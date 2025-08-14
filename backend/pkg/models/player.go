package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Player represents a CS2 player
type Player struct {
	// Basic information
	Name     string `json:"name" binding:"required"`
	SteamID  string `json:"steam_id,omitempty"`
	UserID   int    `json:"user_id,omitempty"`
	Team     string `json:"team"`
	Side     string `json:"side"` // "CT" or "TERRORIST"
	
	// Player configuration
	Role     string `json:"role"` // "entry", "awp", "support", "igl", "lurker"
	
	// Current state
	State    PlayerState `json:"state"`
	
	// Statistics
	Stats    PlayerStats `json:"stats"`
	
	// Economy
	Economy  PlayerEconomy `json:"economy"`
	
	// Performance profile (for realistic generation)
	Profile  PlayerProfile `json:"profile,omitempty"`
}

// PlayerStats represents a player's match statistics
type PlayerStats struct {
	// Basic combat stats
	Kills            int `json:"kills"`
	Deaths           int `json:"deaths"`
	Assists          int `json:"assists"`
	Score            int `json:"score"`
	
	// Damage statistics
	Damage           int `json:"damage"`
	UtilityDamage    int `json:"utility_damage"`
	EnemiesFlashed   int `json:"enemies_flashed"`
	
	// Combat effectiveness
	Headshots        int     `json:"headshots"`
	HeadshotRate     float64 `json:"headshot_rate"`
	Accuracy         float64 `json:"accuracy"`
	
	// Round impact
	FirstKills       int `json:"first_kills"`
	FirstDeaths      int `json:"first_deaths"`
	TradeKills       int `json:"trade_kills"`
	EntryKills       int `json:"entry_kills"`
	
	// Multi-kill rounds
	Multikills2      int `json:"2k_rounds"`
	Multikills3      int `json:"3k_rounds"`
	Multikills4      int `json:"4k_rounds"`
	Multikills5      int `json:"5k_rounds"`
	
	// Objective participation
	BombPlants       int `json:"bomb_plants"`
	BombDefuses      int `json:"bomb_defuses"`
	BombDefuseAttempts int `json:"bomb_defuse_attempts"`
	HostagesRescued  int `json:"hostages_rescued"`
	
	// Equipment usage
	MVPs             int `json:"mvps"`
	MoneySpent       int `json:"money_spent"`
	
	// Utility usage
	GrenadesThrown   map[string]int `json:"grenades_thrown"`
	FlashAssists     int            `json:"flash_assists"`
	
	// Team play
	TeamKills        int `json:"team_kills"`
	TeamDamage       int `json:"team_damage"`
	
	// Performance indicators
	ADR              float64 `json:"adr"` // Average damage per round
	KDRatio          float64 `json:"kd_ratio"`
	Rating           float64 `json:"rating"`
	KAST             float64 `json:"kast"` // Kills, Assists, Survival, Trades percentage
}

// PlayerEconomy represents a player's economic state
type PlayerEconomy struct {
	// Current money
	Money            int `json:"money"`
	MoneySpent       int `json:"money_spent"`
	MoneyEarned      int `json:"money_earned"`
	
	// Equipment value
	EquipmentValue   int `json:"equipment_value"`
	
	// Purchase history
	Purchases        []Purchase `json:"purchases,omitempty"`
	
	// Economic efficiency
	EcoRounds        int     `json:"eco_rounds"`
	ForceBuyRounds   int     `json:"force_buy_rounds"`
	FullBuyRounds    int     `json:"full_buy_rounds"`
	EconomyRating    float64 `json:"economy_rating"`
}

// Purchase represents a single equipment purchase
type Purchase struct {
	Round     int    `json:"round"`
	Item      string `json:"item"`
	Cost      int    `json:"cost"`
	Timestamp string `json:"timestamp,omitempty"`
}

// PlayerProfile represents a player's skill and behavioral profile
type PlayerProfile struct {
	// Skill ratings (0.0 to 1.0)
	AimSkill         float64 `json:"aim_skill"`
	ReflexSpeed      float64 `json:"reflex_speed"`
	GameSense        float64 `json:"game_sense"`
	Positioning      float64 `json:"positioning"`
	Teamwork         float64 `json:"teamwork"`
	UtilityUsage     float64 `json:"utility_usage"`
	
	// Playing style tendencies
	Aggression       float64 `json:"aggression"`       // 0.0 = passive, 1.0 = aggressive
	EconomyDiscipline float64 `json:"economy_discipline"` // Likelihood to save/force buy
	ClutchFactor     float64 `json:"clutch_factor"`    // Performance in clutch situations
	
	// Weapon preferences (0.0 to 1.0)
	RifleSkill       float64 `json:"rifle_skill"`
	AWPSkill         float64 `json:"awp_skill"`
	PistolSkill      float64 `json:"pistol_skill"`
	
	// Role-specific attributes
	EntryFragging    float64 `json:"entry_fragging"`   // Entry fragger effectiveness
	SupportPlay      float64 `json:"support_play"`     // Support role effectiveness
	IGLSkill         float64 `json:"igl_skill"`        // In-game leader abilities
	
	// Consistency factor
	ConsistencyFactor float64 `json:"consistency_factor"` // 0.0 = very inconsistent, 1.0 = very consistent
}

