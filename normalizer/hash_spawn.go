/*
 * HashSpawn is meant to create a unique hash for each bazel action spawn (unique per execlog).
 * This hash will be used to determine that 2 bazel spawns from 2 different execlogs are easily comparable.
 * The 3 things that determine this "easy comparability" are:
 *  - Label
 *  - Input files paths (NOT file content)
 *  - Output files paths (NOT file content)
 */
package normalizer

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"
	"github.com/aspect-build/execlog-diff/diagnostic"
)

func HashSpawn(n *Normalizer) normalizer {
	return &hash_spawn{
		ctx: n.ctx,
	}
}

type hash_spawn struct {
	ctx *context.Context
}

func (hash_spawn *hash_spawn) Run(entry types.ExecLogEntry) types.ExecLogEntry {
	spawn := entry.Spawn
	if spawn != nil {
		hasher := sha512.New384()
		_, err := hasher.Write([]byte(spawn.Label))
		if err != nil {
			hash_spawn.ctx.Warning(&diagnostic.Diagnostic{
				Message: "Failed to hash spawn with Label: " + spawn.Label,
			})
			return entry
		}

		for _, file := range spawn.InputFiles {
			hasher.Write([]byte("input_file:"))
			hasher.Write([]byte(file))
		}

		for _, file := range spawn.OutputFiles {
			hasher.Write([]byte("output_file:"))
			hasher.Write([]byte(file))
		}

		spawn.Hash = hex.EncodeToString(hasher.Sum(nil))

		entry.Spawn = spawn
	}
	return entry
}
