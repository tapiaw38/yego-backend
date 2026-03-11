package order

import (
	"context"
	"time"

	"yego/internal/platform/appcontext"
	apperrors "yego/internal/platform/errors"
	"yego/internal/platform/errors/mappings"
	"yego/internal/usecases/notification"

	"github.com/google/uuid"
)

type AcceptDeliveryOutput struct {
	Data OrderOutputData `json:"data"`
}

type AcceptDeliveryUsecase interface {
	Execute(ctx context.Context, orderID string, deliveryUserID string) (*AcceptDeliveryOutput, apperrors.ApplicationError)
}

type acceptDeliveryUsecase struct {
	contextFactory  appcontext.Factory
	notificationSvc notification.Service
}

func NewAcceptDeliveryUsecase(contextFactory appcontext.Factory, notificationSvc notification.Service) AcceptDeliveryUsecase {
	return &acceptDeliveryUsecase{
		contextFactory:  contextFactory,
		notificationSvc: notificationSvc,
	}
}

func (u *acceptDeliveryUsecase) Execute(ctx context.Context, orderID string, deliveryUserID string) (*AcceptDeliveryOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if _, err := uuid.Parse(orderID); err != nil {
		return nil, apperrors.NewApplicationError(mappings.OrderInvalidIDError, err)
	}

	order, appErr := app.Repositories.Order.AcceptDelivery(ctx, orderID, deliveryUserID)
	if appErr != nil {
		return nil, appErr
	}

	// Notify managers that a delivery user accepted the order
	_ = u.notificationSvc.NotifyManagersDeliveryAccepted(notification.DeliveryAcceptedPayload{
		OrderID:            order.ID,
		DeliveryUserID:     deliveryUserID,
		DeliveryAcceptedAt: time.Now().Format("2006-01-02T15:04:05Z"),
	})

	return &AcceptDeliveryOutput{
		Data: toOrderOutputData(order, false),
	}, nil
}
