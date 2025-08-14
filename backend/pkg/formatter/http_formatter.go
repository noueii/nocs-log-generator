package formatter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/noueii/nocs-log-generator/backend/pkg/models"
)

// HTTPFormatter handles formatting for HTTP endpoints and JSON responses
type HTTPFormatter struct {
	logFormatter *LogFormatter
	config       *models.MatchConfig
}

// NewHTTPFormatter creates a new HTTP formatter
func NewHTTPFormatter(config *models.MatchConfig) *HTTPFormatter {
	return &HTTPFormatter{
		logFormatter: NewLogFormatter(config),
		config:       config,
	}
}

// JSONLogEntry represents a single log entry in JSON format
type JSONLogEntry struct {
	Timestamp   time.Time   `json:"timestamp"`
	Type        string      `json:"type"`
	Tick        int64       `json:"tick"`
	Round       int         `json:"round"`
	LogLine     string      `json:"log_line"`
	RawData     interface{} `json:"raw_data,omitempty"`
	Metadata    *EventMetadata `json:"metadata,omitempty"`
}

// EventMetadata contains additional metadata about the event
type EventMetadata struct {
	Players     []string    `json:"players,omitempty"`
	Teams       []string    `json:"teams,omitempty"`
	Weapon      string      `json:"weapon,omitempty"`
	Location    string      `json:"location,omitempty"`
	Modifiers   []string    `json:"modifiers,omitempty"`
	Damage      int         `json:"damage,omitempty"`
	IsKill      bool        `json:"is_kill,omitempty"`
	IsObjective bool        `json:"is_objective,omitempty"`
}

// HTTPLogResponse represents the complete HTTP response for log data
type HTTPLogResponse struct {
	MatchID     string         `json:"match_id"`
	Map         string         `json:"map"`
	Format      string         `json:"format"`
	Status      string         `json:"status"`
	StartTime   time.Time      `json:"start_time"`
	EndTime     time.Time      `json:"end_time,omitempty"`
	Duration    time.Duration  `json:"duration,omitempty"`
	TotalEvents int            `json:"total_events"`
	Teams       []TeamSummary  `json:"teams"`
	Events      []JSONLogEntry `json:"events"`
	Rounds      []RoundSummary `json:"rounds,omitempty"`
	Statistics  *MatchStats    `json:"statistics,omitempty"`
}

// TeamSummary provides a summary of team performance
type TeamSummary struct {
	Name    string `json:"name"`
	Side    string `json:"side"`
	Score   int    `json:"score"`
	Players []PlayerSummary `json:"players"`
}

// PlayerSummary provides a summary of player performance
type PlayerSummary struct {
	Name     string  `json:"name"`
	UserID   int     `json:"user_id"`
	SteamID  string  `json:"steam_id"`
	Kills    int     `json:"kills"`
	Deaths   int     `json:"deaths"`
	Assists  int     `json:"assists"`
	Rating   float64 `json:"rating"`
	Headshots int    `json:"headshots"`
}

// RoundSummary provides a summary of round data
type RoundSummary struct {
	RoundNumber int           `json:"round_number"`
	Winner      string        `json:"winner"`
	Reason      string        `json:"reason"`
	Duration    time.Duration `json:"duration"`
	MVP         string        `json:"mvp"`
	CTScore     int           `json:"ct_score"`
	TScore      int           `json:"t_score"`
	EventCount  int           `json:"event_count"`
}

// MatchStats provides overall match statistics
type MatchStats struct {
	TotalRounds   int                    `json:"total_rounds"`
	CTWins        int                    `json:"ct_wins"`
	TWins         int                    `json:"t_wins"`
	BombPlants    int                    `json:"bomb_plants"`
	BombDefuses   int                    `json:"bomb_defuses"`
	BombExplosions int                   `json:"bomb_explosions"`
	TotalKills    int                    `json:"total_kills"`
	TotalDamage   int                    `json:"total_damage"`
	EventTypes    map[string]int         `json:"event_types"`
	WeaponStats   map[string]WeaponStat  `json:"weapon_stats"`
}

