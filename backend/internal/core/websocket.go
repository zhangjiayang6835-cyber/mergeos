package core

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocket event types
const (
	EventProjectCreated = "project_created"
	EventProjectFunded  = "project_funded"
)

// WSEvent is the payload sent over WebSocket.
type WSEvent struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// WSClient represents a single WebSocket connection.
type WSClient struct {
	hub  *WSHub
	conn *websocket.Conn
	send chan []byte
	done chan struct{}
}

// WSHub manages all active WebSocket connections and broadcasts events.
type WSHub struct {
	mu      sync.RWMutex
	clients map[*WSClient]bool
}

// NewWSHub creates a new WebSocket hub.
func NewWSHub() *WSHub {
	return &WSHub{
		clients: make(map[*WSClient]bool),
	}
}

// Register adds a client to the hub.
func (h *WSHub) Register(client *WSClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = true
}

// Unregister removes a client from the hub.
func (h *WSHub) Unregister(client *WSClient) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
}

// Broadcast sends an event to all connected clients.
func (h *WSHub) Broadcast(event WSEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("[ws] marshal error: %v", err)
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		select {
		case client.send <- data:
		default:
			// Client's send buffer is full; drop message.
			log.Printf("[ws] dropping message for slow client")
		}
	}
}

// ClientCount returns the number of connected clients.
func (h *WSHub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// HandleWebSocket upgrades an HTTP connection to WebSocket and registers it.
func (h *WSHub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for MVP
		},
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ws] upgrade error: %v", err)
		return
	}

	client := &WSClient{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
		done: make(chan struct{}),
	}
	h.Register(client)

	go client.writePump()
	go client.readPump()
}

// writePump pumps messages from the send channel to the WebSocket connection.
func (c *WSClient) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				// Channel was closed by hub.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("[ws] write error: %v", err)
				return
			}
		case <-c.done:
			return
		}
	}
}

// readPump reads messages from the WebSocket connection (used for keepalive).
func (c *WSClient) readPump() {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			// Client disconnected or error
			break
		}
	}
}
