package blockgeneratormodel

import (
	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/stat/distuv"
)

type BlockGenerator struct {
	TargetTotalEvents    int64
	TargetEventsPerBlock int
	Seed                 int64
}
type Block struct {
	BlockIndex      int   `json:"block_index"`
	BlockEventCount int   `json:"block_event_count"`
	BlockEventStart int64 `json:"block_event_start"`
	BlockEventEnd   int64 `json:"block_event_end"`
}

func (bg BlockGenerator) GenerateBlocks() []Block {
	blocks := make([]Block, 0)
	var EventStart int64 = 1
	var EventEnd int64 = 1
	poisson := distuv.Poisson{}
	poisson.Lambda = float64(bg.TargetEventsPerBlock)
	poisson.Src = rand.NewSource(uint64(bg.Seed))
	Index := 1
	for {
		events := int(poisson.Rand())
		EventEnd += int64(events)
		block := Block{BlockIndex: Index, BlockEventCount: events, BlockEventStart: EventStart, BlockEventEnd: EventEnd}
		blocks = append(blocks, block)
		if EventStart > bg.TargetTotalEvents {
			return blocks
		}
		EventStart += int64(events)
		Index++
	}
}
