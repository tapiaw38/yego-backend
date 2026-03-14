package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type NotificationType string

const (
	OrderClaimedNotification NotificationType = "order_claimed"
	OrderCreatedNotification NotificationType = "order_created"
	OrderUpdatedNotification NotificationType = "order_updated"
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

type OrderCreatedPayload struct {
	OrderID   string `json:"order_id"`
	ProfileID string `json:"profile_id,omitempty"`
	Status    string `json:"status"`
	ETA       string `json:"eta"`
	CreatedAt string `json:"created_at"`
}

type OrderUpdatedPayload struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
	ETA     string `json:"eta"`
}

type Client struct {
	Hub       *Hub
	Conn      *websocket.Conn
	Send      chan []byte
	IsManager bool
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
			log.Printf("WebSocket client registered. Total clients: %d", len(h.clients))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			log.Printf("WebSocket client unregistered. Total clients: %d", len(h.clients))

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
	return h.BroadcastNotification(Notification{Type: OrderClaimedNotification, Payload: payload})
}

func (h *Hub) NotifyOrderCreated(payload OrderCreatedPayload) error {
	return h.BroadcastNotification(Notification{Type: OrderCreatedNotification, Payload: payload})
}

func (h *Hub) NotifyOrderUpdated(payload OrderUpdatedPayload) error {
	return h.BroadcastNotification(Notification{Type: OrderUpdatedNotification, Payload: payload})
}

func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
