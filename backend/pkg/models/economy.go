package models

import (
	"fmt"
	"strings"
)

// EconomyManager manages the CS2 economy system
type EconomyManager struct {
	WeaponPrices    map[string]int
	UtilityPrices   map[string]int
	RoundWinBonus   map[string]int
	KillRewards     map[string]int
	ObjectiveRewards map[string]int
}

// WeaponInfo represents weapon information and pricing
type WeaponInfo struct {
	Name         string  `json:"name"`
	DisplayName  string  `json:"display_name"`
	Type         string  `json:"type"`
	Price        int     `json:"price"`
	KillReward   int     `json:"kill_reward"`
	Damage       int     `json:"damage"`
	Accuracy     float64 `json:"accuracy"`
	ArmorPen     float64 `json:"armor_penetration"`
	RangeModifier float64 `json:"range_modifier"`
	Firerate     float64 `json:"firerate"`
	MovementSpeed float64 `json:"movement_speed"`
	Team         string  `json:"team"` // "both", "ct", "t"
}

// UtilityInfo represents utility/equipment information and pricing
type UtilityInfo struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Type        string `json:"type"` // "grenade", "armor", "utility"
	Price       int    `json:"price"`
	Team        string `json:"team"` // "both", "ct", "t"
}

// EconomyState represents the current economic state
type EconomyState struct {
	Round           int `json:"round"`
	CTLossStreak    int `json:"ct_loss_streak"`
	TLossStreak     int `json:"t_loss_streak"`
	CTConsecutiveWins int `json:"ct_consecutive_wins"`
	TConsecutiveWins  int `json:"t_consecutive_wins"`
}

// NewEconomyManager creates a new economy manager with CS2 prices
func NewEconomyManager() *EconomyManager {
	return &EconomyManager{
		WeaponPrices:     getCS2WeaponPrices(),
		UtilityPrices:    getCS2UtilityPrices(),
		RoundWinBonus:    getRoundWinBonuses(),
		KillRewards:      getKillRewards(),
		ObjectiveRewards: getObjectiveRewards(),
	}
}

// getCS2WeaponPrices returns the current CS2 weapon prices
func getCS2WeaponPrices() map[string]int {
	return map[string]int{
		// Pistols
		"glock":         200,
		"usp_silencer":  200,
		"p250":          300,
		"tec9":          500,
		"fiveseven":     500,
		"cz75a":         500,
		"deagle":        700,
		"revolver":      600,
		
		// SMGs
		"mac10":         1050,
		"mp9":           1250,
		"mp7":           1500,
		"ump45":         1200,
		"p90":           2350,
		"bizon":         1400,
		
		// Rifles
		"famas":         2050,
		"galil":         1800,
		"m4a4":          3100,
		"m4a1_silencer": 2900,
		"ak47":          2700,
		"sg556":         3000,
		"aug":           3300,
		
		// Sniper Rifles
		"ssg08":         1700,
		"awp":           4750,
		"g3sg1":         5000,
		"scar20":        5000,
		
		// Shotguns
		"nova":          1050,
		"xm1014":        2000,
		"sawedoff":      1100,
		"mag7":          1300,
		
		// Machine Guns
		"negev":         1700,
		"m249":          5200,
	}
}

// getCS2UtilityPrices returns the current CS2 utility prices
func getCS2UtilityPrices() map[string]int {
	return map[string]int{
		// Grenades
		"hegrenade":     300,
		"flashbang":     200,
		"smokegrenade":  300,
		"incgrenade":    600,
		"molotov":       400,
		"decoy":         50,
		
		// Armor
		"vest":          650,
		"vesthelm":      1000,
		
		// Utilities
		"defuser":       400,
		"zeus":          200,
	}
}

// getRoundWinBonuses returns round win bonus amounts
func getRoundWinBonuses() map[string]int {
	return map[string]int{
		"elimination":   3250,
		"bomb_defused":  3500,
		"bomb_exploded": 3500,
		"time_expired":  3250,
	}
}

