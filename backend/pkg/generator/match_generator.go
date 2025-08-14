package generator

import (
	"fmt"
	"time"

	"github.com/noueii/nocs-log-generator/backend/pkg/models"
)

// Event structures for WebSocket streaming
type GenerationStartEvent struct {
	MatchID   string    `json:"match_id"`
	Teams     []string  `json:"teams"`
	Map       string    `json:"map"`
	Format    string    `json:"format"`
	MaxRounds int       `json:"max_rounds"`
	StartedAt time.Time `json:"started_at"`
}

type GenerationErrorEvent struct {
	MatchID string    `json:"match_id"`
	Error   string    `json:"error"`
	Time    time.Time `json:"time"`
}

// MatchGenerator handles CS2 match log generation
type MatchGenerator struct {
	economyManager *models.EconomyManager
}

// NewMatchGenerator creates a new match generator instance
func NewMatchGenerator() *MatchGenerator {
	return &MatchGenerator{
		economyManager: models.NewEconomyManager(),
	}
}

// Generate creates a CS2 match log from the given configuration
func (g *MatchGenerator) Generate(req *models.GenerateRequest) (*models.Match, error) {
	if req == nil {
		return nil, fmt.Errorf("generate request cannot be nil")
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Create match configuration from request
	config := models.DefaultMatchConfig()
	config.Format = req.Format
	config.Map = req.Map
	
	// Apply options if provided
	if req.Options.TickRate > 0 {
		config.TickRate = req.Options.TickRate
	}
	if req.Options.Seed > 0 {
		config.Seed = req.Options.Seed
	}
	if req.Options.MaxRounds > 0 {
		config.MaxRounds = req.Options.MaxRounds
	}
	config.Overtime = req.Options.Overtime

	// Prepare teams with proper side assignments
	teams := make([]models.Team, len(req.Teams))
	copy(teams, req.Teams)
	
	// Assign sides (first team CT, second team T)
	teams[0].Side = "CT"
	teams[1].Side = "TERRORIST"
	
	// Update player sides and assign user IDs
	for i := range teams {
		for j := range teams[i].Players {
			teams[i].Players[j].Side = teams[i].Side
			teams[i].Players[j].Team = teams[i].Name
			teams[i].Players[j].UserID = (i * 5) + j + 1 // Simple user ID assignment
		}
	}

	// Create match
	match := models.NewMatch(config, teams)
	match.Status = "generating"
	match.StartTime = time.Now()

	// Create match engine and generate the match
	engine := NewMatchEngine(&config, match)
	if err := engine.GenerateMatch(); err != nil {
		match.Status = "error"
		match.Error = err.Error()
		return match, fmt.Errorf("match generation failed: %w", err)
	}

	return match, nil
}

// GenerateWithStreaming creates a CS2 match log with WebSocket streaming support
func (g *MatchGenerator) GenerateWithStreaming(req *models.GenerateRequest, wsManager WebSocketManager) (*models.Match, error) {
	if req == nil {
		return nil, fmt.Errorf("generate request cannot be nil")
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Create match configuration from request
	config := models.DefaultMatchConfig()
	config.Format = req.Format
	config.Map = req.Map
	
	// Apply options if provided
	if req.Options.TickRate > 0 {
		config.TickRate = req.Options.TickRate
	}
	if req.Options.Seed > 0 {
		config.Seed = req.Options.Seed
	}
	if req.Options.MaxRounds > 0 {
		config.MaxRounds = req.Options.MaxRounds
	}
	config.Overtime = req.Options.Overtime

	// Prepare teams with proper side assignments
	teams := make([]models.Team, len(req.Teams))
	copy(teams, req.Teams)
	
	// Assign sides (first team CT, second team T)
	teams[0].Side = "CT"
	teams[1].Side = "TERRORIST"
	
	// Update player sides and assign user IDs
	for i := range teams {
		for j := range teams[i].Players {
			teams[i].Players[j].Side = teams[i].Side
			teams[i].Players[j].Team = teams[i].Name
			teams[i].Players[j].UserID = (i * 5) + j + 1 // Simple user ID assignment
		}
	}

	// Create match
	match := models.NewMatch(config, teams)
	match.Status = "generating"
	match.StartTime = time.Now()

	// Broadcast generation start event
	if wsManager != nil {
		startEvent := GenerationStartEvent{
			MatchID:   match.ID,
			Teams:     []string{teams[0].Name, teams[1].Name},
			Map:       config.Map,
			Format:    config.Format,
			MaxRounds: config.MaxRounds,
			StartedAt: match.StartTime,
		}
		wsManager.BroadcastMatchEvent(match.ID, "generation_start", startEvent)
	}

	// Create match engine with streaming support and generate the match
	engine := NewMatchEngine(&config, match)
	engine.SetWebSocketManager(wsManager)
	
	if err := engine.GenerateMatchWithStreaming(); err != nil {
		match.Status = "error"
		match.Error = err.Error()
		
		// Broadcast error event
		if wsManager != nil {
			errorEvent := GenerationErrorEvent{
				MatchID: match.ID,
				Error:   err.Error(),
				Time:    time.Now(),
			}
			wsManager.BroadcastMatchEvent(match.ID, "generation_error", errorEvent)
		}
		
		return match, fmt.Errorf("match generation failed: %w", err)
	}

	return match, nil
}