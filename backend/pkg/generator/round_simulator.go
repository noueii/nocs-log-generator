package generator

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/noueii/nocs-log-generator/backend/pkg/models"
)

// RoundSimulator handles individual round simulation
type RoundSimulator struct {
	rng            *rand.Rand
	economyManager *models.EconomyManager
	config         *models.MatchConfig
}

// NewRoundSimulator creates a new round simulator
func NewRoundSimulator(rng *rand.Rand, economyManager *models.EconomyManager, config *models.MatchConfig) *RoundSimulator {
	return &RoundSimulator{
		rng:            rng,
		economyManager: economyManager,
		config:         config,
	}
}

// SimulateRound executes the full round simulation including buy phase and combat
func (rs *RoundSimulator) SimulateRound(match *models.Match, state *models.MatchState, roundNum int) (*RoundResult, []models.GameEvent, error) {
	events := make([]models.GameEvent, 0, 100) // Pre-allocate for ~100 events per round
	
	// Execute buy phase
	buyEvents, err := rs.simulateBuyPhase(match, state, roundNum)
	if err != nil {
		return nil, nil, fmt.Errorf("buy phase simulation failed: %w", err)
	}
	events = append(events, buyEvents...)

	// Reset player states for the round
	rs.resetPlayerStatesForRound(match, state)

	// Determine round strategy and flow
	roundStrategy := rs.determineRoundStrategy(match, state)
	
	// Simulate round based on strategy
	var result *RoundResult
	var combatEvents []models.GameEvent
	
	switch roundStrategy.Type {
	case "bomb_scenario":
		result, combatEvents, err = rs.simulateBombRound(match, state, roundNum, roundStrategy)
	case "elimination":
		result, combatEvents, err = rs.simulateEliminationRound(match, state, roundNum, roundStrategy)
	case "timeout":
		result, combatEvents, err = rs.simulateTimeoutRound(match, state, roundNum, roundStrategy)
	default:
		result, combatEvents, err = rs.simulateEliminationRound(match, state, roundNum, roundStrategy)
	}
	
	if err != nil {
		return nil, nil, fmt.Errorf("round simulation failed: %w", err)
	}
	
	events = append(events, combatEvents...)
	
	// Select MVP
	result.MVP = rs.selectMVP(match, result.Winner, events)

	return result, events, nil
}

// RoundStrategy defines how the round should play out
type RoundStrategy struct {
	Type           string  // "bomb_scenario", "elimination", "timeout"
	Intensity      float64 // 0.0-1.0, affects number of events
	CTAdvantage    float64 // -1.0 to 1.0, team advantage
	ExpectedEvents int     // Target number of events
}

// determineRoundStrategy analyzes the match state and determines round flow
func (rs *RoundSimulator) determineRoundStrategy(match *models.Match, state *models.MatchState) *RoundStrategy {
	ctTeam := rs.getTeamBySide(match, "CT")
	tTeam := rs.getTeamBySide(match, "TERRORIST")
	
	// Calculate team advantages based on economy and skill
	ctEconomy := state.TeamEconomies[ctTeam.Name]
	tEconomy := state.TeamEconomies[tTeam.Name]
	
	economyAdvantage := float64(ctEconomy.AverageMoney-tEconomy.AverageMoney) / 5000.0
	if economyAdvantage > 1.0 {
		economyAdvantage = 1.0
	} else if economyAdvantage < -1.0 {
		economyAdvantage = -1.0
	}

	// Determine round type probabilities
	bombProb := 0.4
	eliminationProb := 0.5
	timeoutProb := 0.1
	
	// Adjust probabilities based on round number and score
	if state.CurrentRound > 15 { // Second half
		bombProb += 0.1 // More tactical play
		timeoutProb += 0.05
		eliminationProb -= 0.15
	}
	
	// Select round type
	randValue := rs.rng.Float64()
	var roundType string
	if randValue < bombProb {
		roundType = "bomb_scenario"
	} else if randValue < bombProb+eliminationProb {
		roundType = "elimination"
	} else {
		roundType = "timeout"
	}

	// Calculate intensity based on economy differential
	intensity := 0.5 + math.Abs(economyAdvantage)*0.3
	if intensity > 1.0 {
		intensity = 1.0
	}

	expectedEvents := int(50 + intensity*50) // 50-100 events per round
	
	return &RoundStrategy{
		Type:           roundType,
		Intensity:      intensity,
		CTAdvantage:    economyAdvantage,
		ExpectedEvents: expectedEvents,
	}
}

