package generator

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/noueii/nocs-log-generator/backend/pkg/models"
)

// EventGenerator creates realistic CS2 events
type EventGenerator struct {
	rng    *rand.Rand
	config *models.MatchConfig
}

// NewEventGenerator creates a new event generator
func NewEventGenerator(rng *rand.Rand, config *models.MatchConfig) *EventGenerator {
	return &EventGenerator{
		rng:    rng,
		config: config,
	}
}

// GenerateRoundEvents creates all events for a round including detailed combat simulation
func (eg *EventGenerator) GenerateRoundEvents(match *models.Match, state *models.MatchState, roundNum int, strategy *RoundStrategy) ([]models.GameEvent, error) {
	var events []models.GameEvent
	
	// Add round start event
	startEvent := eg.createRoundStartEvent(match, state, roundNum)
	events = append(events, startEvent)
	
	// Generate spawn events
	spawnEvents := eg.generateSpawnEvents(match, state, roundNum)
	events = append(events, spawnEvents...)
	
	// Generate buy phase events (these come from round simulator)
	// They should already be created, so we don't duplicate them here
	
	// Generate combat events based on strategy
	combatEvents := eg.generateCombatEvents(match, state, roundNum, strategy)
	events = append(events, combatEvents...)
	
	// Generate utility usage events
	utilityEvents := eg.generateUtilityEvents(match, state, roundNum, strategy)
	events = append(events, utilityEvents...)
	
	// Generate damage events (separate from kills)
	damageEvents := eg.generateDamageEvents(match, state, roundNum, strategy)
	events = append(events, damageEvents...)
	
	// Generate weapon fire events (optional, can be very verbose)
	if eg.config.VerboseLogging {
		fireEvents := eg.generateWeaponFireEvents(match, state, roundNum, strategy)
		events = append(events, fireEvents...)
	}
	
	// Sort events by timestamp/tick
	eg.sortEventsByTime(events)
	
	return events, nil
}

// createRoundStartEvent creates the round start event
func (eg *EventGenerator) createRoundStartEvent(match *models.Match, state *models.MatchState, roundNum int) models.GameEvent {
	ctTeam := eg.getTeamBySide(match, "CT")
	tTeam := eg.getTeamBySide(match, "TERRORIST")
	
	return &models.RoundStartEvent{
		BaseEvent:     models.NewBaseEvent("round_start", 0, roundNum),
		CTScore:       state.Scores[ctTeam.Name],
		TScore:        state.Scores[tTeam.Name],
		CTPlayers:     len(ctTeam.Players),
		TPlayers:      len(tTeam.Players),
		TeamEconomies: eg.copyTeamEconomies(state.TeamEconomies),
	}
}

// generateSpawnEvents creates player spawn events
func (eg *EventGenerator) generateSpawnEvents(match *models.Match, state *models.MatchState, roundNum int) []models.GameEvent {
	var events []models.GameEvent
	tick := int64(eg.config.TickRate) // 1 second after round start
	
	for _, team := range match.Teams {
		for _, player := range team.Players {
			// Players don't have explicit spawn events in standard CS2 logs
			// but we can create them for detailed simulation
			if eg.config.DetailedEvents {
				spawnEvent := &models.ChatEvent{
					BaseEvent: models.NewBaseEvent("player_spawn", tick, roundNum),
					Player:    &player,
					Message:   fmt.Sprintf("%s spawned", player.Name),
					Team:      false,
					Dead:      false,
				}
				events = append(events, spawnEvent)
			}
			tick += int64(eg.config.TickRate / 10) // Slight delay between spawns
		}
	}
	
	return events
}

