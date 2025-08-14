package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// Message represents a WebSocket message
type Message struct {
	Type      string      `json:"type"`
	MatchID   string      `json:"match_id,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

func main() {
	// Connect to WebSocket server
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/api/v1/ws"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer c.Close()

	log.Println("Connected to WebSocket server")

	// Set up signal handling for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Channel for receiving messages
	done := make(chan struct{})

	// Start reading messages
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("Read error: %v", err)
				return
			}

			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				log.Printf("JSON unmarshal error: %v", err)
				continue
			}

			handleMessage(msg)
		}
	}()

	// Subscribe to a test match (you can change this ID)
	testMatchID := "test-match-123"
	subscribeMessage := map[string]interface{}{
		"type":     "subscribe",
		"match_id": testMatchID,
	}

	if err := c.WriteJSON(subscribeMessage); err != nil {
		log.Printf("Failed to send subscribe message: %v", err)
	} else {
		log.Printf("Subscribed to match: %s", testMatchID)
	}

	// Send ping message
	pingMessage := map[string]interface{}{
		"type": "ping",
	}

	if err := c.WriteJSON(pingMessage); err != nil {
		log.Printf("Failed to send ping message: %v", err)
	} else {
		log.Println("Sent ping message")
	}

	// Wait for interrupt signal
	select {
	case <-done:
		return
	case <-interrupt:
		log.Println("Interrupt received, closing connection...")

		// Send close message
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Printf("Write close error: %v", err)
			return
		}

		// Wait for server to close connection or timeout
		select {
		case <-done:
		case <-time.After(time.Second):
		}
	}
}

func handleMessage(msg Message) {
	switch msg.Type {
	case "status":
		log.Printf("📊 Status: %v", msg.Data)

	case "error":
		log.Printf("❌ Error: %v", msg.Data)

	case "event":
		handleEventMessage(msg)

	case "pong":
		log.Printf("🏓 Received pong response")

	default:
		log.Printf("📨 Unknown message type '%s': %v", msg.Type, msg.Data)
	}
}

func handleEventMessage(msg Message) {
	if eventData, ok := msg.Data.(map[string]interface{}); ok {
		eventType, _ := eventData["type"].(string)

		switch eventType {
		case "generation_start":
			log.Printf("🚀 Match generation started: %s", msg.MatchID)
			if teams, ok := eventData["teams"].([]interface{}); ok {
				log.Printf("   Teams: %v", teams)
			}
			if mapName, ok := eventData["map"].(string); ok {
				log.Printf("   Map: %s", mapName)
			}

		case "match_start":
			log.Printf("🏁 Match started: %s", msg.MatchID)

		case "match_progress":
			if data, ok := eventData["data"].(map[string]interface{}); ok {
				if progress, ok := data["progress"].(float64); ok {
					if currentRound, ok := data["current_round"].(float64); ok {
						if totalRounds, ok := data["total_rounds"].(float64); ok {
							log.Printf("⚡ Progress: Round %d/%d (%.1f%%)", 
								int(currentRound), int(totalRounds), progress)
						}
					}
				}
			}

		case "round_start":
			if data, ok := eventData["data"].(map[string]interface{}); ok {
				if roundNum, ok := data["round_number"].(float64); ok {
					log.Printf("🔄 Round %d started", int(roundNum))
				}
			}

		case "round_end":
			if data, ok := eventData["data"].(map[string]interface{}); ok {
				roundNum, _ := data["round_number"].(float64)
				winner, _ := data["winner"].(string)
				reason, _ := data["reason"].(string)
				mvp, _ := data["mvp"].(string)
				log.Printf("✅ Round %d ended: %s won (%s), MVP: %s", 
					int(roundNum), winner, reason, mvp)
			}

		case "player_kill":
			if data, ok := eventData["data"].(map[string]interface{}); ok {
				attacker, _ := data["attacker"].(string)
				victim, _ := data["victim"].(string)
				weapon, _ := data["weapon"].(string)
				headshot, _ := data["headshot"].(bool)
				
				headshotIcon := ""
				if headshot {
					headshotIcon = "💥"
				}
				log.Printf("💀 %s killed %s with %s %s", attacker, victim, weapon, headshotIcon)
			}

		case "bomb_plant":
			if data, ok := eventData["data"].(map[string]interface{}); ok {
				player, _ := data["player"].(string)
				site, _ := data["site"].(string)
				log.Printf("💣 %s planted bomb at site %s", player, site)
			}

		case "bomb_defuse":
			if data, ok := eventData["data"].(map[string]interface{}); ok {
				player, _ := data["player"].(string)
				site, _ := data["site"].(string)
				log.Printf("🛡️ %s defused bomb at site %s", player, site)
			}

		case "bomb_explode":
			if data, ok := eventData["data"].(map[string]interface{}); ok {
				site, _ := data["site"].(string)
				log.Printf("💥 Bomb exploded at site %s", site)
			}

		case "economy_update":
			log.Printf("💰 Economy updated for round")

		case "generation_end":
			if data, ok := eventData["data"].(map[string]interface{}); ok {
				totalRounds, _ := data["total_rounds"].(float64)
				totalEvents, _ := data["total_events"].(float64)
				success, _ := data["success"].(bool)
				
				status := "❌ Failed"
				if success {
					status = "✅ Success"
				}
				log.Printf("🏆 Match generation completed: %s", status)
				log.Printf("   Rounds: %d, Events: %d", int(totalRounds), int(totalEvents))
			}

		case "match_complete":
			log.Printf("🎉 Match completed: %s", msg.MatchID)

		default:
			log.Printf("🎯 Event: %s", eventType)
		}
	}
}