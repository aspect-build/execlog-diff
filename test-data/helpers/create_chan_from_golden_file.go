package helpers

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/types"
)

func CreateSpawnChanFromGoldenFile(goldenFile string, t *testing.T) chan types.Spawn {
	spawnStrings := strings.Split(goldenFile, "\n")
	spawnChan := make(chan types.Spawn, len(spawnStrings))
	for _, spawnString := range spawnStrings {
		var spawn types.Spawn
		err := json.Unmarshal([]byte(spawnString), &spawn)
		if err != nil {
			t.Error(err)
		}

		spawnChan <- spawn
	}

	return spawnChan
}