// WeaponStat tracks statistics for individual weapons
type WeaponStat struct {
	Kills     int     `json:"kills"`
	Headshots int     `json:"headshots"`
	Damage    int     `json:"damage"`
	Accuracy  float64 `json:"accuracy,omitempty"`
}

// FormatAsHTTPLog converts a match to HTTP JSON format
func (f *HTTPFormatter) FormatAsHTTPLog(match *models.Match) (*HTTPLogResponse, error) {
	response := &HTTPLogResponse{
		MatchID:     match.ID,
		Map:         match.Map,
		Format:      match.Format,
		Status:      match.Status,
		StartTime:   match.StartTime,
		EndTime:     match.EndTime,
		Duration:    match.Duration,
		TotalEvents: len(match.Events),
		Teams:       make([]TeamSummary, 0, len(match.Teams)),
		Events:      make([]JSONLogEntry, 0, len(match.Events)),
		Rounds:      make([]RoundSummary, 0, len(match.Rounds)),
	}
	
	// Format teams
	for _, team := range match.Teams {
		teamSummary := TeamSummary{
			Name:    team.Name,
			Side:    team.Side,
			Score:   match.Scores[team.Name],
			Players: make([]PlayerSummary, 0, len(team.Players)),
		}
		
		for _, player := range team.Players {
			playerSummary := PlayerSummary{
				Name:      player.Name,
				UserID:    player.UserID,
				SteamID:   player.SteamID,
				Kills:     player.Stats.Kills,
				Deaths:    player.Stats.Deaths,
				Assists:   player.Stats.Assists,
				Rating:    player.Stats.Rating,
				Headshots: player.Stats.Headshots,
			}
			teamSummary.Players = append(teamSummary.Players, playerSummary)
		}
		
		response.Teams = append(response.Teams, teamSummary)
	}
	
	// Format events
	for _, event := range match.Events {
		jsonEntry, err := f.convertEventToJSON(event)
		if err != nil {
			return nil, fmt.Errorf("error converting event to JSON: %w", err)
		}
		response.Events = append(response.Events, *jsonEntry)
	}
	
	// Format rounds
	for _, round := range match.Rounds {
		roundSummary := RoundSummary{
			RoundNumber: round.RoundNumber,
			Winner:      round.Winner,
			Reason:      round.Reason,
			Duration:    round.EndTime.Sub(round.StartTime),
			MVP:         round.MVP,
			CTScore:     round.Scores["CT"],
			TScore:      round.Scores["TERRORIST"], 
			EventCount:  len(round.Events),
		}
		response.Rounds = append(response.Rounds, roundSummary)
	}
	
	// Generate statistics
	response.Statistics = f.generateMatchStats(match)
	
	return response, nil
}

// FormatEventsAsJSON formats multiple events as JSON array
func (f *HTTPFormatter) FormatEventsAsJSON(events []models.GameEvent) ([]byte, error) {
	jsonEvents := make([]JSONLogEntry, 0, len(events))
	
	for _, event := range events {
		jsonEntry, err := f.convertEventToJSON(event)
		if err != nil {
			return nil, fmt.Errorf("error converting event to JSON: %w", err)
		}
		jsonEvents = append(jsonEvents, *jsonEntry)
	}
	
	return json.Marshal(jsonEvents)
}

// FormatEventAsJSON formats a single event as JSON
func (f *HTTPFormatter) FormatEventAsJSON(event models.GameEvent) ([]byte, error) {
	jsonEntry, err := f.convertEventToJSON(event)
	if err != nil {
		return nil, fmt.Errorf("error converting event to JSON: %w", err)
	}
	
	return json.Marshal(jsonEntry)
}

// BatchFormatEvents formats multiple events in batches for better performance
func (f *HTTPFormatter) BatchFormatEvents(events []models.GameEvent, batchSize int) ([][]byte, error) {
	if batchSize <= 0 {
		batchSize = 100 // Default batch size
	}
	
	var batches [][]byte
	
	for i := 0; i < len(events); i += batchSize {
		end := i + batchSize
		if end > len(events) {
			end = len(events)
		}
		
		batch := events[i:end]
		batchJSON, err := f.FormatEventsAsJSON(batch)
		if err != nil {
			return nil, fmt.Errorf("error formatting batch %d-%d: %w", i, end, err)
		}
		
		batches = append(batches, batchJSON)
	}
	
	return batches, nil
}

