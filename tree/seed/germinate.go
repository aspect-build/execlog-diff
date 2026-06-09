package seed

import (
	"sync"

	"github.com/aspect-build/execlog-diff/types"
)

/*
 * This function will take spawns and begin creating the tree/graph vertices.
 */
func (seed *Seed) Germinate(spawnChan <-chan types.Spawn, vertexChan chan<- *types.Vertex, waitGroup *sync.WaitGroup) {
	for spawn := range spawnChan {
		if !seed.ctx.Continue() {
			break
		}

		vertex := seed.GetVertex(spawn.Hash)
		if vertex == nil {
			vertex = &types.Vertex{
				Mutex: sync.Mutex{},
				Label: spawn.Label,
			}
			seed.SetVertex(spawn.Hash, vertex)
		}

		vertex.Mutex.Lock()
		vertex.SetSpawn(spawn.LogPathHash, &spawn)
		vertex.Mutex.Unlock()

		for _, file := range spawn.OutputFiles {
			seed.SetOutputFileVertex(spawn.LogPathHash, file, vertex)
		}

		vertexChan <- vertex
	}
	waitGroup.Done()
}
