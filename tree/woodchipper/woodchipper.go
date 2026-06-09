/*
 * The primary focus of woodchipper is to reduce memory usage by removing lingering tree/graph branches.
 * We want this to happen in an asynchronous way so that it can be called without blocking the main go routines.
 * Woodchipper will coordinate the functions such as garbage collection and vertex removal
 */

package woodchipper

import (
	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/tree/seed"
)

type Woodchipper struct {
	ctx    *context.Context
	gc     *gc
	walker *walker
}

func NewWoodchipper(ctx *context.Context, seed *seed.Seed, normalizer *normalizer.Normalizer) *Woodchipper {
	return &Woodchipper{
		ctx:    ctx,
		gc:     NewGarbageCollector(ctx),
		walker: NewMapWalker(ctx, seed, normalizer),
	}
}

func (woodchipper *Woodchipper) Start() {
	woodchipper.gc.EnableCollection()
	woodchipper.walker.EnableWalker()
}

func (woodchipper *Woodchipper) Chip() {
	woodchipper.gc.Collect()
	woodchipper.walker.Walk()
}

func (woodchipper *Woodchipper) Stop() {
	woodchipper.gc.Wait()
	woodchipper.walker.Wait()
}
