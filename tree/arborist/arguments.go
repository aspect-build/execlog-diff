package arborist

import (
	"github.com/aspect-build/execlog-diff/types"
)

func (arborist *Arborist) argumentsAreIdentical(vertex *types.Vertex) bool {
	logPathHashes := vertex.GetLogPathHashes()

	logPathHash1 := logPathHashes[0]
	spawn1 := vertex.GetSpawn(logPathHash1)

	logPathHash2 := logPathHashes[1]
	spawn2 := vertex.GetSpawn(logPathHash2)

	if len(spawn1.Args) != len(spawn2.Args) {
		return false
	}

	for i, arg := range spawn1.Args {
		if spawn2.Args[i] == arg {
			continue
		} else {
			return false
		}
	}

	return true
}
