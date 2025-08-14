package generator

import (
	"fmt"
	"strings"

	"github.com/noueii/nocs-log-generator/backend/pkg/models"
)

// LogFormatter handles formatting of game events into CS2 log format
type LogFormatter struct {
	config *models.MatchConfig
}

// NewLogFormatter creates a new log formatter
func NewLogFormatter(config *models.MatchConfig) *LogFormatter {
	return &LogFormatter{
		config: config,
	}
}

// FormatMatch formats a complete match into log format
func (lf *LogFormatter) FormatMatch(match *models.Match) (string, error) {
	var logLines []string
	
	// Add match start header
	logLines = append(logLines, lf.formatMatchStart(match))
	
	// Format all events
	for _, event := range match.Events {
		logLine := event.ToLogLine()
		if logLine != "" {
			logLines = append(logLines, logLine)
		}
	}
	
	// Add match end footer
	logLines = append(logLines, lf.formatMatchEnd(match))
	
	return strings.Join(logLines, "\n"), nil
}

// FormatEvents formats a list of events into log format
func (lf *LogFormatter) FormatEvents(events []models.GameEvent) (string, error) {
	var logLines []string
	
	for _, event := range events {
		logLine := event.ToLogLine()
		if logLine != "" {
			logLines = append(logLines, logLine)
		}
	}
	
	return strings.Join(logLines, "\n"), nil
}

// formatMatchStart creates the match start log entry
func (lf *LogFormatter) formatMatchStart(match *models.Match) string {
	timestamp := match.StartTime.Format(lf.config.TimestampFormat)
	return fmt.Sprintf(`L %s: Log file started (file "logs/match_%s.log") (game "Counter-Strike 2") (version "1.39.8.7")`,
		timestamp, match.ID)
}

// formatMatchEnd creates the match end log entry
func (lf *LogFormatter) formatMatchEnd(match *models.Match) string {
	timestamp := match.EndTime.Format(lf.config.TimestampFormat)
	winningTeam := match.GetWinningTeam()
	
	lines := []string{
		fmt.Sprintf(`L %s: Match ended - %s won %d-%d`, 
			timestamp, winningTeam, match.Scores[winningTeam], 
			match.Scores[lf.getOtherTeam(match, winningTeam)]),
		fmt.Sprintf(`L %s: Log file closed`, timestamp),
	}
	
	return strings.Join(lines, "\n")
}

// getOtherTeam returns the team name that is not the given team
func (lf *LogFormatter) getOtherTeam(match *models.Match, teamName string) string {
	for _, team := range match.Teams {
		if team.Name != teamName {
			return team.Name
		}
	}
	return "Unknown"
}

// FormatEvent formats a single event (helper method)
func (lf *LogFormatter) FormatEvent(event models.GameEvent) string {
	return event.ToLogLine()
}

// FormatRoundData formats round data summary
func (lf *LogFormatter) FormatRoundData(round *models.RoundData) string {
	startTime := round.StartTime.Format(lf.config.TimestampFormat)
	endTime := round.EndTime.Format(lf.config.TimestampFormat)
	duration := round.EndTime.Sub(round.StartTime)
	
	return fmt.Sprintf(`Round %d: %s to %s (%.1fs) - %s won (%s) - MVP: %s`,
		round.RoundNumber, startTime, endTime, duration.Seconds(),
		round.Winner, round.Reason, round.MVP)
}

// GetLogStats returns statistics about the formatted log
func (lf *LogFormatter) GetLogStats(match *models.Match) map[string]interface{} {
	stats := map[string]interface{}{
		"total_events": len(match.Events),
		"total_rounds": len(match.Rounds),
		"match_duration": match.Duration.String(),
		"total_players": len(match.Teams[0].Players) + len(match.Teams[1].Players),
		"format": match.Format,
		"map": match.Map,
	}
	
	// Count event types
	eventTypes := make(map[string]int)
	for _, event := range match.Events {
		eventTypes[event.GetType()]++
	}
	stats["event_breakdown"] = eventTypes
	
	return stats
}