package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"yego/internal/adapters/web/middlewares"
	"yego/internal/platform/config"
	ws "yego/internal/adapters/web/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub *ws.Hub
}

func NewHandler(hub *ws.Hub) *Handler {
	return &Handler{hub: hub}
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	cfg := config.GetInstance()
	claims, err := middlewares.ValidateToken(token, cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	userID := claims.UserID
	role := c.Query("role") // "manager" | "delivery" — sent by frontend

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := &ws.Client{
		Hub:        h.hub,
		Conn:       conn,
		Send:       make(chan []byte, 256),
		IsManager:  role != "delivery",
		IsDelivery: role == "delivery",
		UserID:     userID,
	}

	h.hub.Register <- client

	go clientWritePump(client)
	go clientReadPump(h.hub, client, userID)

	log.Printf("WebSocket connected: user=%s role=%s", userID, role)
}

func clientWritePump(client *ws.Client) {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				client.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}

		case <-ticker.C:
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type incomingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type locationUpdatePayload struct {
	OrderID   string  `json:"order_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func clientReadPump(hub *ws.Hub, client *ws.Client, userID string) {
	defer func() {
		hub.Unregister <- client
		client.Conn.Close()
		log.Printf("WebSocket connection closed for user: %s", userID)
	}()

	client.Conn.SetReadLimit(4096)
	client.Conn.SetPongHandler(func(string) error {
		return nil
	})

	for {
		_, data, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}

		if !client.IsDelivery {
			continue
		}

		var msg incomingMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			log.Printf("WebSocket: failed to parse message from %s: %v", userID, err)
			continue
		}

		switch msg.Type {
		case "location_update":
			var loc locationUpdatePayload
			if err := json.Unmarshal(msg.Payload, &loc); err != nil {
				log.Printf("WebSocket: invalid location_update payload: %v", err)
				continue
			}
			if err := hub.NotifyManagers(ws.Notification{
				Type: ws.DeliveryLocationUpdatedNotification,
				Payload: ws.DeliveryLocationPayload{
					OrderID:   loc.OrderID,
					UserID:    userID,
					Latitude:  loc.Latitude,
					Longitude: loc.Longitude,
					Timestamp: time.Now().UTC().Format(time.RFC3339),
				},
			}); err != nil {
				log.Printf("WebSocket: failed to broadcast location: %v", err)
			}
		}
	}
}
