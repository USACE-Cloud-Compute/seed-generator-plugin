package eventgenerator_test

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/HydrologicEngineeringCenter/go-statistics/statistics"
	"github.com/henrygeorgist/eventgenerator/eventgenerator"
)

func TestWriteModel(t *testing.T) {
	path := "../exampledata/eg.json"
	seeds := make([]eventgenerator.PluginSeed, 3)
	seeds[0] = eventgenerator.PluginSeed{
		Plugin: "pluginA",
		Seed:   234,
	}
	seeds[1] = eventgenerator.PluginSeed{
		Plugin: "pluginB",
		Seed:   345,
	}
	seeds[2] = eventgenerator.PluginSeed{
		Plugin: "pluginC",
		Seed:   456,
	}
	container := statistics.ContinuousDistributionContainer{}
	container.Type = reflect.TypeOf(statistics.TriangularDistribution{}).Name()
	container.Value = statistics.TriangularDistribution{Min: 0, MostLikely: 150, Max: 365}
	model := eventgenerator.Model{
		InitialSeed:                 1234,
		EventsPerRealization:        10,
		Seeds:                       seeds,
		TimeWindowDurationInHours:   100,
		TimeWindowStartDistribution: container,
	}
	b, err := json.Marshal(model)
	if err != nil {
		t.Fail()
	}
	err = os.WriteFile(path, b, 0600)
	if err != nil {
		t.Fail()
	}
}
func TestReadModel(t *testing.T) {
	path := "../exampledata/eg.json"

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fail()
	}
	m := eventgenerator.Model{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		t.Fail()
	}
	val := m.TimeWindowStartDistribution.Value.InvCDF(.25)
	fmt.Printf("value sampled was %v", val)
}
func TestComputeModel(t *testing.T) {
	path := "../exampledata/eg.json"

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fail()
	}
	m := eventgenerator.Model{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		t.Fail()
	}
	r, err := m.Compute(0)
	if err != nil {
		t.Fail()
	}
	fmt.Println(r.EventNumber)
	fmt.Println(r.RealizationNumber)
	outputPath := "../exampledata/result0.json"
	rb, err := json.Marshal(r)
	if err != nil {
		t.Fail()
	}
	err = os.WriteFile(outputPath, rb, 0600)
	if err != nil {
		t.Fail()
	}
	r, err = m.Compute(12)
	if err != nil {
		t.Fail()
	}
	fmt.Println(r.EventNumber)
	fmt.Println(r.RealizationNumber)
	outputPath = "../exampledata/result12.json"
	rb, err = json.Marshal(r)
	if err != nil {
		t.Fail()
	}
	err = os.WriteFile(outputPath, rb, 0600)
	if err != nil {
		t.Fail()
	}
	r, err = m.Compute(32)
	if err != nil {
		t.Fail()
	}
	fmt.Println(r.EventNumber)
	fmt.Println(r.RealizationNumber)
	outputPath = "../exampledata/result32.json"
	rb, err = json.Marshal(r)
	if err != nil {
		t.Fail()
	}
	err = os.WriteFile(outputPath, rb, 0600)
	if err != nil {
		t.Fail()
	}
}