// simulateBuyPhase handles equipment purchasing for all players
func (rs *RoundSimulator) simulateBuyPhase(match *models.Match, state *models.MatchState, roundNum int) ([]models.GameEvent, error) {
	var events []models.GameEvent
	
	for _, team := range match.Teams {
		teamEconomy := state.TeamEconomies[team.Name]
		
		// Determine team buy strategy
		buyType := rs.determineBuyStrategy(teamEconomy, roundNum)
		
		for i, player := range team.Players {
			playerState := state.PlayerStates[player.Name]
			
			// Get optimal buy for this player
			playerBuy := rs.economyManager.GetOptimalBuy(&player, teamEconomy, buyType)
			
			// Process purchases
			for _, item := range playerBuy {
				cost := rs.getItemCost(item)
				if playerState.Money >= cost {
					// Execute purchase
					playerState.Money -= cost
					
					// Apply item to player state
					rs.applyPurchaseToPlayer(playerState, item)
					
					// Create purchase event
					purchaseEvent := &models.ItemPurchaseEvent{
						BaseEvent: models.NewBaseEvent("item_purchase", 0, roundNum),
						Player:    &match.Teams[rs.getTeamIndex(match, team.Name)].Players[i],
						Item:      item,
						Cost:      cost,
					}
					events = append(events, purchaseEvent)
				}
			}
		}
		
		// Update team economy after purchases
		rs.updateTeamEconomyAfterBuy(&team, state)
	}
	
	return events, nil
}

// simulateBombRound simulates a round with bomb plant/defuse scenario
func (rs *RoundSimulator) simulateBombRound(match *models.Match, state *models.MatchState, roundNum int, strategy *RoundStrategy) (*RoundResult, []models.GameEvent, error) {
	var events []models.GameEvent
	currentTick := int64(0)
	
	// Simulate initial engagements (20-40 seconds)
	initialDuration := time.Duration(20+rs.rng.Intn(20)) * time.Second
	initialTicks := int64(initialDuration.Seconds()) * int64(rs.config.TickRate)
	
	// Generate some early kills
	for currentTick < initialTicks && rs.getAliveCount(match, state, "CT") > 0 && rs.getAliveCount(match, state, "TERRORIST") > 0 {
		if rs.rng.Float64() < 0.3 { // 30% chance of engagement per interval
			if killEvent := rs.generateKillEvent(match, state, currentTick, roundNum); killEvent != nil {
				events = append(events, killEvent)
			}
		}
		currentTick += int64(rs.config.TickRate * 2) // Advance 2 seconds
	}
	
	// Check if round should end early
	if rs.getAliveCount(match, state, "CT") == 0 {
		return &RoundResult{
			Winner:   "TERRORIST",
			Reason:   "elimination",
			Duration: time.Duration(currentTick/int64(rs.config.TickRate)) * time.Second,
		}, events, nil
	}
	if rs.getAliveCount(match, state, "TERRORIST") == 0 {
		return &RoundResult{
			Winner:   "CT",
			Reason:   "elimination",
			Duration: time.Duration(currentTick/int64(rs.config.TickRate)) * time.Second,
		}, events, nil
	}
	
	// Bomb plant phase
	if rs.getAliveCount(match, state, "TERRORIST") > 0 {
		plantSuccess := rs.rng.Float64() < 0.7 // 70% bomb plant success
		
		if plantSuccess {
			// Select planter
			aliveTPlayers := rs.getAlivePlayers(match, state, "TERRORIST")
			if len(aliveTPlayers) > 0 {
				planter := aliveTPlayers[rs.rng.Intn(len(aliveTPlayers))]
				bombSite := []string{"A", "B"}[rs.rng.Intn(2)]
				
				plantEvent := &models.BombPlantEvent{
					BaseEvent: models.NewBaseEvent("bomb_plant", currentTick, roundNum),
					Player:    planter,
					Site:      bombSite,
					Position:  rs.getBombSitePosition(bombSite),
				}
				events = append(events, plantEvent)
				currentTick += int64(rs.config.TickRate * 5) // 5 seconds for plant
				
				// Post-plant scenario
				return rs.simulatePostPlant(match, state, roundNum, currentTick, bombSite, events, strategy)
			}
		}
	}
	
	// If no bomb plant, continue until elimination or time
	for currentTick < int64(115*rs.config.TickRate) { // 115 seconds round time
		if killEvent := rs.generateKillEvent(match, state, currentTick, roundNum); killEvent != nil {
			events = append(events, killEvent)
			
			// Check for round end
			if rs.getAliveCount(match, state, "CT") == 0 {
				return &RoundResult{
					Winner:   "TERRORIST",
					Reason:   "elimination",
					Duration: time.Duration(currentTick/int64(rs.config.TickRate)) * time.Second,
				}, events, nil
			}
			if rs.getAliveCount(match, state, "TERRORIST") == 0 {
				return &RoundResult{
					Winner:   "CT",
					Reason:   "elimination",
					Duration: time.Duration(currentTick/int64(rs.config.TickRate)) * time.Second,
				}, events, nil
			}
		}
		currentTick += int64(rs.config.TickRate * 3) // Advance 3 seconds
	}
	
	// Time expired
	return &RoundResult{
		Winner:   "CT",
		Reason:   "time",
		Duration: time.Duration(115) * time.Second,
	}, events, nil
}

