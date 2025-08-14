package websocket

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Client configuration constants
const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// WebSocket upgrader configuration
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin for development
		// TODO: Restrict origins in production
		return true
	},
}

// Client represents a WebSocket client connection
type Client struct {
	// Unique client identifier
	id string

	// The websocket connection
	conn *websocket.Conn

	// The hub this client belongs to
	hub *Hub

	// Buffered channel of outbound messages
	send chan []byte

	// Map of subscribed match IDs
	subscribedMatches map[string]bool
}

// Message types for WebSocket communication
type MessageType string

const (
	MessageTypeSubscribe   MessageType = "subscribe"
	MessageTypeUnsubscribe MessageType = "unsubscribe"
	MessageTypeEvent       MessageType = "event"
	MessageTypeStatus      MessageType = "status"
	MessageTypeError       MessageType = "error"
	MessageTypePing        MessageType = "ping"
	MessageTypePong        MessageType = "pong"
)

// IncomingMessage represents messages received from clients
type IncomingMessage struct {
	Type    MessageType `json:"type"`
	MatchID string      `json:"match_id,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// OutgoingMessage represents messages sent to clients
type OutgoingMessage struct {
	Type      MessageType `json:"type"`
	MatchID   string      `json:"match_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewClient creates a new WebSocket client
func NewClient(conn *websocket.Conn, hub *Hub, clientID string) *Client {
	return &Client{
		id:                clientID,
		conn:              conn,
		hub:               hub,
		send:              make(chan []byte, 256),
		subscribedMatches: make(map[string]bool),
	}
}

// Start begins the client's read and write pumps
func (c *Client) Start() {
	// Register client with hub
	c.hub.RegisterClient(c)

	// Start goroutines for reading and writing
	go c.writePump()
	go c.readPump()
}

// readPump pumps messages from the websocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.UnregisterClient(c)
		c.conn.Close()
	}()

	// Set read deadline and limits
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for client %s: %v", c.id, err)
			}
			break
		}

		// Trim whitespace from message
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))

		// Parse and handle the message
		c.handleMessage(message)
	}
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to current message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming messages from the client
func (c *Client) handleMessage(message []byte) {
	var inMsg IncomingMessage
	if err := json.Unmarshal(message, &inMsg); err != nil {
		log.Printf("Error parsing message from client %s: %v", c.id, err)
		c.sendError("Invalid message format")
		return
	}

	switch inMsg.Type {
	case MessageTypeSubscribe:
		if inMsg.MatchID != "" {
			c.hub.SubscribeToMatch(c, inMsg.MatchID)
			c.sendStatus("subscribed", map[string]string{"match_id": inMsg.MatchID})
		} else {
			c.sendError("Missing match_id for subscription")
		}

	case MessageTypeUnsubscribe:
		if inMsg.MatchID != "" {
			c.hub.UnsubscribeFromMatch(c, inMsg.MatchID)
			c.sendStatus("unsubscribed", map[string]string{"match_id": inMsg.MatchID})
		} else {
			c.sendError("Missing match_id for unsubscription")
		}

	case MessageTypePing:
		c.sendMessage(MessageTypePong, "", "pong")

	default:
		log.Printf("Unknown message type '%s' from client %s", inMsg.Type, c.id)
		c.sendError("Unknown message type")
	}
}

// sendMessage sends a message to the client
func (c *Client) sendMessage(msgType MessageType, matchID string, data interface{}) {
	message := OutgoingMessage{
		Type:      msgType,
		MatchID:   matchID,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message for client %s: %v", c.id, err)
		return
	}

	select {
	case c.send <- messageBytes:
	default:
		// Channel is full, remove client
		close(c.send)
	}
}

// sendError sends an error message to the client
func (c *Client) sendError(errorMsg string) {
	c.sendMessage(MessageTypeError, "", map[string]string{"error": errorMsg})
}

// sendStatus sends a status message to the client
func (c *Client) sendStatus(status string, data interface{}) {
	statusData := map[string]interface{}{
		"status": status,
	}
	
	if data != nil {
		statusData["data"] = data
	}
	
	c.sendMessage(MessageTypeStatus, "", statusData)
}

// SendEvent sends an event message for a specific match to the client
func (c *Client) SendEvent(matchID string, event interface{}) {
	c.sendMessage(MessageTypeEvent, matchID, event)
}

// IsSubscribedToMatch checks if the client is subscribed to a match
func (c *Client) IsSubscribedToMatch(matchID string) bool {
	return c.subscribedMatches[matchID]
}

// GetSubscribedMatches returns a slice of match IDs the client is subscribed to
func (c *Client) GetSubscribedMatches() []string {
	matches := make([]string, 0, len(c.subscribedMatches))
	for matchID := range c.subscribedMatches {
		matches = append(matches, matchID)
	}
	return matches
}

// Close closes the client connection
func (c *Client) Close() {
	c.conn.Close()
}