// convertEventToJSON converts a GameEvent to JSONLogEntry
func (f *HTTPFormatter) convertEventToJSON(event models.GameEvent) (*JSONLogEntry, error) {
	if event == nil {
		return nil, fmt.Errorf("event is nil")
	}
	
	// Get raw JSON data
	rawData, err := event.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("error converting event to raw JSON: %w", err)
	}
	
	// Parse back to get the raw interface
	var eventData interface{}
	if err := json.Unmarshal(rawData, &eventData); err != nil {
		return nil, fmt.Errorf("error parsing event JSON: %w", err)
	}
	
	// Create JSON log entry
	jsonEntry := &JSONLogEntry{
		Timestamp: event.GetTimestamp(),
		Type:      event.GetType(),
		Tick:      event.GetTick(),
		LogLine:   f.logFormatter.FormatEvent(event),
		RawData:   eventData,
		Metadata:  f.extractEventMetadata(event),
	}
	
	// Extract round number if available
	if eventMap, ok := eventData.(map[string]interface{}); ok {
		if roundNum, exists := eventMap["round"]; exists {
			if round, ok := roundNum.(float64); ok {
				jsonEntry.Round = int(round)
			}
		}
	}
	
	return jsonEntry, nil
}

// extractEventMetadata extracts metadata from events for easier filtering/searching
func (f *HTTPFormatter) extractEventMetadata(event models.GameEvent) *EventMetadata {
	metadata := &EventMetadata{}
	
	switch e := event.(type) {
	case *models.KillEvent:
		metadata.Players = []string{e.Attacker.Name, e.Victim.Name}
		metadata.Teams = []string{e.Attacker.Side, e.Victim.Side}
		metadata.Weapon = e.Weapon
		metadata.IsKill = true
		
		var modifiers []string
		if e.Headshot {
			modifiers = append(modifiers, "headshot")
		}
		if e.Penetrated > 0 {
			modifiers = append(modifiers, "penetrated")
		}
		if e.NoScope {
			modifiers = append(modifiers, "noscope")
		}
		if e.AttackerBlind {
			modifiers = append(modifiers, "attackerblind")
		}
		metadata.Modifiers = modifiers
		
	case *models.PlayerHurtEvent:
		metadata.Players = []string{e.Attacker.Name, e.Victim.Name}
		metadata.Teams = []string{e.Attacker.Side, e.Victim.Side}
		metadata.Weapon = e.Weapon
		metadata.Damage = e.Damage
		
	case *models.BombPlantEvent:
		metadata.Players = []string{e.Player.Name}
		metadata.Teams = []string{e.Player.Side}
		metadata.Location = e.Site
		metadata.IsObjective = true
		
	case *models.BombDefuseEvent:
		metadata.Players = []string{e.Player.Name}
		metadata.Teams = []string{e.Player.Side}
		metadata.Location = e.Site
		metadata.IsObjective = true
		if e.WithKit {
			metadata.Modifiers = []string{"with_kit"}
		}
		
	case *models.BombExplodeEvent:
		metadata.Location = e.Site
		metadata.IsObjective = true
		
	case *models.ItemPurchaseEvent:
		metadata.Players = []string{e.Player.Name}
		metadata.Teams = []string{e.Player.Side}
		metadata.Weapon = e.Item // Item could be weapon or equipment
		
	case *models.GrenadeThrowEvent:
		metadata.Players = []string{e.Player.Name}
		metadata.Teams = []string{e.Player.Side}
		metadata.Weapon = e.GrenadeType
		
	case *models.FlashbangEvent:
		metadata.Players = []string{e.Player.Name}
		metadata.Teams = []string{e.Player.Side}
		metadata.Weapon = "flashbang"
		
		// Add flashed players
		for _, flashed := range e.Flashed {
			metadata.Players = append(metadata.Players, flashed.Name)
		}
		
	case *models.ChatEvent:
		if e.Player != nil {
			metadata.Players = []string{e.Player.Name}
			metadata.Teams = []string{e.Player.Side}
		}
		
		var modifiers []string
		if e.Team {
			modifiers = append(modifiers, "team")
		}
		if e.Dead {
			modifiers = append(modifiers, "dead")
		}
		metadata.Modifiers = modifiers
		
	case *models.RoundStartEvent, *models.RoundEndEvent:
		metadata.IsObjective = true
	}
	
	return metadata
}

