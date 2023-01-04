package seedgeneratormodel_test

import (
	"encoding/json"
	"fmt"
	"os"

	"testing"

	"github.com/usace/seed-generator/seedgeneratormodel"
	"github.com/usace/wat-go/plugin"
)

func TestWriteModel(t *testing.T) {
	path := "../exampledata/eg.json"
	seeds := make(map[string]plugin.SeedSet)
	seeds["fc"] = plugin.SeedSet{
		EventSeed:       234,
		RealizationSeed: 987,
	}
	seeds["pluginB"] = plugin.SeedSet{
		EventSeed:       345,
		RealizationSeed: 876,
	}
	seeds["pluginC"] = plugin.SeedSet{
		EventSeed:       456,
		RealizationSeed: 765,
	}
	model := seedgeneratormodel.Model{
		InitialEventSeed:       1234,
		InitialRealizationSeed: 9876,
		EventsPerRealization:   10,
		PluginInitialSeeds:     seeds,
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
	m := seedgeneratormodel.Model{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		t.Fail()
	}
}
func TestComputeModel(t *testing.T) {
	path := "../exampledata/eg.json"

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fail()
	}
	m := seedgeneratormodel.Model{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		t.Fail()
	}
	r, err := m.Compute(0)
	if err != nil {
		t.Fail()
	}
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
	r, err = m.Compute(14)
	if err != nil {
		t.Fail()
	}
	fmt.Println(r.RealizationNumber)
	outputPath = "../exampledata/result14.json"
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
