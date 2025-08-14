package generator

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/noueii/nocs-log-generator/backend/pkg/models"
)

// WebSocketManager interface for broadcasting events (to avoid import cycle)
type WebSocketManager interface {
	BroadcastMatchEvent(matchID string, eventType string, data interface{}) error
	BroadcastMatchStatus(matchID string, status string, data interface{}) error
	BroadcastMatchError(matchID string, errorMsg string) error
}

// MatchEngine handles the core match generation logic
type MatchEngine struct {
	config           *models.MatchConfig
	match            *models.Match
	state            *models.MatchState
	eventFactory     *models.EventFactory
	roundSimulator   *RoundSimulator
	eventGenerator   *EventGenerator
	economyManager   *EconomyManager
	logFormatter     *LogFormatter
	rng              *rand.Rand
	wsManager        WebSocketManager
	
	// Match settings
	roundTime        time.Duration
	freezeTime       time.Duration
	bombTimer        time.Duration
	
	// Economics
	startMoney       int
	maxMoney         int
	killReward       int
	winBonus         int
	lossBonus        []int // Escalating loss bonus
	
	// Simulation state
	currentTick      int64
	tickRate         int
	totalEvents      int64
}

// NewMatchEngine creates a new match engine with the given configuration
func NewMatchEngine(config *models.MatchConfig, match *models.Match) *MatchEngine {
	seed := config.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	
	engine := &MatchEngine{
		config:       config,
		match:        match,
		eventFactory: models.NewEventFactory(),
		rng:          rand.New(rand.NewSource(seed)),
		
		// Standard CS2 settings
		roundTime:    time.Second * 115,
		freezeTime:   time.Second * 15,
		bombTimer:    time.Second * 40,
		
		// Economics
		startMoney:   config.StartMoney,
		maxMoney:     config.MaxMoney,
		killReward:   300,
		winBonus:     3250,
		lossBonus:    []int{1400, 1900, 2400, 2900, 3400}, // CS2 loss bonus progression
		
		// Technical settings
		tickRate:     config.TickRate,
		currentTick:  0,
		totalEvents:  0,
	}
	
	// Initialize subsystems
	engine.roundSimulator = NewRoundSimulator(engine.rng, models.NewEconomyManager(), config)
	engine.eventGenerator = NewEventGenerator(engine.rng, config)
	engine.economyManager = NewEconomyManager(engine.rng)
	engine.logFormatter = NewLogFormatter(config)
	
	// Initialize match state
	engine.initializeMatchState()
	
	return engine
}

// SetWebSocketManager sets the WebSocket manager for streaming events
func (e *MatchEngine) SetWebSocketManager(wsManager WebSocketManager) {
	e.wsManager = wsManager
}

// initializeMatchState sets up the initial match state
func (e *MatchEngine) initializeMatchState() {
	e.state = &models.MatchState{
		CurrentRound:  0,
		Scores:        make(map[string]int),
		TeamEconomies: make(map[string]*models.TeamEconomy),
		PlayerStates:  make(map[string]*models.PlayerState),
		IsLive:        false,
		IsFreezeTime:  true,
		CurrentTick:   0,
	}
	
	// Initialize team scores and economies
	for _, team := range e.match.Teams {
		e.state.Scores[team.Name] = 0
		
		teamEconomy := &models.TeamEconomy{
			TotalMoney:        e.startMoney * 5,
			AverageMoney:      e.startMoney,
			ConsecutiveLosses: 0,
			LossBonus:         e.lossBonus[0],
		}
		e.state.TeamEconomies[team.Name] = teamEconomy
		
		// Initialize player states
		for i, player := range team.Players {
			playerState := &models.PlayerState{
				IsAlive:      true,
				Health:       100,
				Armor:        0,
				HasHelmet:    false,
				HasDefuseKit: false,
				Money:        e.startMoney,
				Position:     e.getSpawnPosition(team.Side, i),
				Grenades:     make([]models.Grenade, 0),
			}
			e.state.PlayerStates[player.Name] = playerState
		}
	}
}

// GenerateMatch executes the complete match generation process
func (e *MatchEngine) GenerateMatch() error {
	e.match.Status = "generating"
	e.match.StartTime = time.Now()
	
	// Generate match events
	for e.state.CurrentRound < e.match.MaxRounds && !e.isMatchFinished() {
		if err := e.playRound(); err != nil {
			return fmt.Errorf("error playing round %d: %w", e.state.CurrentRound+1, err)
		}
	}
	
	// Finalize match
	e.finalizeMatch()
	
	return nil
}

