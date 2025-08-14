package api

import (
	"errors"
	"strings"

	"github.com/noueii/nocs-log-generator/backend/pkg/models"
)

// ValidateGenerateRequest performs additional validation on generate requests
func ValidateGenerateRequest(req *models.GenerateRequest) error {
	// Check team names are not empty and different
	if strings.TrimSpace(req.Teams[0].Name) == "" {
		return errors.New("team 1 name cannot be empty")
	}
	if strings.TrimSpace(req.Teams[1].Name) == "" {
		return errors.New("team 2 name cannot be empty")
	}
	if strings.EqualFold(req.Teams[0].Name, req.Teams[1].Name) {
		return errors.New("team names must be different")
	}

	// Validate team sizes (should have exactly 5 players each)
	if len(req.Teams[0].Players) != 5 {
		return errors.New("team 1 must have exactly 5 players")
	}
	if len(req.Teams[1].Players) != 5 {
		return errors.New("team 2 must have exactly 5 players")
	}

	// Validate player names are unique across all teams
	playerNames := make(map[string]bool)
	for teamIdx, team := range req.Teams {
		for playerIdx, player := range team.Players {
			playerName := strings.TrimSpace(player.Name)
			if playerName == "" {
				return errors.New("all players must have names")
			}
			
			if playerNames[strings.ToLower(playerName)] {
				return errors.New("duplicate player name: " + playerName)
			}
			playerNames[strings.ToLower(playerName)] = true

			// Validate SteamID format if provided
			if player.SteamID != "" && !isValidSteamIDFormat(player.SteamID) {
				return errors.New("invalid SteamID format for player " + playerName)
			}

			// Set team reference if not set
			if req.Teams[teamIdx].Players[playerIdx].Team == "" {
				req.Teams[teamIdx].Players[playerIdx].Team = team.Name
			}
		}
	}

	// Validate match format
	validFormats := []string{"mr12", "mr15"}
	formatValid := false
	for _, format := range validFormats {
		if strings.EqualFold(req.Format, format) {
			formatValid = true
			break
		}
	}
	if !formatValid {
		return errors.New("format must be 'mr12' or 'mr15'")
	}

	// Validate map name
	if !isValidMapName(req.Map) {
		return errors.New("invalid map name: " + req.Map)
	}

	// Validate options if provided
	if req.Options.TickRate != 0 && (req.Options.TickRate < 64 || req.Options.TickRate > 128) {
		return errors.New("tick rate must be between 64 and 128")
	}

	if req.Options.MaxRounds != 0 {
		if req.Options.MaxRounds < 16 || req.Options.MaxRounds > 60 {
			return errors.New("max rounds must be between 16 and 60")
		}
	}

	return nil
}

// isValidSteamIDFormat validates SteamID format
func isValidSteamIDFormat(steamID string) bool {
	// Accept both old format (STEAM_X:Y:Z) and new format (7656119xxxxxxxx)
	if len(steamID) == 17 && strings.HasPrefix(steamID, "7656119") {
		return true // SteamID64 format
	}
	if strings.HasPrefix(steamID, "STEAM_") {
		return true // Legacy format
	}
	if strings.HasPrefix(steamID, "[U:1:") && strings.HasSuffix(steamID, "]") {
		return true // SteamID3 format
	}
	return false
}

// isValidMapName checks if the map name is in our supported list
func isValidMapName(mapName string) bool {
	validMaps := []string{
		"de_mirage", "de_dust2", "de_inferno", "de_cache", "de_overpass",
		"de_train", "de_nuke", "de_vertigo", "de_ancient", "de_anubis",
	}
	
	mapName = strings.ToLower(mapName)
	for _, validMap := range validMaps {
		if mapName == validMap {
			return true
		}
	}
	return false
}

// SanitizeTeamData ensures team data is properly formatted
func SanitizeTeamData(teams []models.Team) []models.Team {
	for i := range teams {
		// Trim and capitalize team names
		teams[i].Name = strings.TrimSpace(teams[i].Name)
		
		// Set default sides
		if i == 0 {
			teams[i].Side = "CT"
		} else {
			teams[i].Side = "TERRORIST" 
		}

		// Initialize team scores
		teams[i].Score = 0
		teams[i].RoundsWon = 0

		// Sanitize player data
		for j := range teams[i].Players {
			teams[i].Players[j].Name = strings.TrimSpace(teams[i].Players[j].Name)
			teams[i].Players[j].Team = teams[i].Name
			teams[i].Players[j].Side = teams[i].Side
			
			// Set default role if not specified
			if teams[i].Players[j].Role == "" {
				teams[i].Players[j].Role = "support"
			}

			// Initialize player state
			teams[i].Players[j].State.IsAlive = true
			teams[i].Players[j].State.Health = 100
			teams[i].Players[j].State.Armor = 0
			teams[i].Players[j].State.HasHelmet = false
			teams[i].Players[j].State.Grenades = make([]models.Grenade, 0)

			// Initialize player economy
			teams[i].Players[j].Economy.Money = 800 // Starting money
			teams[i].Players[j].Economy.Purchases = make([]models.Purchase, 0)

			// Initialize player stats
			teams[i].Players[j].Stats.GrenadesThrown = make(map[string]int)
		}
	}

	return teams
}

// GenerateResponseError creates a standardized error response
func GenerateResponseError(message string, details ...string) map[string]interface{} {
	response := map[string]interface{}{
		"error":   message,
		"success": false,
	}

	if len(details) > 0 {
		response["details"] = details
	}

	return response
}