package arborist

import (
	"github.com/aspect-build/execlog-diff/types"
)

func (arborist *Arborist) inputFilesAreIdentical(vertex *types.Vertex) bool {
	logPathHashes := vertex.GetLogPathHashes()

	logPathHash1 := logPathHashes[0]
	logPathHash2 := logPathHashes[1]

	spawn1 := vertex.GetSpawn(logPathHash1)

	return arborist.filesAreIdentical(spawn1.InputFiles, logPathHash1, logPathHash2)
}

func (arborist *Arborist) outputFilesAreIdentical(vertex *types.Vertex) bool {
	logPathHashes := vertex.GetLogPathHashes()

	logPathHash1 := logPathHashes[0]
	logPathHash2 := logPathHashes[1]

	spawn1 := vertex.GetSpawn(logPathHash1)

	return arborist.filesAreIdentical(spawn1.OutputFiles, logPathHash1, logPathHash2)
}

func (arborist *Arborist) filesAreIdentical(files []string, logPathHash1 string, logPathHash2 string) bool {
	// Spawns in the same vertex will always have the same files (as in file paths).
	// The file contents may be different which is what we want to check here.
	// But because they are the same we can just get the file names from one of the spawns.
	for _, filePath := range files {
		spawn1FileHash := arborist.normalizer.GetFileHash(logPathHash1, filePath)
		spawn2FileHash := arborist.normalizer.GetFileHash(logPathHash2, filePath)

		// Files that are the same will be removed from the File->FileHash map to save memory.
		// So if we get a "" here it means that file is the same and has been pruned from the map.
		if spawn1FileHash == "" || spawn2FileHash == "" {
			continue
		}

		if spawn1FileHash != spawn2FileHash {
			return false
		}
	}

	return true
}
