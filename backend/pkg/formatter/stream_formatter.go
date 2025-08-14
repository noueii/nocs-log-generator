package formatter

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/noueii/nocs-log-generator/backend/pkg/models"
)

// StreamFormatter handles real-time log formatting for WebSocket/SSE streaming
type StreamFormatter struct {
	logFormatter *LogFormatter
	httpFormatter *HTTPFormatter
	config       *models.MatchConfig
	
	// Buffering and batching
	buffer          []models.GameEvent
	bufferMutex     sync.RWMutex
	maxBufferSize   int
	batchTimeout    time.Duration
	
	// Stream management
	subscribers     map[string]*StreamSubscriber
	subscriberMutex sync.RWMutex
	
	// Statistics
	eventsSent      int64
	bytesSent       int64
	activeStreams   int
	statsMutex      sync.RWMutex
}

// StreamSubscriber represents a client subscribed to the stream
type StreamSubscriber struct {
	ID           string
	Channel      chan StreamMessage
	Filter       *StreamFilter
	Format       StreamFormat
	IsActive     bool
	ConnectedAt  time.Time
	LastActivity time.Time
}

// StreamMessage represents a message sent to stream subscribers
type StreamMessage struct {
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
	Metadata  interface{} `json:"metadata,omitempty"`
}

// StreamFilter defines filtering options for streams
type StreamFilter struct {
	EventTypes []string `json:"event_types,omitempty"`
	Players    []string `json:"players,omitempty"`
	Teams      []string `json:"teams,omitempty"`
	Rounds     []int    `json:"rounds,omitempty"`
	MinDamage  int      `json:"min_damage,omitempty"`
	KillsOnly  bool     `json:"kills_only,omitempty"`
	ObjectivesOnly bool `json:"objectives_only,omitempty"`
}

// StreamFormat defines the output format for streams
type StreamFormat string

const (
	StreamFormatText StreamFormat = "text"
	StreamFormatJSON StreamFormat = "json"
	StreamFormatSSE  StreamFormat = "sse"  // Server-Sent Events
)

// StreamConfig contains configuration for the stream formatter
type StreamConfig struct {
	MaxBufferSize int           `json:"max_buffer_size"`
	BatchTimeout  time.Duration `json:"batch_timeout"`
	MaxSubscribers int          `json:"max_subscribers"`
	MessageTimeout time.Duration `json:"message_timeout"`
}

// NewStreamFormatter creates a new stream formatter
func NewStreamFormatter(config *models.MatchConfig, streamConfig *StreamConfig) *StreamFormatter {
	// Set default values
	if streamConfig == nil {
		streamConfig = &StreamConfig{
			MaxBufferSize:  1000,
			BatchTimeout:   time.Millisecond * 100,
			MaxSubscribers: 100,
			MessageTimeout: time.Second * 30,
		}
	}
	
	sf := &StreamFormatter{
		logFormatter:  NewLogFormatter(config),
		httpFormatter: NewHTTPFormatter(config),
		config:        config,
		buffer:        make([]models.GameEvent, 0, streamConfig.MaxBufferSize),
		maxBufferSize: streamConfig.MaxBufferSize,
		batchTimeout:  streamConfig.BatchTimeout,
		subscribers:   make(map[string]*StreamSubscriber),
	}
	
	// Start background processing
	go sf.processBuffer()
	go sf.cleanupInactiveSubscribers()
	
	return sf
}

// StreamFormat formats events for real-time streaming
func (sf *StreamFormatter) StreamFormat(events []models.GameEvent, format StreamFormat) ([]string, error) {
	var lines []string
	
	for _, event := range events {
		switch format {
		case StreamFormatText:
			line := sf.logFormatter.FormatEvent(event)
			if line != "" {
				lines = append(lines, line)
			}
			
		case StreamFormatJSON:
			jsonEntry, err := sf.httpFormatter.convertEventToJSON(event)
			if err != nil {
				return nil, fmt.Errorf("error converting event to JSON: %w", err)
			}
			
			jsonBytes, err := json.Marshal(jsonEntry)
			if err != nil {
				return nil, fmt.Errorf("error marshaling JSON: %w", err)
			}
			
			lines = append(lines, string(jsonBytes))
			
		case StreamFormatSSE:
			jsonEntry, err := sf.httpFormatter.convertEventToJSON(event)
			if err != nil {
				return nil, fmt.Errorf("error converting event to JSON: %w", err)
			}
			
			jsonBytes, err := json.Marshal(jsonEntry)
			if err != nil {
				return nil, fmt.Errorf("error marshaling JSON: %w", err)
			}
			
			// Format as SSE message
			sseLine := fmt.Sprintf("data: %s\n\n", string(jsonBytes))
			lines = append(lines, sseLine)
		}
	}
	
	return lines, nil
}

