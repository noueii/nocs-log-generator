package api

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/noueii/nocs-log-generator/backend/pkg/generator"
	"github.com/noueii/nocs-log-generator/backend/pkg/models"
	"github.com/noueii/nocs-log-generator/backend/pkg/websocket"
)

// Handler contains dependencies for API handlers
type Handler struct {
	generator *generator.MatchGenerator
	wsManager *websocket.Manager
}

// NewHandler creates a new API handler instance
func NewHandler() *Handler {
	return &Handler{
		generator: generator.NewMatchGenerator(),
	}
}

// SetWebSocketManager sets the WebSocket manager for the handler
func (h *Handler) SetWebSocketManager(wsManager *websocket.Manager) {
	h.wsManager = wsManager
}

// RegisterRoutes sets up API routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	// Match generation endpoints
	router.POST("/generate", h.GenerateMatch)
	
	// Configuration endpoints
	router.GET("/config/templates", h.GetConfigTemplates)
	router.GET("/config/maps", h.GetAvailableMaps)
	
	// Demo parsing endpoints (placeholder)
	router.POST("/parse", h.ParseDemo)
	
	// Utility endpoints
	router.GET("/ping", h.Ping)
	router.GET("/sample/request", h.GetSampleRequest)
}

// GenerateMatch handles match generation requests
func (h *Handler) GenerateMatch(c *gin.Context) {
	var req models.GenerateRequest
	
	// Parse and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format: " + err.Error(),
		})
		return
	}
	
	// Validate the request
	if err := req.Validate(); err != nil {
		log.Printf("Basic validation failed: %v", err)
		c.JSON(http.StatusBadRequest, GenerateResponseError("Basic validation failed: "+err.Error()))
		return
	}
	
	// Additional validation
	if err := ValidateGenerateRequest(&req); err != nil {
		log.Printf("Request validation failed: %v", err)
		c.JSON(http.StatusBadRequest, GenerateResponseError("Validation failed: "+err.Error()))
		return
	}
	
	// Sanitize team data
	req.Teams = SanitizeTeamData(req.Teams)
	
	// Broadcast generation start event if WebSocket is available
	if h.wsManager != nil {
		startEvent := websocket.GenerationStartEvent{
			MatchID:   "", // Will be set after match creation
			Teams:     []string{req.Teams[0].Name, req.Teams[1].Name},
			Map:       req.Map,
			Format:    req.Format,
			MaxRounds: getMaxRoundsForFormat(req.Format),
			StartedAt: time.Now().UTC(),
		}
		// We'll broadcast this after we have the match ID
		_ = startEvent
	}

	// Generate the match using the real generator
	match, err := h.generator.GenerateWithStreaming(&req, h.wsManager)
	if err != nil {
		log.Printf("Match generation failed: %v", err)
		
		// Broadcast error if WebSocket is available
		if h.wsManager != nil && match != nil {
			h.wsManager.BroadcastMatchError(match.ID, "Match generation failed: "+err.Error())
		}
		
		c.JSON(http.StatusInternalServerError, GenerateResponseError("Match generation failed: "+err.Error()))
		return
	}
	
	log.Printf("Successfully generated match %s: %s vs %s on %s (%d rounds, %d events)", 
		match.ID, match.Teams[0].Name, match.Teams[1].Name, match.Map, 
		len(match.Rounds), match.TotalEvents)
	
	// Broadcast completion event if WebSocket is available
	if h.wsManager != nil {
		completionEvent := websocket.GenerationCompleteEvent{
			MatchID:     match.ID,
			TotalRounds: len(match.Rounds),
			TotalEvents: match.TotalEvents,
			Duration:    match.Duration,
			CompletedAt: time.Now().UTC(),
			Success:     true,
		}
		h.wsManager.BroadcastMatchEvent(match.ID, websocket.EventTypeGenerationEnd, completionEvent)
	}
	
	// Return successful response
	response := models.GenerateResponse{
		MatchID: match.ID,
		Status:  match.Status,
		LogURL:  fmt.Sprintf("/api/v1/matches/%s/log", match.ID),
	}
	
	c.JSON(http.StatusOK, response)
}

// GetConfigTemplates returns predefined configuration templates
func (h *Handler) GetConfigTemplates(c *gin.Context) {
	templates := map[string]models.MatchConfig{
		"competitive": func() models.MatchConfig {
			config := models.DefaultMatchConfig()
			config.ApplyProfile("competitive")
			return config
		}(),
		"casual": func() models.MatchConfig {
			config := models.DefaultMatchConfig()
			config.ApplyProfile("casual")
			return config
		}(),
		"testing": func() models.MatchConfig {
			config := models.DefaultMatchConfig()
			config.ApplyProfile("testing")
			return config
		}(),
		"minimal": func() models.MatchConfig {
			config := models.DefaultMatchConfig()
			config.ApplyProfile("minimal")
			return config
		}(),
	}
	
	c.JSON(http.StatusOK, gin.H{
		"templates": templates,
	})
}

