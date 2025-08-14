package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// MatchConfig represents the configuration for a match
type MatchConfig struct {
	// Basic match settings
	Format       string `json:"format" binding:"required,oneof=mr12 mr15"`
	Map          string `json:"map" binding:"required"`
	Overtime     bool   `json:"overtime"`
	MaxRounds    int    `json:"max_rounds,omitempty"`
	
	// Server settings
	TickRate     int    `json:"tick_rate"`
	ServerName   string `json:"server_name,omitempty"`
	
	// Simulation settings
	Seed         int64  `json:"seed,omitempty"`
	Duration     time.Duration `json:"duration,omitempty"`
	
	// Rollback settings
	RollbackEnabled     bool    `json:"rollback_enabled"`
	RollbackProbability float64 `json:"rollback_probability"`
	RollbackMinRound    int     `json:"rollback_min_round"`
	RollbackMaxRound    int     `json:"rollback_max_round"`
	
	// Economy settings
	StartMoney          int  `json:"start_money"`
	MaxMoney            int  `json:"max_money"`
	RealisticEconomy    bool `json:"realistic_economy"`
	
	// Advanced settings
	NetworkIssues       bool    `json:"network_issues"`
	AntiCheatEvents     bool    `json:"anti_cheat_events"`
	ChatMessages        bool    `json:"chat_messages"`
	SkillVariance       float64 `json:"skill_variance"`
	
	// Output settings
	LogFormat           string `json:"log_format"`      // "standard", "json", "custom"
	TimestampFormat     string `json:"timestamp_format"`
	OutputVerbosity     string `json:"output_verbosity"` // "minimal", "standard", "verbose"
	IncludePositions    bool   `json:"include_positions"`
	IncludeWeaponFire   bool   `json:"include_weapon_fire"`
}

// SimulationConfig represents configuration for match simulation
type SimulationConfig struct {
	// Performance settings
	EventsPerSecond     int     `json:"events_per_second"`
	MaxConcurrentMatches int    `json:"max_concurrent_matches"`
	BufferSize          int     `json:"buffer_size"`
	
	// Realism settings
	PlayerBehaviorRealism float64 `json:"player_behavior_realism"` // 0.0 to 1.0
	EconomicRealism      float64 `json:"economic_realism"`        // 0.0 to 1.0
	PositionalRealism    float64 `json:"positional_realism"`      // 0.0 to 1.0
	
	// Randomization
	RandomSeed          int64   `json:"random_seed"`
	SkillVariation      float64 `json:"skill_variation"`
	WeaponAccuracy      float64 `json:"weapon_accuracy"`
	
	// Event probabilities
	TeamKillProbability float64 `json:"team_kill_probability"`
	FlashAssistProbability float64 `json:"flash_assist_probability"`
	WallBangProbability float64 `json:"wallbang_probability"`
	
	// Chat and communication
	ChatFrequency       float64 `json:"chat_frequency"`
	RadioCommandFreq    float64 `json:"radio_command_frequency"`
	DeathCamComments    bool    `json:"death_cam_comments"`
	
	// Network simulation
	NetworkDelay        time.Duration `json:"network_delay"`
	PacketLoss          float64       `json:"packet_loss"`
	JitterVariance      time.Duration `json:"jitter_variance"`
}

