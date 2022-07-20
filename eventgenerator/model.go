package eventgenerator

import (
	"math"
	"math/rand"
	"time"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
)

type PluginSeed struct {
	Plugin string `json:"plugin"`
	Seed   int    `json:"seed"`
}
type Model struct {
	InitialSeed                 int                                        `json:"initial_seed"`
	EventsPerRealization        int                                        `json:"events_per_realization"`
	Seeds                       []PluginSeed                               `json:"model_seeds"` //model or plugin name and string
	TimeWindowDurationInHours   int                                        `json:"timewindow_duration"`
	TimeWindowStartDistribution statistics.ContinuousDistributionContainer `json:"timewindow_start_distribution"`
}

type ModelResult struct {
	EventNumber       int          `json:"event_number"`
	RealizationNumber int          `json:"realization_number"`
	Seeds             []PluginSeed `json:"model_seeds"` //model or plugin name and string
	TimeWindowStart   time.Time    `json:"timewindow_start"`
	TimeWindowEnd     time.Time    `json:"timewindow_end"`
}

func (m Model) Compute(eventIndex int) (ModelResult, error) {
	result := ModelResult{}
	eventNumber := math.Mod(float64(eventIndex), float64(m.EventsPerRealization))
	realizationNumber := math.Floor(float64(eventIndex) / float64(m.EventsPerRealization))
	//set event number and realization number
	result.EventNumber = int(math.Floor(eventNumber))
	result.RealizationNumber = int(realizationNumber)
	//compute seeds
	outputSeeds := make([]PluginSeed, len(m.Seeds))
	rng := rand.New(rand.NewSource(int64(eventIndex)))
	for idx, ps := range m.Seeds {
		modelSeed := m.InitialSeed
		modelSeed += ps.Seed
		modelSeed += rng.Int()
		outputSeeds[idx] = PluginSeed{
			Plugin: ps.Plugin,
			Seed:   modelSeed,
		}
	}
	result.Seeds = outputSeeds
	//sample time window start date
	dayOfYear := int(m.TimeWindowStartDistribution.Value.InvCDF(rng.Float64()))
	result.TimeWindowStart = time.Date(1984, time.January, dayOfYear, 0, 0, 0, 0, time.Local)
	days := m.TimeWindowDurationInHours / 24
	hours := m.TimeWindowDurationInHours - (days * 24)
	result.TimeWindowEnd = time.Date(1984, time.January, dayOfYear+days, hours, 0, 0, 0, time.Local)
	return result, nil
}
