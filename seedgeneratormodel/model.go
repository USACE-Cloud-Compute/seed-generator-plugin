package seedgeneratormodel

import (
	"errors"
	"math"
	"math/rand"

	"github.com/usace/cc-go-sdk/plugin"
	"github.com/usace/seed-generator/blockgeneratormodel"
)

type RealizationModel struct {
	InitialEventSeed       int64    `json:"initial_event_seed"`
	InitialRealizationSeed int64    `json:"initial_realization_seed"`
	EventsPerRealization   int      `json:"events_per_realization"`
	Plugins                []string `json:"plugins"` //model or plugin name and string
}

func (m RealizationModel) Compute(eventIndex int) (plugin.EventConfiguration, error) {

	realizationNumber := math.Floor(float64(eventIndex) / float64(m.EventsPerRealization))
	//compute seeds
	result := createEventConfiguration(int(realizationNumber), eventIndex, m.InitialEventSeed, m.InitialRealizationSeed, m.Plugins)
	return result, nil
}

type BlockModel struct {
	InitialEventSeed       int64    `json:"initial_event_seed"`
	InitialRealizationSeed int64    `json:"initial_realization_seed"`
	Plugins                []string `json:"plugins"` //model or plugin name and string
}

func (m BlockModel) Compute(eventIndex int, blocks []blockgeneratormodel.Block) (plugin.EventConfiguration, error) {
	for _, b := range blocks {
		if b.ContainsEvent(eventIndex) {
			realizationNumber := b.RealizationIndex
			result := createEventConfiguration(realizationNumber, eventIndex, m.InitialEventSeed, m.InitialRealizationSeed, m.Plugins)
			return result, nil

		}
	}
	return plugin.EventConfiguration{}, errors.New("event index not found in blocks")
}

func createEventConfiguration(realizationNumber int, eventIndex int, initialEventSeed int64, initialRealizationSeed int64, pluginList []string) plugin.EventConfiguration {
	result := plugin.EventConfiguration{}

	result.RealizationNumber = int(realizationNumber) //set realization number

	outputSeeds := make(map[string]plugin.SeedSet)
	eventrng := advance(eventIndex, len(pluginList), rand.New(rand.NewSource(initialEventSeed))) //unique to each event, spinning off for each plugin

	realrng := advance(result.RealizationNumber, len(pluginList), rand.New(rand.NewSource(initialRealizationSeed))) //unique to each realization and consistent through many events spinning off for each plugin once per realization.
	for _, pluginName := range pluginList {                                                                         //compute seeds
		outputSeeds[pluginName] = plugin.SeedSet{
			EventSeed:       eventrng.Int63(), // unique to each plugin
			RealizationSeed: realrng.Int63(),  // unique to each plugin
		}
	}
	result.Seeds = outputSeeds
	return result
}
func advance(count int, seedsPerCount int, rng *rand.Rand) *rand.Rand {
	for i := 0; i < count; i++ {
		for j := 0; j < seedsPerCount; j++ {
			rng.Int63()
		}
	}
	return rng
}
