package helpers

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/types"
)

func CreateGoldenFileFromExecLogChan(execLogChan <-chan types.ExecLogEntry, t *testing.T) string {
	messages := make([]string, 0)
	for message := range execLogChan {
		out, err := json.Marshal(message)
		if err != nil {
			t.Error(err)
		}

		messages = append(messages, string(out))
	}

	return strings.Join(messages, "\n")
}

func CreateGoldenFileFromSpawnChan(spawnChan <-chan types.Spawn, t *testing.T) string {
	messages := make([]string, 0)
	for message := range spawnChan {
		out, err := json.Marshal(message)
		if err != nil {
			t.Error(err)
		}

		messages = append(messages, string(out))
	}

	return strings.Join(messages, "\n")
}

func CreateGoldenFileFromVertexChan(vertexChan <-chan *types.Vertex, t *testing.T) string {
	messages := make([]string, 0)
	for message := range vertexChan {
		out, err := json.Marshal(&message)
		if err != nil {
			t.Error(err)
		}

		messages = append(messages, string(out))
	}

	return strings.Join(messages, "\n")
}
