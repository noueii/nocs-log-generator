package generator

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/noueii/nocs-log-generator/backend/pkg/models"
)

// EconomyManager handles team and player money management
type EconomyManager struct {
	rng           *rand.Rand
	economySystem *models.EconomyManager
}

// NewEconomyManager creates a new economy manager
func NewEconomyManager(rng *rand.Rand) *EconomyManager {
	return &EconomyManager{
		rng:           rng,
		economySystem: models.NewEconomyManager(),
	}
}

// HandleRoundEnd processes economy changes after a round ends
func (em *EconomyManager) HandleRoundEnd(match *models.Match, state *models.MatchState, result *RoundResult, events []models.GameEvent) error {
	// Convert side winner to team name
	var winningTeamName, losingTeamName string
	
	if result.Winner == "CT" || result.Winner == "COUNTER-TERRORIST" {
		winningTeam := em.getTeamBySide(match, "CT")
		losingTeam := em.getTeamBySide(match, "TERRORIST")
		if winningTeam != nil && losingTeam != nil {
			winningTeamName = winningTeam.Name
			losingTeamName = losingTeam.Name
		} else {
			return fmt.Errorf("could not find CT/T teams for economy processing")
		}
	} else if result.Winner == "TERRORIST" {
		winningTeam := em.getTeamBySide(match, "TERRORIST") 
		losingTeam := em.getTeamBySide(match, "CT")
		if winningTeam != nil && losingTeam != nil {
			winningTeamName = winningTeam.Name
			losingTeamName = losingTeam.Name
		} else {
			return fmt.Errorf("could not find CT/T teams for economy processing")
		}
	} else {
		// Winner is team name directly
		winningTeam := em.getTeamByName(match, result.Winner)
		losingTeam := em.getLosingTeam(match, result.Winner)
		if winningTeam != nil && losingTeam != nil {
			winningTeamName = winningTeam.Name
			losingTeamName = losingTeam.Name
		} else {
			return fmt.Errorf("could not find teams %s for economy processing", result.Winner)
		}
	}
	
	winningTeam := em.getTeamByName(match, winningTeamName)
	losingTeam := em.getTeamByName(match, losingTeamName)
	
	if winningTeam == nil || losingTeam == nil {
		return fmt.Errorf("could not find teams for economy processing: winner=%s, loser=%s", winningTeamName, losingTeamName)
	}
	
	// Process win bonuses
	em.awardWinBonus(winningTeam, state, result.Reason, events)
	
	// Process loss bonuses
	em.awardLossBonus(losingTeam, state)
	
	// Process kill rewards
	em.awardKillRewards(match, events)
	
	// Process objective rewards
	em.awardObjectiveRewards(match, events)
	
	// Cap money at maximum
	em.capPlayerMoney(match, state)
	
	// Update team economy statistics
	em.updateTeamEconomies(match, state)
	
	return nil
}

// PlanTeamBuys determines what each team should buy based on their economy
func (em *EconomyManager) PlanTeamBuys(match *models.Match, state *models.MatchState, roundNum int) (map[string]string, error) {
	teamBuyTypes := make(map[string]string)
	
	for _, team := range match.Teams {
		teamEconomy := state.TeamEconomies[team.Name]
		buyType := em.determineBuyStrategy(teamEconomy, roundNum, team.Side)
		teamBuyTypes[team.Name] = buyType
	}
	
	return teamBuyTypes, nil
}

// ExecutePlayerBuy handles individual player purchases
func (em *EconomyManager) ExecutePlayerBuy(player *models.Player, playerState *models.PlayerState, buyType string, roundNum int) ([]string, error) {
	var purchases []string
	startMoney := playerState.Money
	
	// Determine buy priority based on role and team strategy
	buyList := em.generateBuyList(player, buyType, startMoney)
	
	// Execute purchases in priority order
	for _, item := range buyList {
		cost := em.getItemCost(item)
		if playerState.Money >= cost {
			// Make purchase
			if err := em.purchaseItem(player, playerState, item, cost); err == nil {
				purchases = append(purchases, item)
				playerState.Money -= cost
				
				// Record purchase in player economy
				purchase := models.Purchase{
					Round: roundNum,
					Item:  item,
					Cost:  cost,
				}
				player.Economy.Purchases = append(player.Economy.Purchases, purchase)
				player.Economy.MoneySpent += cost
			}
		}
	}
	
	return purchases, nil
}