// GetAvailableMaps returns the list of available CS2 maps
func (h *Handler) GetAvailableMaps(c *gin.Context) {
	maps := []map[string]interface{}{
		{"name": "de_mirage", "display_name": "Mirage", "type": "defusal"},
		{"name": "de_dust2", "display_name": "Dust II", "type": "defusal"},
		{"name": "de_inferno", "display_name": "Inferno", "type": "defusal"},
		{"name": "de_cache", "display_name": "Cache", "type": "defusal"},
		{"name": "de_overpass", "display_name": "Overpass", "type": "defusal"},
		{"name": "de_train", "display_name": "Train", "type": "defusal"},
		{"name": "de_nuke", "display_name": "Nuke", "type": "defusal"},
		{"name": "de_vertigo", "display_name": "Vertigo", "type": "defusal"},
		{"name": "de_ancient", "display_name": "Ancient", "type": "defusal"},
	}
	
	c.JSON(http.StatusOK, gin.H{
		"maps": maps,
	})
}

// ParseDemo handles demo parsing requests (placeholder)
func (h *Handler) ParseDemo(c *gin.Context) {
	// TODO: Implement demo parsing
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Demo parsing not yet implemented",
	})
}

// GetSampleRequest returns a sample generate request for testing
func (h *Handler) GetSampleRequest(c *gin.Context) {
	sample := GetSampleGenerateRequest()
	c.JSON(http.StatusOK, gin.H{
		"sample_request": sample,
		"description": "Use this sample data to test the /api/v1/generate endpoint",
	})
}

// Ping is a simple ping endpoint for testing
func (h *Handler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"api":     "cs2-log-generator",
		"version": "0.1.0",
	})
}

// createMockMatch creates a mock match with sample data for testing
func (h *Handler) createMockMatch(match *models.Match) *models.Match {
	// Set some basic mock data
	match.Status = "completed"
	match.EndTime = time.Now().Add(45 * time.Minute) // Mock 45-minute match
	match.Duration = match.EndTime.Sub(match.StartTime)
	
	// Mock some rounds and scores
	match.CurrentRound = 24 // Full MR12 match
	match.Scores[match.Teams[0].Name] = 13
	match.Scores[match.Teams[1].Name] = 11
	
	// Mock some basic events
	match.TotalEvents = 1247 // Realistic number of events
	match.FileSize = 85432   // Mock file size in bytes
	
	// Create sample rounds data
	match.Rounds = h.createMockRounds(match)
	
	// Update team scores
	for i := range match.Teams {
		match.Teams[i].Score = match.Scores[match.Teams[i].Name]
		match.Teams[i].RoundsWon = match.Scores[match.Teams[i].Name]
	}
	
	return match
}

// createMockRounds creates mock round data
func (h *Handler) createMockRounds(match *models.Match) []models.RoundData {
	rounds := make([]models.RoundData, 0, match.CurrentRound)
	startTime := match.StartTime
	
	// Create mock rounds
	for i := 1; i <= match.CurrentRound; i++ {
		round := models.RoundData{
			RoundNumber: i,
			StartTime:   startTime,
			EndTime:     startTime.Add(2 * time.Minute), // Mock 2-minute rounds
			Winner:      h.getMockRoundWinner(i, match.Teams),
			Reason:      h.getMockRoundEndReason(i),
			MVP:         h.getMockMVP(match.Teams),
			Scores:      make(map[string]int),
			Economy:     make(map[string]models.TeamEconomy),
		}
		
		// Update scores up to this round
		team1Score := 0
		team2Score := 0
		
		for j := 1; j <= i; j++ {
			if h.getMockRoundWinner(j, match.Teams) == match.Teams[0].Name {
				team1Score++
			} else {
				team2Score++
			}
		}
		
		round.Scores[match.Teams[0].Name] = team1Score
		round.Scores[match.Teams[1].Name] = team2Score
		
		rounds = append(rounds, round)
		startTime = round.EndTime.Add(15 * time.Second) // Mock freeze time
	}
	
	return rounds
}

// Helper functions for mock data generation
func (h *Handler) getMockRoundWinner(round int, teams []models.Team) string {
	// Simple pattern: first team wins slightly more
	if round%3 == 0 {
		return teams[1].Name
	}
	return teams[0].Name
}

// getMaxRoundsForFormat returns the maximum rounds for a given format
func getMaxRoundsForFormat(format string) int {
	switch format {
	case "mr12":
		return 24
	case "mr15":
		return 30
	default:
		return 24
	}
}

func (h *Handler) getMockRoundEndReason(round int) string {
	reasons := []string{"elimination", "bomb_defused", "bomb_exploded", "time"}
	return reasons[round%len(reasons)]
}

func (h *Handler) getMockMVP(teams []models.Team) string {
	// Rotate through players
	allPlayers := append(teams[0].Players, teams[1].Players...)
	if len(allPlayers) > 0 {
		return allPlayers[0].Name
	}
	return "Player1"
}