// Weapon represents a weapon with its properties
type Weapon struct {
	Name         string  `json:"name"`
	Type         string  `json:"type"`         // "rifle", "pistol", "sniper", "smg", "shotgun", "machinegun"
	Damage       int     `json:"damage"`
	Accuracy     float64 `json:"accuracy"`
	RangeModifier float64 `json:"range_modifier"`
	PenetrationPower float64 `json:"penetration_power"`
	Price        int     `json:"price"`
	
	// Ammo state
	Ammo         int     `json:"ammo"`
	AmmoReserve  int     `json:"ammo_reserve"`
	MaxAmmo      int     `json:"max_ammo"`
	
	// Weapon attachments/skins (optional)
	Skin         string  `json:"skin,omitempty"`
	StatTrak     bool    `json:"stat_trak"`
}

// Grenade represents a grenade with its properties
type Grenade struct {
	Type         string  `json:"type"`         // "he", "flash", "smoke", "incendiary", "molotov", "decoy"
	Price        int     `json:"price"`
	Damage       int     `json:"damage,omitempty"`
	EffectRadius float64 `json:"effect_radius,omitempty"`
	Duration     float64 `json:"duration,omitempty"` // For smoke/molotov duration
}

// NewPlayer creates a new player with default values
func NewPlayer(name, steamID string) *Player {
	return &Player{
		Name:    name,
		SteamID: steamID,
		Role:    "support", // Default role
		State: PlayerState{
			IsAlive:   true,
			Health:    100,
			Armor:     0,
			HasHelmet: false,
			Position:  Vector3{X: 0, Y: 0, Z: 0},
			Grenades:  make([]Grenade, 0),
		},
		Stats: PlayerStats{
			GrenadesThrown: make(map[string]int),
		},
		Economy: PlayerEconomy{
			Money:     800, // Starting money in CS2
			Purchases: make([]Purchase, 0),
		},
		Profile: DefaultPlayerProfile(),
	}
}

// DefaultPlayerProfile returns a default player profile with average skills
func DefaultPlayerProfile() PlayerProfile {
	return PlayerProfile{
		AimSkill:          0.5,
		ReflexSpeed:       0.5,
		GameSense:         0.5,
		Positioning:       0.5,
		Teamwork:          0.5,
		UtilityUsage:      0.5,
		Aggression:        0.5,
		EconomyDiscipline: 0.5,
		ClutchFactor:      0.5,
		RifleSkill:        0.5,
		AWPSkill:          0.3,
		PistolSkill:       0.5,
		EntryFragging:     0.5,
		SupportPlay:       0.5,
		IGLSkill:          0.3,
		ConsistencyFactor: 0.5,
	}
}

// Validate validates the player configuration
func (p *Player) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("player name is required")
	}
	
	// Validate SteamID format if provided
	if p.SteamID != "" && !IsValidSteamID(p.SteamID) {
		return fmt.Errorf("invalid SteamID format: %s", p.SteamID)
	}
	
	// Validate role if provided
	if p.Role != "" && !IsValidRole(p.Role) {
		return fmt.Errorf("invalid role: %s", p.Role)
	}
	
	// Validate side if provided
	if p.Side != "" && !IsValidSide(p.Side) {
		return fmt.Errorf("invalid side: %s", p.Side)
	}
	
	return nil
}

// IsValidSteamID validates SteamID format (STEAM_X:Y:Z)
func IsValidSteamID(steamID string) bool {
	// Simple regex for SteamID validation
	steamIDRegex := regexp.MustCompile(`^STEAM_[0-1]:[0-1]:\d+$`)
	return steamIDRegex.MatchString(steamID)
}

// IsValidRole checks if the player role is valid
func IsValidRole(role string) bool {
	validRoles := []string{"entry", "awp", "support", "igl", "lurker", "rifler"}
	role = strings.ToLower(role)
	
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}

// IsAlive returns whether the player is alive
func (p *Player) IsAlive() bool {
	return p.State.IsAlive && p.State.Health > 0
}

// TakeDamage applies damage to the player
func (p *Player) TakeDamage(damage int, hasHelmet bool) int {
	if !p.State.IsAlive {
		return 0
	}
	
	actualDamage := damage
	
	// Apply armor reduction if player has armor
	if p.State.Armor > 0 {
		armorReduction := int(float64(damage) * 0.5) // Simplified armor calculation
		actualDamage = damage - armorReduction
		
		// Reduce armor
		armorDamage := min(armorReduction, p.State.Armor)
		p.State.Armor -= armorDamage
		
		// Remove helmet if armor reaches 0
		if p.State.Armor == 0 {
			p.State.HasHelmet = false
		}
	}
	
	// Apply damage to health
	p.State.Health -= actualDamage
	if p.State.Health <= 0 {
		p.State.Health = 0
		p.State.IsAlive = false
	}
	
	return actualDamage
}

