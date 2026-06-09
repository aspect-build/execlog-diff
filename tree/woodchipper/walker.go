/*
 * As we go through the tree and determine that the bazel action spawns on a given vertex are the same,
 * we remove that vertex and its spawns from memory.
 * However, the files that those spawns create might have already been added as inputs to another spawn / vertex.
 * In an effort to reduce memory usage, we want to occasionally go through the tree and prune these files.
 * That is the purpose of this class, to walk the tree and remove "zombie" input files
 */

package woodchipper

import (
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/tree/seed"
	"github.com/aspect-build/execlog-diff/types"
)

type walker struct {
	ctx        *context.Context
	spawner    chan struct{}
	wait       *sync.WaitGroup
	seed       *seed.Seed
	normalizer *normalizer.Normalizer
}

func NewMapWalker(ctx *context.Context, seed *seed.Seed, normalizer *normalizer.Normalizer) *walker {
	return &walker{
		ctx:        ctx,
		spawner:    make(chan struct{}, 1),
		wait:       &sync.WaitGroup{},
		seed:       seed,
		normalizer: normalizer,
	}
}

// When requested, walk the tree and remove any "zombie" input files from spawns
func (w *walker) walker() {
	for s := range w.spawner {
		if !w.ctx.Continue() {
			break
		}

		_ = s

		w.seed.GetHashToVertexMap().Range(func(k, val interface{}) bool {
			v := val.(*types.Vertex)

			v.Mutex.Lock()
			logPathHashes := v.GetLogPathHashes()
			v.Mutex.Unlock()
			for _, logPathHash := range logPathHashes {
				inputFiles := make([]string, 0)
				v.Mutex.Lock()
				spawn := v.GetSpawn(logPathHash)
				v.Mutex.Unlock()

				for _, filePath := range spawn.InputFiles {
					if hash := w.normalizer.GetFileHash(logPathHash, filePath); hash != "" {
						inputFiles = append(inputFiles, filePath)
					}
				}
				spawn.InputFiles = inputFiles
				v.Mutex.Lock()
				v.SetSpawn(logPathHash, spawn)
				v.Mutex.Unlock()
			}

			if w.ctx.Continue() {
				return true
			}
			return false
		})
	}
	w.wait.Done()
}

// Start the walker go routine
func (walker *walker) EnableWalker() {
	walker.wait.Add(1)
	go walker.walker()
}

// Queue up a walker run if one has not already been requested
func (walker *walker) Walk() {
	if len(walker.spawner) == 0 {
		walker.spawner <- struct{}{}
	}
}

// Wait for the last walker run to complete and then return
func (walker *walker) Wait() {
	close(walker.spawner)
	walker.wait.Wait()
}
