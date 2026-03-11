package order

import (
	"context"

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

	// Notify the assigned delivery user about the status change
	if updated.DeliveryUserID != nil && *updated.DeliveryUserID != "" {
		statusMsg := ""
		if updated.StatusMessage != nil {
			statusMsg = *updated.StatusMessage
		}
		_ = u.notificationSvc.NotifyDeliveryUserOrderUpdated(*updated.DeliveryUserID, notification.OrderStatusUpdatedPayload{
			OrderID:       updated.ID,
			Status:        string(updated.Status),
			StatusMessage: statusMsg,
			ETA:           updated.ETA,
		})
	}

	return &UpdateStatusOutput{
		Data: toOrderOutputData(updated, false),
	}, nil
}

