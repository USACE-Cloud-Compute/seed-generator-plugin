package seedgenerator

import (
	"math"
	"math/rand"

	"github.com/usace/cc-go-sdk/plugin"
)

type Model struct {
	InitialEventSeed       int64                     `json:"initial_event_seed"`
	InitialRealizationSeed int64                     `json:"initial_realization_seed"`
	EventsPerRealization   int                       `json:"events_per_realization"`
	PluginInitialSeeds     map[string]plugin.SeedSet `json:"plugin_initial_seeds"` //model or plugin name and string
}

func (m Model) Compute(eventIndex int) (plugin.EventConfiguration, error) {
	result := plugin.EventConfiguration{}
	realizationNumber := math.Floor(float64(eventIndex) / float64(m.EventsPerRealization))
	//set realization number
	result.RealizationNumber = int(realizationNumber)
	//compute seeds
	outputSeeds := make(map[string]plugin.SeedSet)
	rng := rand.New(rand.NewSource(int64(eventIndex)))
	realrng := rand.New(rand.NewSource(int64(result.RealizationNumber) + m.InitialRealizationSeed))
	for pluginName, ps := range m.PluginInitialSeeds {
		realseed := realrng.Int63()
		realseed += ps.RealizationSeed
		eventSeed := m.InitialEventSeed
		eventSeed += ps.EventSeed
		eventSeed += rng.Int63()
		outputSeeds[pluginName] = plugin.SeedSet{
			EventSeed:       eventSeed,
			RealizationSeed: realseed,
		}
	}
	result.Seeds = outputSeeds
	return result, nil
}
