package arborist

import (
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/types"
	. "github.com/onsi/gomega"
)

func TestArguments(t *testing.T) {
	t.Run("vertex with 2 spawns with the same arguments returns true", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		arborist := New(normalizer)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			Args:        []string{"foo", "bar", "baz"},
			LogPathHash: "path1",
		})

		vertex.SetSpawn("path2", &types.Spawn{
			Args:        []string{"foo", "bar", "baz"},
			LogPathHash: "path2",
		})

		g.Expect(arborist.argumentsAreIdentical(vertex)).To(BeTrue())
	})

	t.Run("vertex with 2 spawns with no arguments returns true", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		arborist := New(normalizer)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			Args:        []string{},
			LogPathHash: "path1",
		})

		vertex.SetSpawn("path2", &types.Spawn{
			Args:        []string{},
			LogPathHash: "path2",
		})

		g.Expect(arborist.argumentsAreIdentical(vertex)).To(BeTrue())
	})

	t.Run("vertex with 2 spawns with unset arguments returns true", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		arborist := New(normalizer)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			LogPathHash: "path1",
		})

		vertex.SetSpawn("path2", &types.Spawn{
			LogPathHash: "path2",
		})

		g.Expect(arborist.argumentsAreIdentical(vertex)).To(BeTrue())
	})

	t.Run("returns false when second spawn is missing an argument", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		arborist := New(normalizer)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			Args:        []string{"foo", "bar", "baz"},
			LogPathHash: "path1",
		})

		vertex.SetSpawn("path2", &types.Spawn{
			Args:        []string{"foo", "bar"},
			LogPathHash: "path2",
		})

		g.Expect(arborist.argumentsAreIdentical(vertex)).To(BeFalse())
	})

	t.Run("returns false when second spawn has an additional argument (1)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		arborist := New(normalizer)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			Args:        []string{"foo", "bar"},
			LogPathHash: "path1",
		})

		vertex.SetSpawn("path2", &types.Spawn{
			Args:        []string{"foo", "bar", "baz"},
			LogPathHash: "path2",
		})

		g.Expect(arborist.argumentsAreIdentical(vertex)).To(BeFalse())
	})

	t.Run("returns false when second spawn has an additional argument (2)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		arborist := New(normalizer)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			Args:        []string{"foo", "bar"},
			LogPathHash: "path1",
		})

		vertex.SetSpawn("path2", &types.Spawn{
			Args:        []string{"foo", "baz", "bar"},
			LogPathHash: "path2",
		})

		g.Expect(arborist.argumentsAreIdentical(vertex)).To(BeFalse())
	})

	t.Run("returns false when arguments do not match (1)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		arborist := New(normalizer)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			Args:        []string{"foo"},
			LogPathHash: "path1",
		})

		vertex.SetSpawn("path2", &types.Spawn{
			Args:        []string{"bar"},
			LogPathHash: "path2",
		})

		g.Expect(arborist.argumentsAreIdentical(vertex)).To(BeFalse())
	})

	t.Run("returns false when arguments do not match (2)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		arborist := New(normalizer)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			Args:        []string{"foo1", "bar", "baz"},
			LogPathHash: "path1",
		})

		vertex.SetSpawn("path2", &types.Spawn{
			Args:        []string{"foo2", "bar", "baz"},
			LogPathHash: "path2",
		})

		g.Expect(arborist.argumentsAreIdentical(vertex)).To(BeFalse())
	})

	t.Run("returns false when arguments do not match (3)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		arborist := New(normalizer)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			Args:        []string{"foo", "bar1", "baz"},
			LogPathHash: "path1",
		})

		vertex.SetSpawn("path2", &types.Spawn{
			Args:        []string{"foo", "bar2", "baz"},
			LogPathHash: "path2",
		})

		g.Expect(arborist.argumentsAreIdentical(vertex)).To(BeFalse())
	})

	t.Run("returns false when arguments do not match (4)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		arborist := New(normalizer)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			Args:        []string{"foo", "bar", "baz1"},
			LogPathHash: "path1",
		})

		vertex.SetSpawn("path2", &types.Spawn{
			Args:        []string{"foo", "bar", "baz2"},
			LogPathHash: "path2",
		})

		g.Expect(arborist.argumentsAreIdentical(vertex)).To(BeFalse())
	})
}