// GenerateMatchWithStreaming executes the complete match generation process with WebSocket streaming
func (e *MatchEngine) GenerateMatchWithStreaming() error {
	e.match.Status = "generating"
	e.match.StartTime = time.Now()
	
	// Broadcast match start event
	if e.wsManager != nil {
		e.wsManager.BroadcastMatchEvent(e.match.ID, "match_start", map[string]interface{}{
			"match_id": e.match.ID,
			"teams": []string{e.match.Teams[0].Name, e.match.Teams[1].Name},
			"map": e.match.Map,
			"max_rounds": e.match.MaxRounds,
			"started_at": e.match.StartTime,
		})
	}
	
	// Generate match events
	for e.state.CurrentRound < e.match.MaxRounds && !e.isMatchFinished() {
		// Broadcast round start
		if e.wsManager != nil {
			progress := float64(e.state.CurrentRound) / float64(e.match.MaxRounds) * 100
			e.wsManager.BroadcastMatchEvent(e.match.ID, "match_progress", map[string]interface{}{
				"match_id": e.match.ID,
				"current_round": e.state.CurrentRound + 1,
				"total_rounds": e.match.MaxRounds,
				"events_generated": e.totalEvents,
				"progress": progress,
			})
		}
		
		if err := e.playRoundWithStreaming(); err != nil {
			if e.wsManager != nil {
				e.wsManager.BroadcastMatchError(e.match.ID, fmt.Sprintf("Error playing round %d: %s", e.state.CurrentRound+1, err.Error()))
			}
			return fmt.Errorf("error playing round %d: %w", e.state.CurrentRound+1, err)
		}
	}
	
	// Finalize match
	e.finalizeMatch()
	
	// Broadcast match completion
	if e.wsManager != nil {
		e.wsManager.BroadcastMatchEvent(e.match.ID, "match_complete", map[string]interface{}{
			"match_id": e.match.ID,
			"total_rounds": len(e.match.Rounds),
			"total_events": e.match.TotalEvents,
			"duration": e.match.Duration.Seconds(),
			"completed_at": e.match.EndTime,
			"success": true,
		})
	}
	
	return nil
}

// playRound executes a single round of the match
func (e *MatchEngine) playRound() error {
	e.state.CurrentRound++
	e.eventFactory.SetRound(e.state.CurrentRound)
	
	// Check for side switch at halftime
	if e.state.CurrentRound == (e.match.MaxRounds/2)+1 {
		e.switchSides()
	}
	
	// Handle pre-round economy
	if err := e.handleBuyPhase(); err != nil {
		return fmt.Errorf("buy phase error: %w", err)
	}
	
	// Start round
	roundStartTime := time.Now()
	e.state.RoundStartTime = roundStartTime
	e.state.IsFreezeTime = false
	e.state.IsLive = true
	
	// Create round start event
	ctTeam := e.getTeamBySide("CT")
	tTeam := e.getTeamBySide("TERRORIST")
	
	startEvent := e.eventFactory.CreateRoundStartEvent(
		e.state.Scores[ctTeam.Name],
		e.state.Scores[tTeam.Name],
		len(ctTeam.Players),
		len(tTeam.Players),
	)
	e.addEvent(startEvent)
	
	// Simulate round events using the round simulator
	roundResult, roundEvents, err := e.roundSimulator.SimulateRound(e.match, e.state, e.state.CurrentRound)
	if err != nil {
		return fmt.Errorf("round simulation error: %w", err)
	}
	
	// Add all round events to the match
	for _, event := range roundEvents {
		e.addEvent(event)
	}
	
	// Handle round end
	if err := e.handleRoundEnd(roundResult, roundEvents); err != nil {
		return fmt.Errorf("round end handling error: %w", err)
	}
	
	// Update match state
	e.updateMatchStatistics()
	
	return nil
}