// simulatePostPlant handles the post-bomb-plant scenario
func (rs *RoundSimulator) simulatePostPlant(match *models.Match, state *models.MatchState, roundNum int, currentTick int64, bombSite string, events []models.GameEvent, strategy *RoundStrategy) (*RoundResult, []models.GameEvent, error) {
	bombTimer := 40 * time.Second // 40 second bomb timer
	bombTicks := int64(bombTimer.Seconds()) * int64(rs.config.TickRate)
	maxTick := currentTick + bombTicks
	
	// Post-plant engagements
	for currentTick < maxTick-int64(10*rs.config.TickRate) { // Leave 10 seconds for defuse
		if killEvent := rs.generateKillEvent(match, state, currentTick, roundNum); killEvent != nil {
			events = append(events, killEvent)
			
			// Check for elimination
			if rs.getAliveCount(match, state, "CT") == 0 {
				// Bomb explodes
				explodeEvent := &models.BombExplodeEvent{
					BaseEvent: models.NewBaseEvent("bomb_explode", maxTick, roundNum),
					Site:      bombSite,
					Position:  rs.getBombSitePosition(bombSite),
				}
				events = append(events, explodeEvent)
				
				return &RoundResult{
					Winner:   "TERRORIST",
					Reason:   "bomb_exploded",
					Duration: time.Duration(maxTick/int64(rs.config.TickRate)) * time.Second,
				}, events, nil
			}
			if rs.getAliveCount(match, state, "TERRORIST") == 0 {
				break // CTs can try to defuse
			}
		}
		currentTick += int64(rs.config.TickRate * 2) // Advance 2 seconds
	}
	
	// Defuse attempt
	aliveCTPlayers := rs.getAlivePlayers(match, state, "CT")
	if len(aliveCTPlayers) > 0 && currentTick < maxTick {
		defuseSuccess := rs.rng.Float64() < 0.4 // 40% defuse success rate
		
		if defuseSuccess {
			defuser := aliveCTPlayers[0]
			hasKit := rs.rng.Float64() < 0.6 // 60% chance of having kit
			defuseTime := 10
			if hasKit {
				defuseTime = 5
			}
			
			defuseEvent := &models.BombDefuseEvent{
				BaseEvent: models.NewBaseEvent("bomb_defuse", currentTick+int64(defuseTime*rs.config.TickRate), roundNum),
				Player:    defuser,
				Site:      bombSite,
				WithKit:   hasKit,
				Position:  rs.getBombSitePosition(bombSite),
			}
			events = append(events, defuseEvent)
			
			return &RoundResult{
				Winner:   "CT",
				Reason:   "bomb_defused",
				Duration: time.Duration((currentTick+int64(defuseTime*rs.config.TickRate))/int64(rs.config.TickRate)) * time.Second,
			}, events, nil
		}
	}
	
	// Bomb explodes
	explodeEvent := &models.BombExplodeEvent{
		BaseEvent: models.NewBaseEvent("bomb_explode", maxTick, roundNum),
		Site:      bombSite,
		Position:  rs.getBombSitePosition(bombSite),
	}
	events = append(events, explodeEvent)
	
	return &RoundResult{
		Winner:   "TERRORIST",
		Reason:   "bomb_exploded",
		Duration: time.Duration(maxTick/int64(rs.config.TickRate)) * time.Second,
	}, events, nil
}

