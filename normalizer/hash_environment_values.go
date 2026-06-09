/*
 * HashCredentials is meant to remove any leaked private information from spawns.
 * Primarily: environment variables.
 * Often environment variables may contain things like auth keys or sensitive values, especially on CI.
 *
 * We do not want to save the original values and we want to make it as secure as possible, so using SHA384
 */
package normalizer

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"
	"github.com/aspect-build/execlog-diff/diagnostic"
)

func HashEnvironmentValues(n *Normalizer) normalizer {
	return &credentials{
		ctx: n.ctx,
	}
}

type credentials struct {
	ctx *context.Context
}

func (credentials *credentials) Run(entry types.ExecLogEntry) types.ExecLogEntry {
	spawn := entry.Spawn
	if spawn != nil {
		normalizedEnv := make(map[string]string)
		hasher := sha512.New384()
		for key, val := range spawn.Env {
			hasher.Reset()
			_, err := hasher.Write([]byte(val))

			if err != nil {
				credentials.ctx.Warning(&diagnostic.Diagnostic{
					Title:   "Unable to hash environment variable",
					Message: "Removing '" + key + "' from the list of environment variables",
				})
				continue
			}
			normalizedEnv[key] = hex.EncodeToString(hasher.Sum(nil))
		}

		spawn.Env = normalizedEnv
		entry.Spawn = spawn
	}
	return entry
}