// getKillRewards returns kill reward amounts by weapon type
func getKillRewards() map[string]int {
	return map[string]int{
		"pistol":    300,
		"smg":       600,
		"rifle":     300,
		"sniper":    100, // AWP, auto-snipers
		"shotgun":   900,
		"knife":     1500,
		"grenade":   300,
		"zeus":      300,
	}
}

// getObjectiveRewards returns objective-based reward amounts
func getObjectiveRewards() map[string]int {
	return map[string]int{
		"bomb_plant":     300,
		"bomb_defuse":    250,
		"hostage_rescue": 200,
	}
}

// GetWeaponInfo returns detailed weapon information
func (em *EconomyManager) GetWeaponInfo() map[string]WeaponInfo {
	return map[string]WeaponInfo{
		"ak47": {
			Name:          "ak47",
			DisplayName:   "AK-47",
			Type:          "rifle",
			Price:         2700,
			KillReward:    300,
			Damage:        36,
			Accuracy:      0.75,
			ArmorPen:      0.775,
			RangeModifier: 0.98,
			Firerate:      600,
			MovementSpeed: 221,
			Team:          "t",
		},
		"m4a4": {
			Name:          "m4a4",
			DisplayName:   "M4A4",
			Type:          "rifle",
			Price:         3100,
			KillReward:    300,
			Damage:        33,
			Accuracy:      0.78,
			ArmorPen:      0.70,
			RangeModifier: 0.97,
			Firerate:      666,
			MovementSpeed: 225,
			Team:          "ct",
		},
		"m4a1_silencer": {
			Name:          "m4a1_silencer",
			DisplayName:   "M4A1-S",
			Type:          "rifle",
			Price:         2900,
			KillReward:    300,
			Damage:        38,
			Accuracy:      0.82,
			ArmorPen:      0.70,
			RangeModifier: 0.99,
			Firerate:      600,
			MovementSpeed: 225,
			Team:          "ct",
		},
		"awp": {
			Name:          "awp",
			DisplayName:   "AWP",
			Type:          "sniper",
			Price:         4750,
			KillReward:    100,
			Damage:        115,
			Accuracy:      0.99,
			ArmorPen:      0.975,
			RangeModifier: 0.99,
			Firerate:      41,
			MovementSpeed: 200,
			Team:          "both",
		},
		"glock": {
			Name:          "glock",
			DisplayName:   "Glock-18",
			Type:          "pistol",
			Price:         200,
			KillReward:    300,
			Damage:        28,
			Accuracy:      0.58,
			ArmorPen:      0.475,
			RangeModifier: 0.75,
			Firerate:      400,
			MovementSpeed: 240,
			Team:          "t",
		},
		"usp_silencer": {
			Name:          "usp_silencer",
			DisplayName:   "USP-S",
			Type:          "pistol",
			Price:         200,
			KillReward:    300,
			Damage:        35,
			Accuracy:      0.75,
			ArmorPen:      0.50,
			RangeModifier: 0.79,
			Firerate:      352,
			MovementSpeed: 240,
			Team:          "ct",
		},
	}
}

// GetUtilityInfo returns detailed utility information
func (em *EconomyManager) GetUtilityInfo() map[string]UtilityInfo {
	return map[string]UtilityInfo{
		"hegrenade": {
			Name:        "hegrenade",
			DisplayName: "HE Grenade",
			Type:        "grenade",
			Price:       300,
			Team:        "both",
		},
		"flashbang": {
			Name:        "flashbang",
			DisplayName: "Flashbang",
			Type:        "grenade",
			Price:       200,
			Team:        "both",
		},
		"smokegrenade": {
			Name:        "smokegrenade",
			DisplayName: "Smoke Grenade",
			Type:        "grenade",
			Price:       300,
			Team:        "both",
		},
		"incgrenade": {
			Name:        "incgrenade",
			DisplayName: "Incendiary Grenade",
			Type:        "grenade",
			Price:       600,
			Team:        "ct",
		},
		"molotov": {
			Name:        "molotov",
			DisplayName: "Molotov Cocktail",
			Type:        "grenade",
			Price:       400,
			Team:        "t",
		},
		"vest": {
			Name:        "vest",
			DisplayName: "Kevlar Vest",
			Type:        "armor",
			Price:       650,
			Team:        "both",
		},
		"vesthelm": {
			Name:        "vesthelm",
			DisplayName: "Kevlar + Helmet",
			Type:        "armor",
			Price:       1000,
			Team:        "both",
		},
		"defuser": {
			Name:        "defuser",
			DisplayName: "Defuse Kit",
			Type:        "utility",
			Price:       400,
			Team:        "ct",
		},
	}
}