// BufferEvents adds events to the stream buffer for batch processing
func (sf *StreamFormatter) BufferEvents(events ...models.GameEvent) error {
	sf.bufferMutex.Lock()
	defer sf.bufferMutex.Unlock()
	
	// Add events to buffer
	sf.buffer = append(sf.buffer, events...)
	
	// If buffer is full, flush immediately
	if len(sf.buffer) >= sf.maxBufferSize {
		go sf.flushBuffer()
	}
	
	return nil
}

// Subscribe creates a new stream subscription
func (sf *StreamFormatter) Subscribe(subscriberID string, filter *StreamFilter, format StreamFormat) (*StreamSubscriber, error) {
	sf.subscriberMutex.Lock()
	defer sf.subscriberMutex.Unlock()
	
	// Check if subscriber already exists
	if _, exists := sf.subscribers[subscriberID]; exists {
		return nil, fmt.Errorf("subscriber %s already exists", subscriberID)
	}
	
	subscriber := &StreamSubscriber{
		ID:           subscriberID,
		Channel:      make(chan StreamMessage, 100), // Buffered channel
		Filter:       filter,
		Format:       format,
		IsActive:     true,
		ConnectedAt:  time.Now(),
		LastActivity: time.Now(),
	}
	
	sf.subscribers[subscriberID] = subscriber
	sf.activeStreams++
	
	return subscriber, nil
}

// Unsubscribe removes a stream subscription
func (sf *StreamFormatter) Unsubscribe(subscriberID string) error {
	sf.subscriberMutex.Lock()
	defer sf.subscriberMutex.Unlock()
	
	subscriber, exists := sf.subscribers[subscriberID]
	if !exists {
		return fmt.Errorf("subscriber %s not found", subscriberID)
	}
	
	subscriber.IsActive = false
	close(subscriber.Channel)
	delete(sf.subscribers, subscriberID)
	sf.activeStreams--
	
	return nil
}

// BroadcastEvent sends an event to all active subscribers
func (sf *StreamFormatter) BroadcastEvent(event models.GameEvent) error {
	sf.subscriberMutex.RLock()
	defer sf.subscriberMutex.RUnlock()
	
	for _, subscriber := range sf.subscribers {
		if !subscriber.IsActive {
			continue
		}
		
		// Apply filter
		if !sf.eventMatchesFilter(event, subscriber.Filter) {
			continue
		}
		
		// Format message based on subscriber's preferred format
		message, err := sf.formatEventForSubscriber(event, subscriber)
		if err != nil {
			continue
		}
		
		// Send message with timeout
		select {
		case subscriber.Channel <- message:
			subscriber.LastActivity = time.Now()
			sf.updateStats(len(fmt.Sprintf("%v", message.Data)))
		case <-time.After(time.Second * 5):
			// Timeout - mark subscriber as inactive
			subscriber.IsActive = false
		}
	}
	
	return nil
}

// BroadcastEvents sends multiple events to all active subscribers
func (sf *StreamFormatter) BroadcastEvents(events []models.GameEvent) error {
	for _, event := range events {
		if err := sf.BroadcastEvent(event); err != nil {
			return fmt.Errorf("error broadcasting event: %w", err)
		}
	}
	return nil
}