// simulateEliminationRound simulates a round ending in elimination
func (rs *RoundSimulator) simulateEliminationRound(match *models.Match, state *models.MatchState, roundNum int, strategy *RoundStrategy) (*RoundResult, []models.GameEvent, error) {
	var events []models.GameEvent
	currentTick := int64(0)
	maxTicks := int64(115 * rs.config.TickRate) // 115 seconds
	
	// Generate kills until one team is eliminated
	for currentTick < maxTicks {
		if killEvent := rs.generateKillEvent(match, state, currentTick, roundNum); killEvent != nil {
			events = append(events, killEvent)
			
			// Check for elimination
			ctAlive := rs.getAliveCount(match, state, "CT")
			tAlive := rs.getAliveCount(match, state, "TERRORIST")
			
			if ctAlive == 0 {
				return &RoundResult{
					Winner:   "TERRORIST",
					Reason:   "elimination",
					Duration: time.Duration(currentTick/int64(rs.config.TickRate)) * time.Second,
				}, events, nil
			}
			if tAlive == 0 {
				return &RoundResult{
					Winner:   "CT",
					Reason:   "elimination",
					Duration: time.Duration(currentTick/int64(rs.config.TickRate)) * time.Second,
				}, events, nil
			}
		}
		
		// Advance time based on intensity
		advanceTime := int64(1 + rs.rng.Intn(4)) // 1-4 seconds
		if strategy.Intensity > 0.7 {
			advanceTime = 1 // Faster paced round
		}
		currentTick += int64(rs.config.TickRate) * advanceTime
	}
	
	// Time expired - CT wins
	return &RoundResult{
		Winner:   "CT",
		Reason:   "time",
		Duration: time.Duration(115) * time.Second,
	}, events, nil
}

// simulateTimeoutRound simulates a round ending in timeout
func (rs *RoundSimulator) simulateTimeoutRound(match *models.Match, state *models.MatchState, roundNum int, strategy *RoundStrategy) (*RoundResult, []models.GameEvent, error) {
	var events []models.GameEvent
	currentTick := int64(0)
	maxTicks := int64(115 * rs.config.TickRate) // 115 seconds
	
	// Generate fewer kills, round times out
	killCount := 1 + rs.rng.Intn(3) // 1-3 kills max
	killInterval := maxTicks / int64(killCount+1)
	
	for i := 0; i < killCount && currentTick < maxTicks; i++ {
		currentTick += killInterval
		if killEvent := rs.generateKillEvent(match, state, currentTick, roundNum); killEvent != nil {
			events = append(events, killEvent)
			
			// Check if elimination occurred anyway
			if rs.getAliveCount(match, state, "CT") == 0 {
				return &RoundResult{
					Winner:   "TERRORIST",
					Reason:   "elimination",
					Duration: time.Duration(currentTick/int64(rs.config.TickRate)) * time.Second,
				}, events, nil
			}
			if rs.getAliveCount(match, state, "TERRORIST") == 0 {
				return &RoundResult{
					Winner:   "CT",
					Reason:   "elimination",
					Duration: time.Duration(currentTick/int64(rs.config.TickRate)) * time.Second,
				}, events, nil
			}
		}
	}
	
	// Time expired - CT wins
	return &RoundResult{
		Winner:   "CT",
		Reason:   "time",
		Duration: time.Duration(115) * time.Second,
	}, events, nil
}

// Helper methods

func (rs *RoundSimulator) resetPlayerStatesForRound(match *models.Match, state *models.MatchState) {
	for _, team := range match.Teams {
		for i, player := range team.Players {
			playerState := state.PlayerStates[player.Name]
			playerState.IsAlive = true
			playerState.Health = 100
			playerState.Position = rs.getSpawnPosition(team.Side, i)
			playerState.IsFlashed = false
			playerState.IsSmoked = false
			playerState.IsDefusing = false
			playerState.IsPlanting = false
			playerState.IsReloading = false
			playerState.HasBomb = false
			playerState.IsLastAlive = false
		}
	}
}

func (rs *RoundSimulator) determineBuyStrategy(economy *models.TeamEconomy, roundNum int) string {
	avgMoney := economy.AverageMoney
	
	if avgMoney >= 5000 {
		return "full_buy"
	} else if avgMoney >= 2500 {
		return "force_buy"
	}
	return "eco"
}

