/*
 * HashFilePaths is meant to hash filepaths for use in other functions, in order to help reduce memory usage.
 * Given the following path as an example:
 * - bazel-out/darwin_arm64-fastbuild/bin/normalizer/foo.bar
 * This path adds up when dealing with 1000's of spawns containing potentially 100s of files each.
 * The above path can instead be stored on the spawn as "165d0f59517e9003067544ad726c8cf8".
 *
 * This function will save the hashed filepath so it can be used by other functions when doing filepath operations.
 */
package normalizer

import (
	"crypto/md5"
	"encoding/hex"
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"
	"github.com/aspect-build/execlog-diff/diagnostic"
)

// In the very unlikely event of a hash clash when hashing filepaths, add this character as a suffix.
const hash_clash_suffix = "_"

func HashFilePaths(n *Normalizer) normalizer {
	return &file_path{
		hash_to_file_path_map: n.hash_to_file_path_map,
		ctx:                   n.ctx,
	}
}

type file_path struct {
	hash_to_file_path_map *sync.Map
	ctx                   *context.Context
}

func (file_path *file_path) Run(entry types.ExecLogEntry) types.ExecLogEntry {
	file := entry.File
	if file != nil {
		filePathHash, err := hashFilepath(file.Path)
		if err != nil {
			file_path.ctx.Warning(&diagnostic.Diagnostic{
				Title:   "Unable to hash: " + file.Path,
				Message: err.Error(),
			})
			return entry
		}

		if len(filePathHash) < len(file.Path) {
			storedFilePathHash := file_path.storeFilePathHash(filePathHash, file.Path)

			file.Path = storedFilePathHash
			entry.File = file
		}
	}
	return entry
}

// In the event of a hash clash, this function will recursively add "hash_clash_suffix"
// to a hash until it is unique. This is done to ensure we can properly lookup our filepaths later.
func (file_path *file_path) storeFilePathHash(hash string, filepath string) string {
	if val, ok := file_path.hash_to_file_path_map.Load(hash); ok {
		storedPath := val.(string)
		if storedPath != filepath {
			return file_path.storeFilePathHash(hash+hash_clash_suffix, filepath)
		}
	} else {
		file_path.hash_to_file_path_map.Store(hash, filepath)
	}

	return hash
}

func hashFilepath(path string) (string, error) {
	hasher := md5.New()
	_, err := hasher.Write([]byte(path))
	if err != nil {
		return path, err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
