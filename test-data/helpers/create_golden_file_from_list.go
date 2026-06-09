package helpers

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/types"
)

func CreateGoldenFileFromVertexList(vertexChan []*types.Vertex, t *testing.T) string {
	sort.Slice(vertexChan, func(i, j int) bool {
		v1 := vertexChan[i]
		v2 := vertexChan[j]

		s1 := v1.GetSpawn(v1.GetLogPathHashes()[0])
		s2 := v2.GetSpawn(v2.GetLogPathHashes()[0])

		return s1.Hash > s2.Hash
	})

	messages := make([]string, 0)
	for _, message := range vertexChan {
		out, err := json.Marshal(message)
		if err != nil {
			t.Error(err)
		}

		messages = append(messages, string(out))
	}

	return strings.Join(messages, "\n")
}