// generateCombatEvents creates combat-related events
func (eg *EventGenerator) generateCombatEvents(match *models.Match, state *models.MatchState, roundNum int, strategy *RoundStrategy) []models.GameEvent {
	var events []models.GameEvent
	
	// Number of combat engagements based on strategy intensity
	baseEngagements := 3
	maxEngagements := int(float64(baseEngagements) * (1.0 + strategy.Intensity))
	numEngagements := eg.rng.Intn(maxEngagements-baseEngagements) + baseEngagements
	
	roundDuration := int64(115 * eg.config.TickRate) // 115 seconds
	
	for i := 0; i < numEngagements; i++ {
		// Distribute engagements throughout the round
		engagementTime := int64(float64(i+1) / float64(numEngagements+1) * float64(roundDuration))
		
		// Add some randomness to engagement timing
		randomOffset := int64(eg.rng.Intn(20*eg.config.TickRate)) - int64(10*eg.config.TickRate) // ±10 seconds
		engagementTime += randomOffset
		
		if engagementTime < 0 {
			engagementTime = int64(5 * eg.config.TickRate) // Minimum 5 seconds into round
		}
		
		// Create engagement
		engagementEvents := eg.generateEngagement(match, state, roundNum, engagementTime, strategy)
		events = append(events, engagementEvents...)
	}
	
	return events
}

// generateEngagement creates a single combat engagement
func (eg *EventGenerator) generateEngagement(match *models.Match, state *models.MatchState, roundNum int, startTick int64, strategy *RoundStrategy) []models.GameEvent {
	var events []models.GameEvent
	
	ctPlayers := eg.getAlivePlayers(match, state, "CT")
	tPlayers := eg.getAlivePlayers(match, state, "TERRORIST")
	
	if len(ctPlayers) == 0 || len(tPlayers) == 0 {
		return events
	}
	
	// Determine engagement participants (1-3 players per side)
	ctParticipants := eg.selectEngagementParticipants(ctPlayers, strategy)
	tParticipants := eg.selectEngagementParticipants(tPlayers, strategy)
	
	// Simulate the engagement
	tick := startTick
	maxEngagementTicks := int64(10 * eg.config.TickRate) // Max 10 seconds per engagement
	
	for tick < startTick+maxEngagementTicks && len(ctParticipants) > 0 && len(tParticipants) > 0 {
		// Determine who shoots first (based on skill, position, etc.)
		if eg.rng.Float64() < 0.5 {
			// CT attacks first
			if len(ctParticipants) > 0 && len(tParticipants) > 0 {
				attacker := ctParticipants[eg.rng.Intn(len(ctParticipants))]
				victim := tParticipants[eg.rng.Intn(len(tParticipants))]
				
				if damageEvent := eg.createDamageEvent(attacker, victim, tick, roundNum); damageEvent != nil {
					events = append(events, damageEvent)
					
					// Check if damage results in death
					if killEvent := eg.checkForKill(attacker, victim, tick, roundNum, damageEvent.(*models.PlayerHurtEvent)); killEvent != nil {
						events = append(events, killEvent)
						
						// Remove dead player from participants
						tParticipants = eg.removePlayerFromList(tParticipants, victim)
						state.PlayerStates[victim.Name].IsAlive = false
						state.PlayerStates[victim.Name].Health = 0
					}
				}
			}
		} else {
			// T attacks first
			if len(tParticipants) > 0 && len(ctParticipants) > 0 {
				attacker := tParticipants[eg.rng.Intn(len(tParticipants))]
				victim := ctParticipants[eg.rng.Intn(len(ctParticipants))]
				
				if damageEvent := eg.createDamageEvent(attacker, victim, tick, roundNum); damageEvent != nil {
					events = append(events, damageEvent)
					
					// Check if damage results in death
					if killEvent := eg.checkForKill(attacker, victim, tick, roundNum, damageEvent.(*models.PlayerHurtEvent)); killEvent != nil {
						events = append(events, killEvent)
						
						// Remove dead player from participants
						ctParticipants = eg.removePlayerFromList(ctParticipants, victim)
						state.PlayerStates[victim.Name].IsAlive = false
						state.PlayerStates[victim.Name].Health = 0
					}
				}
			}
		}
		
		// Advance time (0.5-2 seconds between shots)
		tick += int64(eg.rng.Intn(int(1.5*float64(eg.config.TickRate)))) + int64(0.5*float64(eg.config.TickRate))
	}
	
	return events
}

