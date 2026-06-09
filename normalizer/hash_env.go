/*
 * HashEnv is meant to reduce the size of environment variables on a given bazel spawn while its stored in memory.
 * For a given bazel target type, the environment is often the same (and often contains many key/value pairs).
 * Rather than save the entire thing on each spawn, save it as a hash on the spawn and look it up later if needed.
 */
package normalizer

import (
	"crypto/sha512"
	"encoding/hex"
	"slices"
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"
)

var envHashColon = []byte(":")

func HashEnv(n *Normalizer) normalizer {
	return &hash_env{
		hash_to_env_map: n.hash_to_env_map,
		ctx:             n.ctx,
	}
}

type hash_env struct {
	hash_to_env_map *sync.Map
	ctx             *context.Context
}

func (hash_env *hash_env) Run(entry types.ExecLogEntry) types.ExecLogEntry {
	spawn := entry.Spawn
	if spawn != nil {
		// Pull keys out of env map so they can be sorted before being hashed (ordering matters when hashing)
		envKeys := make([]string, 0, len(spawn.Env))
		for k := range spawn.Env {
			envKeys = append(envKeys, k)
		}

		slices.Sort(envKeys)

		envHasher := sha512.New384()
		for _, key := range envKeys {
			envHasher.Write([]byte(key))
			envHasher.Write(envHashColon)
			envHasher.Write([]byte(spawn.Env[key]))
		}

		envHash := hex.EncodeToString(envHasher.Sum(nil))
		if _, ok := hash_env.hash_to_env_map.Load(envHash); !ok {
			hash_env.hash_to_env_map.Store(envHash, spawn.Env)
		}

		spawn.EnvHash = envHash
		entry.Spawn = spawn
	}
	return entry
}