func (rs *RoundSimulator) getItemCost(item string) int {
	return rs.economyManager.GetWeaponPrice(item) + rs.economyManager.GetUtilityPrice(item)
}

func (rs *RoundSimulator) applyPurchaseToPlayer(state *models.PlayerState, item string) {
	// Apply purchased item to player state
	weaponInfo := rs.economyManager.GetWeaponInfo()
	utilityInfo := rs.economyManager.GetUtilityInfo()
	
	if info, exists := weaponInfo[item]; exists {
		weapon := &models.Weapon{
			Name:  info.Name,
			Type:  info.Type,
			Price: info.Price,
			Ammo:  30, // Default ammo
		}
		
		if info.Type == "pistol" {
			state.SecondaryWeapon = weapon
		} else {
			state.PrimaryWeapon = weapon
		}
	} else if info, exists := utilityInfo[item]; exists {
		switch info.Type {
		case "armor":
			if item == "vesthelm" {
				state.Armor = 100
				state.HasHelmet = true
			} else if item == "vest" {
				state.Armor = 100
			}
		case "utility":
			if item == "defuser" {
				state.HasDefuseKit = true
			}
		case "grenade":
			grenade := models.Grenade{
				Type:  info.Name,
				Price: info.Price,
			}
			state.Grenades = append(state.Grenades, grenade)
		}
	}
}

func (rs *RoundSimulator) generateKillEvent(match *models.Match, state *models.MatchState, tick int64, roundNum int) models.GameEvent {
	ctPlayers := rs.getAlivePlayers(match, state, "CT")
	tPlayers := rs.getAlivePlayers(match, state, "TERRORIST")
	
	if len(ctPlayers) == 0 || len(tPlayers) == 0 {
		return nil
	}
	
	// Select attacker and victim
	var attacker, victim *models.Player
	if rs.rng.Float64() < 0.5 {
		attacker = ctPlayers[rs.rng.Intn(len(ctPlayers))]
		victim = tPlayers[rs.rng.Intn(len(tPlayers))]
	} else {
		attacker = tPlayers[rs.rng.Intn(len(tPlayers))]
		victim = ctPlayers[rs.rng.Intn(len(ctPlayers))]
	}
	
	// Select weapon
	weapon := rs.selectWeaponForKill(attacker, state)
	headshot := rs.rng.Float64() < rs.getHeadshotProbability(attacker, weapon)
	
	// Create kill event
	killEvent := &models.KillEvent{
		BaseEvent:     models.NewBaseEvent("player_death", tick, roundNum),
		Attacker:      attacker,
		Victim:        victim,
		Weapon:        weapon,
		Headshot:      headshot,
		Penetrated:    0,
		NoScope:       false,
		AttackerBlind: false,
		Distance:      float64(5 + rs.rng.Intn(30)), // 5-35 meters
		AttackerPos:   state.PlayerStates[attacker.Name].Position,
		VictimPos:     state.PlayerStates[victim.Name].Position,
	}
	
	// Update player states
	state.PlayerStates[victim.Name].IsAlive = false
	state.PlayerStates[victim.Name].Health = 0
	
	// Update statistics
	attacker.Stats.Kills++
	victim.Stats.Deaths++
	if headshot {
		attacker.Stats.Headshots++
	}
	
	return killEvent
}

func (rs *RoundSimulator) selectMVP(match *models.Match, winner string, events []models.GameEvent) *models.Player {
	// Count kills per player this round
	killCounts := make(map[string]int)
	
	for _, event := range events {
		if killEvent, ok := event.(*models.KillEvent); ok {
			if killEvent.Attacker.Side == winner || (winner == "CT" && killEvent.Attacker.Side == "COUNTER-TERRORIST") {
				killCounts[killEvent.Attacker.Name]++
			}
		}
	}
	
	// Find player with most kills on winning team
	var mvp *models.Player
	maxKills := -1
	
	winningTeam := rs.getTeamBySide(match, winner)
	for _, player := range winningTeam.Players {
		if kills, exists := killCounts[player.Name]; exists && kills > maxKills {
			maxKills = kills
			mvp = &player
		}
	}
	
	// Fallback to first player of winning team
	if mvp == nil && len(winningTeam.Players) > 0 {
		mvp = &winningTeam.Players[0]
	}
	
	return mvp
}

