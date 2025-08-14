package formatter

import (
	"testing"
	"time"

	"github.com/noueii/nocs-log-generator/backend/pkg/models"
)

func TestLogFormatter_FormatEvent(t *testing.T) {
	config := &models.MatchConfig{
		Map:        "de_mirage",
		Format:     "mr12",
		StartMoney: 800,
		MaxMoney:   16000,
		ServerName: "Test Server",
	}
	
	formatter := NewLogFormatter(config)
	
	// Create test players
	attacker := &models.Player{
		Name:   "TestPlayer1",
		UserID: 1,
		SteamID: "STEAM_1:0:123456",
		Side:   "CT",
		Stats: models.PlayerStats{
			Kills: 1,
		},
	}
	
	victim := &models.Player{
		Name:   "TestPlayer2", 
		UserID: 2,
		SteamID: "STEAM_1:0:654321",
		Side:   "TERRORIST",
		Stats: models.PlayerStats{
			Deaths: 1,
		},
	}
	
	// Create a kill event
	killEvent := &models.KillEvent{
		BaseEvent: models.BaseEvent{
			Timestamp: time.Now(),
			Type:      "player_death",
			Tick:      12800,
			Round:     1,
		},
		Attacker: attacker,
		Victim:   victim,
		Weapon:   "ak47",
		Headshot: true,
	}
	
	// Test formatting
	formatted := formatter.FormatEvent(killEvent)
	
	if formatted == "" {
		t.Error("FormatEvent returned empty string")
	}
	
	if !formatter.ValidateLogFormat(formatted) {
		t.Errorf("Generated log format is invalid: %s", formatted)
	}
	
	t.Logf("Generated log line: %s", formatted)
}

func TestLogFormatter_FormatPlayerConnect(t *testing.T) {
	config := &models.MatchConfig{
		Map:        "de_mirage",
		ServerName: "Test Server",
	}
	
	formatter := NewLogFormatter(config)
	
	player := &models.Player{
		Name:    "TestPlayer",
		UserID:  1,
		SteamID: "STEAM_1:0:123456",
		Side:    "CT",
	}
	
	formatted := formatter.FormatPlayerConnect(player, "192.168.1.100:27005", time.Now())
	
	if formatted == "" {
		t.Error("FormatPlayerConnect returned empty string")
	}
	
	if !formatter.ValidateLogFormat(formatted) {
		t.Errorf("Generated connection log format is invalid: %s", formatted)
	}
	
	t.Logf("Generated connection log: %s", formatted)
}

func TestLogFormatter_SanitizePlayerName(t *testing.T) {
	config := &models.MatchConfig{
		Map:        "de_mirage",
		ServerName: "Test Server",
	}
	
	formatter := NewLogFormatter(config)
	
	testCases := []struct {
		input    string
		expected string
	}{
		{"NormalName", "NormalName"},
		{"Name\"WithQuotes", "Name\\\"WithQuotes"},
		{"Name\\WithBackslash", "Name\\\\WithBackslash"},
		{"Name\nWithNewline", "Name_WithNewline"},
		{"VeryLongPlayerNameThatExceedsTheLimit", "VeryLongPlayerNameThatExceedsTh"}, // 31 chars max
	}
	
	for _, tc := range testCases {
		result := formatter.sanitizePlayerName(tc.input)
		if result != tc.expected {
			t.Errorf("sanitizePlayerName(%s) = %s, expected %s", tc.input, result, tc.expected)
		}
	}
}

func TestHTTPFormatter_FormatEventAsJSON(t *testing.T) {
	config := &models.MatchConfig{
		Map:        "de_mirage",
		Format:     "mr12",
		StartMoney: 800,
		MaxMoney:   16000,
		ServerName: "Test Server",
	}
	
	httpFormatter := NewHTTPFormatter(config)
	
	// Create test event
	bombPlant := &models.BombPlantEvent{
		BaseEvent: models.BaseEvent{
			Timestamp: time.Now(),
			Type:      "bomb_plant",
			Tick:      25600,
			Round:     3,
		},
		Player: &models.Player{
			Name:    "Planter",
			UserID:  5,
			SteamID: "STEAM_1:0:987654",
			Side:    "TERRORIST",
		},
		Site: "A",
		Position: models.Vector3{X: 500, Y: 500, Z: 0},
	}
	
	jsonBytes, err := httpFormatter.FormatEventAsJSON(bombPlant)
	if err != nil {
		t.Fatalf("FormatEventAsJSON failed: %v", err)
	}
	
	if len(jsonBytes) == 0 {
		t.Error("FormatEventAsJSON returned empty JSON")
	}
	
	t.Logf("Generated JSON: %s", string(jsonBytes))
}

