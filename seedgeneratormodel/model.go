package seedgeneratormodel

import (
	"errors"
	"math"
	"math/rand"

	"github.com/usace/seed-generator/blockgeneratormodel"
)

// EventConfiguration is a simple structure to support consistency in cc plugins regarding the usage of seeds for natural variability and knowledge uncertainty and realization numbers for indexing
type EventConfiguration struct {
	Seeds map[string]SeedSet `json:"seeds" eventstore:"seeds"`
}

// SeedSet a seed set is a struct to define a natural variability and a knowledge uncertainty
type SeedSet struct {
	EventSeed       int64 `json:"event_seed" eventstore:"event_seed"`
	BlockSeed       int64 `json:"block_seed" eventstore:"block_seed"`
	RealizationSeed int64 `json:"realization_seed" eventstore:"realization_seed"`
}
type RealizationModel struct {
	InitialEventSeed       int64    `json:"initial_event_seed"`
	InitialRealizationSeed int64    `json:"initial_realization_seed"`
	EventsPerRealization   int      `json:"events_per_realization"`
	Plugins                []string `json:"plugins"` //model or plugin name and string
}

func (m RealizationModel) Compute(eventIndex int) (EventConfiguration, error) {

	realizationNumber := math.Floor(float64(eventIndex) / float64(m.EventsPerRealization))
	//compute seeds
	result := createEventConfiguration(int(realizationNumber), 1, eventIndex, m.InitialEventSeed, 1234, m.InitialRealizationSeed, m.Plugins)
	return result, nil
}

type BlockModel struct {
	InitialEventSeed       int64    `json:"initial_event_seed"`
	InitialBlockSeed       int64    `json:"initial_block_seed"`
	InitialRealizationSeed int64    `json:"initial_realization_seed"`
	Plugins                []string `json:"plugins"` //model or plugin name and string
}

func (m BlockModel) Compute(eventIndex int, blocks []blockgeneratormodel.Block) (EventConfiguration, error) {
	blockIndex := 0
	for _, b := range blocks {
		blockIndex++ //blocks are numbered incrimentally within each realization, they must be incrimented across realizations for seed advancement
		if b.ContainsEvent(eventIndex) {
			realizationNumber := b.RealizationIndex
			result := createEventConfiguration(int(realizationNumber), blockIndex, eventIndex, m.InitialEventSeed, m.InitialBlockSeed, m.InitialRealizationSeed, m.Plugins)
			return result, nil

		}
	}
	return EventConfiguration{}, errors.New("event index not found in blocks")
}
func (m BlockModel) ComputeAll(blocks []blockgeneratormodel.Block) ([]EventConfiguration, error) {
	eventrng := rand.New(rand.NewSource(m.InitialEventSeed))
	eventIndex := blocks[0].BlockEventStart
	blockrng := rand.New(rand.NewSource(m.InitialBlockSeed))
	realrng := rand.New(rand.NewSource(m.InitialRealizationSeed))
	var realIndex int32 = 0
	realRandoms := make(map[string]int64)
	configs := []EventConfiguration{}
	for _, b := range blocks {
		blockRandoms := make(map[string]int64)
		for _, pluginName := range m.Plugins {
			blockRandoms[pluginName] = blockrng.Int63()
		}
		if b.RealizationIndex > realIndex { //should happen on the first block of the each realization.
			realIndex = b.RealizationIndex         //update index to avoid problems.
			for _, pluginName := range m.Plugins { //compute seeds
				realRandoms[pluginName] = realrng.Int63() // unique to each plugin
			}
		}
		for eventIndex <= b.BlockEventEnd {
			seeds := make(map[string]SeedSet)
			for _, pluginName := range m.Plugins { //compute seeds
				seeds[pluginName] = SeedSet{
					EventSeed:       eventrng.Int63(), // unique to each plugin
					BlockSeed:       blockRandoms[pluginName],
					RealizationSeed: realRandoms[pluginName], // unique to each plugin
				}
			}
			configs = append(configs, EventConfiguration{Seeds: seeds})
			eventIndex += 1 //update event index
		}
	}

	return configs, nil
}

func createEventConfiguration(realizationNumber int, blockIndex int, eventIndex int, initialEventSeed int64, initialBlockSeed int64, initialRealizationSeed int64, pluginList []string) EventConfiguration {
	result := EventConfiguration{}

	outputSeeds := make(map[string]SeedSet)
	eventrng := advance(eventIndex, len(pluginList), rand.New(rand.NewSource(initialEventSeed))) //unique to each event, spinning off for each plugin
	blockrng := advance(blockIndex, len(pluginList), rand.New(rand.NewSource(initialBlockSeed)))
	realrng := advance(realizationNumber, len(pluginList), rand.New(rand.NewSource(initialRealizationSeed))) //unique to each realization and consistent through many events spinning off for each plugin once per realization.
	for _, pluginName := range pluginList {                                                                  //compute seeds
		outputSeeds[pluginName] = SeedSet{
			EventSeed:       eventrng.Int63(), // unique to each plugin
			BlockSeed:       blockrng.Int63(),
			RealizationSeed: realrng.Int63(), // unique to each plugin
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
