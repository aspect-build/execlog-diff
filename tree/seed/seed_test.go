package seed

import (
	"strings"
	"sync"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/types"
	. "github.com/onsi/gomega"
)

func TestSeed(t *testing.T) {
	t.Run("can save, get and delete vertices", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		seed := NewSeed(ctx, normalizer)

		seed.SetVertex("foo1", &types.Vertex{
			Label: "bar1",
		})

		seed.SetVertex("foo2", &types.Vertex{
			Label: "bar2",
		})

		v1 := seed.GetVertex("foo1")
		v2 := seed.GetVertex("foo2")

		g.Expect(v1.Label).To(Equal("bar1"))
		g.Expect(v2.Label).To(Equal("bar2"))

		counter := 0
		seed.spawn_hash_to_vertex_map.Range(func(key, value any) bool {
			counter++
			return true
		})

		g.Expect(counter).To(Equal(2))

		seed.DeleteVertex("foo1")

		v1 = seed.GetVertex("foo1")
		v2 = seed.GetVertex("foo2")

		g.Expect(v1).To(BeNil())
		g.Expect(v2.Label).To(Equal("bar2"))

		counter = 0
		seed.spawn_hash_to_vertex_map.Range(func(key, value any) bool {
			counter++
			return true
		})

		g.Expect(counter).To(Equal(1))
	})

	t.Run("can save, get and delete output file vertices", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		seed := NewSeed(ctx, normalizer)

		seed.SetOutputFileVertex("foo1", "bar1", &types.Vertex{
			Label: "baz1",
		})

		seed.SetOutputFileVertex("foo2", "bar2", &types.Vertex{
			Label: "baz2",
		})

		v1 := seed.GetOutputFileVertex("foo1", "bar1")
		v2 := seed.GetOutputFileVertex("foo2", "bar2")

		g.Expect(v1.Label).To(Equal("baz1"))
		g.Expect(v2.Label).To(Equal("baz2"))

		counter := 0
		seed.log_path_hash_to_output_file_and_vertex_map.Range(func(key, value any) bool {
			counter++
			return true
		})

		g.Expect(counter).To(Equal(2))

		counter = 0
		m, _ := seed.log_path_hash_to_output_file_and_vertex_map.Load("foo1")
		m.(*sync.Map).Range(func(key, value any) bool {
			counter++
			return true
		})

		g.Expect(counter).To(Equal(1))

		counter = 0
		m, _ = seed.log_path_hash_to_output_file_and_vertex_map.Load("foo2")
		m.(*sync.Map).Range(func(key, value any) bool {
			counter++
			return true
		})

		g.Expect(counter).To(Equal(1))

		seed.DeleteOutputFileVertex("foo1", "bar1")

		v1 = seed.GetOutputFileVertex("foo1", "bar1")
		v2 = seed.GetOutputFileVertex("foo2", "bar2")

		g.Expect(v1).To(BeNil())
		g.Expect(v2.Label).To(Equal("baz2"))

		counter = 0
		seed.log_path_hash_to_output_file_and_vertex_map.Range(func(key, value any) bool {
			counter++
			return true
		})

		counter = 0
		m, _ = seed.log_path_hash_to_output_file_and_vertex_map.Load("foo1")
		m.(*sync.Map).Range(func(key, value any) bool {
			counter++
			return true
		})

		g.Expect(counter).To(Equal(0))

		counter = 0
		m, _ = seed.log_path_hash_to_output_file_and_vertex_map.Load("foo2")
		m.(*sync.Map).Range(func(key, value any) bool {
			counter++
			return true
		})

		g.Expect(counter).To(Equal(1))
	})
}
