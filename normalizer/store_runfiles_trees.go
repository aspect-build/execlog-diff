package normalizer

import (
	"github.com/aspect-build/execlog-diff/types"
)

func StoreRunfilesTrees(n *Normalizer) normalizer {
	return &store_runfiles_trees{
		runfiles_tree_id_to_input_set_id_map: n.runfiles_tree_id_to_input_set_id_map,
	}
}

type store_runfiles_trees struct {
	runfiles_tree_id_to_input_set_id_map map[MapKey]int64
}

func (s *store_runfiles_trees) Run(entry types.ExecLogEntry) types.ExecLogEntry {
	runfilesTree := entry.RunfilesTree
	if runfilesTree != nil {
		key := lookupId(runfilesTree.LogPathHash, runfilesTree.Id)
		s.runfiles_tree_id_to_input_set_id_map[key] = runfilesTree.InputSetId
	}

	return entry
}