// generateUtilityEvents creates grenade and utility usage events
func (eg *EventGenerator) generateUtilityEvents(match *models.Match, state *models.MatchState, roundNum int, strategy *RoundStrategy) []models.GameEvent {
	var events []models.GameEvent
	
	// Generate grenade throws based on team economies
	for _, team := range match.Teams {
		teamEconomy := state.TeamEconomies[team.Name]
		
		// More utilities in full buy rounds
		numUtilities := 0
		if teamEconomy.AverageMoney > 4000 {
			numUtilities = 2 + eg.rng.Intn(3) // 2-4 utilities
		} else if teamEconomy.AverageMoney > 2500 {
			numUtilities = 1 + eg.rng.Intn(2) // 1-2 utilities
		} else {
			numUtilities = eg.rng.Intn(2) // 0-1 utilities
		}
		
		for i := 0; i < numUtilities; i++ {
			// Select random player from team with grenades
			playersWithGrenades := eg.getPlayersWithUtility(&team, state)
			if len(playersWithGrenades) == 0 {
				continue
			}
			
			player := playersWithGrenades[eg.rng.Intn(len(playersWithGrenades))]
			playerState := state.PlayerStates[player.Name]
			
			if len(playerState.Grenades) > 0 {
				grenade := playerState.Grenades[0] // Take first grenade
				
				// Generate throw time (20-90 seconds into round)
				throwTime := int64(20*eg.config.TickRate) + int64(eg.rng.Intn(70*eg.config.TickRate))
				
				// Create grenade throw event
				throwEvent := &models.GrenadeThrowEvent{
					BaseEvent:   models.NewBaseEvent("grenade_throw", throwTime, roundNum),
					Player:      player,
					GrenadeType: grenade.Type,
					Position:    playerState.Position,
					Velocity:    models.Vector3{X: float64(eg.rng.Intn(200) - 100), Y: float64(eg.rng.Intn(200) - 100), Z: 50},
				}
				events = append(events, throwEvent)
				
				// Handle specific grenade effects
				if grenade.Type == "flashbang" {
					flashEvent := eg.createFlashbangEvent(match, state, player, throwTime, roundNum)
					if flashEvent != nil {
						events = append(events, flashEvent)
					}
				}
				
				// Remove grenade from player inventory
				if len(playerState.Grenades) > 1 {
					playerState.Grenades = playerState.Grenades[1:]
				} else {
					playerState.Grenades = []models.Grenade{}
				}
			}
		}
	}
	
	return events
}

// generateDamageEvents creates non-lethal damage events
func (eg *EventGenerator) generateDamageEvents(match *models.Match, state *models.MatchState, roundNum int, strategy *RoundStrategy) []models.GameEvent {
	var events []models.GameEvent
	
	// Generate additional damage events (near misses, body shots that don't kill, etc.)
	numDamageEvents := int(float64(10+eg.rng.Intn(15)) * strategy.Intensity) // 10-25 damage events based on intensity
	roundDuration := int64(115 * eg.config.TickRate)
	
	for i := 0; i < numDamageEvents; i++ {
		// Random time in round
		eventTime := int64(eg.rng.Intn(int(roundDuration)))
		
		ctPlayers := eg.getAlivePlayers(match, state, "CT")
		tPlayers := eg.getAlivePlayers(match, state, "TERRORIST")
		
		if len(ctPlayers) == 0 || len(tPlayers) == 0 {
			continue
		}
		
		// Select random attacker and victim
		var attacker, victim *models.Player
		if eg.rng.Float64() < 0.5 {
			attacker = ctPlayers[eg.rng.Intn(len(ctPlayers))]
			victim = tPlayers[eg.rng.Intn(len(tPlayers))]
		} else {
			attacker = tPlayers[eg.rng.Intn(len(tPlayers))]
			victim = ctPlayers[eg.rng.Intn(len(ctPlayers))]
		}
		
		if damageEvent := eg.createNonLethalDamageEvent(attacker, victim, eventTime, roundNum); damageEvent != nil {
			events = append(events, damageEvent)
		}
	}
	
	return events
}

