package chainsaw

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/converter"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/tree/seed"
	"github.com/aspect-build/execlog-diff/types"

	. "github.com/onsi/gomega"
)

func vertexListContainsLabel(vertices []*types.Vertex, label string) bool {
	for _, v := range vertices {
		if v.Label == label {
			return true
		}
	}

	return false
}

func getVertex(vertices []*types.Vertex, label string) *types.Vertex {
	for _, v := range vertices {
		if v.Label == label {
			return v
		}
	}

	return nil
}

func getVertexList(seed *seed.Seed) []*types.Vertex {
	vertices := make([]*types.Vertex, 0)
	seed.GetHashToVertexMap().Range(func(k, val interface{}) bool {
		v := val.(*types.Vertex)

		vertices = append(vertices, v)

		return true
	})

	return vertices
}

func TestTestHelper(t *testing.T) {
	t.Run("returns true as expected", func(t *testing.T) {
		g := NewGomegaWithT(t)

		g.Expect(vertexListContainsLabel([]*types.Vertex{
			&types.Vertex{
				Label: "//foo",
			},
		}, "//foo")).To(BeTrue())

		g.Expect(vertexListContainsLabel([]*types.Vertex{
			&types.Vertex{
				Label: "//foo",
			},
			&types.Vertex{
				Label: "//bar",
			},
			&types.Vertex{
				Label: "//baz",
			},
		}, "//foo")).To(BeTrue())

		g.Expect(vertexListContainsLabel([]*types.Vertex{
			&types.Vertex{
				Label: "//foo",
			},
			&types.Vertex{
				Label: "//bar",
			},
			&types.Vertex{
				Label: "//baz",
			},
		}, "//bar")).To(BeTrue())

		g.Expect(vertexListContainsLabel([]*types.Vertex{
			&types.Vertex{
				Label: "//foo",
			},
			&types.Vertex{
				Label: "//bar",
			},
			&types.Vertex{
				Label: "//baz",
			},
		}, "//baz")).To(BeTrue())
	})

	t.Run("returns false as expected", func(t *testing.T) {
		g := NewGomegaWithT(t)

		g.Expect(vertexListContainsLabel([]*types.Vertex{
			&types.Vertex{
				Label: "//foo",
			},
		}, "//bar")).To(BeFalse())

		g.Expect(vertexListContainsLabel([]*types.Vertex{
			&types.Vertex{
				Label: "//foo",
			},
			&types.Vertex{
				Label: "//bar",
			},
			&types.Vertex{
				Label: "//baz",
			},
		}, "//foo/bar/baz")).To(BeFalse())
	})
}

