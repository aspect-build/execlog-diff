package types

import (
	"sort"
	"sync"
)

// This vertex struct represents what will be a "node" / "vertex" in our resulting tree / graph.
// A vertex will contain N bazel action spawns from N execlogs.
// The spawns will have been deemed "comparable", meaning they represent the same bazel action from N execlogs.
type Vertex struct {
	Label string

	// DO NOT USE DIRECTLY. THIS IS PUBLIC FOR TESTS ONLY
	// Key: log path hash
	Spawns map[string]*Spawn

	Mutex sync.Mutex

	Same            bool
	SamenessChecked bool
}

func (v *Vertex) GetSpawn(logPathHash string) *Spawn {
	if v.Spawns == nil {
		return nil
	}

	return v.Spawns[logPathHash]
}

func (v *Vertex) SetSpawn(logPathHash string, spawn *Spawn) {
	if v.Spawns == nil {
		v.Spawns = make(map[string]*Spawn)
	}

	v.Spawns[logPathHash] = spawn
}

func (v *Vertex) HasSpawnForLog(logPathHash string) bool {
	if v.Spawns == nil {
		return false
	}

	return v.Spawns[logPathHash] != nil
}

func (v *Vertex) GetLogPathHashes() []string {
	pathHashes := make([]string, 0)
	for _, spawn := range v.Spawns {
		pathHashes = append(pathHashes, spawn.LogPathHash)
	}

	sort.Strings(pathHashes)

	return pathHashes
}