// generateWeaponFireEvents creates weapon fire events (very verbose)
func (eg *EventGenerator) generateWeaponFireEvents(match *models.Match, state *models.MatchState, roundNum int, strategy *RoundStrategy) []models.GameEvent {
	var events []models.GameEvent
	
	// Only generate these for very detailed logging
	if !eg.config.VerboseLogging {
		return events
	}
	
	numFireEvents := int(float64(20+eg.rng.Intn(40)) * strategy.Intensity) // 20-60 fire events
	roundDuration := int64(115 * eg.config.TickRate)
	
	for i := 0; i < numFireEvents; i++ {
		eventTime := int64(eg.rng.Intn(int(roundDuration)))
		
		// Select random alive player
		allPlayers := append(eg.getAlivePlayers(match, state, "CT"), eg.getAlivePlayers(match, state, "TERRORIST")...)
		if len(allPlayers) == 0 {
			continue
		}
		
		player := allPlayers[eg.rng.Intn(len(allPlayers))]
		playerState := state.PlayerStates[player.Name]
		
		weapon := "glock"
		if playerState.PrimaryWeapon != nil {
			weapon = playerState.PrimaryWeapon.Name
		} else if playerState.SecondaryWeapon != nil {
			weapon = playerState.SecondaryWeapon.Name
		}
		
		fireEvent := &models.WeaponFireEvent{
			BaseEvent: models.NewBaseEvent("weapon_fire", eventTime, roundNum),
			Player:    player,
			Weapon:    weapon,
			Position:  playerState.Position,
			Angle:     models.Vector3{X: float64(eg.rng.Intn(360)), Y: float64(eg.rng.Intn(180)), Z: 0},
			Silenced:  weapon == "m4a1_silencer" || weapon == "usp_silencer",
		}
		events = append(events, fireEvent)
	}
	
	return events
}

// Helper methods

func (eg *EventGenerator) createDamageEvent(attacker, victim *models.Player, tick int64, roundNum int) models.GameEvent {
	playerState := eg.getPlayerState(victim)
	if playerState == nil || !playerState.IsAlive {
		return nil
	}
	
	weapon := eg.selectWeaponForAttack(attacker)
	damage := eg.calculateDamage(attacker, victim, weapon)
	damageArmor := eg.calculateArmorDamage(damage, playerState)
	
	// Apply damage to player state
	newHealth := playerState.Health - damage
	newArmor := playerState.Armor - damageArmor
	if newHealth < 0 {
		newHealth = 0
	}
	if newArmor < 0 {
		newArmor = 0
	}
	
	hitgroup := eg.selectHitgroup(attacker, weapon)
	
	damageEvent := &models.PlayerHurtEvent{
		BaseEvent:   models.NewBaseEvent("player_hurt", tick, roundNum),
		Attacker:    attacker,
		Victim:      victim,
		Weapon:      weapon,
		Damage:      damage,
		DamageArmor: damageArmor,
		Health:      newHealth,
		Armor:       newArmor,
		Hitgroup:    hitgroup,
	}
	
	// Update player state
	playerState.Health = newHealth
	playerState.Armor = newArmor
	
	// Update attacker stats
	attacker.Stats.Damage += damage
	
	return damageEvent
}

func (eg *EventGenerator) createNonLethalDamageEvent(attacker, victim *models.Player, tick int64, roundNum int) models.GameEvent {
	playerState := eg.getPlayerState(victim)
	if playerState == nil || !playerState.IsAlive || playerState.Health <= 20 {
		return nil // Don't create damage that would kill or near-kill
	}
	
	weapon := eg.selectWeaponForAttack(attacker)
	damage := 5 + eg.rng.Intn(15) // 5-19 damage (non-lethal)
	damageArmor := eg.calculateArmorDamage(damage, playerState)
	
	newHealth := playerState.Health - damage
	newArmor := playerState.Armor - damageArmor
	if newHealth < 1 {
		newHealth = 1 // Keep alive
	}
	if newArmor < 0 {
		newArmor = 0
	}
	
	hitgroup := eg.selectHitgroup(attacker, weapon)
	
	damageEvent := &models.PlayerHurtEvent{
		BaseEvent:   models.NewBaseEvent("player_hurt", tick, roundNum),
		Attacker:    attacker,
		Victim:      victim,
		Weapon:      weapon,
		Damage:      damage,
		DamageArmor: damageArmor,
		Health:      newHealth,
		Armor:       newArmor,
		Hitgroup:    hitgroup,
	}
	
	// Update player state
	playerState.Health = newHealth
	playerState.Armor = newArmor
	
	// Update attacker stats
	attacker.Stats.Damage += damage
	
	return damageEvent
}

