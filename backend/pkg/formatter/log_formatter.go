package formatter

import (
	"fmt"
	"strings"
	"time"

	"github.com/noueii/nocs-log-generator/backend/pkg/models"
)

// LogFormatter handles conversion of match events to CS2 log format
type LogFormatter struct {
	config       *models.MatchConfig
	timeZone     *time.Location
	serverName   string
	mapName      string
	playerNames  map[string]string // Sanitized player names
}

// NewLogFormatter creates a new log formatter with the given configuration
func NewLogFormatter(config *models.MatchConfig) *LogFormatter {
	// Default to UTC timezone
	tz := time.UTC
	
	return &LogFormatter{
		config:       config,
		timeZone:     tz,
		serverName:   config.ServerName,
		mapName:      config.Map,
		playerNames:  make(map[string]string),
	}
}

// FormatMatch formats an entire match into CS2 log format
func (f *LogFormatter) FormatMatch(match *models.Match) []string {
	var lines []string
	
	// Add log header
	lines = append(lines, f.formatLogHeader(match))
	
	// Format all events
	for _, event := range match.Events {
		formatted := f.FormatEvent(event)
		if formatted != "" {
			// Handle multi-line events
			eventLines := strings.Split(formatted, "\n")
			for _, line := range eventLines {
				if line != "" {
					lines = append(lines, line)
				}
			}
		}
	}
	
	// Add log footer
	lines = append(lines, f.formatLogFooter(match))
	
	return lines
}

// FormatEvent formats a single event into CS2 log format
func (f *LogFormatter) FormatEvent(event models.GameEvent) string {
	if event == nil {
		return ""
	}
	
	// Use the event's built-in ToLogLine method
	return event.ToLogLine()
}

// FormatEventsToString formats multiple events as a single string
func (f *LogFormatter) FormatEventsToString(events []models.GameEvent) string {
	var lines []string
	
	for _, event := range events {
		formatted := f.FormatEvent(event)
		if formatted != "" {
			lines = append(lines, formatted)
		}
	}
	
	return strings.Join(lines, "\n")
}

// FormatMatchToString formats an entire match as a single string
func (f *LogFormatter) FormatMatchToString(match *models.Match) string {
	lines := f.FormatMatch(match)
	return strings.Join(lines, "\n")
}

// FormatRound formats all events from a specific round
func (f *LogFormatter) FormatRound(roundData models.RoundData) []string {
	var lines []string
	
	for _, event := range roundData.Events {
		formatted := f.FormatEvent(event)
		if formatted != "" {
			lines = append(lines, formatted)
		}
	}
	
	return lines
}

// formatLogHeader creates the standard CS2 log header
func (f *LogFormatter) formatLogHeader(match *models.Match) string {
	timestamp := match.StartTime.In(f.timeZone).Format("01/02/2006 - 15:04:05")
	
	header := fmt.Sprintf(`L %s: Log file started (file "logs/L%s.log") (game "%s") (version "%s")`, 
		timestamp, 
		match.StartTime.Format("010206"), 
		"Counter-Strike: Global Offensive",
		"1.38.5.5")
	
	// Add server info
	header += fmt.Sprintf(`\nL %s: server_cvar: "hostname" "%s"`, timestamp, f.serverName)
	header += fmt.Sprintf(`\nL %s: server_cvar: "mp_startmoney" "%d"`, timestamp, f.config.StartMoney)
	header += fmt.Sprintf(`\nL %s: server_cvar: "mp_maxmoney" "%d"`, timestamp, f.config.MaxMoney)
	header += fmt.Sprintf(`\nL %s: server_cvar: "mp_roundtime" "115"`, timestamp)
	header += fmt.Sprintf(`\nL %s: server_cvar: "mp_freezetime" "15"`, timestamp)
	header += fmt.Sprintf(`\nL %s: Loading map "%s"`, timestamp, f.mapName)
	header += fmt.Sprintf(`\nL %s: Started map "%s" (CRC "0")`, timestamp, f.mapName)
	
	return header
}