// ServerConfig represents server-specific configuration
type ServerConfig struct {
	// Server identification
	ServerName          string `json:"server_name"`
	ServerIP            string `json:"server_ip"`
	ServerPort          int    `json:"server_port"`
	Password            string `json:"password,omitempty"`
	
	// Game settings
	TickRate            int    `json:"tick_rate"`
	FPS                 int    `json:"fps"`
	
	// Round settings
	RoundTime           int    `json:"round_time"`           // seconds
	FreezetimeLength    int    `json:"freezetime_length"`    // seconds
	BuyTime             int    `json:"buy_time"`             // seconds
	BombTimer           int    `json:"bomb_timer"`           // seconds
	DefuseTime          int    `json:"defuse_time"`          // seconds (with kit)
	DefuseTimeNoKit     int    `json:"defuse_time_no_kit"`   // seconds (without kit)
	
	// Economy settings
	StartMoney          int    `json:"start_money"`
	MaxMoney            int    `json:"max_money"`
	
	// Gameplay settings
	FriendlyFire        bool   `json:"friendly_fire"`
	AutoBalance         bool   `json:"auto_balance"`
	RestartGame         int    `json:"restart_game"`
	
	// Anti-cheat settings
	VACEnabled          bool   `json:"vac_enabled"`
	PureServer          bool   `json:"pure_server"`
	
	// Communication settings
	AllTalk             bool   `json:"all_talk"`
	TeamTalk            bool   `json:"team_talk"`
	DeadTalk            bool   `json:"dead_talk"`
	
	// Admin settings
	RCONPassword        string `json:"rcon_password,omitempty"`
	AdminPassword       string `json:"admin_password,omitempty"`
	
	// Logging settings
	LogToFile           bool   `json:"log_to_file"`
	LogDetail           int    `json:"log_detail"`
	LogBans             bool   `json:"log_bans"`
}

// ParserConfig represents configuration for demo parsing
type ParserConfig struct {
	// Input settings
	DemoPath            string   `json:"demo_path"`
	DemoURL             string   `json:"demo_url"`
	DemoBase64          string   `json:"demo_base64"`
	
	// Output settings
	OutputFormat        string   `json:"output_format"`        // "http_log", "json", "csv"
	OutputPath          string   `json:"output_path"`
	IncludeRaw          bool     `json:"include_raw"`
	
	// Event filtering
	EventFilter         []string `json:"event_filter"`         // List of event types to include
	PlayerFilter        []string `json:"player_filter"`        // List of players to track
	RoundFilter         []int    `json:"round_filter"`          // List of rounds to include
	
	// Processing settings
	StartTick           int64    `json:"start_tick"`
	EndTick             int64    `json:"end_tick"`
	SkipWarmup          bool     `json:"skip_warmup"`
	SkipKnifing         bool     `json:"skip_knifing"`
	
	// Data extraction
	ExtractPositions    bool     `json:"extract_positions"`
	ExtractGrenadePaths bool     `json:"extract_grenade_paths"`
	ExtractWeaponStates bool     `json:"extract_weapon_states"`
	ExtractChatLog      bool     `json:"extract_chat_log"`
	
	// Performance settings
	BufferSize          int      `json:"buffer_size"`
	MaxMemory           int64    `json:"max_memory"`          // bytes
	
	// Error handling
	SkipErrors          bool     `json:"skip_errors"`
	StrictMode          bool     `json:"strict_mode"`
}

// DefaultMatchConfig returns a default match configuration
func DefaultMatchConfig() MatchConfig {
	return MatchConfig{
		Format:              "mr12",
		Map:                 "de_mirage",
		Overtime:            false,
		TickRate:            64,
		StartMoney:          800,
		MaxMoney:            16000,
		RealisticEconomy:    true,
		RollbackEnabled:     false,
		RollbackProbability: 0.0,
		NetworkIssues:       false,
		AntiCheatEvents:     false,
		ChatMessages:        true,
		SkillVariance:       0.15,
		LogFormat:           "standard",
		TimestampFormat:     "01/02/2006 - 15:04:05",
		OutputVerbosity:     "standard",
		IncludePositions:    false,
		IncludeWeaponFire:   false,
	}
}

// DefaultSimulationConfig returns a default simulation configuration
func DefaultSimulationConfig() SimulationConfig {
	return SimulationConfig{
		EventsPerSecond:          1000,
		MaxConcurrentMatches:     10,
		BufferSize:              10000,
		PlayerBehaviorRealism:    0.8,
		EconomicRealism:          0.9,
		PositionalRealism:        0.7,
		SkillVariation:           0.2,
		WeaponAccuracy:           0.8,
		TeamKillProbability:      0.001,
		FlashAssistProbability:   0.1,
		WallBangProbability:      0.05,
		ChatFrequency:            0.1,
		RadioCommandFreq:         0.05,
		DeathCamComments:         true,
		NetworkDelay:             time.Millisecond * 30,
		PacketLoss:               0.001,
		JitterVariance:           time.Millisecond * 5,
	}
}

