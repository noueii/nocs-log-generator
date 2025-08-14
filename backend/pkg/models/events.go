package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// GameEvent represents a base interface for all game events
type GameEvent interface {
	GetTimestamp() time.Time
	GetType() string
	GetTick() int64
	ToLogLine() string
	ToJSON() ([]byte, error)
}

// BaseEvent provides common fields for all events
type BaseEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	Tick      int64     `json:"tick"`
	Round     int       `json:"round"`
}

// GetTimestamp returns the event timestamp
func (e *BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetType returns the event type
func (e *BaseEvent) GetType() string {
	return e.Type
}

// GetTick returns the server tick
func (e *BaseEvent) GetTick() int64 {
	return e.Tick
}

// KillEvent represents a player kill event
type KillEvent struct {
	BaseEvent
	Attacker      *Player `json:"attacker"`
	Victim        *Player `json:"victim"`
	Assister      *Player `json:"assister,omitempty"`
	Weapon        string  `json:"weapon"`
	Headshot      bool    `json:"headshot"`
	Penetrated    int     `json:"penetrated"`
	NoScope       bool    `json:"no_scope"`
	AttackerBlind bool    `json:"attacker_blind"`
	Distance      float64 `json:"distance"`
	AttackerPos   Vector3 `json:"attacker_pos"`
	VictimPos     Vector3 `json:"victim_pos"`
}

// ToLogLine converts the kill event to CS2 log format
func (e *KillEvent) ToLogLine() string {
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	
	attackerInfo := fmt.Sprintf(`"%s<%d><%s><%s>"`, 
		e.Attacker.Name, e.Attacker.UserID, e.Attacker.SteamID, e.Attacker.Side)
	victimInfo := fmt.Sprintf(`"%s<%d><%s><%s>"`, 
		e.Victim.Name, e.Victim.UserID, e.Victim.SteamID, e.Victim.Side)
	
	logLine := fmt.Sprintf(`L %s: %s killed %s with "%s"`, 
		timestamp, attackerInfo, victimInfo, e.Weapon)
	
	if e.Headshot {
		logLine += " (headshot)"
	}
	if e.Penetrated > 0 {
		logLine += fmt.Sprintf(" (penetrated %d)", e.Penetrated)
	}
	if e.NoScope {
		logLine += " (noscope)"
	}
	if e.AttackerBlind {
		logLine += " (attackerblind)"
	}
	
	return logLine
}

// ToJSON converts the event to JSON
func (e *KillEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// RoundStartEvent represents the start of a round
type RoundStartEvent struct {
	BaseEvent
	CTScore      int                    `json:"ct_score"`
	TScore       int                    `json:"t_score"`
	CTPlayers    int                    `json:"ct_players"`
	TPlayers     int                    `json:"t_players"`
	TeamEconomies map[string]TeamEconomy `json:"team_economies"`
}

// ToLogLine converts the round start event to CS2 log format
func (e *RoundStartEvent) ToLogLine() string {
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	
	lines := []string{
		fmt.Sprintf(`L %s: World triggered "Round_Start"`, timestamp),
		fmt.Sprintf(`L %s: Team "CT" scored "%d" with "%d" players`, 
			timestamp, e.CTScore, e.CTPlayers),
		fmt.Sprintf(`L %s: Team "TERRORIST" scored "%d" with "%d" players`, 
			timestamp, e.TScore, e.TPlayers),
	}
	
	return strings.Join(lines, "\n")
}

// ToJSON converts the event to JSON
func (e *RoundStartEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// RoundEndEvent represents the end of a round
type RoundEndEvent struct {
	BaseEvent
	Winner       string `json:"winner"`       // "CT" or "TERRORIST"
	Reason       string `json:"reason"`       // "elimination", "bomb_defused", "bomb_exploded", "time"
	CTScore      int    `json:"ct_score"`
	TScore       int    `json:"t_score"`
	MVP          *Player `json:"mvp,omitempty"`
}

// ToLogLine converts the round end event to CS2 log format
func (e *RoundEndEvent) ToLogLine() string {
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	
	reasonMap := map[string]string{
		"elimination":   "Terrorists_Win",
		"bomb_exploded": "Target_Bombed",
		"bomb_defused":  "Bomb_Defused",
		"time":          "CTs_Win",
	}
	
	logReason := reasonMap[e.Reason]
	if logReason == "" {
		logReason = e.Reason
	}
	
	logLine := fmt.Sprintf(`L %s: Team "%s" triggered "%s" (CT "%d") (T "%d")`, 
		timestamp, e.Winner, logReason, e.CTScore, e.TScore)
	
	if e.MVP != nil {
		logLine += fmt.Sprintf(`\nL %s: "%s<%d><%s><%s>" triggered "MVP"`, 
			timestamp, e.MVP.Name, e.MVP.UserID, e.MVP.SteamID, e.MVP.Side)
	}
	
	return logLine
}

// ToJSON converts the event to JSON
func (e *RoundEndEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// BombPlantEvent represents a bomb plant event
type BombPlantEvent struct {
	BaseEvent
	Player   *Player `json:"player"`
	Site     string  `json:"site"`     // "A" or "B"
	Position Vector3 `json:"position"`
}

// ToLogLine converts the bomb plant event to CS2 log format
func (e *BombPlantEvent) ToLogLine() string {
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	
	playerInfo := fmt.Sprintf(`"%s<%d><%s><%s>"`, 
		e.Player.Name, e.Player.UserID, e.Player.SteamID, e.Player.Side)
	
	return fmt.Sprintf(`L %s: %s triggered "Planted_The_Bomb" at bombsite %s`, 
		timestamp, playerInfo, e.Site)
}

// ToJSON converts the event to JSON
func (e *BombPlantEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// BombDefuseEvent represents a bomb defuse event
type BombDefuseEvent struct {
	BaseEvent
	Player   *Player `json:"player"`
	Site     string  `json:"site"`
	WithKit  bool    `json:"with_kit"`
	Position Vector3 `json:"position"`
}

// ToLogLine converts the bomb defuse event to CS2 log format
func (e *BombDefuseEvent) ToLogLine() string {
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	
	playerInfo := fmt.Sprintf(`"%s<%d><%s><%s>"`, 
		e.Player.Name, e.Player.UserID, e.Player.SteamID, e.Player.Side)
	
	kitInfo := ""
	if e.WithKit {
		kitInfo = " (with kit)"
	}
	
	return fmt.Sprintf(`L %s: %s triggered "Defused_The_Bomb"%s`, 
		timestamp, playerInfo, kitInfo)
}

// ToJSON converts the event to JSON
func (e *BombDefuseEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// BombExplodeEvent represents a bomb explosion event
type BombExplodeEvent struct {
	BaseEvent
	Site     string  `json:"site"`
	Position Vector3 `json:"position"`
}

// ToLogLine converts the bomb explode event to CS2 log format
func (e *BombExplodeEvent) ToLogLine() string {
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	return fmt.Sprintf(`L %s: World triggered "Target_Bombed"`, timestamp)
}

// ToJSON converts the event to JSON
func (e *BombExplodeEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// PlayerHurtEvent represents a player damage event
type PlayerHurtEvent struct {
	BaseEvent
	Attacker   *Player `json:"attacker"`
	Victim     *Player `json:"victim"`
	Weapon     string  `json:"weapon"`
	Damage     int     `json:"damage"`
	DamageArmor int    `json:"damage_armor"`
	Health     int     `json:"health"`
	Armor      int     `json:"armor"`
	Hitgroup   int     `json:"hitgroup"` // 0=generic, 1=head, 2=chest, 3=stomach, 4=leftarm, 5=rightarm, 6=leftleg, 7=rightleg
}

// ToLogLine converts the player hurt event to CS2 log format
func (e *PlayerHurtEvent) ToLogLine() string {
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	
	attackerInfo := fmt.Sprintf(`"%s<%d><%s><%s>"`, 
		e.Attacker.Name, e.Attacker.UserID, e.Attacker.SteamID, e.Attacker.Side)
	victimInfo := fmt.Sprintf(`"%s<%d><%s><%s>"`, 
		e.Victim.Name, e.Victim.UserID, e.Victim.SteamID, e.Victim.Side)
	
	return fmt.Sprintf(`L %s: %s attacked %s with "%s" (damage "%d") (damage_armor "%d") (health "%d") (armor "%d") (hitgroup "%d")`, 
		timestamp, attackerInfo, victimInfo, e.Weapon, e.Damage, e.DamageArmor, e.Health, e.Armor, e.Hitgroup)
}

// ToJSON converts the event to JSON
func (e *PlayerHurtEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// PlayerConnectEvent represents a player connection event
type PlayerConnectEvent struct {
	BaseEvent
	Player  *Player `json:"player"`
	Address string  `json:"address"`
}

// ToLogLine converts the player connect event to CS2 log format
func (e *PlayerConnectEvent) ToLogLine() string {
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	
	return fmt.Sprintf(`L %s: "%s<%d><%s><>" connected, address "%s"`, 
		timestamp, e.Player.Name, e.Player.UserID, e.Player.SteamID, e.Address)
}

// ToJSON converts the event to JSON
func (e *PlayerConnectEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// PlayerDisconnectEvent represents a player disconnection event
type PlayerDisconnectEvent struct {
	BaseEvent
	Player *Player `json:"player"`
	Reason string  `json:"reason"`
}

// ToLogLine converts the player disconnect event to CS2 log format
func (e *PlayerDisconnectEvent) ToLogLine() string {
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	
	playerInfo := fmt.Sprintf(`"%s<%d><%s><%s>"`, 
		e.Player.Name, e.Player.UserID, e.Player.SteamID, e.Player.Side)
	
	return fmt.Sprintf(`L %s: %s disconnected (reason "%s")`, 
		timestamp, playerInfo, e.Reason)
}

// ToJSON converts the event to JSON
func (e *PlayerDisconnectEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// ItemPurchaseEvent represents an equipment purchase event
type ItemPurchaseEvent struct {
	BaseEvent
	Player *Player `json:"player"`
	Item   string  `json:"item"`
	Cost   int     `json:"cost"`
}

// ToLogLine converts the purchase event to CS2 log format
func (e *ItemPurchaseEvent) ToLogLine() string {
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	
	playerInfo := fmt.Sprintf(`"%s<%d><%s><%s>"`, 
		e.Player.Name, e.Player.UserID, e.Player.SteamID, e.Player.Side)
	
	return fmt.Sprintf(`L %s: %s purchased "%s"`, 
		timestamp, playerInfo, e.Item)
}

// ToJSON converts the event to JSON
func (e *ItemPurchaseEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// GrenadeThrowEvent represents a grenade thrown event
type GrenadeThrowEvent struct {
	BaseEvent
	Player      *Player `json:"player"`
	GrenadeType string  `json:"grenade_type"`
	Position    Vector3 `json:"position"`
	Velocity    Vector3 `json:"velocity"`
}

// ToLogLine converts the grenade throw event to CS2 log format
func (e *GrenadeThrowEvent) ToLogLine() string {
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	
	playerInfo := fmt.Sprintf(`"%s<%d><%s><%s>"`, 
		e.Player.Name, e.Player.UserID, e.Player.SteamID, e.Player.Side)
	
	return fmt.Sprintf(`L %s: %s threw %s`, 
		timestamp, playerInfo, e.GrenadeType)
}

// ToJSON converts the event to JSON
func (e *GrenadeThrowEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// WeaponFireEvent represents a weapon fire event
type WeaponFireEvent struct {
	BaseEvent
	Player   *Player `json:"player"`
	Weapon   string  `json:"weapon"`
	Position Vector3 `json:"position"`
	Angle    Vector3 `json:"angle"`
	Silenced bool    `json:"silenced"`
}

// ToLogLine converts the weapon fire event to CS2 log format
func (e *WeaponFireEvent) ToLogLine() string {
	// Note: Weapon fire events are typically not logged in standard CS2 logs
	// This is more for internal tracking/analysis
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	
	playerInfo := fmt.Sprintf(`"%s<%d><%s><%s>"`, 
		e.Player.Name, e.Player.UserID, e.Player.SteamID, e.Player.Side)
	
	return fmt.Sprintf(`L %s: %s fired %s`, 
		timestamp, playerInfo, e.Weapon)
}

// ToJSON converts the event to JSON
func (e *WeaponFireEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// FlashbangEvent represents a flashbang detonation event
type FlashbangEvent struct {
	BaseEvent
	Player    *Player   `json:"player"`
	Position  Vector3   `json:"position"`
	Flashed   []*Player `json:"flashed"`   // Players that were flashed
	Duration  float64   `json:"duration"`  // Flash duration in seconds
}

// ToLogLine converts the flashbang event to CS2 log format
func (e *FlashbangEvent) ToLogLine() string {
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	
	playerInfo := fmt.Sprintf(`"%s<%d><%s><%s>"`, 
		e.Player.Name, e.Player.UserID, e.Player.SteamID, e.Player.Side)
	
	logLine := fmt.Sprintf(`L %s: %s threw flashbang`, timestamp, playerInfo)
	
	for _, flashed := range e.Flashed {
		flashedInfo := fmt.Sprintf(`"%s<%d><%s><%s>"`, 
			flashed.Name, flashed.UserID, flashed.SteamID, flashed.Side)
		logLine += fmt.Sprintf(`\nL %s: %s blinded %s with flashbang for %.1f`, 
			timestamp, playerInfo, flashedInfo, e.Duration)
	}
	
	return logLine
}

// ToJSON converts the event to JSON
func (e *FlashbangEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// ChatEvent represents a chat message event
type ChatEvent struct {
	BaseEvent
	Player  *Player `json:"player,omitempty"`
	Message string  `json:"message"`
	Team    bool    `json:"team"`    // true for team chat, false for all chat
	Dead    bool    `json:"dead"`    // true if player is dead
}

// ToLogLine converts the chat event to CS2 log format
func (e *ChatEvent) ToLogLine() string {
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	
	if e.Player == nil {
		// Server message
		return fmt.Sprintf(`L %s: Server say "%s"`, timestamp, e.Message)
	}
	
	playerInfo := fmt.Sprintf(`"%s<%d><%s><%s>"`, 
		e.Player.Name, e.Player.UserID, e.Player.SteamID, e.Player.Side)
	
	chatType := "say"
	if e.Team {
		chatType = "say_team"
	}
	if e.Dead {
		chatType += "_dead"
	}
	
	return fmt.Sprintf(`L %s: %s %s "%s"`, 
		timestamp, playerInfo, chatType, e.Message)
}

// ToJSON converts the event to JSON
func (e *ChatEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// TeamSwitchEvent represents a player switching teams
type TeamSwitchEvent struct {
	BaseEvent
	Player  *Player `json:"player"`
	FromTeam string `json:"from_team"`
	ToTeam   string `json:"to_team"`
}

// ToLogLine converts the team switch event to CS2 log format
func (e *TeamSwitchEvent) ToLogLine() string {
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	
	playerInfo := fmt.Sprintf(`"%s<%d><%s><%s>"`, 
		e.Player.Name, e.Player.UserID, e.Player.SteamID, e.FromTeam)
	
	return fmt.Sprintf(`L %s: %s switched from team <%s> to <%s>`, 
		timestamp, playerInfo, e.FromTeam, e.ToTeam)
}

// ToJSON converts the event to JSON
func (e *TeamSwitchEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// ServerCommandEvent represents a server command execution
type ServerCommandEvent struct {
	BaseEvent
	Command string `json:"command"`
	Args    string `json:"args"`
	Result  string `json:"result,omitempty"`
}

// ToLogLine converts the server command event to CS2 log format
func (e *ServerCommandEvent) ToLogLine() string {
	timestamp := e.Timestamp.Format("01/02/2006 - 15:04:05")
	
	return fmt.Sprintf(`L %s: Server cvar "%s" = "%s"`, 
		timestamp, e.Command, e.Args)
}

// ToJSON converts the event to JSON
func (e *ServerCommandEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// NewBaseEvent creates a new base event with current timestamp
func NewBaseEvent(eventType string, tick int64, round int) BaseEvent {
	return BaseEvent{
		Timestamp: time.Now(),
		Type:      eventType,
		Tick:      tick,
		Round:     round,
	}
}

// EventFactory provides factory methods for creating events
type EventFactory struct {
	currentTick  int64
	currentRound int
}

// NewEventFactory creates a new event factory
func NewEventFactory() *EventFactory {
	return &EventFactory{
		currentTick:  0,
		currentRound: 0,
	}
}

// SetTick sets the current server tick
func (f *EventFactory) SetTick(tick int64) {
	f.currentTick = tick
}

// SetRound sets the current round number
func (f *EventFactory) SetRound(round int) {
	f.currentRound = round
}

// CreateKillEvent creates a new kill event
func (f *EventFactory) CreateKillEvent(attacker, victim *Player, weapon string, headshot bool) *KillEvent {
	return &KillEvent{
		BaseEvent: NewBaseEvent("player_death", f.currentTick, f.currentRound),
		Attacker:  attacker,
		Victim:    victim,
		Weapon:    weapon,
		Headshot:  headshot,
	}
}

// CreateRoundStartEvent creates a new round start event
func (f *EventFactory) CreateRoundStartEvent(ctScore, tScore, ctPlayers, tPlayers int) *RoundStartEvent {
	return &RoundStartEvent{
		BaseEvent: NewBaseEvent("round_start", f.currentTick, f.currentRound),
		CTScore:   ctScore,
		TScore:    tScore,
		CTPlayers: ctPlayers,
		TPlayers:  tPlayers,
	}
}

// CreateRoundEndEvent creates a new round end event
func (f *EventFactory) CreateRoundEndEvent(winner, reason string, ctScore, tScore int, mvp *Player) *RoundEndEvent {
	return &RoundEndEvent{
		BaseEvent: NewBaseEvent("round_end", f.currentTick, f.currentRound),
		Winner:    winner,
		Reason:    reason,
		CTScore:   ctScore,
		TScore:    tScore,
		MVP:       mvp,
	}
}

