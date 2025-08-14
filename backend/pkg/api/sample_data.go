package api

import (
	"github.com/noueii/nocs-log-generator/backend/pkg/models"
)

// GetSampleGenerateRequest returns a sample request for testing
func GetSampleGenerateRequest() models.GenerateRequest {
	// Create sample teams
	team1 := models.Team{
		Name: "Astralis",
		Tag:  "AST",
		Country: "Denmark",
		Players: []models.Player{
			{Name: "device", SteamID: "STEAM_1:0:123456", Role: "awp"},
			{Name: "dupreeh", SteamID: "STEAM_1:1:234567", Role: "entry"},
			{Name: "Xyp9x", SteamID: "STEAM_1:0:345678", Role: "support"},
			{Name: "gla1ve", SteamID: "STEAM_1:1:456789", Role: "igl"},
			{Name: "Magisk", SteamID: "STEAM_1:0:567890", Role: "rifler"},
		},
	}

	team2 := models.Team{
		Name: "NAVI",
		Tag:  "NAVI",
		Country: "Ukraine",
		Players: []models.Player{
			{Name: "s1mple", SteamID: "STEAM_1:1:987654", Role: "awp"},
			{Name: "electronic", SteamID: "STEAM_1:0:876543", Role: "entry"},
			{Name: "Perfecto", SteamID: "STEAM_1:1:765432", Role: "support"},
			{Name: "b1t", SteamID: "STEAM_1:0:654321", Role: "rifler"},
			{Name: "Aleksib", SteamID: "STEAM_1:1:543210", Role: "igl"},
		},
	}

	return models.GenerateRequest{
		Teams:  []models.Team{team1, team2},
		Map:    "de_mirage",
		Format: "mr12",
		Options: models.MatchOptions{
			Seed:     12345,
			TickRate: 64,
			Overtime: true,
		},
	}
}

// GetSampleMatchConfig returns a sample match configuration
func GetSampleMatchConfig() models.MatchConfig {
	config := models.DefaultMatchConfig()
	config.Map = "de_mirage"
	config.Format = "mr12"
	config.Seed = 12345
	config.RealisticEconomy = true
	config.ChatMessages = true
	config.SkillVariance = 0.15
	return config
}

// GetValidMapList returns the complete list of valid maps
func GetValidMapList() []string {
	return []string{
		"de_mirage",
		"de_dust2", 
		"de_inferno",
		"de_cache",
		"de_overpass",
		"de_train",
		"de_nuke",
		"de_vertigo",
		"de_ancient",
		"de_anubis",
	}
}