// awardWinBonus gives money to the winning team
func (em *EconomyManager) awardWinBonus(team *models.Team, state *models.MatchState, reason string, events []models.GameEvent) {
	bonus := em.economySystem.CalculateWinBonus(reason)
	
	for i := range team.Players {
		playerState := state.PlayerStates[team.Players[i].Name]
		playerState.Money += bonus
		team.Players[i].Economy.MoneyEarned += bonus
	}
	
	// Reset loss streak
	teamEconomy := state.TeamEconomies[team.Name]
	teamEconomy.ConsecutiveLosses = 0
}

// awardLossBonus gives loss bonus to the losing team
func (em *EconomyManager) awardLossBonus(team *models.Team, state *models.MatchState) {
	teamEconomy := state.TeamEconomies[team.Name]
	teamEconomy.ConsecutiveLosses++
	
	lossBonus := em.economySystem.CalculateLossBonus(teamEconomy.ConsecutiveLosses)
	teamEconomy.LossBonus = lossBonus
	
	for i := range team.Players {
		playerState := state.PlayerStates[team.Players[i].Name]
		playerState.Money += lossBonus
		team.Players[i].Economy.MoneyEarned += lossBonus
	}
}

// awardKillRewards gives money for kills
func (em *EconomyManager) awardKillRewards(match *models.Match, events []models.GameEvent) {
	for _, event := range events {
		if killEvent, ok := event.(*models.KillEvent); ok {
			reward := em.economySystem.CalculateKillReward(killEvent.Weapon)
			
			// Find the attacker in the match and award money
			attacker := em.findPlayerInMatch(match, killEvent.Attacker.Name)
			if attacker != nil {
				// Money is already managed in player state, but track in economy
				attacker.Economy.MoneyEarned += reward
			}
		}
	}
}

// awardObjectiveRewards gives money for objectives
func (em *EconomyManager) awardObjectiveRewards(match *models.Match, events []models.GameEvent) {
	for _, event := range events {
		switch e := event.(type) {
		case *models.BombPlantEvent:
			// Award bomb plant money
			planter := em.findPlayerInMatch(match, e.Player.Name)
			if planter != nil {
				reward := em.economySystem.ObjectiveRewards["bomb_plant"]
				planter.Economy.MoneyEarned += reward
			}
			
		case *models.BombDefuseEvent:
			// Award bomb defuse money
			defuser := em.findPlayerInMatch(match, e.Player.Name)
			if defuser != nil {
				reward := em.economySystem.ObjectiveRewards["bomb_defuse"]
				defuser.Economy.MoneyEarned += reward
			}
		}
	}
}

// determineBuyStrategy decides what type of buy the team should make
func (em *EconomyManager) determineBuyStrategy(economy *models.TeamEconomy, roundNum int, side string) string {
	avgMoney := economy.AverageMoney
	
	// Consider various factors
	isImportantRound := em.isImportantRound(roundNum)
	hasGoodEconomy := avgMoney >= 4000
	hasOkayEconomy := avgMoney >= 2500
	consecutiveLosses := economy.ConsecutiveLosses
	
	// Anti-eco after enemy eco
	if consecutiveLosses >= 2 && avgMoney >= 2000 && isImportantRound {
		return "anti_eco"
	}
	
	// Full buy conditions
	if hasGoodEconomy || (hasOkayEconomy && isImportantRound) {
		return "full_buy"
	}
	
	// Force buy conditions
	if hasOkayEconomy || (avgMoney >= 1500 && isImportantRound) {
		return "force_buy"
	}
	
	// Semi-eco (light buy)
	if avgMoney >= 1000 {
		return "semi_eco"
	}
	
	// Pure eco
	return "eco"
}

// generateBuyList creates a prioritized buy list for a player
func (em *EconomyManager) generateBuyList(player *models.Player, buyType string, money int) []string {
	var buyList []string
	
	side := strings.ToUpper(player.Side)
	role := player.Role
	
	switch buyType {
	case "full_buy":
		buyList = em.generateFullBuy(side, role, money)
	case "force_buy":
		buyList = em.generateForceBuy(side, role, money)
	case "anti_eco":
		buyList = em.generateAntiEcoBuy(side, role, money)
	case "semi_eco":
		buyList = em.generateSemiEcoBuy(side, role, money)
	case "eco":
		buyList = em.generateEcoBuy(side, role, money)
	default:
		buyList = em.generateDefaultBuy(side, role, money)
	}
	
	return buyList
}

