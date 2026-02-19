package order

import (
	"context"
	"errors"

	"wappi/internal/platform/appcontext"
	apperrors "wappi/internal/platform/errors"
	"wappi/internal/platform/errors/mappings"
	settingsUsecase "wappi/internal/usecases/settings"
)

// PayForOrderInput represents the input for paying an order
type PayForOrderInput struct {
	OrderID      string
	UserID       string
	AuthToken    string
	SecurityCode string
}

// PayForOrderOutput represents the output after paying an order
type PayForOrderOutput struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

// PayForOrderUsecase defines the interface for paying an order
type PayForOrderUsecase interface {
	Execute(ctx context.Context, input PayForOrderInput) (*PayForOrderOutput, apperrors.ApplicationError)
}

type payForOrderUsecase struct {
	contextFactory          appcontext.Factory
	calculateDeliveryFeeUse settingsUsecase.CalculateDeliveryFeeUsecase
}

// NewPayForOrderUsecase creates a new instance of PayForOrderUsecase
func NewPayForOrderUsecase(contextFactory appcontext.Factory, calculateDeliveryFeeUse settingsUsecase.CalculateDeliveryFeeUsecase) PayForOrderUsecase {
	return &payForOrderUsecase{
		contextFactory:          contextFactory,
		calculateDeliveryFeeUse: calculateDeliveryFeeUse,
	}
}

// Execute processes the payment for an order and moves it from CREATED to CONFIRMED
func (u *payForOrderUsecase) Execute(ctx context.Context, input PayForOrderInput) (*PayForOrderOutput, apperrors.ApplicationError) {
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

	paymentErr := ProcessPaymentForOrder(ctx, app, order, input.AuthToken, input.SecurityCode, u.calculateDeliveryFeeUse)
	if paymentErr != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderPaymentFailedError, paymentErr)
	}

	_, _ = app.Repositories.Order.UpdateStatus(ctx, input.OrderID, "CONFIRMED")

	return &PayForOrderOutput{
		OrderID: input.OrderID,
		Status:  "CONFIRMED",
	}, nil
}
