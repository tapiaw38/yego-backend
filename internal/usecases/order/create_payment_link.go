package order

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"

	"wappi/internal/adapters/web/integrations/payments"
	"wappi/internal/platform/appcontext"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
	settingsUsecase "wappi/internal/usecases/settings"
)

// CreatePaymentLinkInput represents the input for creating a payment link
type CreatePaymentLinkInput struct {
	OrderID    string
	UserID     string
	AuthToken  string
	FrontendURL string
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
			payerEmail = fmt.Sprintf("%s@wappi.local", profile.UserID)
		}
	} else {
		payerEmail = fmt.Sprintf("%s@wappi.local", profile.UserID)
	}

	// Calculate order total
	orderTotal, calcErr := calculateOrderTotal(ctx, app, order, profile, u.calculateDeliveryFeeUse)
	if calcErr != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderPaymentFailedError, fmt.Errorf("failed to calculate order total: %w", calcErr))
	}
	if orderTotal <= 0 {
		return nil, apperrors.NewApplicationError(mappings.OrderPaymentFailedError, errors.New("order total is zero or negative"))
	}
	orderTotal = math.Round(orderTotal*100) / 100

	// Build preference items from order items
	var prefItems []payments.PreferenceItem
	if order.Data != nil {
		for _, item := range order.Data.Items {
			prefItems = append(prefItems, payments.PreferenceItem{
				Title:      item.Name,
				Quantity:   item.Quantity,
				UnitPrice:  math.Round(item.Price*100) / 100,
				CurrencyID: "ARS",
			})
		}
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
		frontendURL = "https://wappi.app"
	}
	orderURL := fmt.Sprintf("%s/order/%s", frontendURL, order.ID)

	prefResp, prefErr := app.Integrations.Payments.CreatePreference(
		prefItems,
		payerEmail,
		order.ID,
		orderURL,
		orderURL,
		orderURL,
		"",
	)
	if prefErr != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderPaymentFailedError, fmt.Errorf("failed to create preference: %w", prefErr))
	}

	return &CreatePaymentLinkOutput{
		InitPoint:       prefResp.InitPoint,
		SandboxInitPoint: prefResp.SandboxInitPoint,
	}, nil
}
