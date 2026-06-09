package normalizer

import (
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"

	. "github.com/onsi/gomega"
)

func TestHashSpawn(t *testing.T) {
	t.Run("Files are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		credentials := HashSpawn(normalizer)

		filesAreNotUpdated(t, credentials, ctx)
	})

	t.Run("Input Sets are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		credentials := HashSpawn(normalizer)

		inputSetsAreNotUpdated(t, credentials, ctx)
	})

	t.Run("Updates spawn hash", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		credentials := HashSpawn(normalizer)

		spawn := &types.IncomingSpawn{
			Label:       "//foo/bar:baz",
			Hash:        "",
			InputFiles:  []string{"foo", "bar"},
			OutputFiles: []string{"baz"},
		}

		result := credentials.Run(types.ExecLogEntry{
			Spawn: spawn,
		})
		g.Expect(len(result.Spawn.InputFiles)).To(Equal(2))
		g.Expect(result.Spawn.InputFiles[0]).To(Equal("foo"))
		g.Expect(result.Spawn.InputFiles[1]).To(Equal("bar"))
		g.Expect(len(result.Spawn.OutputFiles)).To(Equal(1))
		g.Expect(result.Spawn.OutputFiles[0]).To(Equal("baz"))
		g.Expect(result.Spawn.Hash).To(Equal("c199f5351ac1a58e2c340d6c64c521dc2cf1684aa8446bd796da7347d124bfe4878dd06cff3c1d5b35cbf6de3f6c6fc8"))
		g.Expect(result.InputSet).To(BeNil())
		g.Expect(result.File).To(BeNil())
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})
}
