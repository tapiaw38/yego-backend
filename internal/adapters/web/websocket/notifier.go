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

func (n *Notifier) NotifyOrderCreated(payload notification.OrderCreatedPayload) error {
	return n.hub.NotifyOrderCreated(OrderCreatedPayload{
		OrderID:   payload.OrderID,
		ProfileID: payload.ProfileID,
		Status:    payload.Status,
		ETA:       payload.ETA,
		CreatedAt: payload.CreatedAt,
	})
}

func (n *Notifier) NotifyOrderUpdated(payload notification.OrderUpdatedPayload) error {
	return n.hub.NotifyOrderUpdated(OrderUpdatedPayload{
		OrderID: payload.OrderID,
		Status:  payload.Status,
		ETA:     payload.ETA,
	})
}

var _ notification.Service = (*Notifier)(nil)