// generateFullBuy creates a full buy list
func (em *EconomyManager) generateFullBuy(side, role string, money int) []string {
	var buyList []string
	
	// Armor first
	if money >= 1000 {
		buyList = append(buyList, "vesthelm")
	} else if money >= 650 {
		buyList = append(buyList, "vest")
	}
	
	// Primary weapon based on side and role
	if role == "awp" && money >= 4750 {
		buyList = append(buyList, "awp")
	} else if side == "CT" {
		if money >= 3100 {
			buyList = append(buyList, "m4a4")
		} else if money >= 2900 {
			buyList = append(buyList, "m4a1_silencer")
		}
	} else { // Terrorist
		if money >= 2700 {
			buyList = append(buyList, "ak47")
		}
	}
	
	// Utilities
	buyList = append(buyList, "smokegrenade")
	buyList = append(buyList, "flashbang")
	buyList = append(buyList, "hegrenade")
	
	// Defuse kit for CT
	if side == "CT" {
		buyList = append(buyList, "defuser")
	}
	
	return buyList
}

// generateForceBuy creates a force buy list
func (em *EconomyManager) generateForceBuy(side, role string, money int) []string {
	var buyList []string
	
	// Armor
	if money >= 650 {
		buyList = append(buyList, "vest")
	}
	
	// Cheaper weapons
	if side == "CT" {
		if money >= 2050 {
			buyList = append(buyList, "famas")
		} else if money >= 1250 {
			buyList = append(buyList, "mp9")
		}
	} else { // Terrorist
		if money >= 1800 {
			buyList = append(buyList, "galil")
		} else if money >= 1050 {
			buyList = append(buyList, "mac10")
		}
	}
	
	// Minimal utility
	buyList = append(buyList, "flashbang")
	
	return buyList
}

// generateAntiEcoBuy creates an anti-eco buy list
func (em *EconomyManager) generateAntiEcoBuy(side, role string, money int) []string {
	var buyList []string
	
	// Light armor
	if money >= 650 {
		buyList = append(buyList, "vest")
	}
	
	// SMGs for anti-eco
	if money >= 1200 {
		if side == "CT" {
			buyList = append(buyList, "mp9")
		} else {
			buyList = append(buyList, "mac10")
		}
	}
	
	// More grenades for anti-eco
	buyList = append(buyList, "hegrenade")
	buyList = append(buyList, "flashbang")
	
	return buyList
}

// generateSemiEcoBuy creates a semi-eco buy list
func (em *EconomyManager) generateSemiEcoBuy(side, role string, money int) []string {
	var buyList []string
	
	// Upgraded pistol
	if money >= 700 {
		buyList = append(buyList, "deagle")
	} else if money >= 500 {
		if side == "CT" {
			buyList = append(buyList, "fiveseven")
		} else {
			buyList = append(buyList, "tec9")
		}
	}
	
	// Single utility
	if money >= 200 {
		buyList = append(buyList, "flashbang")
	}
	
	return buyList
}

// generateEcoBuy creates an eco buy list
func (em *EconomyManager) generateEcoBuy(side, role string, money int) []string {
	var buyList []string
	
	// Maybe upgrade pistol if very cheap
	if money >= 500 && em.rng.Float64() < 0.3 { // 30% chance
		if side == "CT" {
			buyList = append(buyList, "p250")
		} else {
			buyList = append(buyList, "p250")
		}
	}
	
	return buyList
}

// generateDefaultBuy creates a default buy list
func (em *EconomyManager) generateDefaultBuy(side, role string, money int) []string {
	return em.generateForceBuy(side, role, money)
}

// purchaseItem applies a purchased item to the player
func (em *EconomyManager) purchaseItem(player *models.Player, playerState *models.PlayerState, item string, cost int) error {
	// Get item information
	weaponInfo := em.economySystem.GetWeaponInfo()
	utilityInfo := em.economySystem.GetUtilityInfo()
	
	if info, exists := weaponInfo[item]; exists {
		weapon := &models.Weapon{
			Name:     info.Name,
			Type:     info.Type,
			Price:    info.Price,
			Damage:   info.Damage,
			Accuracy: info.Accuracy,
			Ammo:     30, // Default ammo count
		}
		
		switch info.Type {
		case "pistol":
			playerState.SecondaryWeapon = weapon
		default:
			playerState.PrimaryWeapon = weapon
		}
		
	} else if info, exists := utilityInfo[item]; exists {
		switch info.Type {
		case "armor":
			if item == "vesthelm" {
				playerState.Armor = 100
				playerState.HasHelmet = true
			} else if item == "vest" {
				playerState.Armor = 100
			}
		case "utility":
			if item == "defuser" {
				playerState.HasDefuseKit = true
			}
		case "grenade":
			if len(playerState.Grenades) < 4 { // Max 4 grenades
				grenade := models.Grenade{
					Type:  info.Name,
					Price: info.Price,
				}
				playerState.Grenades = append(playerState.Grenades, grenade)
			}
		}
	} else {
		return fmt.Errorf("unknown item: %s", item)
	}
	
	return nil
}

// Utility methods