func (eg *EventGenerator) checkForKill(attacker, victim *models.Player, tick int64, roundNum int, damageEvent *models.PlayerHurtEvent) models.GameEvent {
	if damageEvent.Health <= 0 {
		headshot := damageEvent.Hitgroup == 1
		
		killEvent := &models.KillEvent{
			BaseEvent:     models.NewBaseEvent("player_death", tick, roundNum),
			Attacker:      attacker,
			Victim:        victim,
			Weapon:        damageEvent.Weapon,
			Headshot:      headshot,
			Penetrated:    0,
			NoScope:       false,
			AttackerBlind: false,
			Distance:      float64(5 + eg.rng.Intn(30)),
			AttackerPos:   eg.getPlayerState(attacker).Position,
			VictimPos:     eg.getPlayerState(victim).Position,
		}
		
		// Update stats
		attacker.Stats.Kills++
		victim.Stats.Deaths++
		if headshot {
			attacker.Stats.Headshots++
		}
		
		return killEvent
	}
	return nil
}

func (eg *EventGenerator) createFlashbangEvent(match *models.Match, state *models.MatchState, thrower *models.Player, tick int64, roundNum int) models.GameEvent {
	// Find players who might be flashed (simplified logic)
	var flashed []*models.Player
	
	// Get all alive players from the opposite team
	oppositeTeam := "TERRORIST"
	if thrower.Side == "TERRORIST" {
		oppositeTeam = "CT"
	}
	
	potentialVictims := eg.getAlivePlayers(match, state, oppositeTeam)
	
	// Randomly select 0-3 players to be flashed
	numFlashed := eg.rng.Intn(4) // 0-3 players
	if numFlashed > len(potentialVictims) {
		numFlashed = len(potentialVictims)
	}
	
	for i := 0; i < numFlashed; i++ {
		flashed = append(flashed, potentialVictims[i])
	}
	
	if len(flashed) > 0 {
		flashEvent := &models.FlashbangEvent{
			BaseEvent: models.NewBaseEvent("flashbang_detonate", tick+int64(2*eg.config.TickRate), roundNum),
			Player:    thrower,
			Position:  eg.getPlayerState(thrower).Position,
			Flashed:   flashed,
			Duration:  1.0 + eg.rng.Float64()*3.0, // 1-4 seconds flash
		}
		
		// Update thrower stats
		thrower.Stats.EnemiesFlashed += len(flashed)
		thrower.Stats.FlashAssists += len(flashed)
		
		return flashEvent
	}
	
	return nil
}

// Utility methods

func (eg *EventGenerator) getTeamBySide(match *models.Match, side string) *models.Team {
	for i := range match.Teams {
		if match.Teams[i].Side == side {
			return &match.Teams[i]
		}
	}
	return nil
}

func (eg *EventGenerator) getAlivePlayers(match *models.Match, state *models.MatchState, side string) []*models.Player {
	var alive []*models.Player
	team := eg.getTeamBySide(match, side)
	if team != nil {
		for i, player := range team.Players {
			if playerState := state.PlayerStates[player.Name]; playerState != nil && playerState.IsAlive {
				alive = append(alive, &team.Players[i])
			}
		}
	}
	return alive
}

func (eg *EventGenerator) getPlayerState(player *models.Player) *models.PlayerState {
	// This would need access to the match state - should be passed as parameter
	// For now, return a mock state or handle this differently
	return &models.PlayerState{
		IsAlive: true,
		Health:  100,
		Armor:   0,
	}
}