// playRoundWithStreaming executes a single round of the match with WebSocket streaming
func (e *MatchEngine) playRoundWithStreaming() error {
	e.state.CurrentRound++
	e.eventFactory.SetRound(e.state.CurrentRound)
	
	// Broadcast round start event
	if e.wsManager != nil {
		e.wsManager.BroadcastMatchEvent(e.match.ID, "round_start", map[string]interface{}{
			"match_id": e.match.ID,
			"round_number": e.state.CurrentRound,
			"ct_score": e.state.Scores[e.getTeamBySide("CT").Name],
			"t_score": e.state.Scores[e.getTeamBySide("TERRORIST").Name],
		})
	}
	
	// Check for side switch at halftime
	if e.state.CurrentRound == (e.match.MaxRounds/2)+1 {
		e.switchSides()
		
		// Broadcast side switch event
		if e.wsManager != nil {
			e.wsManager.BroadcastMatchEvent(e.match.ID, "side_switch", map[string]interface{}{
				"match_id": e.match.ID,
				"round_number": e.state.CurrentRound,
				"message": "Teams switched sides",
			})
		}
	}
	
	// Handle pre-round economy
	if err := e.handleBuyPhase(); err != nil {
		return fmt.Errorf("buy phase error: %w", err)
	}
	
	// Broadcast economy update
	if e.wsManager != nil {
		economyData := make(map[string]map[string]int)
		for _, team := range e.match.Teams {
			teamEconomy := make(map[string]int)
			for _, player := range team.Players {
				teamEconomy[player.Name] = e.state.PlayerStates[player.Name].Money
			}
			economyData[team.Name] = teamEconomy
		}
		
		e.wsManager.BroadcastMatchEvent(e.match.ID, "economy_update", map[string]interface{}{
			"match_id": e.match.ID,
			"round": e.state.CurrentRound,
			"economy": economyData,
		})
	}
	
	// Start round
	roundStartTime := time.Now()
	e.state.RoundStartTime = roundStartTime
	e.state.IsFreezeTime = false
	e.state.IsLive = true
	
	// Create round start event
	ctTeam := e.getTeamBySide("CT")
	tTeam := e.getTeamBySide("TERRORIST")
	
	startEvent := e.eventFactory.CreateRoundStartEvent(
		e.state.Scores[ctTeam.Name],
		e.state.Scores[tTeam.Name],
		len(ctTeam.Players),
		len(tTeam.Players),
	)
	e.addEvent(startEvent)
	
	// Simulate round events using the round simulator
	roundResult, roundEvents, err := e.roundSimulator.SimulateRound(e.match, e.state, e.state.CurrentRound)
	if err != nil {
		return fmt.Errorf("round simulation error: %w", err)
	}
	
	// Add all round events to the match and broadcast them
	for _, event := range roundEvents {
		e.addEvent(event)
		
		// Broadcast significant events
		if e.wsManager != nil {
			e.broadcastGameEvent(event)
		}
	}
	
	// Handle round end
	if err := e.handleRoundEnd(roundResult, roundEvents); err != nil {
		return fmt.Errorf("round end handling error: %w", err)
	}
	
	// Broadcast round end event
	if e.wsManager != nil {
		e.wsManager.BroadcastMatchEvent(e.match.ID, "round_end", map[string]interface{}{
			"match_id": e.match.ID,
			"round_number": e.state.CurrentRound,
			"winner": roundResult.Winner,
			"reason": roundResult.Reason,
			"mvp": roundResult.MVP.Name,
			"ct_score": e.state.Scores[ctTeam.Name],
			"t_score": e.state.Scores[tTeam.Name],
			"duration": roundResult.Duration.Seconds(),
		})
	}
	
	// Update match state
	e.updateMatchStatistics()
	
	return nil
}

// broadcastGameEvent broadcasts specific game events via WebSocket
func (e *MatchEngine) broadcastGameEvent(event models.GameEvent) {
	if e.wsManager == nil {
		return
	}
	
	switch evt := event.(type) {
	case *models.KillEvent:
		e.wsManager.BroadcastMatchEvent(e.match.ID, "player_kill", map[string]interface{}{
			"match_id": e.match.ID,
			"round": e.state.CurrentRound,
			"attacker": evt.Attacker.Name,
			"victim": evt.Victim.Name,
			"weapon": evt.Weapon,
			"headshot": evt.Headshot,
			"distance": evt.Distance,
		})
	
	case *models.BombPlantEvent:
		e.wsManager.BroadcastMatchEvent(e.match.ID, "bomb_plant", map[string]interface{}{
			"match_id": e.match.ID,
			"round": e.state.CurrentRound,
			"player": evt.Player.Name,
			"site": evt.Site,
		})
	
	case *models.BombDefuseEvent:
		e.wsManager.BroadcastMatchEvent(e.match.ID, "bomb_defuse", map[string]interface{}{
			"match_id": e.match.ID,
			"round": e.state.CurrentRound,
			"player": evt.Player.Name,
			"site": evt.Site,
			"with_kit": evt.WithKit,
		})
	
	case *models.BombExplodeEvent:
		e.wsManager.BroadcastMatchEvent(e.match.ID, "bomb_explode", map[string]interface{}{
			"match_id": e.match.ID,
			"round": e.state.CurrentRound,
			"site": evt.Site,
		})
	}
}

