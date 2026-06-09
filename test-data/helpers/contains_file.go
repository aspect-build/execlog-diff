package helpers

import (
	"github.com/aspect-build/execlog-diff/types"
)

func ContainsFile(execLogList []types.ExecLogEntry, path string) bool {
	for _, execlog := range execLogList {
		if execlog.File != nil && execlog.File.Path == path {
			return true
		}
	}
	return false
}
