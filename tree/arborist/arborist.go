/*
 * Arborists main purpose is to take a given tree vertex and determine if its spawns are the same
 */
package arborist

import (
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/types"
)

type Arborist struct {
	normalizer *normalizer.Normalizer
}

func New(normalizer *normalizer.Normalizer) *Arborist {
	return &Arborist{
		normalizer: normalizer,
	}
}

func (arborist *Arborist) VertexIsIdentical(vertex *types.Vertex) bool {
	logPathHashes := vertex.GetLogPathHashes()
	if len(logPathHashes) != 2 {
		return false
	}

	return arborist.inputFilesAreIdentical(vertex) &&
		arborist.outputFilesAreIdentical(vertex) &&
		arborist.argumentsAreIdentical(vertex) &&
		arborist.environmentIsIdentical(vertex)
}
