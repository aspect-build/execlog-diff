package sawmill

import (
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/result/proto"
	"github.com/aspect-build/execlog-diff/types"
)

type SawMill struct {
	normalizer *normalizer.Normalizer
}

type differ interface {
	Run(spawn1 *types.Spawn, spawn2 *types.Spawn) []*proto.InputItem
}

func New(normalizer *normalizer.Normalizer) *SawMill {
	return &SawMill{
		normalizer: normalizer,
	}
}

func (sawmill *SawMill) DiffSpawns(spawn1 *types.Spawn, spawn2 *types.Spawn) *proto.ActionGraphDiffNode {
	differs := []differ{
		// TODO: Add Environment
		Arguments(sawmill),
		InputFiles(sawmill),
		OutputFiles(sawmill),
	}

	changes := make([]*proto.InputItem, 0)
	for _, differ := range differs {
		changes = append(changes, differ.Run(spawn1, spawn2)...)
	}

	return &proto.ActionGraphDiffNode{
		Reason:         proto.Reason_CHANGED,
		ActionHash:     spawn1.Hash,
		TargetLabel:    spawn1.Label,
		Changes:        changes,
		RemoteCachable: spawn1.RemoteCachable,
	}
}
