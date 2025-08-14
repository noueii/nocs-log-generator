package websocket

import (
	"log"
	"sync"
)

// Hub maintains active client connections and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Channel for new client registration
	register chan *Client

	// Channel for client unregistration
	unregister chan *Client

	// Channel for broadcasting messages to all clients
	broadcast chan []byte

	// Channel for broadcasting messages to specific match subscribers
	matchBroadcast chan *MatchMessage

	// Map of match ID to subscribed clients
	matchClients map[string]map[*Client]bool

	// Mutex for thread safety
	mu sync.RWMutex

	// Channel to stop the hub
	stop chan struct{}
}

// MatchMessage represents a message targeted at specific match subscribers
type MatchMessage struct {
	MatchID string
	Data    []byte
}

// NewHub creates a new WebSocket hub instance
func NewHub() *Hub {
	return &Hub{
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan []byte),
		matchBroadcast: make(chan *MatchMessage),
		matchClients:   make(map[string]map[*Client]bool),
		stop:           make(chan struct{}),
	}
}

// Run starts the WebSocket hub and handles client management
func (h *Hub) Run() {
	log.Println("WebSocket hub started")
	
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastToAll(message)

		case matchMsg := <-h.matchBroadcast:
			h.broadcastToMatch(matchMsg)

		case <-h.stop:
			log.Println("WebSocket hub stopping")
			return
		}
	}
}

// Stop gracefully shuts down the hub
func (h *Hub) Stop() {
	close(h.stop)
}

// RegisterClient adds a new client to the hub
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient removes a client from the hub
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// BroadcastToAll sends a message to all connected clients
func (h *Hub) BroadcastToAll(message []byte) {
	h.broadcast <- message
}

// BroadcastToMatch sends a message to all clients subscribed to a specific match
func (h *Hub) BroadcastToMatch(matchID string, message []byte) {
	h.matchBroadcast <- &MatchMessage{
		MatchID: matchID,
		Data:    message,
	}
}

// SubscribeToMatch subscribes a client to match-specific messages
func (h *Hub) SubscribeToMatch(client *Client, matchID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.matchClients[matchID] == nil {
		h.matchClients[matchID] = make(map[*Client]bool)
	}
	
	h.matchClients[matchID][client] = true
	client.subscribedMatches[matchID] = true
	
	log.Printf("Client %s subscribed to match %s", client.id, matchID)
}

// UnsubscribeFromMatch unsubscribes a client from match-specific messages
func (h *Hub) UnsubscribeFromMatch(client *Client, matchID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.matchClients[matchID] != nil {
		delete(h.matchClients[matchID], client)
		
		// Clean up empty match subscription map
		if len(h.matchClients[matchID]) == 0 {
			delete(h.matchClients, matchID)
		}
	}
	
	delete(client.subscribedMatches, matchID)
	
	log.Printf("Client %s unsubscribed from match %s", client.id, matchID)
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetMatchSubscribers returns the number of clients subscribed to a match
func (h *Hub) GetMatchSubscribers(matchID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	if matchClients, exists := h.matchClients[matchID]; exists {
		return len(matchClients)
	}
	return 0
}

// registerClient handles client registration
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	h.clients[client] = true
	
	log.Printf("Client %s connected. Total clients: %d", client.id, len(h.clients))
}

// unregisterClient handles client unregistration
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	
	if _, ok := h.clients[client]; ok {
		// Remove client from general clients list
		delete(h.clients, client)
		
		// Remove client from all match subscriptions
		for matchID := range client.subscribedMatches {
			if h.matchClients[matchID] != nil {
				delete(h.matchClients[matchID], client)
				
				// Clean up empty match subscription map
				if len(h.matchClients[matchID]) == 0 {
					delete(h.matchClients, matchID)
				}
			}
		}
		
		// Close client's send channel
		close(client.send)
		
		log.Printf("Client %s disconnected. Total clients: %d", client.id, len(h.clients))
	}
}

// broadcastToAll sends a message to all connected clients
func (h *Hub) broadcastToAll(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			// Client's send channel is full or closed
			// Remove client and close channel
			delete(h.clients, client)
			close(client.send)
		}
	}
}

// broadcastToMatch sends a message to clients subscribed to a specific match
func (h *Hub) broadcastToMatch(matchMsg *MatchMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	matchClients, exists := h.matchClients[matchMsg.MatchID]
	if !exists {
		log.Printf("No clients subscribed to match %s", matchMsg.MatchID)
		return
	}
	
	for client := range matchClients {
		select {
		case client.send <- matchMsg.Data:
		default:
			// Client's send channel is full or closed
			// Remove client from match subscription
			delete(matchClients, client)
			delete(client.subscribedMatches, matchMsg.MatchID)
			
			// Clean up empty match subscription map
			if len(matchClients) == 0 {
				delete(h.matchClients, matchMsg.MatchID)
			}
		}
	}
}