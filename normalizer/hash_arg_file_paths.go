/*
 * HashArgFilePaths is meant to reduce the size of the arguments on a given bazel spawn while its stored in memory.
 * The argument list can potentially contain MANY long filepaths such as:
 * - bazel-out/darwin_arm64-fastbuild/bin/execlog-diff/normalizer/foo.bar
 * This adds up when dealing with 1000's of spawns containing potentially 100s of files each.
 * The above path can instead be stored on the spawn as "165d0f59517e9003067544ad726c8cf8".
 * Which is more than 3 times smaller in this case.
 *
 * The mapping of "hash" -> "file path" should already stored in the normalizer and
 * so this function should not need to store them again for later lookup.
 */
package normalizer

import (
	"strings"
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"
	"github.com/aspect-build/execlog-diff/diagnostic"
)

func HashArgFilePaths(n *Normalizer) normalizer {
	return &arg_file_paths{
		hash_to_file_path_map: n.hash_to_file_path_map,
		ctx:                   n.ctx,
	}
}

type arg_file_paths struct {
	hash_to_file_path_map *sync.Map
	ctx                   *context.Context
}

func (arg_file_paths *arg_file_paths) Run(entry types.ExecLogEntry) types.ExecLogEntry {
	spawn := entry.Spawn
	if spawn != nil {
		args := make([]string, len(spawn.Args))
		for i, arg := range spawn.Args {
			splitArg := splitOnEquals(arg)
			argPieces := make([]string, len(splitArg))
			for j, a := range splitArg {
				if arg != "" {
					argPieces[j] = arg_file_paths.getArgHash(a)
				}
			}

			args[i] = strings.Join(argPieces, "=")
		}

		spawn.Args = args
		entry.Spawn = spawn
	}
	return entry
}

func (arg_file_paths *arg_file_paths) getArgHash(arg string) string {
	filePathHash, err := hashFilepath(arg)
	if err != nil {
		arg_file_paths.ctx.Warning(&diagnostic.Diagnostic{
			Title:   "Unable to hash: " + arg,
			Message: err.Error(),
		})

		// In the event that there is an issue, the original value is preferable to no value.
		// The purpose of HashArgFilePaths is to save space on the spawn args when they are stored in memory.
		// It is therefore not a critical failure if we fail to hash, so just log a warning and move on.
		// We will still want to do comparisons later using these args so we dont want to lose track of the original value.
		return arg
	}

	return arg_file_paths.getArgHashRecursive(arg, filePathHash)
}

// This function is recursive to handle the very unlikely possibility of hash clashes.
// Should 2 filepaths have a hash clash, one will have been saved with our "hash_clash_suffix" added to the end.
// This function will recursively check the hash to ensure we use the right one so that our arg can be restored later.
func (arg_file_paths *arg_file_paths) getArgHashRecursive(arg string, hash string) string {
	if val, ok := arg_file_paths.hash_to_file_path_map.Load(hash); ok {
		storedArg := val.(string)
		if storedArg != arg {
			return arg_file_paths.getArgHashRecursive(arg+hash_clash_suffix, hash)
		} else {
			return hash
		}
	}

	return arg
}

func splitOnEquals(arg string) []string {
	if strings.Contains(arg, "=") && !strings.Contains(arg, "==") && !strings.HasPrefix(arg, "=") && !strings.HasSuffix(arg, "=") {
		return strings.Split(arg, "=")
	} else {
		return []string{arg}
	}
}
