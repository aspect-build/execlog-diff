package tree

import (
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/converter"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/result/proto"
	"github.com/aspect-build/execlog-diff/tree/arborist"
	"github.com/aspect-build/execlog-diff/tree/chainsaw"
	"github.com/aspect-build/execlog-diff/tree/sawmill"
	"github.com/aspect-build/execlog-diff/tree/seed"
	"github.com/aspect-build/execlog-diff/types"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/onsi/gomega/matchers"
)

func getVertexList(seed *seed.Seed) []types.Vertex {
	vertices := make([]types.Vertex, 0)
	seed.GetHashToVertexMap().Range(func(k, val interface{}) bool {
		v := val.(*types.Vertex)

		vertices = append(vertices, types.Vertex{
			Label: v.Label,
		})

		return true
	})

	return vertices
}

func TestExecLogDiff(t *testing.T) {
	t.Run("a single execlog does not register any changes", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		execLogChan := make(chan types.ExecLogEntry, 100)
		spawnChan := make(chan types.Spawn, 100)
		var compactWait sync.WaitGroup
		compactWait.Add(1)
		go converter.Convert("../test-data/changed_file_7.3.0_1.compact.log", execLogChan, &compactWait, ctx)

		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)
		normalizer := normalizer.New(ctx)
		go normalizer.Normalize(execLogChan, spawnChan, &normalizerWait)

		seed := seed.NewSeed(ctx, normalizer)
		deletionChan := make(chan *types.Vertex, 1)
		var chainsawWait sync.WaitGroup
		chainsawWait.Add(1)

		vertexChan := make(chan *types.Vertex)
		vertexList := make([]*types.Vertex, 0)

		diff_tree := &diff_tree{
			tree: tree{
				arborist:   arborist.New(normalizer),
				seed:       seed,
				ctx:        ctx,
				vertexChan: vertexChan,
			},
			changedVertices: vertexList,
			deletionChan:    deletionChan,
			chainsawWait:    &chainsawWait,
			chainsaw: chainsaw.NewChainsaw(
				ctx,
				normalizer,
				seed,
			),
			sawmill: sawmill.New(normalizer),
		}

		var growWait sync.WaitGroup
		growWait.Add(1)

		go diff_tree.Grow(spawnChan, &growWait)

		compactWait.Wait()
		close(execLogChan)
		normalizerWait.Wait()
		close(spawnChan)

		growWait.Wait()

		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
		g.Expect(len(diff_tree.changedVertices)).To(Equal(0))

		seedVertexList := getVertexList(seed)
		g.Expect(len(seedVertexList)).To(Equal(4))
		g.Expect(seedVertexList).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Label": &matchers.EqualMatcher{Expected: "//foo:foo"},
		})))
		g.Expect(seedVertexList).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Label": &matchers.EqualMatcher{Expected: "//bar:bar"},
		})))
		g.Expect(seedVertexList).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Label": &matchers.EqualMatcher{Expected: "//baz:baz"},
		})))
		g.Expect(seedVertexList).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Label": &matchers.EqualMatcher{Expected: "//:static"},
		})))
	})

	t.Run("identical vertices are removed from seed", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		execLogChan := make(chan types.ExecLogEntry, 100)
		spawnChan := make(chan types.Spawn, 100)
		var compactWait sync.WaitGroup
		compactWait.Add(2)
		go converter.Convert("../test-data/changed_file_7.3.0_1.compact.log", execLogChan, &compactWait, ctx)
		go converter.Convert("../test-data/changed_file_7.3.0_2.compact.log", execLogChan, &compactWait, ctx)

		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)
		normalizer := normalizer.New(ctx)
		go normalizer.Normalize(execLogChan, spawnChan, &normalizerWait)

		seed := seed.NewSeed(ctx, normalizer)
		deletionChan := make(chan *types.Vertex, 1)
		var chainsawWait sync.WaitGroup
		chainsawWait.Add(1)

		vertexChan := make(chan *types.Vertex)
		vertexList := make([]*types.Vertex, 0)

		diff_tree := &diff_tree{
			tree: tree{
				arborist:   arborist.New(normalizer),
				seed:       seed,
				ctx:        ctx,
				vertexChan: vertexChan,
			},
			changedVertices: vertexList,
			deletionChan:    deletionChan,
			chainsawWait:    &chainsawWait,
			chainsaw: chainsaw.NewChainsaw(
				ctx,
				normalizer,
				seed,
			),
			sawmill: sawmill.New(normalizer),
		}

		var growWait sync.WaitGroup
		growWait.Add(1)

		go diff_tree.Grow(spawnChan, &growWait)

		compactWait.Wait()
		close(execLogChan)
		normalizerWait.Wait()
		close(spawnChan)

		growWait.Wait()

		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
		g.Expect(len(diff_tree.changedVertices)).To(Equal(0))

		seedVertexList := getVertexList(seed)
		g.Expect(len(seedVertexList)).To(Equal(2))
		g.Expect(seedVertexList).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Label": &matchers.EqualMatcher{Expected: "//foo:foo"},
		})))
		g.Expect(seedVertexList).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Label": &matchers.EqualMatcher{Expected: "//bar:bar"},
		})))
	})

	t.Run("mill extracts changed vertices from seed and generates proto output", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		execLogChan := make(chan types.ExecLogEntry, 100)
		spawnChan := make(chan types.Spawn, 100)
		var compactWait sync.WaitGroup
		compactWait.Add(2)
		go converter.Convert("../test-data/changed_file_7.3.0_1.compact.log", execLogChan, &compactWait, ctx)
		go converter.Convert("../test-data/changed_file_7.3.0_2.compact.log", execLogChan, &compactWait, ctx)

		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)
		normalizer := normalizer.New(ctx)
		go normalizer.Normalize(execLogChan, spawnChan, &normalizerWait)

		seed := seed.NewSeed(ctx, normalizer)
		deletionChan := make(chan *types.Vertex, 1)
		var chainsawWait sync.WaitGroup
		chainsawWait.Add(1)

		vertexChan := make(chan *types.Vertex)
		vertexList := make([]*types.Vertex, 0)

		diff_tree := &diff_tree{
			tree: tree{
				arborist:   arborist.New(normalizer),
				seed:       seed,
				ctx:        ctx,
				vertexChan: vertexChan,
			},
			changedVertices: vertexList,
			deletionChan:    deletionChan,
			chainsawWait:    &chainsawWait,
			chainsaw: chainsaw.NewChainsaw(
				ctx,
				normalizer,
				seed,
			),
			sawmill: sawmill.New(normalizer),
		}

		var growWait sync.WaitGroup
		growWait.Add(1)

		go diff_tree.Grow(spawnChan, &growWait)

		compactWait.Wait()
		close(execLogChan)
		normalizerWait.Wait()
		close(spawnChan)

		growWait.Wait()

		millResult := diff_tree.CreateResult()

		vertices := make([]types.Vertex, 0)
		for _, vertex := range diff_tree.changedVertices {
			// Unfortunately, the matchers used below do not like pointers
			vertices = append(vertices, types.Vertex{
				Label: vertex.Label,
			})
		}

		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
		g.Expect(len(diff_tree.changedVertices)).To(Equal(2))

		g.Expect(vertices).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Label": &matchers.EqualMatcher{Expected: "//foo:foo"},
		})))
		g.Expect(vertices).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"Label": &matchers.EqualMatcher{Expected: "//bar:bar"},
		})))

		nodes := make([]proto.ActionGraphDiffNode, 0)
		for _, node := range millResult.GetExecDiff().Nodes {
			// Unfortunately, the matchers used below do not like pointers
			nodes = append(nodes, proto.ActionGraphDiffNode{
				ActionHash:  node.ActionHash,
				TargetLabel: node.TargetLabel,
				Children:    node.Children,
			})
		}

		g.Expect(len(nodes)).To(Equal(2))

		g.Expect(nodes).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"ActionHash":  &matchers.EqualMatcher{Expected: "9366b90f37a5eeaed3528b11539487d8a23781e87522601f5bb26f3cd02350916405d622b2400f333ff95931aaf3eaf3"},
			"TargetLabel": &matchers.EqualMatcher{Expected: "//foo:foo"},
			"Children":    &matchers.HaveLenMatcher{Count: 0},
		})))

		g.Expect(nodes).To(ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			"ActionHash":  &matchers.EqualMatcher{Expected: "887da361e121526fc152cbf27796a099b4521636d63984bf849533ab492399237ad227397567501982f4a2af6d435cab"},
			"TargetLabel": &matchers.EqualMatcher{Expected: "//bar:bar"},
			"Children":    &matchers.HaveLenMatcher{Count: 1},
		})))

		for _, node := range millResult.GetExecDiff().Nodes {
			if node.TargetLabel == "//bar:bar" {
				g.Expect(node.Children[0].Reason)

				g.Expect(node.Children[0].Reason).To(Equal(proto.Reason_CHANGED))
				g.Expect(len(node.Children[0].ChangedFiles)).To(Equal(1))
				g.Expect(node.Children[0].ChangedFiles[0].Reason).To(Equal(proto.Reason_CHANGED))
				g.Expect(reflect.TypeOf(node.Children[0].ChangedFiles[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
				g.Expect(node.Children[0].ChangedFiles[0].GetInputFileListChanges().Path).To(Equal("bazel-out/darwin_arm64-fastbuild/bin/foo/foo.txt"))
			}
		}
	})
}