// formatLogFooter creates the standard CS2 log footer
func (f *LogFormatter) formatLogFooter(match *models.Match) string {
	timestamp := match.EndTime.In(f.timeZone).Format("01/02/2006 - 15:04:05")
	
	footer := fmt.Sprintf(`L %s: Log file closed`, timestamp)
	
	return footer
}

// formatPlayerInfo formats player information in CS2 log format
func (f *LogFormatter) formatPlayerInfo(player *models.Player) string {
	sanitizedName := f.sanitizePlayerName(player.Name)
	return fmt.Sprintf(`"%s<%d><%s><%s>"`, 
		sanitizedName, 
		player.UserID, 
		player.SteamID, 
		player.Side)
}

// sanitizePlayerName ensures player names are safe for log format
func (f *LogFormatter) sanitizePlayerName(name string) string {
	if sanitized, exists := f.playerNames[name]; exists {
		return sanitized
	}
	
	// Remove quotes and backslashes that could break log format
	// First escape backslashes, then escape quotes to avoid double escaping
	sanitized := strings.ReplaceAll(name, `\`, `\\`)
	sanitized = strings.ReplaceAll(sanitized, `"`, `\"`)
	
	// Remove control characters and non-printable characters
	var result strings.Builder
	for _, r := range sanitized {
		if r >= 32 && r <= 126 { // ASCII printable characters
			result.WriteRune(r)
		} else {
			result.WriteRune('_') // Replace with underscore
		}
	}
	
	sanitized = result.String()
	
	// Limit length to reasonable size
	if len(sanitized) > 31 {
		sanitized = sanitized[:31]
	}
	
	// Cache the sanitized name
	f.playerNames[name] = sanitized
	
	return sanitized
}

// formatTimestamp formats a timestamp in CS2 log format
func (f *LogFormatter) formatTimestamp(t time.Time) string {
	return t.In(f.timeZone).Format("01/02/2006 - 15:04:05")
}

// FormatPlayerConnect formats a player connection in standard CS2 format
func (f *LogFormatter) FormatPlayerConnect(player *models.Player, address string, timestamp time.Time) string {
	ts := f.formatTimestamp(timestamp)
	return fmt.Sprintf(`L %s: "%s<%d><%s><>" connected, address "%s"`, 
		ts, 
		f.sanitizePlayerName(player.Name), 
		player.UserID, 
		player.SteamID, 
		address)
}

// FormatPlayerDisconnect formats a player disconnection in standard CS2 format
func (f *LogFormatter) FormatPlayerDisconnect(player *models.Player, reason string, timestamp time.Time) string {
	ts := f.formatTimestamp(timestamp)
	return fmt.Sprintf(`L %s: %s disconnected (reason "%s")`, 
		ts, 
		f.formatPlayerInfo(player), 
		reason)
}

// FormatTeamSwitch formats a team switch event in standard CS2 format
func (f *LogFormatter) FormatTeamSwitch(player *models.Player, fromTeam, toTeam string, timestamp time.Time) string {
	ts := f.formatTimestamp(timestamp)
	return fmt.Sprintf(`L %s: %s switched from team <%s> to <%s>`, 
		ts, 
		f.formatPlayerInfo(player), 
		fromTeam, 
		toTeam)
}

// FormatKill formats a kill event with all modifiers
func (f *LogFormatter) FormatKill(attacker, victim *models.Player, weapon string, headshot bool, penetrated int, blind bool, noScope bool, distance float64, timestamp time.Time) string {
	ts := f.formatTimestamp(timestamp)
	
	baseLine := fmt.Sprintf(`L %s: %s killed %s with "%s"`, 
		ts, 
		f.formatPlayerInfo(attacker), 
		f.formatPlayerInfo(victim), 
		weapon)
	
	var modifiers []string
	if headshot {
		modifiers = append(modifiers, "headshot")
	}
	if penetrated > 0 {
		modifiers = append(modifiers, fmt.Sprintf("penetrated %d", penetrated))
	}
	if blind {
		modifiers = append(modifiers, "attackerblind")
	}
	if noScope {
		modifiers = append(modifiers, "noscope")
	}
	
	if len(modifiers) > 0 {
		baseLine += " (" + strings.Join(modifiers, ") (") + ")"
	}
	
	return baseLine
}

// FormatDamage formats a damage event
func (f *LogFormatter) FormatDamage(attacker, victim *models.Player, weapon string, damage, damageArmor, health, armor, hitgroup int, timestamp time.Time) string {
	ts := f.formatTimestamp(timestamp)
	
	return fmt.Sprintf(`L %s: %s attacked %s with "%s" (damage "%d") (damage_armor "%d") (health "%d") (armor "%d") (hitgroup "%d")`, 
		ts, 
		f.formatPlayerInfo(attacker), 
		f.formatPlayerInfo(victim), 
		weapon, 
		damage, 
		damageArmor, 
		health, 
		armor, 
		hitgroup)
}

// FormatPurchase formats an item purchase event
func (f *LogFormatter) FormatPurchase(player *models.Player, item string, cost int, timestamp time.Time) string {
	ts := f.formatTimestamp(timestamp)
	
	return fmt.Sprintf(`L %s: %s purchased "%s"`, 
		ts, 
		f.formatPlayerInfo(player), 
		item)
}

// FormatBombPlant formats a bomb plant event
func (f *LogFormatter) FormatBombPlant(player *models.Player, site string, timestamp time.Time) string {
	ts := f.formatTimestamp(timestamp)
	
	return fmt.Sprintf(`L %s: %s triggered "Planted_The_Bomb" at bombsite %s`, 
		ts, 
		f.formatPlayerInfo(player), 
		site)
}

// FormatBombDefuse formats a bomb defuse event
func (f *LogFormatter) FormatBombDefuse(player *models.Player, withKit bool, timestamp time.Time) string {
	ts := f.formatTimestamp(timestamp)
	
	kitInfo := ""
	if withKit {
		kitInfo = " (with kit)"
	}
	
	return fmt.Sprintf(`L %s: %s triggered "Defused_The_Bomb"%s`, 
		ts, 
		f.formatPlayerInfo(player), 
		kitInfo)
}

// FormatBombExplode formats a bomb explosion event
func (f *LogFormatter) FormatBombExplode(timestamp time.Time) string {
	ts := f.formatTimestamp(timestamp)
	
	return fmt.Sprintf(`L %s: World triggered "Target_Bombed"`, ts)
}

// FormatRoundStart formats a round start event
func (f *LogFormatter) FormatRoundStart(ctScore, tScore, ctPlayers, tPlayers int, timestamp time.Time) string {
	ts := f.formatTimestamp(timestamp)
	
	lines := []string{
		fmt.Sprintf(`L %s: World triggered "Round_Start"`, ts),
		fmt.Sprintf(`L %s: Team "CT" scored "%d" with "%d" players`, ts, ctScore, ctPlayers),
		fmt.Sprintf(`L %s: Team "TERRORIST" scored "%d" with "%d" players`, ts, tScore, tPlayers),
	}
	
	return strings.Join(lines, "\n")
}

// FormatRoundEnd formats a round end event
func (f *LogFormatter) FormatRoundEnd(winner, reason string, ctScore, tScore int, mvp *models.Player, timestamp time.Time) string {
	ts := f.formatTimestamp(timestamp)
	
	// Map reasons to CS2 log format
	reasonMap := map[string]string{
		"elimination":   "Terrorists_Win",
		"bomb_exploded": "Target_Bombed", 
		"bomb_defused":  "Bomb_Defused",
		"time":          "CTs_Win",
	}
	
	logReason := reasonMap[reason]
	if logReason == "" {
		logReason = reason
	}
	
	// Adjust for correct team name based on reason
	if reason == "elimination" && winner == "CT" {
		logReason = "CTs_Win"
	} else if reason == "time" && winner == "TERRORIST" {
		logReason = "Terrorists_Win"
	}
	
	baseLine := fmt.Sprintf(`L %s: Team "%s" triggered "%s" (CT "%d") (T "%d")`, 
		ts, winner, logReason, ctScore, tScore)
	
	if mvp != nil {
		mvpLine := fmt.Sprintf(`L %s: %s triggered "MVP"`, 
			ts, f.formatPlayerInfo(mvp))
		return baseLine + "\n" + mvpLine
	}
	
	return baseLine
}

// FormatChat formats a chat message event
func (f *LogFormatter) FormatChat(player *models.Player, message string, team, dead bool, timestamp time.Time) string {
	ts := f.formatTimestamp(timestamp)
	
	if player == nil {
		// Server message
		return fmt.Sprintf(`L %s: Server say "%s"`, ts, message)
	}
	
	chatType := "say"
	if team {
		chatType = "say_team"
	}
	if dead {
		chatType = chatType + "_dead"
	}
	
	return fmt.Sprintf(`L %s: %s %s "%s"`, 
		ts, 
		f.formatPlayerInfo(player), 
		chatType, 
		message)
}

// FormatGrenade formats a grenade event
func (f *LogFormatter) FormatGrenade(player *models.Player, grenadeType string, timestamp time.Time) string {
	ts := f.formatTimestamp(timestamp)
	
	return fmt.Sprintf(`L %s: %s threw %s`, 
		ts, 
		f.formatPlayerInfo(player), 
		grenadeType)
}

// FormatFlashbang formats a flashbang event with blinded players
func (f *LogFormatter) FormatFlashbang(player *models.Player, flashed []*models.Player, duration float64, timestamp time.Time) string {
	ts := f.formatTimestamp(timestamp)
	
	lines := []string{
		fmt.Sprintf(`L %s: %s threw flashbang`, ts, f.formatPlayerInfo(player)),
	}
	
	for _, flashedPlayer := range flashed {
		line := fmt.Sprintf(`L %s: %s blinded %s with flashbang for %.1f`, 
			ts, 
			f.formatPlayerInfo(player), 
			f.formatPlayerInfo(flashedPlayer), 
			duration)
		lines = append(lines, line)
	}
	
	return strings.Join(lines, "\n")
}

// FormatServerCommand formats a server command event
func (f *LogFormatter) FormatServerCommand(command, args string, timestamp time.Time) string {
	ts := f.formatTimestamp(timestamp)
	
	return fmt.Sprintf(`L %s: Server cvar "%s" = "%s"`, ts, command, args)
}

// ValidateLogFormat validates that the formatted log follows CS2 standards
func (f *LogFormatter) ValidateLogFormat(logLine string) bool {
	// Basic validation - must start with "L " and have proper timestamp format
	if !strings.HasPrefix(logLine, "L ") {
		return false
	}
	
	// Check minimum length for "L MM/DD/YYYY - HH:MM:SS: "
	if len(logLine) < 25 {
		return false
	}
	
	// Check for colon and space after timestamp (positions 23 and 24)
	if logLine[23] != ':' || logLine[24] != ' ' {
		return false
	}
	
	// Parse timestamp (MM/DD/YYYY - HH:MM:SS)
	timestampPart := logLine[2:23]
	_, err := time.Parse("01/02/2006 - 15:04:05", timestampPart)
	return err == nil
}

// GetFormatterStats returns statistics about the formatter usage
func (f *LogFormatter) GetFormatterStats() map[string]interface{} {
	return map[string]interface{}{
		"server_name":        f.serverName,
		"map_name":           f.mapName,
		"timezone":           f.timeZone.String(),
		"cached_players":     len(f.playerNames),
		"sanitized_names":    f.playerNames,
	}
}