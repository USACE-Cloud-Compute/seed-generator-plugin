package seedgeneratormodel

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
	eventrng := rand.New(rand.NewSource(int64(eventIndex) + m.InitialEventSeed))                    //unique to each event (though offset by 1, which is risky.)
	realrng := rand.New(rand.NewSource(int64(result.RealizationNumber) + m.InitialRealizationSeed)) //unique to each realization and consistent through many events (though offset by 1 which is risky)
	realSeed := realrng.Int63()                                                                     //do not sample random numbers in a range over a map, order matters in an RNG sampling.
	eventSeed := eventrng.Int63()                                                                   //do not sample random numbers in a range over a map, order matters in an RNG sampling.
	for pluginName, ps := range m.PluginInitialSeeds {
		realSeed += ps.RealizationSeed //unique to each plugin
		eventSeed += ps.EventSeed      // unique to each plugin
		outputSeeds[pluginName] = plugin.SeedSet{
			EventSeed:       eventSeed,
			RealizationSeed: realSeed,
		}
	}
	result.Seeds = outputSeeds
	return result, nil
}
