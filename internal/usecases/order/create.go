package order

import (
	"context"
	"fmt"
	"log"
	"time"

	"yego/internal/domain"
	"yego/internal/platform/appcontext"
	apperrors "yego/internal/platform/errors"
	"yego/internal/platform/errors/mappings"
	"yego/internal/usecases/notification"
	settingsUsecase "yego/internal/usecases/settings"
)

// CreateInput represents the input for creating an order
type CreateInput struct {
	ProfileID    string `json:"profile_id" binding:"required"`
	ETA          string `json:"eta"`
	SecurityCode string `json:"security_code"`
	Token        string
}

// CreateOutput represents the output after creating an order
type CreateOutput struct {
	Data OrderOutputData `json:"data"`
}

// CreateUsecase defines the interface for creating orders
type CreateUsecase interface {
	Execute(ctx context.Context, input CreateInput) (*CreateOutput, apperrors.ApplicationError)
}

type createUsecase struct {
	contextFactory          appcontext.Factory
	calculateDeliveryFeeUse settingsUsecase.CalculateDeliveryFeeUsecase
	notificationSvc         notification.Service
}

// NewCreateUsecase creates a new instance of CreateUsecase
func NewCreateUsecase(contextFactory appcontext.Factory, calculateDeliveryFeeUse settingsUsecase.CalculateDeliveryFeeUsecase, notificationSvc notification.Service) CreateUsecase {
	return &createUsecase{
		contextFactory:          contextFactory,
		calculateDeliveryFeeUse: calculateDeliveryFeeUse,
		notificationSvc:         notificationSvc,
	}
}

// Execute creates a new order and processes payment immediately
func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	newOrder := &domain.Order{
		ProfileID: &input.ProfileID,
		ETA:       input.ETA,
		Status:    domain.StatusCreated,
	}

	created, err := app.Repositories.Order.Create(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	// Process payment immediately if security code is provided
	if input.SecurityCode != "" {
		paymentErr := ProcessPaymentForOrder(ctx, app, created, input.Token, input.SecurityCode, u.calculateDeliveryFeeUse)
		if paymentErr != nil {
			log.Printf("Payment failed for order %s: %v", created.ID, paymentErr)
			return nil, apperrors.NewApplicationError(mappings.OrderPaymentFailedError, fmt.Errorf("payment failed: %w", paymentErr))
		}
		created.Status = domain.StatusConfirmed
		created, _ = app.Repositories.Order.Update(ctx, created)
	}

	// Notify managers of the new order
	if u.notificationSvc != nil {
		profileID := ""
		if created.ProfileID != nil {
			profileID = *created.ProfileID
		}
		payload := notification.OrderCreatedPayload{
			OrderID:   created.ID,
			ProfileID: profileID,
			Status:    string(created.Status),
			ETA:       created.ETA,
			CreatedAt: time.Now().Format("2006-01-02T15:04:05Z"),
		}
		go func() {
			if notifyErr := u.notificationSvc.NotifyOrderCreated(payload); notifyErr != nil {
				log.Printf("Warning: failed to notify order created %s: %v", created.ID, notifyErr)
			}
		}()
	}

	return &CreateOutput{
		Data: toOrderOutputData(created, false),
	}, nil
}
