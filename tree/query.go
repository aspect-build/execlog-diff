package tree

import (
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/result/proto"
	"github.com/aspect-build/execlog-diff/tree/arborist"
	"github.com/aspect-build/execlog-diff/tree/chainsaw"
	"github.com/aspect-build/execlog-diff/tree/sawmill"
	"github.com/aspect-build/execlog-diff/tree/seed"
	"github.com/aspect-build/execlog-diff/types"
)

func NewQueryTree(ctx *context.Context, normalizer *normalizer.Normalizer) Tree {
	seed := seed.NewSeed(ctx, normalizer)
	deletionChan := make(chan *types.Vertex, 1)
	var chainsawWait sync.WaitGroup
	chainsawWait.Add(1)

	return &query_tree{
		tree: tree{
			arborist:   arborist.New(normalizer),
			seed:       seed,
			ctx:        ctx,
			normalizer: normalizer,
			vertexChan: make(chan *types.Vertex),
		},
		vertices:     make([]*types.Vertex, 0),
		deletionChan: deletionChan,
		chainsawWait: &chainsawWait,
		chainsaw: chainsaw.NewChainsaw(
			ctx,
			normalizer,
			seed,
		),
		sawmill: sawmill.New(normalizer),
	}
}

type query_tree struct {
	tree
	deletionChan chan *types.Vertex
	chainsaw     *chainsaw.Chainsaw
	chainsawWait *sync.WaitGroup
	vertices     []*types.Vertex
	sawmill      *sawmill.SawMill
}

func (tree *query_tree) Grow(spawnChan <-chan types.Spawn, waitGroup *sync.WaitGroup) {
	var germinateWait sync.WaitGroup
	germinateWait.Add(1)
	go tree.seed.Germinate(spawnChan, tree.vertexChan, &germinateWait)

	var collectorWait sync.WaitGroup
	collectorWait.Add(1)
	go func() {
		for vertex := range tree.vertexChan {
			tree.vertices = append(tree.vertices, vertex)
		}
		collectorWait.Done()
	}()

	germinateWait.Wait()
	close(tree.vertexChan)
	collectorWait.Wait()

	waitGroup.Done()
}

func (tree *query_tree) CreateResult() *proto.Result {
	switch tree.ctx.Mode {
	case context.ListLabels:
		return tree.listLabels()
	}

	return &proto.Result{}
}

func (tree *query_tree) listLabels() *proto.Result {
	labelMap := make(map[string]*proto.Label)

	for _, vertex := range tree.vertices {
		spawn := tree.getSpawn(vertex)

		if _, exists := labelMap[spawn.Label]; !exists {
			labelMap[spawn.Label] = &proto.Label{
				TargetLabel:     spawn.Label,
				CacheHit:        spawn.CacheHit,
				Remoteable:      spawn.Remoteable,
				Cacheable:       spawn.Cacheable,
				RemoteCacheable: spawn.RemoteCachable,
			}
		}
	}

	pbLabels := make([]*proto.Label, 0, len(labelMap))
	for _, label := range labelMap {
		pbLabels = append(pbLabels, label)
	}

	return &proto.Result{
		Result: &proto.Result_LabelQuery{
			LabelQuery: &proto.ListLabelsQuery{
				Labels: pbLabels,
			},
		},
	}
}

func (tree *query_tree) getSpawn(vertex *types.Vertex) *types.Spawn {
	spawns := []*types.Spawn{}
	for _, logPathHash := range vertex.GetLogPathHashes() {
		spawn := vertex.GetSpawn(logPathHash)

		spawns = append(spawns, spawn)
	}
	return spawns[0]
}
