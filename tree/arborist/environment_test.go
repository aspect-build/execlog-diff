package arborist

import (
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/types"
	. "github.com/onsi/gomega"
)

func TestEnvironment(t *testing.T) {
	t.Run("vertex with 2 spawns with the same environment returns true", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		arborist := New(normalizer)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			EnvHash:     "somehash",
			LogPathHash: "path1",
		})

		vertex.SetSpawn("path2", &types.Spawn{
			EnvHash:     "somehash",
			LogPathHash: "path2",
		})

		g.Expect(arborist.environmentIsIdentical(vertex)).To(BeTrue())
	})

	t.Run("vertex with 2 spawns with different environment returns false", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		arborist := New(normalizer)

		vertex := &types.Vertex{
			Label: "//foo/bar:baz",
		}

		vertex.SetSpawn("path1", &types.Spawn{
			EnvHash:     "somehash",
			LogPathHash: "path1",
		})

		vertex.SetSpawn("path2", &types.Spawn{
			EnvHash:     "someotherhash",
			LogPathHash: "path2",
		})

		g.Expect(arborist.environmentIsIdentical(vertex)).To(BeFalse())
	})
}
