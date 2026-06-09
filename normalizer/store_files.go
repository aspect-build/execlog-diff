/*
 * StoreFiles takes incoming file paths and their hashes to save them so they can be looked up later.
 */

package normalizer

import (
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"
)

func StoreFiles(n *Normalizer) normalizer {
	return &store_files{
		file_path_to_file_hash_map: n.file_path_to_file_hash_map,
		file_id_to_path_map:        n.file_id_to_path_map,
		ctx:                        n.ctx,
	}
}

type store_files struct {
	file_path_to_file_hash_map *sync.Map
	file_id_to_path_map        map[MapKey]string
	ctx                        *context.Context
}

func (store_files *store_files) Run(entry types.ExecLogEntry) types.ExecLogEntry {
	file := entry.File
	if file != nil {
		fileHashLookupId := lookupId(file.LogPathHash, file.Path)
		if _, ok := store_files.file_path_to_file_hash_map.Load(fileHashLookupId); !ok {
			store_files.file_path_to_file_hash_map.Store(fileHashLookupId, file.Hash)
			if file.Id != 0 {
				fileIdLookupId := lookupId(file.LogPathHash, file.Id)
				store_files.file_id_to_path_map[fileIdLookupId] = file.Path
			}
		}
	}
	return entry
}