// simulateRoundEvents generates events for a single round (legacy method, now unused)
func (e *MatchEngine) simulateRoundEvents() (*RoundResult, error) {
	roundStartTick := e.currentTick
	maxRoundTicks := int64(e.roundTime.Seconds()) * int64(e.tickRate)
	
	// Initialize round state
	e.resetPlayerStates()
	
	// Determine round outcome probability based on team economies and skill
	ctTeam := e.getTeamBySide("CT")
	_ = e.getTeamBySide("TERRORIST") // tTeam unused in legacy method
	
	// Simple round outcome simulation
	for e.currentTick-roundStartTick < maxRoundTicks {
		// Simulate bomb plant scenario
		if e.rng.Float64() < 0.3 && e.currentTick-roundStartTick > int64(20*e.tickRate) {
			if bombPlantResult := e.simulateBombPlant(); bombPlantResult != nil {
				return bombPlantResult, nil
			}
		}
		
		// Simulate elimination rounds
		if eliminationResult := e.simulateElimination(); eliminationResult != nil {
			return eliminationResult, nil
		}
		
		// Advance tick
		e.currentTick += int64(e.tickRate) // Advance by 1 second
	}
	
	// Time expired - CT wins
	return &RoundResult{
		Winner:    "CT",
		Reason:    "time",
		MVP:       e.selectMVP(ctTeam),
		Duration:  e.roundTime,
		EndTick:   e.currentTick,
	}, nil
}

// simulateBombPlant handles bomb plant scenarios
func (e *MatchEngine) simulateBombPlant() *RoundResult {
	tTeam := e.getTeamBySide("TERRORIST")
	ctTeam := e.getTeamBySide("CT")
	
	// Select random T player for bomb plant
	aliveTPlayers := e.getAlivePlayers(tTeam)
	if len(aliveTPlayers) == 0 {
		return nil
	}
	
	planter := aliveTPlayers[e.rng.Intn(len(aliveTPlayers))]
	bombSite := []string{"A", "B"}[e.rng.Intn(2)]
	
	// Create bomb plant event
	plantEvent := &models.BombPlantEvent{
		BaseEvent: models.NewBaseEvent("bomb_plant", e.currentTick, e.state.CurrentRound),
		Player:    planter,
		Site:      bombSite,
		Position:  e.getBombSitePosition(bombSite),
	}
	e.addEvent(plantEvent)
	
	// Simulate post-plant scenario
	defuseTime := time.Second * 10 // Default defuse time
	bombTimer := e.bombTimer
	
	// Simple probability: 60% bomb explodes, 40% defused
	if e.rng.Float64() < 0.4 && len(e.getAlivePlayers(ctTeam)) > 0 {
		// Bomb defused
		defuser := e.getAlivePlayers(ctTeam)[0]
		hasKit := e.rng.Float64() < 0.7 // 70% chance of defuse kit
		if hasKit {
			defuseTime = time.Second * 5
		}
		
		defuseEvent := &models.BombDefuseEvent{
			BaseEvent: models.NewBaseEvent("bomb_defuse", e.currentTick+int64(defuseTime.Seconds())*int64(e.tickRate), e.state.CurrentRound),
			Player:    defuser,
			Site:      bombSite,
			WithKit:   hasKit,
			Position:  e.getBombSitePosition(bombSite),
		}
		e.addEvent(defuseEvent)
		
		return &RoundResult{
			Winner:   "CT",
			Reason:   "bomb_defused",
			MVP:      defuser,
			Duration: time.Duration(e.currentTick/int64(e.tickRate)) * time.Second,
			EndTick:  e.currentTick,
		}
	} else {
		// Bomb explodes
		explodeEvent := &models.BombExplodeEvent{
			BaseEvent: models.NewBaseEvent("bomb_explode", e.currentTick+int64(bombTimer.Seconds())*int64(e.tickRate), e.state.CurrentRound),
			Site:      bombSite,
			Position:  e.getBombSitePosition(bombSite),
		}
		e.addEvent(explodeEvent)
		
		return &RoundResult{
			Winner:   "TERRORIST",
			Reason:   "bomb_exploded",
			MVP:      planter,
			Duration: time.Duration(e.currentTick/int64(e.tickRate)) * time.Second,
			EndTick:  e.currentTick,
		}
	}
}

