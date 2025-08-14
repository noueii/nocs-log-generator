package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Manager manages WebSocket connections and message broadcasting
type Manager struct {
	hub *Hub
}

// NewManager creates a new WebSocket manager
func NewManager() *Manager {
	hub := NewHub()
	go hub.Run() // Start the hub in a goroutine
	
	return &Manager{
		hub: hub,
	}
}

// GetHub returns the underlying hub for direct access
func (m *Manager) GetHub() *Hub {
	return m.hub
}

// HandleWebSocketUpgrade handles WebSocket connection upgrades
func (m *Manager) HandleWebSocketUpgrade(c *gin.Context) {
	// Generate unique client ID
	clientID := generateClientID()
	
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "WebSocket upgrade failed",
		})
		return
	}
	
	// Create new client and start it
	client := NewClient(conn, m.hub, clientID)
	client.Start()
	
	log.Printf("WebSocket connection established for client %s from %s", 
		clientID, c.ClientIP())
}

// BroadcastMatchEvent broadcasts an event to all clients subscribed to a match
func (m *Manager) BroadcastMatchEvent(matchID string, eventType string, data interface{}) error {
	event := MatchEvent{
		Type:      eventType,
		MatchID:   matchID,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
	
	message, err := json.Marshal(OutgoingMessage{
		Type:      MessageTypeEvent,
		MatchID:   matchID,
		Data:      event,
		Timestamp: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal match event: %w", err)
	}
	
	m.hub.BroadcastToMatch(matchID, message)
	return nil
}

// BroadcastMatchStatus broadcasts a status update to all clients subscribed to a match
func (m *Manager) BroadcastMatchStatus(matchID string, status string, data interface{}) error {
	statusUpdate := MatchStatus{
		Status:    status,
		MatchID:   matchID,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
	
	message, err := json.Marshal(OutgoingMessage{
		Type:      MessageTypeStatus,
		MatchID:   matchID,
		Data:      statusUpdate,
		Timestamp: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal match status: %w", err)
	}
	
	m.hub.BroadcastToMatch(matchID, message)
	return nil
}

// BroadcastMatchError broadcasts an error to all clients subscribed to a match
func (m *Manager) BroadcastMatchError(matchID string, errorMsg string) error {
	errorData := MatchError{
		Error:     errorMsg,
		MatchID:   matchID,
		Timestamp: time.Now().UTC(),
	}
	
	message, err := json.Marshal(OutgoingMessage{
		Type:      MessageTypeError,
		MatchID:   matchID,
		Data:      errorData,
		Timestamp: time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal match error: %w", err)
	}
	
	m.hub.BroadcastToMatch(matchID, message)
	return nil
}

// GetConnectionStats returns statistics about WebSocket connections
func (m *Manager) GetConnectionStats() ConnectionStats {
	return ConnectionStats{
		TotalClients:   m.hub.GetClientCount(),
		ActiveMatches:  len(m.hub.matchClients),
		Timestamp:      time.Now().UTC(),
	}
}

// GetMatchStats returns statistics for a specific match
func (m *Manager) GetMatchStats(matchID string) MatchStats {
	return MatchStats{
		MatchID:     matchID,
		Subscribers: m.hub.GetMatchSubscribers(matchID),
		Timestamp:   time.Now().UTC(),
	}
}

// Shutdown gracefully shuts down the WebSocket manager
func (m *Manager) Shutdown() {
	log.Println("Shutting down WebSocket manager")
	m.hub.Stop()
}

// Event and message structures for WebSocket communication

// MatchEvent represents an event that occurred in a match
type MatchEvent struct {
	Type      string      `json:"type"`
	MatchID   string      `json:"match_id"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// MatchStatus represents a status update for a match
type MatchStatus struct {
	Status    string      `json:"status"`
	MatchID   string      `json:"match_id"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// MatchError represents an error for a match
type MatchError struct {
	Error     string    `json:"error"`
	MatchID   string    `json:"match_id"`
	Timestamp time.Time `json:"timestamp"`
}

// ConnectionStats provides statistics about WebSocket connections
type ConnectionStats struct {
	TotalClients  int       `json:"total_clients"`
	ActiveMatches int       `json:"active_matches"`
	Timestamp     time.Time `json:"timestamp"`
}

// MatchStats provides statistics for a specific match
type MatchStats struct {
	MatchID     string    `json:"match_id"`
	Subscribers int       `json:"subscribers"`
	Timestamp   time.Time `json:"timestamp"`
}

// generateClientID generates a unique client identifier
func generateClientID() string {
	return uuid.New().String()[:8] // Use first 8 characters of UUID for brevity
}

// Event types for match generation streaming
const (
	EventTypeMatchStart      = "match_start"
	EventTypeRoundStart      = "round_start"
	EventTypeRoundEnd        = "round_end"
	EventTypePlayerEvent     = "player_event"
	EventTypeEconomyUpdate   = "economy_update"
	EventTypeMatchProgress   = "match_progress"
	EventTypeMatchComplete   = "match_complete"
	EventTypeMatchError      = "match_error"
	EventTypeGenerationStart = "generation_start"
	EventTypeGenerationEnd   = "generation_end"
)

// Status types for match generation
const (
	StatusGenerating = "generating"
	StatusPaused     = "paused"
	StatusResumed    = "resumed"
	StatusCompleted  = "completed"
	StatusError      = "error"
)

// Predefined event data structures

// GenerationStartEvent represents the start of match generation
type GenerationStartEvent struct {
	MatchID     string    `json:"match_id"`
	Teams       []string  `json:"teams"`
	Map         string    `json:"map"`
	Format      string    `json:"format"`
	MaxRounds   int       `json:"max_rounds"`
	StartedAt   time.Time `json:"started_at"`
}

// MatchProgressEvent represents progress during match generation
type MatchProgressEvent struct {
	MatchID         string  `json:"match_id"`
	CurrentRound    int     `json:"current_round"`
	TotalRounds     int     `json:"total_rounds"`
	EventsGenerated int     `json:"events_generated"`
	Progress        float64 `json:"progress"` // Percentage (0-100)
}

// RoundEvent represents round-specific events
type RoundEvent struct {
	MatchID     string                 `json:"match_id"`
	RoundNumber int                    `json:"round_number"`
	EventType   string                 `json:"event_type"`
	Data        map[string]interface{} `json:"data"`
}

// PlayerEvent represents player-specific events
type PlayerEvent struct {
	MatchID    string                 `json:"match_id"`
	PlayerName string                 `json:"player_name"`
	EventType  string                 `json:"event_type"`
	Data       map[string]interface{} `json:"data"`
}

// EconomyUpdateEvent represents economy updates
type EconomyUpdateEvent struct {
	MatchID string                        `json:"match_id"`
	Round   int                           `json:"round"`
	Economy map[string]map[string]int `json:"economy"` // team -> player -> money
}

// GenerationCompleteEvent represents the completion of match generation
type GenerationCompleteEvent struct {
	MatchID       string        `json:"match_id"`
	TotalRounds   int           `json:"total_rounds"`
	TotalEvents   int           `json:"total_events"`
	Duration      time.Duration `json:"duration"`
	CompletedAt   time.Time     `json:"completed_at"`
	Success       bool          `json:"success"`
}