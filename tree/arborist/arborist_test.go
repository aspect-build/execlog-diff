package arborist

import (
	"strings"
	"sync"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/test-data/helpers"
	"github.com/aspect-build/execlog-diff/types"
	. "github.com/onsi/gomega"
)

func TestArboristFiles(t *testing.T) {
	t.Run("identical vertex should return true", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		arborist := New(normalizer)

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log1File("foo.txt", "abc")
		th.Log2File("foo.txt", "abc")
		th.Log1File("bar.txt", "abcd")
		th.Log2File("bar.txt", "abcd")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			Args:        []string{"foo", "bar", "baz"},
			EnvHash:     "somehash",
			InputFiles:  []string{"foo.txt"},
			OutputFiles: []string{"bar.txt"},
			LogPathHash: th.GetLog1Path(),
		})

		vertex.SetSpawn("path2", &types.Spawn{
			Args:        []string{"foo", "bar", "baz"},
			EnvHash:     "somehash",
			InputFiles:  []string{"foo.txt"},
			OutputFiles: []string{"bar.txt"},
			LogPathHash: th.GetLog2Path(),
		})

		g.Expect(arborist.VertexIsIdentical(vertex)).To(BeTrue())
	})

	t.Run("returns false if vertex has only 1 spawn", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		arborist := New(normalizer)

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log1File("foo.txt", "abc")
		th.Log1File("bar.txt", "abcd")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			Args:        []string{"foo", "bar", "baz"},
			EnvHash:     "somehash",
			InputFiles:  []string{"foo.txt"},
			OutputFiles: []string{"bar.txt"},
			LogPathHash: th.GetLog1Path(),
		})

		g.Expect(arborist.VertexIsIdentical(vertex)).To(BeFalse())
	})

	t.Run("returns false if vertex has 3 spawns", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		arborist := New(normalizer)

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log1File("foo.txt", "abc")
		th.Log1File("bar.txt", "abcd")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			Args:        []string{"foo", "bar", "baz"},
			EnvHash:     "somehash",
			InputFiles:  []string{"foo.txt"},
			OutputFiles: []string{"bar.txt"},
			LogPathHash: th.GetLog1Path(),
		})

		vertex.SetSpawn("path2", &types.Spawn{
			Args:        []string{"foo", "bar", "baz"},
			EnvHash:     "somehash",
			InputFiles:  []string{"foo.txt"},
			OutputFiles: []string{"bar.txt"},
			LogPathHash: th.GetLog2Path(),
		})

		vertex.SetSpawn("path3", &types.Spawn{
			Args:        []string{"foo", "bar", "baz"},
			EnvHash:     "somehash",
			InputFiles:  []string{"foo.txt"},
			OutputFiles: []string{"bar.txt"},
			LogPathHash: "path3",
		})

		g.Expect(arborist.VertexIsIdentical(vertex)).To(BeFalse())
	})
}