// simulateElimination handles elimination scenarios
func (e *MatchEngine) simulateElimination() *RoundResult {
	ctTeam := e.getTeamBySide("CT")
	tTeam := e.getTeamBySide("TERRORIST")
	
	// Generate some kill events based on team skill and economy
	for i := 0; i < e.rng.Intn(3)+1; i++ {
		if killEvent := e.generateKillEvent(); killEvent != nil {
			e.addEvent(killEvent)
		}
	}
	
	// Check if one team is eliminated
	ctAlive := len(e.getAlivePlayers(ctTeam))
	tAlive := len(e.getAlivePlayers(tTeam))
	
	if ctAlive == 0 {
		return &RoundResult{
			Winner:   "TERRORIST",
			Reason:   "elimination",
			MVP:      e.selectMVP(tTeam),
			Duration: time.Duration(e.currentTick/int64(e.tickRate)) * time.Second,
			EndTick:  e.currentTick,
		}
	} else if tAlive == 0 {
		return &RoundResult{
			Winner:   "CT",
			Reason:   "elimination",
			MVP:      e.selectMVP(ctTeam),
			Duration: time.Duration(e.currentTick/int64(e.tickRate)) * time.Second,
			EndTick:  e.currentTick,
		}
	}
	
	return nil
}

// generateKillEvent creates a realistic kill event
func (e *MatchEngine) generateKillEvent() *models.KillEvent {
	ctTeam := e.getTeamBySide("CT")
	tTeam := e.getTeamBySide("TERRORIST")
	
	ctAlive := e.getAlivePlayers(ctTeam)
	tAlive := e.getAlivePlayers(tTeam)
	
	if len(ctAlive) == 0 || len(tAlive) == 0 {
		return nil
	}
	
	// Randomly select attacker and victim from different teams
	var attacker, victim *models.Player
	if e.rng.Float64() < 0.5 {
		attacker = ctAlive[e.rng.Intn(len(ctAlive))]
		victim = tAlive[e.rng.Intn(len(tAlive))]
	} else {
		attacker = tAlive[e.rng.Intn(len(tAlive))]
		victim = ctAlive[e.rng.Intn(len(ctAlive))]
	}
	
	// Select weapon based on economy and round
	weapon := e.selectWeapon(attacker)
	headshot := e.rng.Float64() < 0.25 // 25% headshot rate
	
	// Create kill event
	killEvent := &models.KillEvent{
		BaseEvent:   models.NewBaseEvent("player_death", e.currentTick, e.state.CurrentRound),
		Attacker:    attacker,
		Victim:      victim,
		Weapon:      weapon,
		Headshot:    headshot,
		Penetrated:  0,
		NoScope:     false,
		AttackerBlind: false,
		Distance:    float64(e.rng.Intn(30) + 5), // 5-35 meters
		AttackerPos: e.state.PlayerStates[attacker.Name].Position,
		VictimPos:   e.state.PlayerStates[victim.Name].Position,
	}
	
	// Update player states
	e.state.PlayerStates[victim.Name].IsAlive = false
	e.state.PlayerStates[victim.Name].Health = 0
	
	// Update statistics
	attacker.Stats.Kills++
	victim.Stats.Deaths++
	if headshot {
		attacker.Stats.Headshots++
	}
	
	return killEvent
}

