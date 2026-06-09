package normalizer

import (
	"fmt"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"
)

func StoreInputSets(n *Normalizer) normalizer {
	return &store_input_sets{
		file_id_to_path_map:                  n.file_id_to_path_map,
		input_set_id_to_files_map:            n.input_set_id_to_files_map,
		runfiles_tree_id_to_input_set_id_map: n.runfiles_tree_id_to_input_set_id_map,
		ctx:                                  n.ctx,
	}
}

type store_input_sets struct {
	// Updated map types to use MapKey
	file_id_to_path_map                  map[MapKey]string
	input_set_id_to_files_map            map[MapKey][]string
	runfiles_tree_id_to_input_set_id_map map[MapKey]int64
	ctx                                  *context.Context
}

func (store_input_sets *store_input_sets) Run(entry types.ExecLogEntry) types.ExecLogEntry {
	inputSet := entry.InputSet
	if inputSet != nil {
		key := lookupId(inputSet.LogPathHash, inputSet.Id)
		fileNames := inputSet.FileNames

		for _, fileId := range inputSet.FileIds {
			fileKey := lookupId(inputSet.LogPathHash, fileId)
			fileName, ok := store_input_sets.file_id_to_path_map[fileKey]
			if ok {
				fileNames = append(fileNames, fileName)
			} else {
				// Check if this ID is a runfiles tree; if so, add a nested ref to its input set
				if runfilesInputSetId, isRunfiles := store_input_sets.runfiles_tree_id_to_input_set_id_map[fileKey]; isRunfiles {
					nestedKeyString := fmt.Sprintf("%s%s:%d", nestedPrefix, inputSet.LogPathHash, runfilesInputSetId)
					fileNames = append(fileNames, nestedKeyString)
				}
			}
		}

		for _, nestedSetId := range inputSet.NestedInputSets {
			// Store nested sets as strings so they can fit in the []string value list.
			// Format: "nested:LogPathHash:ID"
			// This matches the parsing logic in normalizer.go
			nestedKeyString := fmt.Sprintf("%s%s:%d", nestedPrefix, inputSet.LogPathHash, nestedSetId)
			fileNames = append(fileNames, nestedKeyString)
		}

		store_input_sets.input_set_id_to_files_map[key] = fileNames
	}

	return entry
}