func (eg *EventGenerator) selectEngagementParticipants(players []*models.Player, strategy *RoundStrategy) []*models.Player {
	if len(players) == 0 {
		return players
	}
	
	maxParticipants := 3
	if strategy.Intensity > 0.7 {
		maxParticipants = int(math.Min(float64(len(players)), 4)) // More chaotic fights
	} else {
		maxParticipants = int(math.Min(float64(len(players)), 2)) // Smaller engagements
	}
	
	numParticipants := 1 + eg.rng.Intn(maxParticipants)
	if numParticipants > len(players) {
		numParticipants = len(players)
	}
	
	// Randomly select participants
	selected := make([]*models.Player, numParticipants)
	selectedIndices := eg.rng.Perm(len(players))[:numParticipants]
	
	for i, idx := range selectedIndices {
		selected[i] = players[idx]
	}
	
	return selected
}

func (eg *EventGenerator) removePlayerFromList(players []*models.Player, toRemove *models.Player) []*models.Player {
	var result []*models.Player
	for _, player := range players {
		if player.Name != toRemove.Name {
			result = append(result, player)
		}
	}
	return result
}

func (eg *EventGenerator) getPlayersWithUtility(team *models.Team, state *models.MatchState) []*models.Player {
	var players []*models.Player
	for i, player := range team.Players {
		if playerState := state.PlayerStates[player.Name]; playerState != nil && len(playerState.Grenades) > 0 {
			players = append(players, &team.Players[i])
		}
	}
	return players
}

func (eg *EventGenerator) selectWeaponForAttack(attacker *models.Player) string {
	// Default based on side
	if attacker.Side == "CT" {
		return "m4a4"
	}
	return "ak47"
}

func (eg *EventGenerator) calculateDamage(attacker, victim *models.Player, weapon string) int {
	baseDamage := 25
	
	// Weapon-specific damage
	switch weapon {
	case "awp":
		baseDamage = 115
	case "ak47":
		baseDamage = 36
	case "m4a4", "m4a1_silencer":
		baseDamage = 33
	case "deagle":
		baseDamage = 63
	case "glock", "usp_silencer":
		baseDamage = 28
	}
	
	// Add randomness (±20%)
	variation := int(float64(baseDamage) * 0.2)
	damage := baseDamage + eg.rng.Intn(variation*2) - variation
	
	if damage < 1 {
		damage = 1
	}
	
	return damage
}

func (eg *EventGenerator) calculateArmorDamage(damage int, playerState *models.PlayerState) int {
	if playerState.Armor <= 0 {
		return 0
	}
	
	armorDamage := int(float64(damage) * 0.5) // Simplified armor calculation
	if armorDamage > playerState.Armor {
		armorDamage = playerState.Armor
	}
	
	return armorDamage
}

func (eg *EventGenerator) selectHitgroup(attacker *models.Player, weapon string) int {
	// Hitgroup probabilities
	// 0=generic, 1=head, 2=chest, 3=stomach, 4=leftarm, 5=rightarm, 6=leftleg, 7=rightleg
	
	rand := eg.rng.Float64()
	
	// Adjusted for skill level
	headshotChance := 0.15
	if attacker.Profile.AimSkill > 0.8 {
		headshotChance = 0.25
	} else if attacker.Profile.AimSkill < 0.3 {
		headshotChance = 0.08
	}
	
	if rand < headshotChance {
		return 1 // Head
	} else if rand < 0.5 {
		return 2 // Chest
	} else if rand < 0.65 {
		return 3 // Stomach
	} else if rand < 0.8 {
		// Arms
		if eg.rng.Float64() < 0.5 {
			return 4 // Left arm
		}
		return 5 // Right arm
	} else {
		// Legs
		if eg.rng.Float64() < 0.5 {
			return 6 // Left leg
		}
		return 7 // Right leg
	}
}

func (eg *EventGenerator) copyTeamEconomies(economies map[string]*models.TeamEconomy) map[string]models.TeamEconomy {
	copied := make(map[string]models.TeamEconomy)
	for name, economy := range economies {
		copied[name] = *economy
	}
	return copied
}

func (eg *EventGenerator) sortEventsByTime(events []models.GameEvent) {
	// Simple bubble sort by tick - for production use sort.Slice
	n := len(events)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if events[j].GetTick() > events[j+1].GetTick() {
				events[j], events[j+1] = events[j+1], events[j]
			}
		}
	}
}