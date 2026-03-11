package admin

import (
	"context"
	"time"

	"yego/internal/platform/appcontext"
	apperrors "yego/internal/platform/errors"
	"yego/internal/platform/errors/mappings"
	"yego/internal/usecases/notification"

	"github.com/google/uuid"
)

type AssignDeliveryInput struct {
	DeliveryUserID string `json:"delivery_user_id"`
}

type AssignDeliveryUsecase interface {
	Execute(ctx context.Context, orderID string, input AssignDeliveryInput) (*OrderOutput, apperrors.ApplicationError)
}

type assignDeliveryUsecase struct {
	contextFactory  appcontext.Factory
	notificationSvc notification.Service
}

func NewAssignDeliveryUsecase(contextFactory appcontext.Factory, notificationSvc notification.Service) AssignDeliveryUsecase {
	return &assignDeliveryUsecase{
		contextFactory:  contextFactory,
		notificationSvc: notificationSvc,
	}
}

func (u *assignDeliveryUsecase) Execute(ctx context.Context, orderID string, input AssignDeliveryInput) (*OrderOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if _, err := uuid.Parse(orderID); err != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderInvalidIDError, err)
	}

	if input.DeliveryUserID == "" {
		return nil, apperrors.NewApplicationError(mappings.RequestBodyParsingError, nil)
	}

	order, appErr := app.Repositories.Order.AssignDelivery(ctx, orderID, input.DeliveryUserID)
	if appErr != nil {
		return nil, appErr
	}

	// Notify all connected delivery users about the new available order
	_ = u.notificationSvc.NotifyDeliveryUsers(notification.OrderAssignedToDeliveryPayload{
		OrderID:    order.ID,
		Status:     string(order.Status),
		ETA:        order.ETA,
		AssignedAt: time.Now().Format("2006-01-02T15:04:05Z"),
	})

	output := toOrderOutput(order)
	return &output, nil
}
