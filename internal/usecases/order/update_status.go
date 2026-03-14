package order

import (
	"context"
	"log"

	"yego/internal/domain"
	"yego/internal/platform/appcontext"
	apperrors "yego/internal/platform/errors"
	"yego/internal/platform/errors/mappings"
	"yego/internal/usecases/notification"
	settingsUsecase "yego/internal/usecases/settings"

	"github.com/google/uuid"
)

// UpdateStatusInput represents the input for updating order status
type UpdateStatusInput struct {
	Status string `json:"status" binding:"required"`
	Token  string
}

// UpdateStatusOutput represents the output after updating order status
type UpdateStatusOutput struct {
	Data OrderOutputData `json:"data"`
}

// UpdateStatusUsecase defines the interface for updating order status
type UpdateStatusUsecase interface {
	Execute(ctx context.Context, id string, input UpdateStatusInput) (*UpdateStatusOutput, apperrors.ApplicationError)
}

type updateStatusUsecase struct {
	contextFactory          appcontext.Factory
	calculateDeliveryFeeUse settingsUsecase.CalculateDeliveryFeeUsecase
	notificationSvc         notification.Service
}

// NewUpdateStatusUsecase creates a new instance of UpdateStatusUsecase
func NewUpdateStatusUsecase(contextFactory appcontext.Factory, calculateDeliveryFeeUse settingsUsecase.CalculateDeliveryFeeUsecase, notificationSvc notification.Service) UpdateStatusUsecase {
	return &updateStatusUsecase{
		contextFactory:          contextFactory,
		calculateDeliveryFeeUse: calculateDeliveryFeeUse,
		notificationSvc:         notificationSvc,
	}
}

// Execute updates the status of an order
func (u *updateStatusUsecase) Execute(ctx context.Context, id string, input UpdateStatusInput) (*UpdateStatusOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if _, err := uuid.Parse(id); err != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderInvalidIDError, err)
	}

	if !domain.IsValidStatus(input.Status) {
		return nil, apperrors.NewApplicationError(mappings.OrderInvalidStatusError, nil)
	}

	updated, err := app.Repositories.Order.UpdateStatus(ctx, id, domain.OrderStatus(input.Status))
	if err != nil {
		return nil, err
	}

	// Notify managers of the status change
	if u.notificationSvc != nil {
		payload := notification.OrderUpdatedPayload{
			OrderID: updated.ID,
			Status:  string(updated.Status),
			ETA:     updated.ETA,
		}
		go func() {
			if notifyErr := u.notificationSvc.NotifyOrderUpdated(payload); notifyErr != nil {
				log.Printf("Warning: failed to notify order updated %s: %v", updated.ID, notifyErr)
			}
		}()
	}

	return &UpdateStatusOutput{
		Data: toOrderOutputData(updated, false),
	}, nil
}
