package blockgeneratormodel

import (
	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/stat/distuv"
)

type BlockGenerator struct {
	TargetTotalEvents    int64
	BlocksPerRealization int
	TargetEventsPerBlock int
	Seed                 int64
}
type Block struct {
	RealizationIndex int32 `json:"realization_index" eventstore:"realization_index"`
	BlockIndex       int32 `json:"block_index" eventstore:"block_index"`
	BlockEventCount  int32 `json:"block_event_count" eventstore:"block_event_count"`
	BlockEventStart  int64 `json:"block_event_start" eventstore:"block_event_start"` //inclusive - will be one greater than previous event end
	BlockEventEnd    int64 `json:"block_event_end" eventstore:"block_event_end"`     //inclusive - will be one less than event start if event count is 0.
}

func (b Block) ContainsEvent(eventIndex int) bool {
	if b.BlockEventStart <= int64(eventIndex) {
		if b.BlockEventEnd >= int64(eventIndex) {
			return true
		}
	}
	return false
}
func (bg BlockGenerator) GenerateBlocks() []Block {
	blocks := make([]Block, 0)
	var EventStart int64 = 1
	var EventEnd int64 = 1
	poisson := distuv.Poisson{}
	poisson.Lambda = float64(bg.TargetEventsPerBlock)
	poisson.Src = rand.NewSource(uint64(bg.Seed))
	Index := 1
	Realization := 1
	for {
		events := int(poisson.Rand())
		EventEnd = EventStart + (int64(events) - 1)
		block := Block{BlockIndex: int32(Index), RealizationIndex: int32(Realization), BlockEventCount: int32(events), BlockEventStart: EventStart, BlockEventEnd: EventEnd}
		blocks = append(blocks, block)
		if Index == bg.BlocksPerRealization {
			Realization++
			Index = 0 //always will be adding at the last line of the method.
			if EventEnd >= bg.TargetTotalEvents {
				return blocks
			}
		}
		EventStart = EventEnd + 1
		Index++
	}
}