func TestChainsaw(t *testing.T) {
	t.Run("vertices are removed as expected", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		execLogChan := make(chan types.ExecLogEntry, 100)
		spawnChan := make(chan types.Spawn, 100)
		var compactWait sync.WaitGroup
		compactWait.Add(1)
		go converter.Convert("../../test-data/changed_directory_7.3.0_1.compact.log", execLogChan, &compactWait, ctx)

		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)
		normalizer := normalizer.New(ctx)
		go normalizer.Normalize(execLogChan, spawnChan, &normalizerWait)

		seed := seed.NewSeed(ctx, normalizer)

		var vertexWait sync.WaitGroup
		vertexWait.Add(1)
		vertexChan := make(chan *types.Vertex, 10)

		go seed.Germinate(spawnChan, vertexChan, &vertexWait)

		var chainsawWait sync.WaitGroup
		chainsawVertexChan := make(chan *types.Vertex, 10)
		chainsaw := NewChainsaw(ctx, normalizer, seed)

		chainsawWait.Add(1)
		go chainsaw.RemoveVertices(chainsawVertexChan, &chainsawWait)

		compactWait.Wait()
		close(execLogChan)
		normalizerWait.Wait()
		close(spawnChan)
		vertexWait.Wait()
		close(vertexChan)

		// Baseline
		vertices := getVertexList(seed)
		g.Expect(len(vertices)).To(Equal(5))
		g.Expect(vertexListContainsLabel(vertices, "//:static")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//:create_dir")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//foo:foo")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//bar:bar")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//baz:baz")).To(BeTrue())

		// remove
		chainsawVertexChan <- getVertex(vertices, "//:static")

		// Let the deletion process happen async
		time.Sleep(200 * time.Millisecond)

		// Verify
		vertices = getVertexList(seed)
		g.Expect(len(vertices)).To(Equal(4))
		g.Expect(vertexListContainsLabel(vertices, "//:static")).To(BeFalse())
		g.Expect(vertexListContainsLabel(vertices, "//:create_dir")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//foo:foo")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//bar:bar")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//baz:baz")).To(BeTrue())

		// remove
		chainsawVertexChan <- getVertex(vertices, "//:create_dir")
		chainsawVertexChan <- getVertex(vertices, "//foo:foo")
		chainsawVertexChan <- getVertex(vertices, "//baz:baz")

		// Let the deletion process happen async
		time.Sleep(600 * time.Millisecond)

		// Verify
		vertices = getVertexList(seed)
		g.Expect(len(vertices)).To(Equal(1))
		g.Expect(vertexListContainsLabel(vertices, "//:static")).To(BeFalse())
		g.Expect(vertexListContainsLabel(vertices, "//:create_dir")).To(BeFalse())
		g.Expect(vertexListContainsLabel(vertices, "//foo:foo")).To(BeFalse())
		g.Expect(vertexListContainsLabel(vertices, "//bar:bar")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//baz:baz")).To(BeFalse())

		// remove
		chainsawVertexChan <- getVertex(vertices, "//bar:bar")

		// Let the deletion process happen async
		time.Sleep(200 * time.Millisecond)

		// Verify
		vertices = getVertexList(seed)
		g.Expect(len(vertices)).To(Equal(0))

		close(chainsawVertexChan)
		chainsawWait.Wait()

		// Verify post close
		vertices = getVertexList(seed)
		g.Expect(len(vertices)).To(Equal(0))

		g.Expect(stderr.String()).To(Equal(""))
		g.Expect(ctx.Continue()).To(BeTrue())
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})

	t.Run("requesting that a nil vertex is removed is handled", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		execLogChan := make(chan types.ExecLogEntry, 100)
		spawnChan := make(chan types.Spawn, 100)
		var compactWait sync.WaitGroup
		compactWait.Add(1)
		go converter.Convert("../../test-data/changed_directory_7.3.0_1.compact.log", execLogChan, &compactWait, ctx)

		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)
		normalizer := normalizer.New(ctx)
		go normalizer.Normalize(execLogChan, spawnChan, &normalizerWait)

		seed := seed.NewSeed(ctx, normalizer)

		var vertexWait sync.WaitGroup
		vertexWait.Add(1)
		vertexChan := make(chan *types.Vertex, 10)

		go seed.Germinate(spawnChan, vertexChan, &vertexWait)

		var chainsawWait sync.WaitGroup
		chainsawVertexChan := make(chan *types.Vertex, 10)
		chainsaw := NewChainsaw(ctx, normalizer, seed)

		chainsawWait.Add(1)
		go chainsaw.RemoveVertices(chainsawVertexChan, &chainsawWait)

		compactWait.Wait()
		close(execLogChan)
		normalizerWait.Wait()
		close(spawnChan)
		vertexWait.Wait()
		close(vertexChan)

		// Baseline
		vertices := getVertexList(seed)
		g.Expect(len(vertices)).To(Equal(5))
		g.Expect(vertexListContainsLabel(vertices, "//:static")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//:create_dir")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//foo:foo")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//bar:bar")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//baz:baz")).To(BeTrue())

		// remove
		chainsawVertexChan <- nil

		// Let the deletion process happen async
		time.Sleep(200 * time.Millisecond)

		// Verify
		vertices = getVertexList(seed)
		g.Expect(len(vertices)).To(Equal(5))
		g.Expect(vertexListContainsLabel(vertices, "//:static")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//:create_dir")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//foo:foo")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//bar:bar")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//baz:baz")).To(BeTrue())

		close(chainsawVertexChan)
		chainsawWait.Wait()

		// Verify post close
		vertices = getVertexList(seed)
		g.Expect(len(vertices)).To(Equal(5))
		g.Expect(vertexListContainsLabel(vertices, "//:static")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//:create_dir")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//foo:foo")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//bar:bar")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//baz:baz")).To(BeTrue())

		g.Expect(stderr.String()).To(Equal(""))
		g.Expect(ctx.Continue()).To(BeTrue())
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})

	t.Run("requesting that a non existant vertex be removed does not fail", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		execLogChan := make(chan types.ExecLogEntry, 100)
		spawnChan := make(chan types.Spawn, 100)
		var compactWait sync.WaitGroup
		compactWait.Add(1)
		go converter.Convert("../../test-data/changed_directory_7.3.0_1.compact.log", execLogChan, &compactWait, ctx)

		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)
		normalizer := normalizer.New(ctx)
		go normalizer.Normalize(execLogChan, spawnChan, &normalizerWait)

		seed := seed.NewSeed(ctx, normalizer)

		var vertexWait sync.WaitGroup
		vertexWait.Add(1)
		vertexChan := make(chan *types.Vertex, 10)

		go seed.Germinate(spawnChan, vertexChan, &vertexWait)

		var chainsawWait sync.WaitGroup
		chainsawVertexChan := make(chan *types.Vertex, 10)
		chainsaw := NewChainsaw(ctx, normalizer, seed)

		chainsawWait.Add(1)
		go chainsaw.RemoveVertices(chainsawVertexChan, &chainsawWait)

		compactWait.Wait()
		close(execLogChan)
		normalizerWait.Wait()
		close(spawnChan)
		vertexWait.Wait()
		close(vertexChan)

		// remove
		vertices := getVertexList(seed)
		vertexToRemove := getVertex(vertices, "//:static")
		chainsawVertexChan <- vertexToRemove

		// Let the deletion process happen async
		time.Sleep(200 * time.Millisecond)

		// Verify
		vertices = getVertexList(seed)
		g.Expect(len(vertices)).To(Equal(4))
		g.Expect(vertexListContainsLabel(vertices, "//:static")).To(BeFalse())
		g.Expect(vertexListContainsLabel(vertices, "//:create_dir")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//foo:foo")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//bar:bar")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//baz:baz")).To(BeTrue())

		// remove same vertex again
		chainsawVertexChan <- vertexToRemove

		// Let the deletion process happen async
		time.Sleep(200 * time.Millisecond)

		// Verify
		vertices = getVertexList(seed)
		g.Expect(len(vertices)).To(Equal(4))
		g.Expect(vertexListContainsLabel(vertices, "//:static")).To(BeFalse())
		g.Expect(vertexListContainsLabel(vertices, "//:create_dir")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//foo:foo")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//bar:bar")).To(BeTrue())
		g.Expect(vertexListContainsLabel(vertices, "//baz:baz")).To(BeTrue())

		close(chainsawVertexChan)
		chainsawWait.Wait()

		g.Expect(stderr.String()).To(Equal(""))
		g.Expect(ctx.Continue()).To(BeTrue())
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})
}
