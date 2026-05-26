package core

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
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

// PublicProjectPayload is the safe, public subset of project data
// broadcast to all connected clients (no private/customer fields).
type PublicProjectPayload struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	ClientName      string `json:"client_name"`
	CompanyName     string `json:"company_name"`
	SiteType        string `json:"site_type"`
	PackageTier     string `json:"package_tier"`
	Timeline        string `json:"timeline"`
	BudgetCents     int64  `json:"budget_cents"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
}

// WSClient represents a single WebSocket connection.
type WSClient struct {
	hub    *WSHub
	conn   *websocket.Conn
	userID string // empty string for unauthenticated/anonymous clients
	send   chan []byte
	done   chan struct{}
}

// WSHub manages all active WebSocket connections and broadcasts events.
type WSHub struct {
	mu            sync.RWMutex
	clients       map[*WSClient]bool
	allowedOrigins []string // list of allowed origin domains
}

// NewWSHub creates a new WebSocket hub.
func NewWSHub() *WSHub {
	return &WSHub{
		clients:        make(map[*WSClient]bool),
		allowedOrigins: make([]string, 0),
	}
}

// SetAllowedOrigins configures the list of allowed CORS origins.
func (h *WSHub) SetAllowedOrigins(origins []string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.allowedOrigins = origins
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

// BroadcastPublic sends an event to ALL connected clients (anonymous + authenticated).
// Only use this for public-safe payloads that contain no private user/project data.
func (h *WSHub) BroadcastPublic(event WSEvent) {
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
			log.Printf("[ws] dropping message for slow client")
		}
	}
}

// BroadcastToUser sends an event only to clients authenticated as the given userID.
// If userID is empty, no clients receive the message.
func (h *WSHub) BroadcastToUser(userID string, event WSEvent) {
	if userID == "" {
		return
	}
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("[ws] marshal error: %v", err)
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.userID != userID {
			continue
		}
		select {
		case client.send <- data:
		default:
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

// HandleWebSocket upgrades an HTTP connection to WebSocket, authenticates
// the client via token query param or Authorization header, and registers it.
//
// The token can be provided as:
//   - Query parameter: /api/ws?token=<bearer_token>
//   - Authorization header on the initial HTTP upgrade request
//
// Unauthenticated clients are registered as anonymous (empty userID) and
// will only receive public broadcast events.
func (h *WSHub) HandleWebSocket(store *Store, w http.ResponseWriter, r *http.Request) {
	// --- Authentication ---
	token := r.URL.Query().Get("token")
	if token == "" {
		token = r.Header.Get("Authorization")
	}

	var userID string
	if token != "" {
		if user, ok := store.UserByToken(token); ok {
			userID = user.ID
		}
		// Invalid tokens are silently treated as anonymous; we don't reject
		// the connection so anonymous users can still receive public events.
	}

	// --- Origin check ---
	upgrader := websocket.Upgrader{
		CheckOrigin: h.checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ws] upgrade error: %v", err)
		return
	}

	client := &WSClient{
		hub:    h,
		conn:   conn,
		userID: userID,
		send:   make(chan []byte, 256),
		done:   make(chan struct{}),
	}

	if userID != "" {
		log.Printf("[ws] authenticated client: user=%s", userID)
	} else {
		log.Printf("[ws] anonymous client connected")
	}

	h.Register(client)

	go client.writePump()
	go client.readPump()
}

// checkOrigin validates the Origin header against allowed origins.
// It parses the Origin URL and compares the hostname (not a substring match)
// to prevent domain-squatting attacks (e.g. mergeos.shop.evil.com).
// If no origins are configured, falls back to allowing the request
// (dev-mode behaviour). In production, configure allowed origins
// via PrimaryDomain, AdminDomain, and ScanDomain config values.
func (h *WSHub) checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}

	h.mu.RLock()
	origins := h.allowedOrigins
	h.mu.RUnlock()

	if len(origins) == 0 {
		// No origins configured — allow all (dev mode).
		return true
	}

	// Parse the origin as a URL to extract just the hostname.
	u, err := url.Parse(origin)
	if err != nil {
		log.Printf("[ws] failed to parse origin: %s, err: %v", origin, err)
		return false
	}
	host := u.Hostname()

	for _, allowed := range origins {
		if host == allowed || strings.HasSuffix(host, "."+allowed) {
			return true
		}
	}

	log.Printf("[ws] origin not allowed: %s", origin)
	return false
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
