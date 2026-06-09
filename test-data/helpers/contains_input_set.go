package helpers

import (
	"github.com/aspect-build/execlog-diff/types"
)

func ContainsInputSet(execLogList []types.ExecLogEntry, paths []string) bool {
	for _, execlog := range execLogList {
		if execlog.InputSet != nil {
			if len(execlog.InputSet.FileNames) != len(paths) {
				continue
			}

			containsPaths := make(map[string]bool)
			for _, path := range paths {
				containsPaths[path] = false
			}

			for _, path := range paths {
				for _, inputSetPath := range execlog.InputSet.FileNames {
					if path == inputSetPath {
						containsPaths[path] = true
					}
				}
			}

			inputSetMatches := true
			for _, contains := range containsPaths {
				if !contains {
					inputSetMatches = false
				}
			}

			if inputSetMatches {
				return true
			}
		}

	}
	return false
}