// handleBuyPhase manages the economy and equipment purchases
func (e *MatchEngine) handleBuyPhase() error {
	for _, team := range e.match.Teams {
		teamEconomy := e.state.TeamEconomies[team.Name]
		
		// Simple buy logic based on team economy
		avgMoney := teamEconomy.AverageMoney
		
		for i, player := range team.Players {
			playerState := e.state.PlayerStates[player.Name]
			
			// Buy armor if affordable
			if playerState.Money >= 650 && playerState.Armor == 0 {
				playerState.Armor = 100
				playerState.HasHelmet = true
				playerState.Money -= 1000 // Helmet + armor
				
				purchaseEvent := &models.ItemPurchaseEvent{
					BaseEvent: models.NewBaseEvent("item_purchase", e.currentTick, e.state.CurrentRound),
					Player:    &team.Players[i],
					Item:      "item_assaultsuit",
					Cost:      1000,
				}
				e.addEvent(purchaseEvent)
			}
			
			// Buy primary weapon based on economy
			if playerState.PrimaryWeapon == nil {
				weapon := e.selectBuyWeapon(avgMoney, player.Role)
				if weapon != nil && playerState.Money >= weapon.Price {
					playerState.PrimaryWeapon = weapon
					playerState.Money -= weapon.Price
					
					purchaseEvent := &models.ItemPurchaseEvent{
						BaseEvent: models.NewBaseEvent("item_purchase", e.currentTick, e.state.CurrentRound),
						Player:    &team.Players[i],
						Item:      weapon.Name,
						Cost:      weapon.Price,
					}
					e.addEvent(purchaseEvent)
				}
			}
			
			// Buy grenades
			if playerState.Money >= 300 && len(playerState.Grenades) < 2 {
				grenadeType := e.selectGrenade(team.Side)
				grenade := models.Grenade{Type: grenadeType, Price: 300}
				playerState.Grenades = append(playerState.Grenades, grenade)
				playerState.Money -= 300
				
				purchaseEvent := &models.ItemPurchaseEvent{
					BaseEvent: models.NewBaseEvent("item_purchase", e.currentTick, e.state.CurrentRound),
					Player:    &team.Players[i],
					Item:      grenadeType,
					Cost:      300,
				}
				e.addEvent(purchaseEvent)
			}
			
			// Buy defuse kit for CTs
			if team.Side == "CT" && !playerState.HasDefuseKit && playerState.Money >= 400 {
				playerState.HasDefuseKit = true
				playerState.Money -= 400
				
				purchaseEvent := &models.ItemPurchaseEvent{
					BaseEvent: models.NewBaseEvent("item_purchase", e.currentTick, e.state.CurrentRound),
					Player:    &team.Players[i],
					Item:      "item_defuser",
					Cost:      400,
				}
				e.addEvent(purchaseEvent)
			}
		}
		
		// Update team economy
		e.updateTeamEconomy(&team)
	}
	
	return nil
}

// handleRoundEnd processes the end of a round
func (e *MatchEngine) handleRoundEnd(result *RoundResult, roundEvents []models.GameEvent) error {
	// Update scores
	e.state.Scores[result.Winner]++
	e.match.Scores[result.Winner]++
	
	// Handle economy rewards using the economy manager
	if err := e.economyManager.HandleRoundEnd(e.match, e.state, result, roundEvents); err != nil {
		return fmt.Errorf("failed to handle round end economy: %w", err)
	}
	
	// Create round end event
	ctScore := e.state.Scores[e.getTeamBySide("CT").Name]
	tScore := e.state.Scores[e.getTeamBySide("TERRORIST").Name]
	
	endEvent := &models.RoundEndEvent{
		BaseEvent: models.NewBaseEvent("round_end", e.currentTick, e.state.CurrentRound),
		Winner:    result.Winner,
		Reason:    result.Reason,
		CTScore:   ctScore,
		TScore:    tScore,
		MVP:       result.MVP,
	}
	e.addEvent(endEvent)
	
	// Create round data
	roundData := models.RoundData{
		RoundNumber: e.state.CurrentRound,
		StartTime:   e.state.RoundStartTime,
		EndTime:     time.Now(),
		Winner:      result.Winner,
		Reason:      result.Reason,
		MVP:         result.MVP.Name,
		Scores:      make(map[string]int),
		Economy:     make(map[string]models.TeamEconomy),
	}
	
	// Copy scores and economies
	for teamName, score := range e.state.Scores {
		roundData.Scores[teamName] = score
	}
	for teamName, economy := range e.state.TeamEconomies {
		roundData.Economy[teamName] = *economy
	}
	
	e.match.Rounds = append(e.match.Rounds, roundData)
	return nil
}

