package helpers

import (
	"github.com/aspect-build/execlog-diff/types"
)

func ContainsSpawn(execLogList []types.ExecLogEntry, label string) bool {
	for _, execlog := range execLogList {
		if execlog.Spawn != nil && execlog.Spawn.Label == label {
			return true
		}
	}
	return false
}
