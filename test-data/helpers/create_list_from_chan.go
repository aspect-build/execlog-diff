package helpers

import (
	"testing"

	"github.com/aspect-build/execlog-diff/types"
)

func CreateListFromExecLogChan(execLogChan <-chan types.ExecLogEntry, t *testing.T) []types.ExecLogEntry {
	execlogs := make([]types.ExecLogEntry, 0)
	for execlog := range execLogChan {
		execlogs = append(execlogs, execlog)
	}

	return execlogs
}

func CreateListFromSpawnChan(execLogChan <-chan types.Spawn, t *testing.T) []types.Spawn {
	spawns := make([]types.Spawn, 0)
	for spawn := range execLogChan {
		spawns = append(spawns, spawn)
	}

	return spawns
}

func CreateListFromVertexChan(execLogChan <-chan *types.Vertex, t *testing.T) []*types.Vertex {
	vertices := make([]*types.Vertex, 0)
	for vertex := range execLogChan {
		vertices = append(vertices, vertex)
	}

	return vertices
}