// CalculateLossBonus calculates the loss bonus for a team
func (em *EconomyManager) CalculateLossBonus(consecutiveLosses int) int {
	baseLossBonus := 1400
	bonusIncrement := 500
	maxLossBonus := 3400
	
	lossBonus := baseLossBonus + (consecutiveLosses-1)*bonusIncrement
	if lossBonus > maxLossBonus {
		lossBonus = maxLossBonus
	}
	
	return lossBonus
}

// CalculateWinBonus calculates the win bonus for a team
func (em *EconomyManager) CalculateWinBonus(winReason string) int {
	if bonus, exists := em.RoundWinBonus[winReason]; exists {
		return bonus
	}
	return em.RoundWinBonus["elimination"] // Default bonus
}

// CalculateKillReward calculates the kill reward for a weapon
func (em *EconomyManager) CalculateKillReward(weaponName string) int {
	// First try to get exact weapon reward
	if reward, exists := em.KillRewards[weaponName]; exists {
		return reward
	}
	
	// Try to get reward by weapon type
	weaponInfo := em.GetWeaponInfo()
	if info, exists := weaponInfo[weaponName]; exists {
		if reward, exists := em.KillRewards[info.Type]; exists {
			return reward
		}
	}
	
	// Default kill reward
	return 300
}

// GetWeaponPrice returns the price of a weapon
func (em *EconomyManager) GetWeaponPrice(weaponName string) int {
	if price, exists := em.WeaponPrices[weaponName]; exists {
		return price
	}
	return 0
}

// GetUtilityPrice returns the price of utility/equipment
func (em *EconomyManager) GetUtilityPrice(utilityName string) int {
	if price, exists := em.UtilityPrices[utilityName]; exists {
		return price
	}
	return 0
}

// CanAfford checks if a player can afford an item
func (em *EconomyManager) CanAfford(playerMoney int, itemName string) bool {
	price := em.GetWeaponPrice(itemName)
	if price == 0 {
		price = em.GetUtilityPrice(itemName)
	}
	return playerMoney >= price
}

// GetOptimalBuy suggests an optimal buy for a player
func (em *EconomyManager) GetOptimalBuy(player *Player, teamEconomy *TeamEconomy, roundType string) []string {
	money := player.Economy.Money
	var buy []string
	
	// Determine buy type based on money and team economy
	avgMoney := teamEconomy.AverageMoney
	
	if avgMoney >= 5000 {
		// Full buy round
		buy = em.getFullBuy(player, money)
	} else if avgMoney >= 2500 {
		// Force buy round
		buy = em.getForceBuy(player, money)
	} else {
		// Eco round
		buy = em.getEcoBuy(player, money)
	}
	
	return buy
}