// GetSubscriber returns a subscriber by ID
func (sf *StreamFormatter) GetSubscriber(subscriberID string) (*StreamSubscriber, error) {
	sf.subscriberMutex.RLock()
	defer sf.subscriberMutex.RUnlock()
	
	subscriber, exists := sf.subscribers[subscriberID]
	if !exists {
		return nil, fmt.Errorf("subscriber %s not found", subscriberID)
	}
	
	return subscriber, nil
}

// processBuffer processes buffered events in batches
func (sf *StreamFormatter) processBuffer() {
	ticker := time.NewTicker(sf.batchTimeout)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			sf.flushBuffer()
		}
	}
}

// flushBuffer sends all buffered events to subscribers
func (sf *StreamFormatter) flushBuffer() {
	sf.bufferMutex.Lock()
	
	if len(sf.buffer) == 0 {
		sf.bufferMutex.Unlock()
		return
	}
	
	// Copy buffer and reset
	events := make([]models.GameEvent, len(sf.buffer))
	copy(events, sf.buffer)
	sf.buffer = sf.buffer[:0] // Reset slice but keep capacity
	
	sf.bufferMutex.Unlock()
	
	// Broadcast all buffered events
	sf.BroadcastEvents(events)
}

// cleanupInactiveSubscribers removes inactive subscribers periodically
func (sf *StreamFormatter) cleanupInactiveSubscribers() {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			sf.subscriberMutex.Lock()
			
			for id, subscriber := range sf.subscribers {
				// Remove subscribers inactive for more than 30 minutes
				if !subscriber.IsActive || time.Since(subscriber.LastActivity) > time.Minute*30 {
					close(subscriber.Channel)
					delete(sf.subscribers, id)
					sf.activeStreams--
				}
			}
			
			sf.subscriberMutex.Unlock()
		}
	}
}

