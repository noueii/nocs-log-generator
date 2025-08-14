package models

import (
	"errors"
	"fmt"
	"strings"
)

// Team represents a CS2 team
type Team struct {
	// Basic information
	Name        string `json:"name" binding:"required"`
	Tag         string `json:"tag,omitempty"`
	Country     string `json:"country,omitempty"`
	Ranking     int    `json:"ranking,omitempty"`
	
	// Players
	Players     []Player `json:"players" binding:"required,len=5"`
	
	// Match state
	Side        string `json:"side"`         // "CT" or "TERRORIST"
	Score       int    `json:"score"`
	RoundsWon   int    `json:"rounds_won"`
	
	// Economy state
	Economy     TeamEconomy `json:"economy"`
	
	// Statistics
	Stats       TeamStats `json:"stats"`
}

// TeamEconomy represents the team's economic state
type TeamEconomy struct {
	TotalMoney     int `json:"total_money"`
	AverageMoney   int `json:"average_money"`
	EquipmentValue int `json:"equipment_value"`
	
	// Loss bonus tracking
	ConsecutiveLosses int `json:"consecutive_losses"`
	LossBonus        int `json:"loss_bonus"`
	
	// Round economy
	MoneySpent       int `json:"money_spent"`
	MoneyEarned      int `json:"money_earned"`
	
	// Equipment counts
	Rifles           int `json:"rifles"`
	SMGs             int `json:"smgs"`
	Pistols          int `json:"pistols"`
	Snipers          int `json:"snipers"`
	Grenades         int `json:"grenades"`
	
	// Armor and utilities
	Armor            int `json:"armor"`
	Helmets          int `json:"helmets"`
	DefuseKits       int `json:"defuse_kits"`
}

// TeamStats represents aggregate team statistics
type TeamStats struct {
	// Basic stats
	Kills            int `json:"kills"`
	Deaths           int `json:"deaths"`
	Assists          int `json:"assists"`
	Damage           int `json:"damage"`
	
	// Combat effectiveness
	Headshots        int     `json:"headshots"`
	HeadshotRate     float64 `json:"headshot_rate"`
	FirstKills       int     `json:"first_kills"`
	FirstDeaths      int     `json:"first_deaths"`
	
	// Objective stats
	BombPlants       int `json:"bomb_plants"`
	BombDefuses      int `json:"bomb_defuses"`
	
	// Round stats
	RoundsPlayed     int `json:"rounds_played"`
	RoundsWonCT      int `json:"rounds_won_ct"`
	RoundsWonT       int `json:"rounds_won_t"`
	
	// Economic efficiency
	MoneyPerRound    int     `json:"money_per_round"`
	EconomyRating    float64 `json:"economy_rating"`
}

// PlayerState represents a player's current state during the match
type PlayerState struct {
	// Basic state
	IsAlive      bool    `json:"is_alive"`
	Health       int     `json:"health"`
	Armor        int     `json:"armor"`
	HasHelmet    bool    `json:"has_helmet"`
	HasDefuseKit bool    `json:"has_defuse_kit"`
	
	// Position and movement
	Position     Vector3 `json:"position"`
	ViewAngle    Vector3 `json:"view_angle"`
	Velocity     Vector3 `json:"velocity"`
	
	// Weapons and equipment
	PrimaryWeapon   *Weapon `json:"primary_weapon,omitempty"`
	SecondaryWeapon *Weapon `json:"secondary_weapon,omitempty"`
	Grenades        []Grenade `json:"grenades"`
	
	// Economy
	Money        int `json:"money"`
	
	// Temporary states
	IsFlashed    bool    `json:"is_flashed"`
	IsSmoked     bool    `json:"is_smoked"`
	IsDefusing   bool    `json:"is_defusing"`
	IsPlanting   bool    `json:"is_planting"`
	IsReloading  bool    `json:"is_reloading"`
	
	// Round-specific
	HasBomb      bool    `json:"has_bomb"`
	IsLastAlive  bool    `json:"is_last_alive"`
}

// Vector3 represents a 3D position or direction
type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// NewTeam creates a new team with the given name and players
func NewTeam(name string, players []Player) *Team {
	team := &Team{
		Name:    name,
		Players: players,
		Side:    "CT", // Default to CT
		Score:   0,
		Economy: TeamEconomy{
			ConsecutiveLosses: 0,
			LossBonus:        1400, // Starting loss bonus
		},
		Stats: TeamStats{},
	}
	
	// Initialize player economies
	for i := range team.Players {
		team.Players[i].Economy.Money = 800 // Starting money
		team.Players[i].Team = name
	}
	
	return team
}

// Validate validates the team configuration
func (t *Team) Validate() error {
	if strings.TrimSpace(t.Name) == "" {
		return errors.New("team name is required")
	}
	
	if len(t.Players) != 5 {
		return fmt.Errorf("team must have exactly 5 players, got %d", len(t.Players))
	}
	
	// Validate each player
	playerNames := make(map[string]bool)
	for i, player := range t.Players {
		if err := player.Validate(); err != nil {
			return fmt.Errorf("player %d validation failed: %w", i+1, err)
		}
		
		// Check for duplicate player names
		if playerNames[player.Name] {
			return fmt.Errorf("duplicate player name: %s", player.Name)
		}
		playerNames[player.Name] = true
	}
	
	// Validate side if set
	if t.Side != "" && !IsValidSide(t.Side) {
		return fmt.Errorf("invalid side: %s (must be 'CT' or 'TERRORIST')", t.Side)
	}
	
	return nil
}

