package seedgeneratormodel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"testing"
)

func TestWriteRealizationModel(t *testing.T) {
	path := "../exampledata/eg.json"
	seeds := make([]string, 0)
	seeds = append(seeds, "fc")
	seeds = append(seeds, "pluginB")
	seeds = append(seeds, "pluginC")

	model := RealizationModel{
		InitialEventSeed:       1234,
		InitialRealizationSeed: 9876,
		EventsPerRealization:   10,
		Plugins:                seeds,
	}
	b, err := json.Marshal(model)
	if err != nil {
		t.Fail()
	}
	err = ioutil.WriteFile(path, b, 0600)
	if err != nil {
		t.Fail()
	}
}
func TestReadRealizationModel(t *testing.T) {
	path := "../exampledata/eg.json"

	b, err := ioutil.ReadFile(path) //ReadFile(path)
	if err != nil {
		t.Fail()
	}
	m := RealizationModel{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		t.Fail()
	}
}
func TestComputeModel(t *testing.T) {
	path := "../exampledata/eg.json"

	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fail()
	}
	m := RealizationModel{}
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
	err = ioutil.WriteFile(outputPath, rb, 0600)
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
	err = ioutil.WriteFile(outputPath, rb, 0600)
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
	err = ioutil.WriteFile(outputPath, rb, 0600)
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
	err = ioutil.WriteFile(outputPath, rb, 0600)
	if err != nil {
		t.Fail()
	}
}
func Test_GenerateSeedList(t *testing.T) {
	seed := rand.Int63()
	r := rand.New(rand.NewSource(seed))
	for i := 0; i < 100; i++ {
		fmt.Printf("%v,%v\n", i, r.Int63())
	}

}
func Test_Advance(t *testing.T) {
	seed := rand.Int63()
	r := rand.New(rand.NewSource(seed))
	now := time.Now()

	advance(1000000, 25, r)
	since := time.Since(now)
	fmt.Println(since)
}
