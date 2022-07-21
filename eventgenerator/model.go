package eventgenerator

import (
	"math"
	"math/rand"
)

type PluginSeeds struct {
	Plugin          string `json:"plugin"`
	EventSeed       int64  `json:"event_seed"`
	RealizationSeed int64  `json:"realization_seed"`
}
type Model struct {
	InitialEventSeed       int64         `json:"initial_event_seed"`
	InitialRealizationSeed int64         `json:"initial_realization_seed"`
	EventsPerRealization   int           `json:"events_per_realization"`
	Seeds                  []PluginSeeds `json:"model_seeds"` //model or plugin name and string
}

type EventConfiguration struct {
	EventNumber       int           `json:"event_number"`
	RealizationNumber int           `json:"realization_number"`
	Seeds             []PluginSeeds `json:"model_seeds"` //model or plugin name and string
}

func (m Model) Compute(eventIndex int) (EventConfiguration, error) {
	result := EventConfiguration{}
	eventNumber := math.Mod(float64(eventIndex), float64(m.EventsPerRealization))
	realizationNumber := math.Floor(float64(eventIndex) / float64(m.EventsPerRealization))
	//set event number and realization number
	result.EventNumber = int(math.Floor(eventNumber))
	result.RealizationNumber = int(realizationNumber)
	//compute seeds
	outputSeeds := make([]PluginSeeds, len(m.Seeds))
	rng := rand.New(rand.NewSource(int64(eventIndex)))
	realrng := rand.New(rand.NewSource(int64(result.RealizationNumber) + m.InitialRealizationSeed))
	for idx, ps := range m.Seeds {
		realseed := realrng.Int63()
		realseed += ps.RealizationSeed
		eventSeed := m.InitialEventSeed
		eventSeed += ps.EventSeed
		eventSeed += rng.Int63()
		outputSeeds[idx] = PluginSeeds{
			Plugin:          ps.Plugin,
			EventSeed:       eventSeed,
			RealizationSeed: realseed,
		}
	}
	result.Seeds = outputSeeds
	return result, nil
}
