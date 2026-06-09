package seed

import (
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/test-data/helpers"
	"github.com/aspect-build/execlog-diff/types"
	"github.com/aspect-build/execlog-diff/diagnostic"
	. "github.com/onsi/gomega"
	"gotest.tools/v3/golden"
)

func TestGerminate(t *testing.T) {
	t.Run("stops when context says to", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		seed := NewSeed(ctx, normalizer)

		var wait sync.WaitGroup
		wait.Add(1)
		spawnChan := make(chan types.Spawn, 10)
		vertexChan := make(chan *types.Vertex, 10)

		go seed.Germinate(spawnChan, vertexChan, &wait)

		spawnChan <- types.Spawn{
			Hash:  "foo",
			Label: "//foo",
		}

		spawnChan <- types.Spawn{
			Hash:  "bar",
			Label: "//bar",
		}

		time.Sleep(time.Millisecond * 200)
		ctx.Error(&diagnostic.Diagnostic{Message: "Stop!"})

		spawnChan <- types.Spawn{
			Hash:  "baz",
			Label: "//baz",
		}

		close(spawnChan)
		wait.Wait()
		close(vertexChan)

		g.Expect(len(vertexChan)).To(Equal(2))
	})

	t.Run("spawns from spawn chan are converted and put onto the vertex chan", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		seed := NewSeed(ctx, normalizer)

		var wait sync.WaitGroup
		wait.Add(1)
		spawnChan := make(chan types.Spawn, 10)
		vertexChan := make(chan *types.Vertex, 10)

		go seed.Germinate(spawnChan, vertexChan, &wait)

		spawnChan <- types.Spawn{
			Hash:  "foo",
			Label: "//foo",
		}

		spawnChan <- types.Spawn{
			Hash:  "bar",
			Label: "//bar",
		}

		close(spawnChan)
		wait.Wait()
		close(vertexChan)

		g.Expect(len(vertexChan)).To(Equal(2))

		vertexList := helpers.CreateListFromVertexChan(vertexChan, t)

		g.Expect(len(vertexList)).To(Equal(2))
		g.Expect(vertexList[0].Label).To(Equal("//foo"))
		g.Expect(vertexList[1].Label).To(Equal("//bar"))
	})

	t.Run("vertices are saved", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		seed := NewSeed(ctx, normalizer)

		var wait sync.WaitGroup
		wait.Add(1)
		spawnChan := make(chan types.Spawn, 10)
		vertexChan := make(chan *types.Vertex, 10)

		go seed.Germinate(spawnChan, vertexChan, &wait)

		spawnChan <- types.Spawn{
			Hash:        "foo",
			Label:       "//foo",
			LogPathHash: "1",
			EnvHash:     "someEnvHash1",
		}

		spawnChan <- types.Spawn{
			Hash:        "bar",
			Label:       "//bar",
			LogPathHash: "1",
			EnvHash:     "someEnvHash2",
		}

		close(spawnChan)
		wait.Wait()
		close(vertexChan)

		v1 := seed.GetVertex("foo")
		v2 := seed.GetVertex("bar")

		g.Expect(v1.Label).To(Equal("//foo"))
		g.Expect(v2.Label).To(Equal("//bar"))
	})

	t.Run("vertices are saved to the output file map", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		seed := NewSeed(ctx, normalizer)

		var wait sync.WaitGroup
		wait.Add(1)
		spawnChan := make(chan types.Spawn, 10)
		vertexChan := make(chan *types.Vertex, 10)

		go seed.Germinate(spawnChan, vertexChan, &wait)

		spawnChan <- types.Spawn{
			Hash:        "foo",
			Label:       "//foo",
			OutputFiles: []string{"foo.txt"},
			LogPathHash: "logPathHash1",
			EnvHash:     "someEnvHash1",
		}

		spawnChan <- types.Spawn{
			Hash:        "bar",
			Label:       "//bar",
			OutputFiles: []string{"bar.txt"},
			LogPathHash: "logPathHash2",
			EnvHash:     "someEnvHash2",
		}

		close(spawnChan)
		wait.Wait()
		close(vertexChan)

		v1 := seed.GetOutputFileVertex("logPathHash1", "foo.txt")
		v2 := seed.GetOutputFileVertex("logPathHash2", "bar.txt")

		g.Expect(v1.Label).To(Equal("//foo"))
		g.Expect(v2.Label).To(Equal("//bar"))
	})

	t.Run("spawns with the same hash are added to the same vertex", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		seed := NewSeed(ctx, normalizer)

		var wait sync.WaitGroup
		wait.Add(1)
		spawnChan := make(chan types.Spawn, 10)
		vertexChan := make(chan *types.Vertex, 10)

		go seed.Germinate(spawnChan, vertexChan, &wait)

		spawnChan <- types.Spawn{
			Hash:        "foo",
			Label:       "//foo",
			LogPathHash: "1",
			EnvHash:     "someEnvHash1",
		}

		spawnChan <- types.Spawn{
			Hash:        "bar",
			Label:       "//bar",
			LogPathHash: "1",
			EnvHash:     "someEnvHash2",
		}

		spawnChan <- types.Spawn{
			Hash:        "foo",
			Label:       "//foo",
			LogPathHash: "2",
			EnvHash:     "someEnvHash3",
		}

		close(spawnChan)
		wait.Wait()
		close(vertexChan)

		g.Expect(len(vertexChan)).To(Equal(3))

		vertexList := helpers.CreateListFromVertexChan(vertexChan, t)

		g.Expect(len(vertexList)).To(Equal(3))

		g.Expect(vertexList[0].Label).To(Equal("//foo"))
		g.Expect(vertexList[0].GetSpawn("1")).NotTo(BeNil())
		g.Expect(vertexList[0].GetSpawn("1").EnvHash).To(Equal("someEnvHash1"))
		g.Expect(vertexList[0].GetSpawn("2")).NotTo(BeNil())
		g.Expect(vertexList[0].GetSpawn("2").EnvHash).To(Equal("someEnvHash3"))

		g.Expect(vertexList[1].Label).To(Equal("//bar"))
		g.Expect(vertexList[1].GetSpawn("1")).NotTo(BeNil())
		g.Expect(vertexList[1].GetSpawn("1").EnvHash).To(Equal("someEnvHash2"))
		g.Expect(vertexList[1].GetSpawn("2")).To(BeNil())

		g.Expect(vertexList[2].Label).To(Equal("//foo"))
		g.Expect(vertexList[2].GetSpawn("1")).NotTo(BeNil())
		g.Expect(vertexList[2].GetSpawn("1").EnvHash).To(Equal("someEnvHash1"))
		g.Expect(vertexList[2].GetSpawn("2")).NotTo(BeNil())
		g.Expect(vertexList[2].GetSpawn("2").EnvHash).To(Equal("someEnvHash3"))

		v1 := seed.GetVertex("foo")
		v2 := seed.GetVertex("bar")

		g.Expect(v1.Label).To(Equal("//foo"))
		g.Expect(v1.GetSpawn("1")).NotTo(BeNil())
		g.Expect(v1.GetSpawn("1").EnvHash).To(Equal("someEnvHash1"))
		g.Expect(v1.GetSpawn("2")).NotTo(BeNil())
		g.Expect(v1.GetSpawn("2").EnvHash).To(Equal("someEnvHash3"))

		g.Expect(v2.Label).To(Equal("//bar"))
		g.Expect(v2.GetSpawn("1")).NotTo(BeNil())
		g.Expect(v2.GetSpawn("1").EnvHash).To(Equal("someEnvHash2"))
		g.Expect(v2.GetSpawn("2")).To(BeNil())
	})

	t.Run("golden files", func(t *testing.T) {
		g := NewGomegaWithT(t)

		files := map[string]string{
			"golden_files/changed_directory.txt": "../../normalizer/testdata/golden_files/changed_directory.txt",
			"golden_files/changed_file.txt":      "../../normalizer/testdata/golden_files/changed_file.txt",
		}

		for goldenFilepath, normalizedGoldenFile := range files {
			// Run this 1000 times to try and weed out any potential async issues
			for i := 0; i < 1000; i++ {
				b, err := os.ReadFile(normalizedGoldenFile)
				if err != nil {
					t.Error(err)
				}

				spawnChan := helpers.CreateSpawnChanFromGoldenFile(string(b), t)

				var stderr strings.Builder
				ctx := context.New(&stderr, context.Diff)

				normalizer := normalizer.New(ctx)

				seed := NewSeed(ctx, normalizer)

				var wait sync.WaitGroup
				wait.Add(1)
				vertexChan := make(chan *types.Vertex, 10)

				go seed.Germinate(spawnChan, vertexChan, &wait)

				close(spawnChan)
				wait.Wait()
				close(vertexChan)

				g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))

				output := helpers.CreateGoldenFileFromVertexChan(vertexChan, t)
				golden.Assert(t, output, goldenFilepath)
			}
		}
	})

	t.Run("smoketest with changed_file execlog", func(t *testing.T) {
		g := NewGomegaWithT(t)

		// Run this 1000 times to try and weed out any potential async issues
		for i := 0; i < 1000; i++ {
			b, err := os.ReadFile("../../normalizer/testdata/golden_files/changed_file.txt")
			if err != nil {
				t.Error(err)
			}

			spawnChan := helpers.CreateSpawnChanFromGoldenFile(string(b), t)

			var stderr strings.Builder
			ctx := context.New(&stderr, context.Diff)

			normalizer := normalizer.New(ctx)

			seed := NewSeed(ctx, normalizer)

			var wait sync.WaitGroup
			wait.Add(1)
			vertexChan := make(chan *types.Vertex, 10)

			go seed.Germinate(spawnChan, vertexChan, &wait)

			close(spawnChan)
			wait.Wait()
			close(vertexChan)

			g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))

			vertexList := helpers.CreateListFromVertexChan(vertexChan, t)

			g.Expect(len(vertexList)).To(Equal(4))
			g.Expect(vertexList[0].Label).To(Equal("//:static"))
			g.Expect(vertexList[1].Label).To(Equal("//foo:foo"))
			g.Expect(vertexList[2].Label).To(Equal("//bar:bar"))
			g.Expect(vertexList[3].Label).To(Equal("//baz:baz"))

			g.Expect(seed.GetOutputFileVertex("795af60dff57435a74ad7edd1ce4a3fa", "594453f86972a9984602031992cb4ea4").Label).To(Equal("//:static"))

			g.Expect(len(vertexList[3].GetLogPathHashes())).To(Equal(1))
			g.Expect(vertexList[3].GetLogPathHashes()[0]).To(Equal("795af60dff57435a74ad7edd1ce4a3fa"))
			t.Log("Expected Value:", vertexList[3].GetSpawn("795af60dff57435a74ad7edd1ce4a3fa").Hash)
			g.Expect(vertexList[3].GetSpawn("795af60dff57435a74ad7edd1ce4a3fa").Hash).To(Equal("d8627614ef1186cdeba52fa38aa23654645acc124cf7e793458783cdc1a50bd6145a9bd157d798267f6272a1cd2dd9aa"))
			g.Expect(seed.GetVertex("d8627614ef1186cdeba52fa38aa23654645acc124cf7e793458783cdc1a50bd6145a9bd157d798267f6272a1cd2dd9aa").Label).To(Equal("//baz:baz"))
		}
	})
}
