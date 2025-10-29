package internal

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestGenerateBlocks(t *testing.T) {
	blockGenerator := BlockGenerator{
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
	compareBlocks := []Block{}
	err = json.Unmarshal(b, &compareBlocks)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(blocks, compareBlocks) {
		t.Error("unable to marshall and unmarshall blocks")
	}
}

func TestReadBlocks(t *testing.T) {
	bytes, err := os.ReadFile("../testdata/blocks.json")
	if err != nil {
		t.Error(err)
	}
	var blocks []Block
	err = json.Unmarshal(bytes, &blocks)
	if err != nil {
		t.Error(err)
	}
	expectedLen := 20000
	actualLen := len(blocks)
	if expectedLen != actualLen {
		t.Errorf("expected %d records but found %d", expectedLen, actualLen)
	}
	expectedBlock := Block{
		RealizationIndex: 40,
		BlockIndex:       500,
		BlockEventCount:  5,
		BlockEventStart:  100153,
		BlockEventEnd:    100157,
	}
	actualBlock := blocks[expectedLen-1]
	if !reflect.DeepEqual(expectedBlock, actualBlock) {
		t.Errorf("Expected %v got %v", expectedBlock, actualBlock)
	}
}
