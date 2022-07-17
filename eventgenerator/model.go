package eventgenerator

import (
	"math"
	"time"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
	"github.com/usace/wat-go-sdk/plugin"
)

type Model struct {
	InitialSeed                 int                               `json:"initial_seed"`
	EventsPerRealization        int                               `json:"events_per_realization"`
	Seeds                       map[string]int                    `json:"model_seeds"` //model or plugin name and string
	TimeWindowDurationInHours   int                               `json:"timewindow_duration"`
	TimeWindowStartDistribution statistics.ContinuousDistribution `json:"timewindow_start_distribution"`
}
type ModelResult struct {
	EventNumber       int            `json:"event_number"`
	RealizationNumber int            `json:"realization_number"`
	Seeds             map[string]int `json:"model_seeds"` //model or plugin name and string
	TimeWindowStart   time.Time      `json:"timewindow_start"`
	TimeWindowEnd     time.Time      `json:"timewindow_end"`
}

func Init(input plugin.ResourceInfo) (Model, error) {
	return Model{}, nil
}
func (m Model) Compute(eventIndex int) (ModelResult, error) {
	result := ModelResult{}
	eventNumber := math.Mod(float64(eventIndex), float64(m.EventsPerRealization))
	realizationNumber := math.Floor(float64(eventIndex) / float64(m.EventsPerRealization))

	result.EventNumber = int(math.Floor(eventNumber))
	result.RealizationNumber = int(realizationNumber)
	//compute seeds

	return result, nil
}