// DefaultServerConfig returns a default server configuration
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		ServerName:          "CS2 Log Generator Server",
		ServerIP:            "127.0.0.1",
		ServerPort:          27015,
		TickRate:            64,
		FPS:                 300,
		RoundTime:           115,
		FreezetimeLength:    15,
		BuyTime:             20,
		BombTimer:           40,
		DefuseTime:          5,
		DefuseTimeNoKit:     10,
		StartMoney:          800,
		MaxMoney:            16000,
		FriendlyFire:        true,
		AutoBalance:         false,
		VACEnabled:          true,
		PureServer:          true,
		AllTalk:             false,
		TeamTalk:            true,
		DeadTalk:            false,
		LogToFile:           true,
		LogDetail:           3,
		LogBans:             true,
	}
}

// DefaultParserConfig returns a default parser configuration
func DefaultParserConfig() ParserConfig {
	return ParserConfig{
		OutputFormat:        "http_log",
		IncludeRaw:          false,
		SkipWarmup:          true,
		SkipKnifing:         true,
		ExtractPositions:    false,
		ExtractGrenadePaths: false,
		ExtractWeaponStates: false,
		ExtractChatLog:      true,
		BufferSize:          10000,
		MaxMemory:           500 * 1024 * 1024, // 500MB
		SkipErrors:          false,
		StrictMode:          false,
	}
}

// Validate validates the match configuration
func (c *MatchConfig) Validate() error {
	if c.Format != "mr12" && c.Format != "mr15" {
		return errors.New("format must be 'mr12' or 'mr15'")
	}
	
	if strings.TrimSpace(c.Map) == "" {
		return errors.New("map is required")
	}
	
	if c.TickRate != 0 && (c.TickRate < 64 || c.TickRate > 128) {
		return errors.New("tick rate must be between 64 and 128")
	}
	
	if c.RollbackProbability < 0 || c.RollbackProbability > 1 {
		return errors.New("rollback probability must be between 0 and 1")
	}
	
	if c.SkillVariance < 0 || c.SkillVariance > 1 {
		return errors.New("skill variance must be between 0 and 1")
	}
	
	if c.StartMoney < 0 || c.StartMoney > c.MaxMoney {
		return errors.New("start money must be between 0 and max money")
	}
	
	return nil
}

// Validate validates the simulation configuration
func (c *SimulationConfig) Validate() error {
	if c.EventsPerSecond <= 0 {
		return errors.New("events per second must be positive")
	}
	
	if c.MaxConcurrentMatches <= 0 {
		return errors.New("max concurrent matches must be positive")
	}
	
	if c.BufferSize <= 0 {
		return errors.New("buffer size must be positive")
	}
	
	// Validate realism values (0.0 to 1.0)
	realism := []struct {
		name  string
		value float64
	}{
		{"player behavior realism", c.PlayerBehaviorRealism},
		{"economic realism", c.EconomicRealism},
		{"positional realism", c.PositionalRealism},
		{"skill variation", c.SkillVariation},
		{"weapon accuracy", c.WeaponAccuracy},
	}
	
	for _, r := range realism {
		if r.value < 0 || r.value > 1 {
			return fmt.Errorf("%s must be between 0.0 and 1.0", r.name)
		}
	}
	
	// Validate probabilities
	probabilities := []struct {
		name  string
		value float64
	}{
		{"team kill probability", c.TeamKillProbability},
		{"flash assist probability", c.FlashAssistProbability},
		{"wallbang probability", c.WallBangProbability},
		{"chat frequency", c.ChatFrequency},
		{"radio command frequency", c.RadioCommandFreq},
		{"packet loss", c.PacketLoss},
	}
	
	for _, p := range probabilities {
		if p.value < 0 || p.value > 1 {
			return fmt.Errorf("%s must be between 0.0 and 1.0", p.name)
		}
	}
	
	return nil
}

