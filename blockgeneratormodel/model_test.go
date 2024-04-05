package blockgeneratormodel

import (
	"encoding/json"
	"fmt"
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
	fmt.Println(string(b))
}
