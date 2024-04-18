package blockgeneratormodel_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/usace/seed-generator/blockgeneratormodel"
)

func ExampleBlockGenerator() {
	blockGenerator := blockgeneratormodel.BlockGenerator{
		BlocksPerRealization: 500,
		TargetTotalEvents:    5000,
		TargetEventsPerBlock: 5,
		Seed:                 1234,
	}
	blocks := blockGenerator.GenerateBlocks()
	for idx, block := range blocks {
		if block.ContainsEvent(1) {
			fmt.Println(idx)
		}
	}
	// Output: 0
}
func TestGenerateBlocks(t *testing.T) {
	blockGenerator := blockgeneratormodel.BlockGenerator{
		BlocksPerRealization: 500,
		TargetTotalEvents:    5000,
		TargetEventsPerBlock: 5,
		Seed:                 1234,
	}
	blocks := blockGenerator.GenerateBlocks()
	b, err := json.Marshal(blocks)
	if err != nil {
		t.Fail()
	}
	fmt.Println(string(b))
}
func TestReadBlocks(t *testing.T) {
	bytes, err := os.ReadFile("/workspaces/seedgenerator/exampledata/blocks.json")
	if err != nil {
		t.Fail()
		fmt.Println(err)
	}
	var blocks []blockgeneratormodel.Block
	err = json.Unmarshal(bytes, &blocks)
	if err != nil {
		t.Fail()
		fmt.Println(err)
	}
	fmt.Println(blocks)
}
func TestReaderBlocks(t *testing.T) {
	file, err := os.Open("/workspaces/seedgenerator/exampledata/blocks.json")

	if err != nil {
		t.Fail()
		fmt.Println(err)
	}
	var blocks []blockgeneratormodel.Block
	err = json.NewDecoder(file).Decode(&blocks)
	if err != nil {
		t.Fail()
		fmt.Println(err)
	}
	fmt.Println(blocks)
}
