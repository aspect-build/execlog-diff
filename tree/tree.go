package tree

import (
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/result/proto"
	"github.com/aspect-build/execlog-diff/tree/arborist"
	"github.com/aspect-build/execlog-diff/tree/seed"
	"github.com/aspect-build/execlog-diff/types"
)

type Tree interface {
	Grow(spawnChan <-chan types.Spawn, waitGroup *sync.WaitGroup)
	CreateResult() *proto.Result
}

type tree struct {
	arborist   *arborist.Arborist
	seed       *seed.Seed
	ctx        *context.Context
	normalizer *normalizer.Normalizer
	vertexChan chan *types.Vertex
}