func (em *EconomyManager) getItemCost(item string) int {
	cost := em.economySystem.GetWeaponPrice(item)
	if cost == 0 {
		cost = em.economySystem.GetUtilityPrice(item)
	}
	return cost
}

func (em *EconomyManager) capPlayerMoney(match *models.Match, state *models.MatchState) {
	maxMoney := 16000 // CS2 money cap
	
	for _, team := range match.Teams {
		for _, player := range team.Players {
			if playerState := state.PlayerStates[player.Name]; playerState != nil {
				if playerState.Money > maxMoney {
					playerState.Money = maxMoney
				}
			}
		}
	}
}

func (em *EconomyManager) updateTeamEconomies(match *models.Match, state *models.MatchState) {
	for _, team := range match.Teams {
		teamEconomy := state.TeamEconomies[team.Name]
		totalMoney := 0
		equipmentValue := 0
		
		for _, player := range team.Players {
			if playerState := state.PlayerStates[player.Name]; playerState != nil {
				totalMoney += playerState.Money
				equipmentValue += em.calculateEquipmentValue(playerState)
			}
		}
		
		teamEconomy.TotalMoney = totalMoney
		teamEconomy.AverageMoney = totalMoney / len(team.Players)
		teamEconomy.EquipmentValue = equipmentValue
	}
}

func (em *EconomyManager) calculateEquipmentValue(playerState *models.PlayerState) int {
	value := 0
	
	if playerState.PrimaryWeapon != nil {
		value += playerState.PrimaryWeapon.Price
	}
	if playerState.SecondaryWeapon != nil {
		value += playerState.SecondaryWeapon.Price
	}
	for _, grenade := range playerState.Grenades {
		value += grenade.Price
	}
	if playerState.Armor > 0 {
		if playerState.HasHelmet {
			value += 1000
		} else {
			value += 650
		}
	}
	if playerState.HasDefuseKit {
		value += 400
	}
	
	return value
}

func (em *EconomyManager) getTeamByName(match *models.Match, name string) *models.Team {
	for i := range match.Teams {
		if match.Teams[i].Name == name {
			return &match.Teams[i]
		}
	}
	return nil
}

func (em *EconomyManager) getTeamBySide(match *models.Match, side string) *models.Team {
	for i := range match.Teams {
		if match.Teams[i].Side == side {
			return &match.Teams[i]
		}
	}
	return nil
}

func (em *EconomyManager) getLosingTeam(match *models.Match, winnerName string) *models.Team {
	for i := range match.Teams {
		if match.Teams[i].Name != winnerName {
			return &match.Teams[i]
		}
	}
	return nil
}

func (em *EconomyManager) findPlayerInMatch(match *models.Match, playerName string) *models.Player {
	for i := range match.Teams {
		for j := range match.Teams[i].Players {
			if match.Teams[i].Players[j].Name == playerName {
				return &match.Teams[i].Players[j]
			}
		}
	}
	return nil
}

func (em *EconomyManager) isImportantRound(roundNum int) bool {
	// Pistol rounds (1st and 16th in MR12)
	if roundNum == 1 || roundNum == 13 {
		return true
	}
	
	// Anti-eco rounds (2nd, 3rd, 17th, 18th)
	if roundNum == 2 || roundNum == 3 || roundNum == 14 || roundNum == 15 {
		return true
	}
	
	// Match point rounds
	if roundNum >= 12 || roundNum >= 24 { // Near end of half or match
		return true
	}
	
	return false
}

// CalculateTeamEconomyRating calculates an overall economy rating for a team
func (em *EconomyManager) CalculateTeamEconomyRating(team *models.Team, teamEconomy *models.TeamEconomy) float64 {
	// Factors: average money, equipment value, recent spending efficiency
	avgMoney := float64(teamEconomy.AverageMoney)
	equipValue := float64(teamEconomy.EquipmentValue)
	
	// Normalize values (0.0 to 1.0)
	moneyRating := avgMoney / 16000.0 // Max money
	if moneyRating > 1.0 {
		moneyRating = 1.0
	}
	
	equipRating := equipValue / 25000.0 // Rough estimate of max equipment value
	if equipRating > 1.0 {
		equipRating = 1.0
	}
	
	// Combine ratings
	overallRating := (moneyRating*0.6 + equipRating*0.4)
	
	return overallRating
}

// GetBuyTypeDistribution returns statistics on buy types used
func (em *EconomyManager) GetBuyTypeDistribution(match *models.Match) map[string]int {
	distribution := map[string]int{
		"full_buy":  0,
		"force_buy": 0,
		"eco":       0,
		"semi_eco":  0,
		"anti_eco":  0,
	}
	
	// This would be populated during match generation
	// For now, return empty distribution
	return distribution
}