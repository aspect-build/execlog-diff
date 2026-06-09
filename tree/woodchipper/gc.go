/*
 * When running execlog-diff, we are processing large amounts of data.
 * Much of that data is stored in memory and then pruned.
 * We often hit a situation where the garbage collector was not running often enough and caused us to OOM.
 * This class allows us to asynchronously queue up a garbage collector run, once at a time.
 */

package woodchipper

import (
	"runtime"
	"sync"

	"github.com/aspect-build/execlog-diff/context"
)

type gc struct {
	ctx                 *context.Context
	spawner             chan struct{}
	wait                *sync.WaitGroup
	runGarbageCollector func()
}

// Wrapper that can be overwritten during tests to avoid actually running the garbage collector
func runGarbageCollector() {
	runtime.GC()
}

func NewGarbageCollector(ctx *context.Context) *gc {
	return &gc{
		ctx:                 ctx,
		spawner:             make(chan struct{}, 1),
		wait:                &sync.WaitGroup{},
		runGarbageCollector: runGarbageCollector,
	}
}

func (gc *gc) gc() {
	for s := range gc.spawner {
		if !gc.ctx.Continue() {
			break
		}

		_ = s

		gc.runGarbageCollector()
	}
	gc.wait.Done()
}

// Start the go routine
func (gc *gc) EnableCollection() {
	gc.wait.Add(1)
	go gc.gc()
}

// Queue up a garbage collection run if one has not already been requested
func (gc *gc) Collect() {
	if len(gc.spawner) == 0 {
		gc.spawner <- struct{}{}
	}
}

// Wait for the last garbage collection run to complete and then return
func (gc *gc) Wait() {
	close(gc.spawner)
	gc.wait.Wait()
}
