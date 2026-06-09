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

func NewDiffTree(ctx *context.Context, normalizer *normalizer.Normalizer) Tree {
	seed := seed.NewSeed(ctx, normalizer)
	deletionChan := make(chan *types.Vertex, 1)
	var chainsawWait sync.WaitGroup
	chainsawWait.Add(1)

	return &diff_tree{
		tree: tree{
			arborist:   arborist.New(normalizer),
			seed:       seed,
			ctx:        ctx,
			normalizer: normalizer,
			vertexChan: make(chan *types.Vertex),
		},
		changedVertices: make([]*types.Vertex, 0),
		deletionChan:    deletionChan,
		chainsawWait:    &chainsawWait,
		chainsaw: chainsaw.NewChainsaw(
			ctx,
			normalizer,
			seed,
		),
		sawmill: sawmill.New(normalizer),
	}
}

type diff_tree struct {
	tree
	deletionChan    chan *types.Vertex
	chainsaw        *chainsaw.Chainsaw
	chainsawWait    *sync.WaitGroup
	changedVertices []*types.Vertex
	sawmill         *sawmill.SawMill
}

func (tree *diff_tree) Grow(spawnChan <-chan types.Spawn, waitGroup *sync.WaitGroup) {
	var germinateWait sync.WaitGroup
	germinateWait.Add(1)
	go tree.seed.Germinate(spawnChan, tree.vertexChan, &germinateWait)

	var trimBranchesWait sync.WaitGroup
	trimBranchesWait.Add(1)
	go tree.removeSameVertices(&trimBranchesWait)

	germinateWait.Wait()
	close(tree.vertexChan)

	trimBranchesWait.Wait()

	waitGroup.Done()
}

func (tree *diff_tree) removeSameVertices(waitGroup *sync.WaitGroup) {
	go tree.chainsaw.RemoveVertices(tree.deletionChan, tree.chainsawWait)
	for vertex := range tree.vertexChan {
		if !tree.ctx.Continue() {
			break
		}

		vertex.Mutex.Lock()

		logPathHashes := vertex.GetLogPathHashes()

		if len(logPathHashes) == 2 {
			vertex.SamenessChecked = true
			if tree.arborist.VertexIsIdentical(vertex) {
				vertex.Same = true
				tree.deletionChan <- vertex
			} else {
				vertex.Same = false
			}
		}

		vertex.Mutex.Unlock()
	}

	close(tree.deletionChan)
	tree.chainsawWait.Wait()

	waitGroup.Done()
}

func (tree *diff_tree) CreateResult() *proto.Result {
	tree.seed.GetHashToVertexMap().Range(func(k, val interface{}) bool {
		v := val.(*types.Vertex)

		tree.changedVertices = append(tree.changedVertices, v)

		return true
	})

	execDiff := &proto.Result_ExecDiff{
		ExecDiff: &proto.ExecDiff{
			Nodes:     make([]*proto.ActionGraphDiffNode, 0),
			RootNodes: make([]string, 0),
		},
	}
	for _, vertex := range tree.changedVertices {
		if len(vertex.Spawns) == 2 {
			spawns := make([]*types.Spawn, 0, 2)
			for _, spawn := range vertex.Spawns {
				spawns = append(spawns, spawn)
			}
			execDiff.ExecDiff.Nodes = append(execDiff.ExecDiff.Nodes, tree.sawmill.DiffSpawns(spawns[0], spawns[1]))
		}
	}

	isRootNodesMap := make(map[string]bool)
	for _, node := range execDiff.ExecDiff.Nodes {
		isRootNodesMap[node.ActionHash] = true
	}

	for _, node := range execDiff.ExecDiff.Nodes {
		children := make(map[string]*proto.ActionGraphDiffEdge, 0)
		for _, change := range node.Changes {
			switch item := change.Item.(type) {
			// TODO: add added & removed
			case *proto.InputItem_InputFileListChanges:
				for _, outputNode := range execDiff.ExecDiff.Nodes {
					for _, outputChange := range outputNode.Changes {
						switch outputItem := outputChange.Item.(type) {
						case *proto.InputItem_OutputFileListChanges:
							if outputItem.OutputFileListChanges.Path == item.InputFileListChanges.Path {
								edge := children[outputNode.ActionHash]
								if edge == nil {
									edge = &proto.ActionGraphDiffEdge{
										Child:  outputNode.ActionHash,
										Reason: proto.Reason_CHANGED,
									}
								}

								edge.ChangedFiles = append(edge.ChangedFiles, change)
								children[outputNode.ActionHash] = edge
							}
						}
					}
				}
			}
		}

		childrenList := make([]*proto.ActionGraphDiffEdge, 0)
		for key, value := range children {
			childrenList = append(childrenList, value)
			isRootNodesMap[key] = false
		}
		node.Children = childrenList
	}

	for target, isRootNode := range isRootNodesMap {
		if isRootNode {
			execDiff.ExecDiff.RootNodes = append(execDiff.ExecDiff.RootNodes, target)
		}
	}

	return &proto.Result{
		Result: execDiff,
	}
}