func TestStreamFormatter_Subscribe(t *testing.T) {
	config := &models.MatchConfig{
		Map:        "de_mirage",
		ServerName: "Test Server",
	}
	
	streamConfig := &StreamConfig{
		MaxBufferSize:  100,
		BatchTimeout:   time.Millisecond * 50,
		MaxSubscribers: 10,
		MessageTimeout: time.Second * 5,
	}
	
	streamFormatter := NewStreamFormatter(config, streamConfig)
	defer streamFormatter.Shutdown()
	
	// Test subscription
	filter := &StreamFilter{
		EventTypes: []string{"player_death", "bomb_plant"},
		KillsOnly:  false,
	}
	
	subscriber, err := streamFormatter.Subscribe("test_client", filter, StreamFormatJSON)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}
	
	if subscriber.ID != "test_client" {
		t.Errorf("Expected subscriber ID 'test_client', got %s", subscriber.ID)
	}
	
	if subscriber.Format != StreamFormatJSON {
		t.Errorf("Expected format JSON, got %s", subscriber.Format)
	}
	
	// Test unsubscribe
	err = streamFormatter.Unsubscribe("test_client")
	if err != nil {
		t.Fatalf("Unsubscribe failed: %v", err)
	}
}

func TestStreamFormatter_BroadcastEvent(t *testing.T) {
	config := &models.MatchConfig{
		Map:        "de_mirage",
		ServerName: "Test Server",
	}
	
	streamConfig := &StreamConfig{
		MaxBufferSize:  100,
		BatchTimeout:   time.Millisecond * 50,
		MaxSubscribers: 10,
		MessageTimeout: time.Second * 5,
	}
	
	streamFormatter := NewStreamFormatter(config, streamConfig)
	defer streamFormatter.Shutdown()
	
	// Subscribe to events
	filter := &StreamFilter{
		EventTypes: []string{"round_start"},
	}
	
	subscriber, err := streamFormatter.Subscribe("test_client", filter, StreamFormatJSON)
	if err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}
	
	// Create and broadcast event
	roundStart := &models.RoundStartEvent{
		BaseEvent: models.BaseEvent{
			Timestamp: time.Now(),
			Type:      "round_start",
			Tick:      0,
			Round:     1,
		},
		CTScore:   0,
		TScore:    0,
		CTPlayers: 5,
		TPlayers:  5,
	}
	
	err = streamFormatter.BroadcastEvent(roundStart)
	if err != nil {
		t.Fatalf("BroadcastEvent failed: %v", err)
	}
	
	// Check if message was received (with timeout)
	select {
	case message := <-subscriber.Channel:
		if message.Type != "event" {
			t.Errorf("Expected message type 'event', got %s", message.Type)
		}
		t.Logf("Received message: %+v", message)
	case <-time.After(time.Second):
		t.Error("No message received within timeout")
	}
}

func BenchmarkLogFormatter_FormatEvent(b *testing.B) {
	config := &models.MatchConfig{
		Map:        "de_mirage",
		ServerName: "Benchmark Server",
	}
	
	formatter := NewLogFormatter(config)
	
	killEvent := &models.KillEvent{
		BaseEvent: models.BaseEvent{
			Timestamp: time.Now(),
			Type:      "player_death",
			Tick:      12800,
			Round:     1,
		},
		Attacker: &models.Player{
			Name:    "BenchPlayer1",
			UserID:  1,
			SteamID: "STEAM_1:0:111111",
			Side:    "CT",
		},
		Victim: &models.Player{
			Name:    "BenchPlayer2",
			UserID:  2,
			SteamID: "STEAM_1:0:222222",
			Side:    "TERRORIST",
		},
		Weapon:   "ak47",
		Headshot: true,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = formatter.FormatEvent(killEvent)
	}
}