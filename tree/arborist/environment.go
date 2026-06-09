package arborist

import (
	"github.com/aspect-build/execlog-diff/types"
)

func (arborist *Arborist) environmentIsIdentical(vertex *types.Vertex) bool {
	logPathHashes := vertex.GetLogPathHashes()

	logPathHash1 := logPathHashes[0]
	spawn1 := vertex.GetSpawn(logPathHash1)

	logPathHash2 := logPathHashes[1]
	spawn2 := vertex.GetSpawn(logPathHash2)

	return spawn1.EnvHash == spawn2.EnvHash
}