// getFullBuy returns a full buy recommendation
func (em *EconomyManager) getFullBuy(player *Player, money int) []string {
	var buy []string
	remaining := money
	
	// Primary weapon based on side and role
	var primary string
	if strings.EqualFold(player.Side, "CT") {
		if player.Role == "awp" && remaining >= 4750 {
			primary = "awp"
		} else if remaining >= 3100 {
			primary = "m4a4"
		} else if remaining >= 2900 {
			primary = "m4a1_silencer"
		}
	} else { // Terrorist
		if player.Role == "awp" && remaining >= 4750 {
			primary = "awp"
		} else if remaining >= 2700 {
			primary = "ak47"
		}
	}
	
	if primary != "" {
		buy = append(buy, primary)
		remaining -= em.GetWeaponPrice(primary)
	}
	
	// Armor
	if remaining >= 1000 {
		buy = append(buy, "vesthelm")
		remaining -= 1000
	} else if remaining >= 650 {
		buy = append(buy, "vest")
		remaining -= 650
	}
	
	// Grenades
	if remaining >= 300 {
		buy = append(buy, "smokegrenade")
		remaining -= 300
	}
	if remaining >= 200 {
		buy = append(buy, "flashbang")
		remaining -= 200
	}
	if remaining >= 300 {
		buy = append(buy, "hegrenade")
		remaining -= 300
	}
	
	// Defuse kit for CT
	if strings.EqualFold(player.Side, "CT") && remaining >= 400 {
		buy = append(buy, "defuser")
		remaining -= 400
	}
	
	return buy
}

// getForceBuy returns a force buy recommendation
func (em *EconomyManager) getForceBuy(player *Player, money int) []string {
	var buy []string
	remaining := money
	
	// Cheaper primary weapons
	var primary string
	if strings.EqualFold(player.Side, "CT") {
		if remaining >= 2050 {
			primary = "famas"
		} else if remaining >= 1250 {
			primary = "mp9"
		}
	} else { // Terrorist
		if remaining >= 1800 {
			primary = "galil"
		} else if remaining >= 1050 {
			primary = "mac10"
		}
	}
	
	if primary != "" {
		buy = append(buy, primary)
		remaining -= em.GetWeaponPrice(primary)
	}
	
	// Armor - prioritize vest
	if remaining >= 650 {
		buy = append(buy, "vest")
		remaining -= 650
	}
	
	// One utility
	if remaining >= 200 {
		buy = append(buy, "flashbang")
		remaining -= 200
	}
	
	return buy
}

// getEcoBuy returns an eco round recommendation
func (em *EconomyManager) getEcoBuy(player *Player, money int) []string {
	var buy []string
	remaining := money
	
	// Upgraded pistol or cheap SMG
	if remaining >= 700 {
		buy = append(buy, "deagle")
		remaining -= 700
	} else if remaining >= 500 {
		if strings.EqualFold(player.Side, "CT") {
			buy = append(buy, "fiveseven")
		} else {
			buy = append(buy, "tec9")
		}
		remaining -= 500
	}
	
	// Minimal utility
	if remaining >= 200 {
		buy = append(buy, "flashbang")
		remaining -= 200
	}
	
	return buy
}

// CalculateEquipmentValue calculates the total value of equipment
func (em *EconomyManager) CalculateEquipmentValue(weapons []string, utilities []string) int {
	total := 0
	
	for _, weapon := range weapons {
		total += em.GetWeaponPrice(weapon)
	}
	
	for _, utility := range utilities {
		total += em.GetUtilityPrice(utility)
	}
	
	return total
}

// IsValidWeaponForSide checks if a weapon is available for the specified side
func (em *EconomyManager) IsValidWeaponForSide(weaponName, side string) bool {
	weaponInfo := em.GetWeaponInfo()
	if info, exists := weaponInfo[weaponName]; exists {
		return info.Team == "both" || strings.EqualFold(info.Team, side)
	}
	return false
}

// GetWeaponsByType returns weapons filtered by type
func (em *EconomyManager) GetWeaponsByType(weaponType string) []WeaponInfo {
	var weapons []WeaponInfo
	weaponInfo := em.GetWeaponInfo()
	
	for _, info := range weaponInfo {
		if strings.EqualFold(info.Type, weaponType) {
			weapons = append(weapons, info)
		}
	}
	
	return weapons
}

// FormatMoney formats money amount for display
func (em *EconomyManager) FormatMoney(amount int) string {
	if amount >= 1000 {
		return fmt.Sprintf("$%dk", amount/1000)
	}
	return fmt.Sprintf("$%d", amount)
}