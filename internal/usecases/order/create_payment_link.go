package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"

	"yego/internal/adapters/web/integrations/payments"
	"yego/internal/platform/appcontext"
	apperrors "yego/internal/platform/errors"
	"yego/internal/platform/errors/mappings"
	settingsUsecase "yego/internal/usecases/settings"
)

// CreatePaymentLinkInput represents the input for creating a payment link
type CreatePaymentLinkInput struct {
	OrderID     string
	UserID      string
	AuthToken   string
	FrontendURL string
	BackendURL  string
}

// CreatePaymentLinkOutput represents the output with the payment link
type CreatePaymentLinkOutput struct {
	InitPoint       string `json:"init_point"`
	SandboxInitPoint string `json:"sandbox_init_point"`
}

// CreatePaymentLinkUsecase defines the interface for creating a payment link
type CreatePaymentLinkUsecase interface {
	Execute(ctx context.Context, input CreatePaymentLinkInput) (*CreatePaymentLinkOutput, apperrors.ApplicationError)
}

type createPaymentLinkUsecase struct {
	contextFactory          appcontext.Factory
	calculateDeliveryFeeUse settingsUsecase.CalculateDeliveryFeeUsecase
}

// NewCreatePaymentLinkUsecase creates a new instance of CreatePaymentLinkUsecase
func NewCreatePaymentLinkUsecase(contextFactory appcontext.Factory, calculateDeliveryFeeUse settingsUsecase.CalculateDeliveryFeeUsecase) CreatePaymentLinkUsecase {
	return &createPaymentLinkUsecase{
		contextFactory:          contextFactory,
		calculateDeliveryFeeUse: calculateDeliveryFeeUse,
	}
}

// Execute creates a MercadoPago Checkout Pro preference and returns the payment link
func (u *createPaymentLinkUsecase) Execute(ctx context.Context, input CreatePaymentLinkInput) (*CreatePaymentLinkOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	order, err := app.Repositories.Order.GetByID(ctx, input.OrderID)
	if err != nil {
		return nil, err
	}

	// Verify the user owns this order
	if order.UserID == nil || *order.UserID != input.UserID {
		return nil, apperrors.NewApplicationError(mappings.UnauthorizedError, errors.New("order does not belong to this user"))
	}

	// Only allow payment when order is in CREATED status
	if string(order.Status) != "CREATED" {
		return nil, apperrors.NewApplicationError(mappings.OrderAlreadyAssignedError, errors.New("order has already been paid"))
	}

	if order.ProfileID == nil {
		return nil, apperrors.NewApplicationError(mappings.OrderPaymentFailedError, errors.New("order has no profile"))
	}

	profile, profileErr := app.Repositories.Profile.GetByID(ctx, *order.ProfileID)
	if profileErr != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderPaymentFailedError, fmt.Errorf("failed to get profile: %w", profileErr))
	}

	// Get payer email
	var payerEmail string
	if input.AuthToken != "" {
		var emailErr error
		payerEmail, emailErr = app.Integrations.Auth.GetUserEmail(profile.UserID, input.AuthToken)
		if emailErr != nil {
			log.Printf("Warning: Failed to get user email for user %s: %v", profile.UserID, emailErr)
			payerEmail = fmt.Sprintf("%s@yego.local", profile.UserID)
		}
	} else {
		payerEmail = fmt.Sprintf("%s@yego.local", profile.UserID)
	}

	// Calculate items subtotal
	var itemsTotal float64
	if order.Data != nil {
		for _, item := range order.Data.Items {
			itemsTotal += item.Price * float64(item.Quantity)
		}
	}

	// Calculate delivery fee separately
	var deliveryFee float64
	if profile.LocationID != nil {
		location, locErr := app.Repositories.Profile.GetLocationByID(ctx, *profile.LocationID)
		if locErr == nil && location != nil {
			deliveryInput := settingsUsecase.CalculateDeliveryFeeInput{
				UserLatitude:  location.Latitude,
				UserLongitude: location.Longitude,
				Items: make([]struct {
					Quantity int  `json:"quantity"`
					Weight   *int `json:"weight,omitempty"`
				}, len(order.Data.Items)),
			}
			if order.Data != nil {
				for i, item := range order.Data.Items {
					deliveryInput.Items[i].Quantity = item.Quantity
					deliveryInput.Items[i].Weight = item.Weight
				}
			}
			if feeOutput, feeErr := u.calculateDeliveryFeeUse.Execute(ctx, deliveryInput); feeErr == nil && feeOutput != nil {
				deliveryFee = feeOutput.TotalPrice
			}
		}
	}

	orderTotal := math.Round((itemsTotal+deliveryFee)*100) / 100
	if orderTotal <= 0 {
		return nil, apperrors.NewApplicationError(mappings.OrderPaymentFailedError, errors.New("order total is zero or negative"))
	}

	// Build preference items: products + delivery fee if applicable
	var prefItems []payments.PreferenceItem
	if order.Data != nil && len(order.Data.Items) > 0 {
		for _, item := range order.Data.Items {
			prefItems = append(prefItems, payments.PreferenceItem{
				Title:      item.Name,
				Quantity:   item.Quantity,
				UnitPrice:  math.Round(item.Price*100) / 100,
				CurrencyID: "ARS",
			})
		}
	}
	if deliveryFee > 0 {
		prefItems = append(prefItems, payments.PreferenceItem{
			Title:      "Envío",
			Quantity:   1,
			UnitPrice:  math.Round(deliveryFee*100) / 100,
			CurrencyID: "ARS",
		})
	}
	if len(prefItems) == 0 {
		prefItems = append(prefItems, payments.PreferenceItem{
			Title:      fmt.Sprintf("Pedido %s", order.ID),
			Quantity:   1,
			UnitPrice:  orderTotal,
			CurrencyID: "ARS",
		})
	}

	frontendURL := input.FrontendURL
	if frontendURL == "" {
		frontendURL = "https://yego.app"
	}
	orderURL := fmt.Sprintf("%s/order/%s", frontendURL, order.ID)

	notificationURL := ""
	if input.BackendURL != "" {
		notificationURL = fmt.Sprintf("%s/api/orders/webhook/mp", input.BackendURL)
	}

	prefResp, prefErr := app.Integrations.Payments.CreatePreference(
		prefItems,
		payerEmail,
		order.ID,
		orderURL,
		orderURL,
		orderURL,
		notificationURL,
	)
	if prefErr != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderPaymentFailedError, fmt.Errorf("failed to create preference: %w", prefErr))
	}

	return &CreatePaymentLinkOutput{
		InitPoint:       prefResp.InitPoint,
		SandboxInitPoint: prefResp.SandboxInitPoint,
	}, nil
}