// Validate validates the server configuration
func (c *ServerConfig) Validate() error {
	if strings.TrimSpace(c.ServerName) == "" {
		return errors.New("server name is required")
	}
	
	if c.ServerPort <= 0 || c.ServerPort > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	
	if c.TickRate != 64 && c.TickRate != 128 {
		return errors.New("tick rate must be 64 or 128")
	}
	
	if c.RoundTime <= 0 || c.RoundTime > 300 {
		return errors.New("round time must be between 1 and 300 seconds")
	}
	
	if c.FreezetimeLength < 0 || c.FreezetimeLength > 60 {
		return errors.New("freezetime must be between 0 and 60 seconds")
	}
	
	if c.BuyTime <= 0 || c.BuyTime > c.FreezetimeLength {
		return errors.New("buy time must be positive and not exceed freezetime")
	}
	
	return nil
}

// Validate validates the parser configuration
func (c *ParserConfig) Validate() error {
	// Must have one input source
	inputSources := 0
	if c.DemoPath != "" {
		inputSources++
	}
	if c.DemoURL != "" {
		inputSources++
	}
	if c.DemoBase64 != "" {
		inputSources++
	}
	
	if inputSources != 1 {
		return errors.New("exactly one demo input source must be specified")
	}
	
	// Validate output format
	validFormats := []string{"http_log", "json", "csv"}
	validFormat := false
	for _, format := range validFormats {
		if c.OutputFormat == format {
			validFormat = true
			break
		}
	}
	if !validFormat {
		return fmt.Errorf("output format must be one of: %s", strings.Join(validFormats, ", "))
	}
	
	if c.BufferSize <= 0 {
		return errors.New("buffer size must be positive")
	}
	
	if c.MaxMemory <= 0 {
		return errors.New("max memory must be positive")
	}
	
	if c.StartTick < 0 {
		return errors.New("start tick must be non-negative")
	}
	
	if c.EndTick > 0 && c.EndTick <= c.StartTick {
		return errors.New("end tick must be greater than start tick")
	}
	
	return nil
}

// GetMaxRounds returns the maximum number of rounds for the format
func (c *MatchConfig) GetMaxRounds() int {
	if c.MaxRounds > 0 {
		return c.MaxRounds
	}
	
	switch c.Format {
	case "mr12":
		return 24
	case "mr15":
		return 30
	default:
		return 24
	}
}

// GetWinThreshold returns the number of rounds needed to win
func (c *MatchConfig) GetWinThreshold() int {
	return (c.GetMaxRounds() / 2) + 1
}

// IsValidMap checks if a map name is valid
func (c *MatchConfig) IsValidMap() bool {
	validMaps := []string{
		"de_mirage", "de_dust2", "de_inferno", "de_cache", "de_overpass",
		"de_train", "de_nuke", "de_cbble", "de_vertigo", "de_ancient",
	}
	
	for _, validMap := range validMaps {
		if strings.EqualFold(c.Map, validMap) {
			return true
		}
	}
	
	return false
}

// ApplyProfile applies a predefined configuration profile
func (c *MatchConfig) ApplyProfile(profileName string) {
	switch strings.ToLower(profileName) {
	case "competitive":
		c.RealisticEconomy = true
		c.SkillVariance = 0.1
		c.NetworkIssues = false
		c.AntiCheatEvents = true
		c.ChatMessages = false
		
	case "casual":
		c.RealisticEconomy = false
		c.SkillVariance = 0.3
		c.NetworkIssues = true
		c.AntiCheatEvents = false
		c.ChatMessages = true
		
	case "testing":
		c.RollbackEnabled = true
		c.RollbackProbability = 0.1
		c.NetworkIssues = true
		c.IncludePositions = true
		c.IncludeWeaponFire = true
		c.OutputVerbosity = "verbose"
		
	case "minimal":
		c.ChatMessages = false
		c.IncludePositions = false
		c.IncludeWeaponFire = false
		c.OutputVerbosity = "minimal"
	}
}

// Clone creates a deep copy of the match configuration
func (c *MatchConfig) Clone() *MatchConfig {
	clone := *c
	return &clone
}

// Merge merges another configuration into this one (non-zero values override)
func (c *MatchConfig) Merge(other *MatchConfig) {
	if other.Format != "" {
		c.Format = other.Format
	}
	if other.Map != "" {
		c.Map = other.Map
	}
	if other.TickRate != 0 {
		c.TickRate = other.TickRate
	}
	if other.Seed != 0 {
		c.Seed = other.Seed
	}
	// Add more fields as needed...
}