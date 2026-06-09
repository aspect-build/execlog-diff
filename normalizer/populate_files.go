/*
 * PopulateFiles takes a spawn's input set and converts it to files (on the spawn)
 */

package normalizer

import (
	"slices"

	"github.com/aspect-build/execlog-diff/types"
)

func PopulateFiles(n *Normalizer) normalizer {
	return &populate_files{
		file_id_to_path_map: n.file_id_to_path_map,
		n:                   n,
	}
}

type populate_files struct {
	n                   *Normalizer
	file_id_to_path_map map[MapKey]string
}

func (populate_files *populate_files) Run(entry types.ExecLogEntry) types.ExecLogEntry {
	spawn := entry.Spawn
	if spawn != nil {
		// Input Files
		inputSetKey := lookupId(spawn.LogPathHash, spawn.InputSetId)
		spawn.InputFiles = populate_files.n.getInputsFromInputSet(inputSetKey)
		slices.Sort(spawn.InputFiles)

		// Output Files
		outputFiles := spawn.OutputFiles

		for _, fileId := range spawn.OutputFileIds {
			fileIdKey := lookupId(spawn.LogPathHash, fileId)
			if path, ok := populate_files.file_id_to_path_map[fileIdKey]; ok {
				outputFiles = append(outputFiles, path)
			}
		}

		for _, setId := range spawn.OutputSetIds {
			setIdKey := lookupId(spawn.LogPathHash, setId)
			for _, file := range populate_files.n.getInputsFromInputSet(setIdKey) {
				outputFiles = append(outputFiles, file)
			}
		}

		spawn.OutputFiles = removeDuplicates(outputFiles)

		slices.Sort(spawn.OutputFiles)

		entry.Spawn = spawn
	}

	return entry
}

func removeDuplicates(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
