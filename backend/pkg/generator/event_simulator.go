package generator

import (
	"math/rand"

	"github.com/noueii/nocs-log-generator/backend/pkg/models"
)

// EventSimulator simulates game events (legacy, now using EventGenerator)
type EventSimulator struct {
	rng           *rand.Rand
	config        *models.MatchConfig
	eventGen      *EventGenerator
}

// NewEventSimulator creates a new event simulator
func NewEventSimulator(rng *rand.Rand, config *models.MatchConfig) *EventSimulator {
	return &EventSimulator{
		rng:      rng,
		config:   config,
		eventGen: NewEventGenerator(rng, config),
	}
}

// GenerateEvents generates events for compatibility (delegates to EventGenerator)
func (es *EventSimulator) GenerateEvents(match *models.Match, state *models.MatchState, roundNum int, strategy *RoundStrategy) ([]models.GameEvent, error) {
	return es.eventGen.GenerateRoundEvents(match, state, roundNum, strategy)
}