package order

import (
	"context"
	"fmt"
	"log"
	"math"

	"yego/internal/adapters/web/integrations/payments"
	"yego/internal/domain"
	"yego/internal/platform/appcontext"
	settingsUsecase "yego/internal/usecases/settings"
)

// ProcessPaymentForOrder processes the payment when an order is delivered.
// It resolves the user's internal UUID, checks for a payment method, calculates
// the order total, charges the user, and records the transaction.
func ProcessPaymentForOrder(ctx context.Context, app *appcontext.Context, order *domain.Order, token string, securityCode string, calculateDeliveryFeeUse settingsUsecase.CalculateDeliveryFeeUsecase) error {
	if order.ProfileID == nil {
		// Profile may have been created after claiming — try to find and assign it now
		if order.UserID != nil {
			userProfile, profileLookupErr := app.Repositories.Profile.GetByUserID(ctx, *order.UserID)
			if profileLookupErr == nil && userProfile != nil {
				_ = app.Repositories.Order.AssignProfile(ctx, order.ID, userProfile.ID)
				order.ProfileID = &userProfile.ID
			}
		}
		if order.ProfileID == nil {
			return fmt.Errorf("order has no profile_id")
		}
	}

	profile, profileErr := app.Repositories.Profile.GetByID(ctx, *order.ProfileID)
	if profileErr != nil {
		return fmt.Errorf("failed to get profile: %w", profileErr)
	}

	// Payment methods are stored under profile.UserID (auth username).
	// Use it directly to avoid the GetUserIDByUsername UUID mismatch.
	paymentUserID := profile.UserID

	hasPaymentMethod, paymentErr := app.Integrations.Payments.HasPaymentMethod(paymentUserID)
	if paymentErr != nil {
		return fmt.Errorf("failed to check payment method: %w", paymentErr)
	}
	if !hasPaymentMethod {
		return fmt.Errorf("user has no payment method configured")
	}

	orderTotal, calcErr := calculateOrderTotal(ctx, app, order, profile, calculateDeliveryFeeUse)
	if calcErr != nil {
		return fmt.Errorf("failed to calculate order total: %w", calcErr)
	}
	if orderTotal <= 0 {
		return fmt.Errorf("order total is zero or negative")
	}
	// Round to 2 decimal places to avoid floating point issues with MercadoPago
	orderTotal = math.Round(orderTotal*100) / 100

	var userEmail string
	if token != "" {
		var emailErr error
		userEmail, emailErr = app.Integrations.Auth.GetUserEmail(profile.UserID, token)
		if emailErr != nil {
			log.Printf("Warning: Failed to get user email for user %s: %v", profile.UserID, emailErr)
			userEmail = fmt.Sprintf("%s@yego.local", profile.UserID)
		}
	} else {
		userEmail = fmt.Sprintf("%s@yego.local", profile.UserID)
		log.Printf("Warning: No token provided, using placeholder email for user %s", profile.UserID)
	}
	if userEmail == "" {
		return fmt.Errorf("user email not found")
	}

	var collectorID string
	settings, _ := app.Repositories.Settings.Get(ctx)
	if settings != nil && settings.ManagerCollectorID != nil {
		collectorID = *settings.ManagerCollectorID
	}

	var paymentResponse *payments.ProcessPaymentResponse
	paymentResponse, paymentErr = app.Integrations.Payments.ProcessPaymentWithSavedMethod(
		paymentUserID,
		orderTotal,
		fmt.Sprintf("Pago por pedido %s", order.ID),
		order.ID,
		userEmail,
		collectorID,
		securityCode,
	)
	if paymentErr != nil {
		return fmt.Errorf("failed to process payment: %w", paymentErr)
	}

	log.Printf("Payment processed for order %s: Payment ID %d, Gateway ID %s, Status %s",
		order.ID, paymentResponse.PaymentID, paymentResponse.GatewayPaymentID, paymentResponse.Status)

	if paymentResponse.Status == "rejected" {
		return fmt.Errorf("payment rejected by gateway (status: %s)", paymentResponse.Status)
	}

	description := fmt.Sprintf("Pago por pedido %s", order.ID)
	transaction := &domain.Transaction{
		OrderID:          order.ID,
		UserID:           paymentUserID,
		ProfileID:        order.ProfileID,
		Amount:           orderTotal,
		Currency:         "ARS",
		Status:           paymentResponse.Status,
		PaymentID:        &paymentResponse.PaymentID,
		GatewayPaymentID: &paymentResponse.GatewayPaymentID,
		CollectorID:      &collectorID,
		Description:      &description,
	}

	_, transErr := app.Repositories.Transaction.Create(ctx, transaction)
	if transErr != nil {
		log.Printf("Warning: Failed to create transaction record for order %s: %v", order.ID, transErr)
	}

	return nil
}

// calculateOrderTotal calculates the total amount for an order (items + delivery fee).
func calculateOrderTotal(ctx context.Context, app *appcontext.Context, order *domain.Order, profile *domain.Profile, calculateDeliveryFeeUse settingsUsecase.CalculateDeliveryFeeUsecase) (float64, error) {
	if order.Data == nil || len(order.Data.Items) == 0 {
		return 0, nil
	}

	var itemsTotal float64
	for _, item := range order.Data.Items {
		itemsTotal += item.Price * float64(item.Quantity)
	}

	var deliveryFee float64
	if profile.LocationID != nil {
		location, err := app.Repositories.Profile.GetLocationByID(ctx, *profile.LocationID)
		if err == nil && location != nil {
			deliveryFeeInput := settingsUsecase.CalculateDeliveryFeeInput{
				UserLatitude:  location.Latitude,
				UserLongitude: location.Longitude,
				Items: make([]struct {
					Quantity int  `json:"quantity"`
					Weight   *int `json:"weight,omitempty"`
				}, len(order.Data.Items)),
			}

			for i, item := range order.Data.Items {
				deliveryFeeInput.Items[i].Quantity = item.Quantity
				deliveryFeeInput.Items[i].Weight = item.Weight
			}

			deliveryFeeOutput, err := calculateDeliveryFeeUse.Execute(ctx, deliveryFeeInput)
			if err == nil && deliveryFeeOutput != nil {
				deliveryFee = deliveryFeeOutput.TotalPrice
			}
		}
	}

	return itemsTotal + deliveryFee, nil
}
