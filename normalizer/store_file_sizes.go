/*
 * StoreFiles takes incoming file paths and their hashes to save them so they can be looked up later.
 */

package normalizer

import (
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"
)

func StoreFileSizes(n *Normalizer) normalizer {
	return &store_file_sizes{
		file_path_to_file_size_map: n.file_path_to_file_size_map,
		file_id_to_size_map:        n.file_id_to_size_map,
		ctx:                        n.ctx,
	}
}

type store_file_sizes struct {
	file_path_to_file_size_map *sync.Map
	file_id_to_size_map        map[MapKey]int64
	ctx                        *context.Context
}

func (store_file_sizes *store_file_sizes) Run(entry types.ExecLogEntry) types.ExecLogEntry {
	file := entry.File
	if file != nil {
		fileHashLookupId := lookupId(file.LogPathHash, file.Path)
		if _, ok := store_file_sizes.file_path_to_file_size_map.Load(fileHashLookupId); !ok {
			store_file_sizes.file_path_to_file_size_map.Store(fileHashLookupId, file.SizeBytes)
			if file.Id != 0 {
				fileIdLookupId := lookupId(file.LogPathHash, file.Id)
				store_file_sizes.file_id_to_size_map[fileIdLookupId] = file.SizeBytes
			}
		}
	}
	return entry
}
