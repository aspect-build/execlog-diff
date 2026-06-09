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

func TestInputFiles(t *testing.T) {
	t.Run("vertex with 2 spawns with the same single input file returns true", func(t *testing.T) {
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

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			InputFiles:  []string{"foo.txt"},
			LogPathHash: th.GetLog1Path(),
		})

		vertex.SetSpawn("path2", &types.Spawn{
			InputFiles:  []string{"foo.txt"},
			LogPathHash: th.GetLog2Path(),
		})

		g.Expect(arborist.inputFilesAreIdentical(vertex)).To(BeTrue())
	})

	t.Run("vertex with 2 spawns with the same input files returns true", func(t *testing.T) {
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

		th.Log1File("baz.txt", "abcde")
		th.Log2File("baz.txt", "abcde")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog1Path(),
		})

		vertex.SetSpawn("path2", &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog2Path(),
		})

		g.Expect(arborist.inputFilesAreIdentical(vertex)).To(BeTrue())
	})

	t.Run("vertex with 2 spawns with the same input files returns true when a file is not present in the file list", func(t *testing.T) {
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

		th.Log1File("baz.txt", "abcde")
		th.Log2File("baz.txt", "abcde")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog1Path(),
		})

		vertex.SetSpawn("path2", &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog2Path(),
		})
		arborist.inputFilesAreIdentical(vertex)

		g.Expect(arborist.inputFilesAreIdentical(vertex)).To(BeTrue())
	})

	t.Run("vertex with 2 spawns with the same input files returns true when no files are present in the file list", func(t *testing.T) {
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

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: "path1",
		})

		vertex.SetSpawn("path2", &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: "path2",
		})
		arborist.inputFilesAreIdentical(vertex)

		g.Expect(arborist.inputFilesAreIdentical(vertex)).To(BeTrue())
	})

	t.Run("vertex with 2 spawns with different input file hashes returns false", func(t *testing.T) {
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
		th.Log2File("bar.txt", "abcdef")

		th.Log1File("baz.txt", "abcde")
		th.Log2File("baz.txt", "abcde")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog1Path(),
		})

		vertex.SetSpawn("path2", &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog2Path(),
		})

		g.Expect(arborist.inputFilesAreIdentical(vertex)).To(BeFalse())
	})
}

func TestOutputFiles(t *testing.T) {
	t.Run("vertex with 2 spawns with the same single output file returns true", func(t *testing.T) {
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

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			OutputFiles: []string{"foo.txt"},
			LogPathHash: th.GetLog1Path(),
		})

		vertex.SetSpawn("path2", &types.Spawn{
			OutputFiles: []string{"foo.txt"},
			LogPathHash: th.GetLog2Path(),
		})

		g.Expect(arborist.outputFilesAreIdentical(vertex)).To(BeTrue())
	})

	t.Run("vertex with 2 spawns with the same output files returns true", func(t *testing.T) {
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

		th.Log1File("baz.txt", "abcde")
		th.Log2File("baz.txt", "abcde")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			OutputFiles: []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog1Path(),
		})

		vertex.SetSpawn("path2", &types.Spawn{
			OutputFiles: []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog2Path(),
		})

		g.Expect(arborist.outputFilesAreIdentical(vertex)).To(BeTrue())
	})

	t.Run("vertex with 2 spawns with the same output files returns true when a file is not present in the file list", func(t *testing.T) {
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

		th.Log1File("baz.txt", "abcde")
		th.Log2File("baz.txt", "abcde")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			OutputFiles: []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog1Path(),
		})

		vertex.SetSpawn("path2", &types.Spawn{
			OutputFiles: []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog2Path(),
		})
		arborist.outputFilesAreIdentical(vertex)

		g.Expect(arborist.outputFilesAreIdentical(vertex)).To(BeTrue())
	})

	t.Run("vertex with 2 spawns with the same output files returns true when no files are present in the file list", func(t *testing.T) {
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

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			OutputFiles: []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: "path1",
		})

		vertex.SetSpawn("path2", &types.Spawn{
			OutputFiles: []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: "path2",
		})
		arborist.outputFilesAreIdentical(vertex)

		g.Expect(arborist.outputFilesAreIdentical(vertex)).To(BeTrue())
	})

	t.Run("vertex with 2 spawns with different output file hashes returns false", func(t *testing.T) {
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
		th.Log2File("bar.txt", "abcdef")

		th.Log1File("baz.txt", "abcde")
		th.Log2File("baz.txt", "abcde")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			OutputFiles: []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog1Path(),
		})

		vertex.SetSpawn("path2", &types.Spawn{
			OutputFiles: []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog2Path(),
		})

		g.Expect(arborist.outputFilesAreIdentical(vertex)).To(BeFalse())
	})
}