// IsValidSide checks if the side is valid
func IsValidSide(side string) bool {
	return strings.EqualFold(side, "CT") || 
		   strings.EqualFold(side, "TERRORIST") ||
		   strings.EqualFold(side, "COUNTER-TERRORIST")
}

// GetAlivePlayers returns all living players on the team
func (t *Team) GetAlivePlayers() []Player {
	var alive []Player
	for _, player := range t.Players {
		if player.State.IsAlive {
			alive = append(alive, player)
		}
	}
	return alive
}

// GetAliveCount returns the number of living players
func (t *Team) GetAliveCount() int {
	count := 0
	for _, player := range t.Players {
		if player.State.IsAlive {
			count++
		}
	}
	return count
}

// GetPlayerByName returns a player by name
func (t *Team) GetPlayerByName(name string) *Player {
	for i := range t.Players {
		if strings.EqualFold(t.Players[i].Name, name) {
			return &t.Players[i]
		}
	}
	return nil
}

// UpdateEconomy updates the team's economic state based on player economies
func (t *Team) UpdateEconomy() {
	totalMoney := 0
	equipmentValue := 0
	rifles := 0
	smgs := 0
	pistols := 0
	snipers := 0
	grenades := 0
	armor := 0
	helmets := 0
	defuseKits := 0
	
	for _, player := range t.Players {
		totalMoney += player.Economy.Money
		equipmentValue += player.Economy.EquipmentValue
		
		// Count equipment
		if player.State.PrimaryWeapon != nil {
			switch player.State.PrimaryWeapon.Type {
			case "rifle":
				rifles++
			case "smg":
				smgs++
			case "sniper":
				snipers++
			}
		}
		
		if player.State.SecondaryWeapon != nil && player.State.SecondaryWeapon.Type == "pistol" {
			pistols++
		}
		
		grenades += len(player.State.Grenades)
		
		if player.State.Armor > 0 {
			armor++
		}
		if player.State.HasHelmet {
			helmets++
		}
		if player.State.HasDefuseKit {
			defuseKits++
		}
	}
	
	// Update team economy
	t.Economy.TotalMoney = totalMoney
	t.Economy.AverageMoney = totalMoney / len(t.Players)
	t.Economy.EquipmentValue = equipmentValue
	t.Economy.Rifles = rifles
	t.Economy.SMGs = smgs
	t.Economy.Pistols = pistols
	t.Economy.Snipers = snipers
	t.Economy.Grenades = grenades
	t.Economy.Armor = armor
	t.Economy.Helmets = helmets
	t.Economy.DefuseKits = defuseKits
}

// HandleRoundWin handles the economic and statistical impact of winning a round
func (t *Team) HandleRoundWin(reason string) {
	t.Score++
	t.RoundsWon++
	
	// Reset consecutive losses
	t.Economy.ConsecutiveLosses = 0
	
	// Award win bonus
	winBonus := 3250 // Standard win bonus
	if reason == "bomb_exploded" || reason == "bomb_defused" {
		winBonus = 3500
	}
	
	for i := range t.Players {
		t.Players[i].Economy.Money += winBonus
		t.Players[i].Economy.MoneyEarned += winBonus
	}
	
	t.Economy.MoneyEarned += winBonus * len(t.Players)
}

// HandleRoundLoss handles the economic impact of losing a round
func (t *Team) HandleRoundLoss() {
	t.Economy.ConsecutiveLosses++
	
	// Calculate loss bonus
	lossBonus := 1400 + (t.Economy.ConsecutiveLosses-1)*500
	if lossBonus > 3400 {
		lossBonus = 3400 // Max loss bonus
	}
	t.Economy.LossBonus = lossBonus
	
	// Award loss bonus to players
	for i := range t.Players {
		t.Players[i].Economy.Money += lossBonus
		t.Players[i].Economy.MoneyEarned += lossBonus
	}
	
	t.Economy.MoneyEarned += lossBonus * len(t.Players)
}

// SwitchSides switches the team to the opposite side
func (t *Team) SwitchSides() {
	if strings.EqualFold(t.Side, "CT") {
		t.Side = "TERRORIST"
	} else {
		t.Side = "CT"
	}
	
	// Update all players' sides
	for i := range t.Players {
		t.Players[i].Side = t.Side
	}
}

// CalculateRating calculates a simple team rating based on performance
func (t *Team) CalculateRating() float64 {
	if t.Stats.RoundsPlayed == 0 {
		return 0.0
	}
	
	// Simple rating calculation
	kdr := float64(t.Stats.Kills) / float64(max(1, t.Stats.Deaths))
	winRate := float64(t.RoundsWon) / float64(t.Stats.RoundsPlayed)
	
	return (kdr * 0.6) + (winRate * 0.4)
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}