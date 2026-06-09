package woodchipper

import (
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/test-data/helpers"
	"github.com/aspect-build/execlog-diff/tree/seed"
	"github.com/aspect-build/execlog-diff/types"
	"github.com/aspect-build/execlog-diff/diagnostic"
	. "github.com/onsi/gomega"
)

func TestWalker(t *testing.T) {
	t.Run("Walker removes input files that should be removed", func(t *testing.T) {
		g := NewGomegaWithT(t)

		b, err := os.ReadFile("../../normalizer/testdata/golden_files/changed_file.txt")
		if err != nil {
			t.Error(err)
		}

		spawnChan := helpers.CreateSpawnChanFromGoldenFile(string(b), t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		seed := seed.NewSeed(ctx, normalizer)

		var wait sync.WaitGroup
		wait.Add(1)
		vertexChan := make(chan *types.Vertex, 10)

		go seed.Germinate(spawnChan, vertexChan, &wait)

		close(spawnChan)
		wait.Wait()
		close(vertexChan)

		spawnThatIsDependedOn := "9b197a821a06c5765cd6f22d5947207de472dcb042ff3b1db1bdf3ddd0b429acf8f4151763ac36ba5fb42b8a257ad2f5"
		spawnThatDependsOn := "9366b90f37a5eeaed3528b11539487d8a23781e87522601f5bb26f3cd02350916405d622b2400f333ff95931aaf3eaf3"
		logPathHash := "795af60dff57435a74ad7edd1ce4a3fa"
		filePathHash := "c08d251aff44b6eb190987e839104f2c"

		vertex := seed.GetVertex(spawnThatDependsOn)

		for _, spawn := range vertex.Spawns {
			g.Expect(spawn.InputFiles).To(ContainElement(ContainSubstring(filePathHash)))
		}

		seed.DeleteOutputFileVertex(logPathHash, filePathHash)
		normalizer.DeleteFileHash(logPathHash, filePathHash)
		seed.DeleteVertex(spawnThatIsDependedOn)

		for _, spawn := range vertex.Spawns {
			g.Expect(spawn.InputFiles).To(ContainElement(ContainSubstring(filePathHash)))
		}

		w := &walker{
			ctx:        ctx,
			spawner:    make(chan struct{}, 1),
			wait:       &sync.WaitGroup{},
			seed:       seed,
			normalizer: normalizer,
		}

		w.EnableWalker()
		w.Walk()
		w.Wait()

		vertex = seed.GetVertex(spawnThatDependsOn)
		for _, spawn := range vertex.Spawns {
			g.Expect(spawn.InputFiles).NotTo(ContainElement(ContainSubstring(filePathHash)))
		}

		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})

	t.Run("Walker stops when requested by the context", func(t *testing.T) {
		g := NewGomegaWithT(t)

		b, err := os.ReadFile("../../normalizer/testdata/golden_files/changed_file.txt")
		if err != nil {
			t.Error(err)
		}

		spawnChan := helpers.CreateSpawnChanFromGoldenFile(string(b), t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		seed := seed.NewSeed(ctx, normalizer)

		var wait sync.WaitGroup
		wait.Add(1)
		vertexChan := make(chan *types.Vertex, 10)

		go seed.Germinate(spawnChan, vertexChan, &wait)

		close(spawnChan)
		wait.Wait()
		close(vertexChan)

		spawnThatIsDependedOn := "9b197a821a06c5765cd6f22d5947207de472dcb042ff3b1db1bdf3ddd0b429acf8f4151763ac36ba5fb42b8a257ad2f5"
		spawnThatDependsOn := "9366b90f37a5eeaed3528b11539487d8a23781e87522601f5bb26f3cd02350916405d622b2400f333ff95931aaf3eaf3"
		logPathHash := "795af60dff57435a74ad7edd1ce4a3fa"
		filePathHash := "c08d251aff44b6eb190987e839104f2c"

		vertex := seed.GetVertex(spawnThatDependsOn)

		for _, spawn := range vertex.Spawns {
			g.Expect(spawn.InputFiles).To(ContainElement(ContainSubstring(filePathHash)))
		}

		seed.DeleteOutputFileVertex(logPathHash, filePathHash)
		normalizer.DeleteFileHash(logPathHash, filePathHash)
		seed.DeleteVertex(spawnThatIsDependedOn)

		for _, spawn := range vertex.Spawns {
			g.Expect(spawn.InputFiles).To(ContainElement(ContainSubstring(filePathHash)))
		}

		w := &walker{
			ctx:        ctx,
			spawner:    make(chan struct{}, 1),
			wait:       &sync.WaitGroup{},
			seed:       seed,
			normalizer: normalizer,
		}

		w.EnableWalker()

		ctx.Error(&diagnostic.Diagnostic{
			Message: "STOP!",
		})

		time.Sleep(10 * time.Millisecond)

		w.Walk()
		w.Wait()

		vertex = seed.GetVertex(spawnThatDependsOn)
		for _, spawn := range vertex.Spawns {
			g.Expect(spawn.InputFiles).To(ContainElement(ContainSubstring(filePathHash)))
		}

		g.Expect(len(ctx.GetDiagnostics())).To(Equal(1))
	})
}
