package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type NotificationType string

const (
	OrderClaimedNotification            NotificationType = "order_claimed"
	OrderUpdatedNotification            NotificationType = "order_updated"
	OrderAssignedToDeliveryNotification NotificationType = "order_assigned_to_delivery"
	DeliveryAcceptedNotification        NotificationType = "delivery_accepted"
	DeliveryLocationUpdatedNotification NotificationType = "delivery_location_updated"
	OrderStatusUpdatedNotification      NotificationType = "order_status_updated"
)

type Notification struct {
	Type    NotificationType `json:"type"`
	Payload interface{}      `json:"payload"`
}

type OrderClaimedPayload struct {
	OrderID   string `json:"order_id"`
	UserID    string `json:"user_id"`
	ProfileID string `json:"profile_id,omitempty"`
	Status    string `json:"status"`
	ETA       string `json:"eta"`
	ClaimedAt string `json:"claimed_at"`
}

type OrderAssignedToDeliveryPayload struct {
	OrderID   string `json:"order_id"`
	Status    string `json:"status"`
	ETA       string `json:"eta"`
	AssignedAt string `json:"assigned_at"`
}

type DeliveryAcceptedPayload struct {
	OrderID            string `json:"order_id"`
	DeliveryUserID     string `json:"delivery_user_id"`
	DeliveryAcceptedAt string `json:"delivery_accepted_at"`
}

type OrderStatusUpdatedPayload struct {
	OrderID       string `json:"order_id"`
	Status        string `json:"status"`
	StatusMessage string `json:"status_message,omitempty"`
	ETA           string `json:"eta"`
}

type DeliveryLocationPayload struct {
	OrderID   string  `json:"order_id"`
	UserID    string  `json:"user_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
}

type Client struct {
	Hub        *Hub
	Conn       *websocket.Conn
	Send       chan []byte
	IsManager  bool
	IsDelivery bool
	UserID     string
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("WebSocket client registered (manager=%v, delivery=%v). Total: %d", client.IsManager, client.IsDelivery, len(h.clients))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket client unregistered. Total: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				if client.IsManager {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(h.clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) BroadcastNotification(notification Notification) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	h.broadcast <- data
	return nil
}

func (h *Hub) NotifyOrderClaimed(payload OrderClaimedPayload) error {
	return h.BroadcastNotification(Notification{
		Type:    OrderClaimedNotification,
		Payload: payload,
	})
}

// NotifyDeliveryUsers sends a notification to all connected delivery users.
func (h *Hub) NotifyDeliveryUsers(payload OrderAssignedToDeliveryPayload) error {
	data, err := json.Marshal(Notification{
		Type:    OrderAssignedToDeliveryNotification,
		Payload: payload,
	})
	if err != nil {
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.IsDelivery {
			select {
			case client.Send <- data:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
	return nil
}

// NotifyManagers sends a notification only to manager clients.
func (h *Hub) NotifyManagers(notification Notification) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.IsManager {
			select {
			case client.Send <- data:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
	return nil
}

// SendToUser sends a notification to a specific user by userID.
func (h *Hub) SendToUser(userID string, notification Notification) error {
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.UserID == userID {
			select {
			case client.Send <- data:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
	return nil
}

func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