// handleEconomyRewards manages money rewards after round end
func (e *MatchEngine) handleEconomyRewards(result *RoundResult) {
	winningTeamName := result.Winner
	losingTeamName := ""
	
	// Identify losing team
	for _, team := range e.match.Teams {
		if team.Name != winningTeamName {
			losingTeamName = team.Name
			break
		}
	}
	
	// Award win bonus
	winningTeam := e.getTeamByName(winningTeamName)
	winBonus := e.winBonus
	if result.Reason == "bomb_exploded" {
		winBonus = 3500 // Bomb plant bonus
	}
	
	for i := range winningTeam.Players {
		playerState := e.state.PlayerStates[winningTeam.Players[i].Name]
		playerState.Money = e.capMoney(playerState.Money + winBonus)
	}
	
	// Reset winning team loss streak
	e.state.TeamEconomies[winningTeamName].ConsecutiveLosses = 0
	
	// Award loss bonus
	losingTeam := e.getTeamByName(losingTeamName)
	teamEconomy := e.state.TeamEconomies[losingTeamName]
	teamEconomy.ConsecutiveLosses++
	
	lossIndex := teamEconomy.ConsecutiveLosses - 1
	if lossIndex >= len(e.lossBonus) {
		lossIndex = len(e.lossBonus) - 1
	}
	lossBonus := e.lossBonus[lossIndex]
	teamEconomy.LossBonus = lossBonus
	
	for i := range losingTeam.Players {
		playerState := e.state.PlayerStates[losingTeam.Players[i].Name]
		playerState.Money = e.capMoney(playerState.Money + lossBonus)
	}
}

// Helper functions

// getTeamBySide returns the team playing on the specified side
func (e *MatchEngine) getTeamBySide(side string) *models.Team {
	for i := range e.match.Teams {
		if e.match.Teams[i].Side == side {
			return &e.match.Teams[i]
		}
	}
	return nil
}

// getTeamByName returns the team with the specified name
func (e *MatchEngine) getTeamByName(name string) *models.Team {
	for i := range e.match.Teams {
		if e.match.Teams[i].Name == name {
			return &e.match.Teams[i]
		}
	}
	return nil
}

// getAlivePlayers returns all living players from a team
func (e *MatchEngine) getAlivePlayers(team *models.Team) []*models.Player {
	var alive []*models.Player
	for i := range team.Players {
		if e.state.PlayerStates[team.Players[i].Name].IsAlive {
			alive = append(alive, &team.Players[i])
		}
	}
	return alive
}

// resetPlayerStates resets player states for a new round
func (e *MatchEngine) resetPlayerStates() {
	for _, team := range e.match.Teams {
		for i, player := range team.Players {
			playerState := e.state.PlayerStates[player.Name]
			playerState.IsAlive = true
			playerState.Health = 100
			playerState.Position = e.getSpawnPosition(team.Side, i)
			playerState.IsFlashed = false
			playerState.IsSmoked = false
			playerState.IsDefusing = false
			playerState.IsPlanting = false
			playerState.IsReloading = false
			playerState.HasBomb = false
		}
	}
}

// selectMVP selects the MVP for a team based on performance
func (e *MatchEngine) selectMVP(team *models.Team) *models.Player {
	// Simple MVP selection - player with most kills this round
	// In a real implementation, this would consider damage, assists, objective play, etc.
	mvp := &team.Players[0]
	for i := range team.Players {
		if team.Players[i].Stats.Kills > mvp.Stats.Kills {
			mvp = &team.Players[i]
		}
	}
	return mvp
}

// selectWeapon selects a weapon for a kill event
func (e *MatchEngine) selectWeapon(player *models.Player) string {
	weapons := []string{
		"ak47", "m4a4", "m4a1_silencer", "awp", "ssg08",
		"aug", "sg556", "famas", "galil", "mp9", "mac10",
		"ump45", "p90", "bizon", "deagle", "glock", "usp_silencer",
	}
	return weapons[e.rng.Intn(len(weapons))]
}