// Utility methods

func (rs *RoundSimulator) getTeamBySide(match *models.Match, side string) *models.Team {
	for i := range match.Teams {
		if match.Teams[i].Side == side {
			return &match.Teams[i]
		}
	}
	return nil
}

func (rs *RoundSimulator) getTeamIndex(match *models.Match, teamName string) int {
	for i, team := range match.Teams {
		if team.Name == teamName {
			return i
		}
	}
	return 0
}

func (rs *RoundSimulator) getAliveCount(match *models.Match, state *models.MatchState, side string) int {
	count := 0
	team := rs.getTeamBySide(match, side)
	if team != nil {
		for _, player := range team.Players {
			if playerState := state.PlayerStates[player.Name]; playerState != nil && playerState.IsAlive {
				count++
			}
		}
	}
	return count
}

func (rs *RoundSimulator) getAlivePlayers(match *models.Match, state *models.MatchState, side string) []*models.Player {
	var alive []*models.Player
	team := rs.getTeamBySide(match, side)
	if team != nil {
		for i, player := range team.Players {
			if playerState := state.PlayerStates[player.Name]; playerState != nil && playerState.IsAlive {
				alive = append(alive, &team.Players[i])
			}
		}
	}
	return alive
}

func (rs *RoundSimulator) updateTeamEconomyAfterBuy(team *models.Team, state *models.MatchState) {
	economy := state.TeamEconomies[team.Name]
	totalMoney := 0
	equipmentValue := 0
	
	for _, player := range team.Players {
		playerState := state.PlayerStates[player.Name]
		totalMoney += playerState.Money
		equipmentValue += rs.calculateEquipmentValue(playerState)
	}
	
	economy.TotalMoney = totalMoney
	economy.AverageMoney = totalMoney / len(team.Players)
	economy.EquipmentValue = equipmentValue
}

func (rs *RoundSimulator) calculateEquipmentValue(state *models.PlayerState) int {
	value := 0
	
	if state.PrimaryWeapon != nil {
		value += state.PrimaryWeapon.Price
	}
	if state.SecondaryWeapon != nil {
		value += state.SecondaryWeapon.Price
	}
	for _, grenade := range state.Grenades {
		value += grenade.Price
	}
	if state.Armor > 0 {
		if state.HasHelmet {
			value += 1000
		} else {
			value += 650
		}
	}
	if state.HasDefuseKit {
		value += 400
	}
	
	return value
}

func (rs *RoundSimulator) selectWeaponForKill(attacker *models.Player, state *models.MatchState) string {
	playerState := state.PlayerStates[attacker.Name]
	
	// Prefer primary weapon if available
	if playerState.PrimaryWeapon != nil {
		return playerState.PrimaryWeapon.Name
	}
	
	// Fall back to secondary
	if playerState.SecondaryWeapon != nil {
		return playerState.SecondaryWeapon.Name
	}
	
	// Default weapons based on side
	if attacker.Side == "CT" {
		return "usp_silencer"
	}
	return "glock"
}

func (rs *RoundSimulator) getHeadshotProbability(attacker *models.Player, weapon string) float64 {
	baseRate := 0.25 // 25% base headshot rate
	
	// Adjust based on player skill
	if attacker.Profile.AimSkill > 0.8 {
		baseRate += 0.15
	} else if attacker.Profile.AimSkill < 0.3 {
		baseRate -= 0.10
	}
	
	// Adjust based on weapon
	if weapon == "awp" {
		baseRate = 0.95 // AWP headshots are usually one-shot kills
	} else if weapon == "ak47" {
		baseRate += 0.05 // AK47 rewards headshots
	}
	
	if baseRate > 0.9 {
		baseRate = 0.9
	} else if baseRate < 0.1 {
		baseRate = 0.1
	}
	
	return baseRate
}

func (rs *RoundSimulator) getSpawnPosition(side string, playerIndex int) models.Vector3 {
	baseX := float64(playerIndex * 100)
	if side == "CT" {
		return models.Vector3{X: baseX, Y: 0, Z: 0}
	}
	return models.Vector3{X: baseX, Y: 1000, Z: 0}
}

func (rs *RoundSimulator) getBombSitePosition(site string) models.Vector3 {
	if site == "A" {
		return models.Vector3{X: 500, Y: 500, Z: 0}
	}
	return models.Vector3{X: 1500, Y: 500, Z: 0}
}