// generateMatchStats generates comprehensive match statistics
func (f *HTTPFormatter) generateMatchStats(match *models.Match) *MatchStats {
	stats := &MatchStats{
		TotalRounds:   len(match.Rounds),
		EventTypes:    make(map[string]int),
		WeaponStats:   make(map[string]WeaponStat),
	}
	
	// Count wins
	for _, round := range match.Rounds {
		if round.Winner == "CT" {
			stats.CTWins++
		} else {
			stats.TWins++
		}
	}
	
	// Analyze events
	for _, event := range match.Events {
		eventType := event.GetType()
		stats.EventTypes[eventType]++
		
		switch e := event.(type) {
		case *models.KillEvent:
			stats.TotalKills++
			
			weaponStat := stats.WeaponStats[e.Weapon]
			weaponStat.Kills++
			if e.Headshot {
				weaponStat.Headshots++
			}
			stats.WeaponStats[e.Weapon] = weaponStat
			
		case *models.PlayerHurtEvent:
			stats.TotalDamage += e.Damage
			
			weaponStat := stats.WeaponStats[e.Weapon]
			weaponStat.Damage += e.Damage
			stats.WeaponStats[e.Weapon] = weaponStat
			
		case *models.BombPlantEvent:
			stats.BombPlants++
			
		case *models.BombDefuseEvent:
			stats.BombDefuses++
			
		case *models.BombExplodeEvent:
			stats.BombExplosions++
		}
	}
	
	return stats
}

// FormatTimestamp formats a timestamp for HTTP responses
func (f *HTTPFormatter) FormatTimestamp(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FilterEventsByType filters events by type for HTTP responses
func (f *HTTPFormatter) FilterEventsByType(events []models.GameEvent, eventType string) []JSONLogEntry {
	var filtered []JSONLogEntry
	
	for _, event := range events {
		if event.GetType() == eventType {
			if jsonEntry, err := f.convertEventToJSON(event); err == nil {
				filtered = append(filtered, *jsonEntry)
			}
		}
	}
	
	return filtered
}

// FilterEventsByPlayer filters events by player name for HTTP responses
func (f *HTTPFormatter) FilterEventsByPlayer(events []models.GameEvent, playerName string) []JSONLogEntry {
	var filtered []JSONLogEntry
	
	for _, event := range events {
		jsonEntry, err := f.convertEventToJSON(event)
		if err != nil {
			continue
		}
		
		// Check if player is involved in the event
		if jsonEntry.Metadata != nil {
			for _, player := range jsonEntry.Metadata.Players {
				if player == playerName {
					filtered = append(filtered, *jsonEntry)
					break
				}
			}
		}
	}
	
	return filtered
}

// FilterEventsByRound filters events by round number for HTTP responses
func (f *HTTPFormatter) FilterEventsByRound(events []models.GameEvent, roundNumber int) []JSONLogEntry {
	var filtered []JSONLogEntry
	
	for _, event := range events {
		if jsonEntry, err := f.convertEventToJSON(event); err == nil && jsonEntry.Round == roundNumber {
			filtered = append(filtered, *jsonEntry)
		}
	}
	
	return filtered
}

// GetHTTPFormatterStats returns formatter statistics for HTTP endpoints
func (f *HTTPFormatter) GetHTTPFormatterStats() map[string]interface{} {
	baseStats := f.logFormatter.GetFormatterStats()
	
	httpStats := map[string]interface{}{
		"formatter_type":    "http",
		"json_support":      true,
		"batch_support":     true,
		"filter_support":    true,
		"metadata_support":  true,
	}
	
	// Merge with base formatter stats
	for k, v := range baseStats {
		httpStats[k] = v
	}
	
	return httpStats
}