// selectBuyWeapon selects a weapon to buy based on economy
func (e *MatchEngine) selectBuyWeapon(money int, role string) *models.Weapon {
	if money >= 4700 && role == "awp" {
		return &models.Weapon{Name: "awp", Type: "sniper", Price: 4750}
	} else if money >= 2700 {
		return &models.Weapon{Name: "ak47", Type: "rifle", Price: 2700}
	} else if money >= 1300 {
		return &models.Weapon{Name: "ump45", Type: "smg", Price: 1200}
	}
	return nil
}

// selectGrenade selects a grenade type to buy
func (e *MatchEngine) selectGrenade(side string) string {
	grenades := []string{"hegrenade", "flashbang", "smokegrenade"}
	if side == "TERRORIST" {
		grenades = append(grenades, "molotov")
	} else {
		grenades = append(grenades, "incgrenade")
	}
	return grenades[e.rng.Intn(len(grenades))]
}

// getSpawnPosition returns a spawn position for a player
func (e *MatchEngine) getSpawnPosition(side string, playerIndex int) models.Vector3 {
	// Simple spawn positions - in a real implementation these would be map-specific
	baseX := float64(playerIndex * 100)
	if side == "CT" {
		return models.Vector3{X: baseX, Y: 0, Z: 0}
	}
	return models.Vector3{X: baseX, Y: 1000, Z: 0}
}

// getBombSitePosition returns the position of a bomb site
func (e *MatchEngine) getBombSitePosition(site string) models.Vector3 {
	if site == "A" {
		return models.Vector3{X: 500, Y: 500, Z: 0}
	}
	return models.Vector3{X: 1500, Y: 500, Z: 0}
}

// updateTeamEconomy updates a team's economic statistics
func (e *MatchEngine) updateTeamEconomy(team *models.Team) {
	economy := e.state.TeamEconomies[team.Name]
	totalMoney := 0
	equipmentValue := 0
	
	for _, player := range team.Players {
		playerState := e.state.PlayerStates[player.Name]
		totalMoney += playerState.Money
		equipmentValue += e.calculateEquipmentValue(playerState)
	}
	
	economy.TotalMoney = totalMoney
	economy.AverageMoney = totalMoney / len(team.Players)
	economy.EquipmentValue = equipmentValue
}

// calculateEquipmentValue calculates the value of a player's equipment
func (e *MatchEngine) calculateEquipmentValue(state *models.PlayerState) int {
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

// capMoney ensures money doesn't exceed the maximum
func (e *MatchEngine) capMoney(money int) int {
	if money > e.maxMoney {
		return e.maxMoney
	}
	return money
}

// switchSides switches team sides at halftime
func (e *MatchEngine) switchSides() {
	for i := range e.match.Teams {
		if e.match.Teams[i].Side == "CT" {
			e.match.Teams[i].Side = "TERRORIST"
		} else {
			e.match.Teams[i].Side = "CT"
		}
		
		// Update all players in the team
		for j := range e.match.Teams[i].Players {
			e.match.Teams[i].Players[j].Side = e.match.Teams[i].Side
		}
	}
}

// isMatchFinished checks if the match is complete
func (e *MatchEngine) isMatchFinished() bool {
	winThreshold := (e.match.MaxRounds / 2) + 1
	for _, score := range e.state.Scores {
		if score >= winThreshold {
			return true
		}
	}
	return false
}

// updateMatchStatistics updates overall match statistics
func (e *MatchEngine) updateMatchStatistics() {
	// Update player statistics
	for _, team := range e.match.Teams {
		for i := range team.Players {
			player := &team.Players[i]
			player.CalculateRating(e.state.CurrentRound)
		}
	}
}

// finalizeMatch completes the match generation
func (e *MatchEngine) finalizeMatch() {
	e.match.Status = "completed"
	e.match.EndTime = time.Now()
	e.match.Duration = e.match.EndTime.Sub(e.match.StartTime)
	e.match.CurrentRound = e.state.CurrentRound
	e.match.TotalEvents = e.totalEvents
	
	// Set final scores
	for teamName, score := range e.state.Scores {
		e.match.Scores[teamName] = score
	}
}

// addEvent adds an event to the match and increments counters
func (e *MatchEngine) addEvent(event models.GameEvent) {
	e.match.Events = append(e.match.Events, event)
	e.totalEvents++
	e.eventFactory.SetTick(e.currentTick)
}

// RoundResult represents the outcome of a round
type RoundResult struct {
	Winner   string
	Reason   string
	MVP      *models.Player
	Duration time.Duration
	EndTick  int64
}