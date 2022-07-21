package eventgenerator

import (
	"math"
	"math/rand"

	"github.com/usace/wat-go-sdk/plugin"
)

type Model struct {
	InitialEventSeed       int64            `json:"initial_event_seed"`
	InitialRealizationSeed int64            `json:"initial_realization_seed"`
	EventsPerRealization   int              `json:"events_per_realization"`
	Seeds                  []plugin.SeedSet `json:"model_seeds"` //model or plugin name and string
}

func (m Model) Compute(eventIndex int) (plugin.EventConfiguration, error) {
	result := plugin.EventConfiguration{}
	eventNumber := math.Mod(float64(eventIndex), float64(m.EventsPerRealization))
	realizationNumber := math.Floor(float64(eventIndex) / float64(m.EventsPerRealization))
	//set event number and realization number
	result.EventNumber = int(math.Floor(eventNumber))
	result.RealizationNumber = int(realizationNumber)
	//compute seeds
	outputSeeds := make([]plugin.SeedSet, len(m.Seeds))
	rng := rand.New(rand.NewSource(int64(eventIndex)))
	realrng := rand.New(rand.NewSource(int64(result.RealizationNumber) + m.InitialRealizationSeed))
	for idx, ps := range m.Seeds {
		realseed := realrng.Int63()
		realseed += ps.RealizationSeed
		eventSeed := m.InitialEventSeed
		eventSeed += ps.EventSeed
		eventSeed += rng.Int63()
		outputSeeds[idx] = plugin.SeedSet{
			Identifier:      ps.Identifier,
			EventSeed:       eventSeed,
			RealizationSeed: realseed,
		}
	}
	result.Seeds = outputSeeds
	return result, nil
}