// Heal restores health to the player
func (p *Player) Heal(amount int) {
	p.State.Health = min(100, p.State.Health+amount)
}

// GiveArmor gives armor to the player
func (p *Player) GiveArmor(amount int, withHelmet bool) {
	p.State.Armor = min(100, p.State.Armor+amount)
	if withHelmet {
		p.State.HasHelmet = true
	}
}

// AddWeapon gives a weapon to the player
func (p *Player) AddWeapon(weapon Weapon) {
	switch weapon.Type {
	case "rifle", "smg", "sniper", "shotgun", "machinegun":
		p.State.PrimaryWeapon = &weapon
	case "pistol":
		p.State.SecondaryWeapon = &weapon
	}
}

// AddGrenade gives a grenade to the player
func (p *Player) AddGrenade(grenade Grenade) {
	// Check if player can carry more grenades (max 4 in CS2)
	if len(p.State.Grenades) < 4 {
		p.State.Grenades = append(p.State.Grenades, grenade)
	}
}

// RemoveGrenade removes a grenade of the specified type
func (p *Player) RemoveGrenade(grenadeType string) bool {
	for i, grenade := range p.State.Grenades {
		if grenade.Type == grenadeType {
			p.State.Grenades = append(p.State.Grenades[:i], p.State.Grenades[i+1:]...)
			return true
		}
	}
	return false
}

// Purchase handles equipment purchase
func (p *Player) Purchase(item string, cost int, round int) error {
	if p.Economy.Money < cost {
		return fmt.Errorf("insufficient funds: need %d, have %d", cost, p.Economy.Money)
	}
	
	p.Economy.Money -= cost
	p.Economy.MoneySpent += cost
	
	// Record purchase
	purchase := Purchase{
		Round: round,
		Item:  item,
		Cost:  cost,
	}
	p.Economy.Purchases = append(p.Economy.Purchases, purchase)
	
	return nil
}

// Spawn resets the player state for a new round
func (p *Player) Spawn(position Vector3) {
	p.State.IsAlive = true
	p.State.Health = 100
	p.State.Position = position
	p.State.IsFlashed = false
	p.State.IsSmoked = false
	p.State.IsDefusing = false
	p.State.IsPlanting = false
	p.State.IsReloading = false
	p.State.HasBomb = false
	p.State.IsLastAlive = false
}

// Kill marks the player as dead
func (p *Player) Kill() {
	p.State.IsAlive = false
	p.State.Health = 0
	p.Stats.Deaths++
}

// AddKill records a kill for this player
func (p *Player) AddKill(headshot bool, weapon string) {
	p.Stats.Kills++
	p.Stats.Score += 100 // Standard kill score
	
	if headshot {
		p.Stats.Headshots++
	}
	
	// Update headshot rate
	if p.Stats.Kills > 0 {
		p.Stats.HeadshotRate = float64(p.Stats.Headshots) / float64(p.Stats.Kills)
	}
}

// AddAssist records an assist for this player
func (p *Player) AddAssist() {
	p.Stats.Assists++
	p.Stats.Score += 50 // Standard assist score
}

// AddDamage records damage dealt by this player
func (p *Player) AddDamage(damage int) {
	p.Stats.Damage += damage
}

// CalculateRating calculates a simplified player rating
func (p *Player) CalculateRating(roundsPlayed int) float64 {
	if roundsPlayed == 0 {
		return 0.0
	}
	
	// Calculate ADR (Average Damage per Round)
	p.Stats.ADR = float64(p.Stats.Damage) / float64(roundsPlayed)
	
	// Calculate K/D ratio
	deaths := p.Stats.Deaths
	if deaths == 0 {
		deaths = 1 // Avoid division by zero
	}
	p.Stats.KDRatio = float64(p.Stats.Kills) / float64(deaths)
	
	// Simple rating calculation (similar to HLTV 1.0 rating)
	killRating := float64(p.Stats.Kills) / float64(roundsPlayed)
	survivalRating := float64(roundsPlayed-p.Stats.Deaths) / float64(roundsPlayed)
	damageRating := p.Stats.ADR / 100.0
	
	p.Stats.Rating = (killRating + 0.7*survivalRating + damageRating) / 2.7
	return p.Stats.Rating
}

// GetEquipmentValue calculates the total value of player's equipment
func (p *Player) GetEquipmentValue() int {
	total := 0
	
	if p.State.PrimaryWeapon != nil {
		total += p.State.PrimaryWeapon.Price
	}
	
	if p.State.SecondaryWeapon != nil {
		total += p.State.SecondaryWeapon.Price
	}
	
	for _, grenade := range p.State.Grenades {
		total += grenade.Price
	}
	
	// Add armor value
	if p.State.Armor > 0 {
		if p.State.HasHelmet {
			total += 1000 // Helmet + armor
		} else {
			total += 650 // Armor only
		}
	}
	
	// Add defuse kit value
	if p.State.HasDefuseKit {
		total += 400
	}
	
	p.Economy.EquipmentValue = total
	return total
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}