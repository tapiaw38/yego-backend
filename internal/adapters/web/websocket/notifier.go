package websocket

import (
	"yego/internal/usecases/notification"
)

type Notifier struct {
	hub *Hub
}

func NewNotifier(hub *Hub) *Notifier {
	return &Notifier{hub: hub}
}

func (n *Notifier) NotifyOrderClaimed(payload notification.OrderClaimedPayload) error {
	return n.hub.NotifyOrderClaimed(OrderClaimedPayload{
		OrderID:   payload.OrderID,
		UserID:    payload.UserID,
		ProfileID: payload.ProfileID,
		Status:    payload.Status,
		ETA:       payload.ETA,
		ClaimedAt: payload.ClaimedAt,
	})
}

func (n *Notifier) NotifyDeliveryUsers(payload notification.OrderAssignedToDeliveryPayload) error {
	return n.hub.NotifyDeliveryUsers(OrderAssignedToDeliveryPayload{
		OrderID:    payload.OrderID,
		Status:     payload.Status,
		ETA:        payload.ETA,
		AssignedAt: payload.AssignedAt,
	})
}

func (n *Notifier) NotifyManagersDeliveryAccepted(payload notification.DeliveryAcceptedPayload) error {
	return n.hub.NotifyManagers(Notification{
		Type: DeliveryAcceptedNotification,
		Payload: DeliveryAcceptedPayload{
			OrderID:            payload.OrderID,
			DeliveryUserID:     payload.DeliveryUserID,
			DeliveryAcceptedAt: payload.DeliveryAcceptedAt,
		},
	})
}

func (n *Notifier) NotifyDeliveryUserOrderUpdated(deliveryUserID string, payload notification.OrderStatusUpdatedPayload) error {
	return n.hub.SendToUser(deliveryUserID, Notification{
		Type: OrderStatusUpdatedNotification,
		Payload: OrderStatusUpdatedPayload{
			OrderID:       payload.OrderID,
			Status:        payload.Status,
			StatusMessage: payload.StatusMessage,
			ETA:           payload.ETA,
		},
	})
}

var _ notification.Service = (*Notifier)(nil)
