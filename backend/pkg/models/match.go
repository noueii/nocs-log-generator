package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Match represents a CS2 match configuration and state
type Match struct {
	// Basic information
	ID          string    `json:"id"`
	Title       string    `json:"title,omitempty"`
	Map         string    `json:"map"`
	Format      string    `json:"format"` // "mr12" or "mr15"
	Status      string    `json:"status"` // "pending", "generating", "completed", "error"
	StartTime   time.Time `json:"start_time,omitempty"`
	EndTime     time.Time `json:"end_time,omitempty"`
	LogURL      string    `json:"log_url,omitempty"`
	Error       string    `json:"error,omitempty"`
	
	// Match configuration
	Config      MatchConfig `json:"config"`
	
	// Teams and players
	Teams       []Team    `json:"teams"`
	
	// Match state
	CurrentRound int       `json:"current_round"`
	MaxRounds    int       `json:"max_rounds"`
	Overtime     bool      `json:"overtime"`
	Scores       map[string]int `json:"scores"`
	
	// Round history and events
	Rounds       []RoundData `json:"rounds,omitempty"`
	Events       []GameEvent `json:"events,omitempty"`
	
	// Statistics
	TotalEvents  int64     `json:"total_events"`
	FileSize     int64     `json:"file_size,omitempty"`
	Duration     time.Duration `json:"duration,omitempty"`
}

// RoundData represents the state and events of a single round
type RoundData struct {
	RoundNumber  int         `json:"round_number"`
	StartTime    time.Time   `json:"start_time"`
	EndTime      time.Time   `json:"end_time"`
	Winner       string      `json:"winner"`      // "CT", "TERRORIST"
	Reason       string      `json:"reason"`      // "elimination", "bomb_defused", "bomb_exploded", "time"
	MVP          string      `json:"mvp"`         // Player name
	Events       []GameEvent `json:"events"`
	Economy      map[string]TeamEconomy `json:"economy"`
	Scores       map[string]int `json:"scores"`
}

// MatchState represents the current state during match generation
type MatchState struct {
	CurrentRound  int
	Scores        map[string]int
	TeamEconomies map[string]*TeamEconomy
	PlayerStates  map[string]*PlayerState
	BombCarrier   *Player
	IsLive        bool
	IsFreezeTime  bool
	RoundStartTime time.Time
	CurrentTick   int64
}

// GenerateRequest represents the request body for match generation
type GenerateRequest struct {
	Teams     []Team       `json:"teams" binding:"required,len=2"`
	Map       string       `json:"map" binding:"required"`
	Format    string       `json:"format" binding:"required,oneof=mr12 mr15"`
	Options   MatchOptions `json:"options"`
}

// MatchOptions contains additional configuration for match generation
type MatchOptions struct {
	Seed       int64 `json:"seed,omitempty"`       // Random seed for reproducible generation
	TickRate   int   `json:"tick_rate,omitempty"`  // Default: 64
	Overtime   bool  `json:"overtime,omitempty"`   // Allow overtime
	MaxRounds  int   `json:"max_rounds,omitempty"` // Override default based on format
}

// GenerateResponse represents the response from match generation
type GenerateResponse struct {
	MatchID string `json:"match_id"`
	Status  string `json:"status"`
	LogURL  string `json:"log_url,omitempty"`
	Error   string `json:"error,omitempty"`
}

// NewMatch creates a new match with the given configuration
func NewMatch(config MatchConfig, teams []Team) *Match {
	match := &Match{
		ID:           generateMatchID(),
		Title:        fmt.Sprintf("%s vs %s", teams[0].Name, teams[1].Name),
		Map:          config.Map,
		Format:       config.Format,
		Status:       "pending",
		Config:       config,
		Teams:        teams,
		CurrentRound: 0,
		Scores:       make(map[string]int),
		Rounds:       make([]RoundData, 0),
		Events:       make([]GameEvent, 0),
	}
	
	// Set max rounds based on format
	switch config.Format {
	case "mr12":
		match.MaxRounds = 24
	case "mr15":
		match.MaxRounds = 30
	default:
		match.MaxRounds = 24
	}
	
	// Initialize scores
	for _, team := range teams {
		match.Scores[team.Name] = 0
	}
	
	return match
}

// IsFinished returns true if the match is complete
func (m *Match) IsFinished() bool {
	if m.Status == "completed" {
		return true
	}
	
	// Check if any team has won
	winThreshold := (m.MaxRounds / 2) + 1
	for _, score := range m.Scores {
		if score >= winThreshold {
			return true
		}
	}
	
	// Check overtime conditions
	if m.CurrentRound >= m.MaxRounds && !m.Overtime {
		return true
	}
	
	return false
}

// GetWinningTeam returns the name of the winning team, or empty string if no winner
func (m *Match) GetWinningTeam() string {
	winThreshold := (m.MaxRounds / 2) + 1
	highestScore := 0
	winningTeam := ""
	
	for teamName, score := range m.Scores {
		if score >= winThreshold && score > highestScore {
			highestScore = score
			winningTeam = teamName
		}
	}
	
	return winningTeam
}

// AddEvent adds a game event to the match
func (m *Match) AddEvent(event GameEvent) {
	m.Events = append(m.Events, event)
	m.TotalEvents++
}

// Validate validates the match configuration
func (m *Match) Validate() error {
	if m.ID == "" {
		return errors.New("match ID is required")
	}
	
	if len(m.Teams) != 2 {
		return errors.New("exactly 2 teams are required")
	}
	
	if m.Map == "" {
		return errors.New("map is required")
	}
	
	if m.Format != "mr12" && m.Format != "mr15" {
		return errors.New("format must be 'mr12' or 'mr15'")
	}
	
	// Validate teams
	for i, team := range m.Teams {
		if err := team.Validate(); err != nil {
			return fmt.Errorf("team %d validation failed: %w", i+1, err)
		}
	}
	
	return nil
}

// Validate validates the generate request
func (r *GenerateRequest) Validate() error {
	if len(r.Teams) != 2 {
		return errors.New("exactly 2 teams are required")
	}
	
	if r.Map == "" {
		return errors.New("map is required")
	}
	
	if r.Format != "mr12" && r.Format != "mr15" {
		return errors.New("format must be 'mr12' or 'mr15'")
	}
	
	// Validate teams
	for i, team := range r.Teams {
		if err := team.Validate(); err != nil {
			return fmt.Errorf("team %d validation failed: %w", i+1, err)
		}
	}
	
	// Validate options
	if r.Options.TickRate != 0 && (r.Options.TickRate < 64 || r.Options.TickRate > 128) {
		return errors.New("tick rate must be between 64 and 128")
	}
	
	return nil
}

// generateMatchID generates a unique match ID
func generateMatchID() string {
	// Simple timestamp-based ID for MVP
	return fmt.Sprintf("match_%d", time.Now().Unix())
}

// GetTeamBySide returns the team playing on the specified side
func (m *Match) GetTeamBySide(side string) *Team {
	for i := range m.Teams {
		if strings.EqualFold(m.Teams[i].Side, side) {
			return &m.Teams[i]
		}
	}
	return nil
}

// GetPlayerByName returns a player by name from any team
func (m *Match) GetPlayerByName(name string) *Player {
	for _, team := range m.Teams {
		for i := range team.Players {
			if strings.EqualFold(team.Players[i].Name, name) {
				return &team.Players[i]
			}
		}
	}
	return nil
}