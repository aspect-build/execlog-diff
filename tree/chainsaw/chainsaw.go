package chainsaw

import (
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/tree/seed"
	"github.com/aspect-build/execlog-diff/tree/woodchipper"
	"github.com/aspect-build/execlog-diff/types"
)

type Chainsaw struct {
	ctx         *context.Context
	normalizer  *normalizer.Normalizer
	seed        *seed.Seed
	woodchipper *woodchipper.Woodchipper
}

func NewChainsaw(ctx *context.Context, normalizer *normalizer.Normalizer, seed *seed.Seed) *Chainsaw {
	return &Chainsaw{
		ctx:         ctx,
		normalizer:  normalizer,
		seed:        seed,
		woodchipper: woodchipper.NewWoodchipper(ctx, seed, normalizer),
	}
}

func (c *Chainsaw) RemoveVertices(vertices <-chan *types.Vertex, wait *sync.WaitGroup) {
	c.woodchipper.Start()
	count := 0

	for vertex := range vertices {
		if !c.ctx.Continue() {
			break
		}

		if vertex == nil {
			// Should never hit this case, added just in case
			continue
		}

		vertex.Mutex.Lock()
		logPathHashes := vertex.GetLogPathHashes()
		vertex.Mutex.Unlock()

		for i, logPathHash := range logPathHashes {
			vertex.Mutex.Lock()
			spawn := vertex.GetSpawn(logPathHash)
			vertex.Mutex.Unlock()

			for _, outputFilePath := range spawn.OutputFiles {
				c.seed.DeleteOutputFileVertex(logPathHash, outputFilePath)
				c.normalizer.DeleteFileHash(logPathHash, outputFilePath)
				c.normalizer.DeleteFileSize(logPathHash, outputFilePath)
			}

			if i == 0 {
				// We only need this to be run once
				c.seed.DeleteVertex(spawn.Hash)
			}
		}

		// Remove this vertex from parents & memory
		// c.woodchipper.Chip()
		count++
		// Run cleanup every 1000 items.
		// This keeps memory usage in check without killing performance.
		if count%1000 == 0 {
			c.woodchipper.Chip()
		}
	}
	c.woodchipper.Stop()
	wait.Done()
}