// eventMatchesFilter checks if an event matches a subscriber's filter
func (sf *StreamFormatter) eventMatchesFilter(event models.GameEvent, filter *StreamFilter) bool {
	if filter == nil {
		return true // No filter means accept all
	}
	
	eventType := event.GetType()
	
	// Check event type filter
	if len(filter.EventTypes) > 0 {
		found := false
		for _, allowedType := range filter.EventTypes {
			if eventType == allowedType {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	
	// Check kills only filter
	if filter.KillsOnly && eventType != "player_death" {
		return false
	}
	
	// Check objectives only filter
	if filter.ObjectivesOnly {
		objectiveEvents := map[string]bool{
			"bomb_plant":   true,
			"bomb_defuse":  true,
			"bomb_explode": true,
			"round_start":  true,
			"round_end":    true,
		}
		if !objectiveEvents[eventType] {
			return false
		}
	}
	
	// Check specific event filters
	switch e := event.(type) {
	case *models.KillEvent:
		// Player filter
		if len(filter.Players) > 0 {
			found := false
			for _, playerName := range filter.Players {
				if e.Attacker.Name == playerName || e.Victim.Name == playerName {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		
		// Team filter
		if len(filter.Teams) > 0 {
			found := false
			for _, teamName := range filter.Teams {
				if e.Attacker.Side == teamName || e.Victim.Side == teamName {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
		
	case *models.PlayerHurtEvent:
		// Min damage filter
		if filter.MinDamage > 0 && e.Damage < filter.MinDamage {
			return false
		}
		
		// Player filter
		if len(filter.Players) > 0 {
			found := false
			for _, playerName := range filter.Players {
				if e.Attacker.Name == playerName || e.Victim.Name == playerName {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}
	
	return true
}

// formatEventForSubscriber formats an event according to subscriber preferences
func (sf *StreamFormatter) formatEventForSubscriber(event models.GameEvent, subscriber *StreamSubscriber) (StreamMessage, error) {
	message := StreamMessage{
		Type:      "event",
		Timestamp: event.GetTimestamp(),
	}
	
	switch subscriber.Format {
	case StreamFormatText:
		message.Data = sf.logFormatter.FormatEvent(event)
		
	case StreamFormatJSON, StreamFormatSSE:
		jsonEntry, err := sf.httpFormatter.convertEventToJSON(event)
		if err != nil {
			return message, fmt.Errorf("error converting event to JSON: %w", err)
		}
		message.Data = jsonEntry
		
	default:
		return message, fmt.Errorf("unsupported format: %s", subscriber.Format)
	}
	
	return message, nil
}

// updateStats updates streaming statistics
func (sf *StreamFormatter) updateStats(bytes int) {
	sf.statsMutex.Lock()
	defer sf.statsMutex.Unlock()
	
	sf.eventsSent++
	sf.bytesSent += int64(bytes)
}

// GetStreamStats returns streaming statistics
func (sf *StreamFormatter) GetStreamStats() map[string]interface{} {
	sf.statsMutex.RLock()
	sf.subscriberMutex.RLock()
	defer sf.statsMutex.RUnlock()
	defer sf.subscriberMutex.RUnlock()
	
	stats := map[string]interface{}{
		"active_streams":   sf.activeStreams,
		"total_subscribers": len(sf.subscribers),
		"events_sent":      sf.eventsSent,
		"bytes_sent":       sf.bytesSent,
		"buffer_size":      len(sf.buffer),
		"max_buffer_size":  sf.maxBufferSize,
		"batch_timeout":    sf.batchTimeout.String(),
	}
	
	// Add subscriber details
	subscriberStats := make(map[string]interface{})
	for id, subscriber := range sf.subscribers {
		subscriberStats[id] = map[string]interface{}{
			"format":        string(subscriber.Format),
			"is_active":     subscriber.IsActive,
			"connected_at":  subscriber.ConnectedAt,
			"last_activity": subscriber.LastActivity,
			"channel_size":  len(subscriber.Channel),
		}
	}
	stats["subscribers"] = subscriberStats
	
	return stats
}

// StreamLiveMatch streams events from a live match generation
func (sf *StreamFormatter) StreamLiveMatch(ctx context.Context, matchEngine interface{}) error {
	// This would integrate with the match engine to stream events in real-time
	// For now, it's a placeholder for the live streaming functionality
	
	// In a real implementation, this would:
	// 1. Hook into the match engine's event generation
	// 2. Stream events as they are created
	// 3. Handle backpressure and client disconnections
	// 4. Provide match status updates
	
	return fmt.Errorf("live match streaming not implemented yet")
}

// FormatForWebSocket formats events for WebSocket streaming
func (sf *StreamFormatter) FormatForWebSocket(events []models.GameEvent) ([]byte, error) {
	jsonEvents, err := sf.httpFormatter.FormatEventsAsJSON(events)
	if err != nil {
		return nil, fmt.Errorf("error formatting events for WebSocket: %w", err)
	}
	
	// Wrap in WebSocket message format
	message := map[string]interface{}{
		"type":      "events",
		"timestamp": time.Now(),
		"data":      json.RawMessage(jsonEvents),
	}
	
	return json.Marshal(message)
}

// FormatForSSE formats events for Server-Sent Events
func (sf *StreamFormatter) FormatForSSE(events []models.GameEvent) ([]string, error) {
	var sseLines []string
	
	for _, event := range events {
		jsonEntry, err := sf.httpFormatter.convertEventToJSON(event)
		if err != nil {
			return nil, fmt.Errorf("error converting event to JSON: %w", err)
		}
		
		jsonBytes, err := json.Marshal(jsonEntry)
		if err != nil {
			return nil, fmt.Errorf("error marshaling JSON: %w", err)
		}
		
		// Format as SSE
		sseLine := fmt.Sprintf("data: %s\n\n", string(jsonBytes))
		sseLines = append(sseLines, sseLine)
	}
	
	return sseLines, nil
}

// Shutdown gracefully shuts down the stream formatter
func (sf *StreamFormatter) Shutdown() error {
	sf.subscriberMutex.Lock()
	defer sf.subscriberMutex.Unlock()
	
	// Close all subscriber channels
	for id, subscriber := range sf.subscribers {
		subscriber.IsActive = false
		close(subscriber.Channel)
		delete(sf.subscribers, id)
	}
	
	sf.activeStreams = 0
	
	// Clear buffer
	sf.bufferMutex.Lock()
	sf.buffer = sf.buffer[:0]
	sf.bufferMutex.Unlock()
	
	return nil
}