package normalizer

import (
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"

	. "github.com/onsi/gomega"
)

func TestPopulateInputFiles(t *testing.T) {
	t.Run("Files are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		storeInputSets := PopulateFiles(normalizer)

		filesAreNotUpdated(t, storeInputSets, ctx)
	})

	t.Run("Input sets are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		storeInputSets := PopulateFiles(normalizer)

		inputSetsAreNotUpdated(t, storeInputSets, ctx)
	})

	t.Run("Files are added to spawn", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		filePath := StoreFiles(normalizer)

		filePath.Run(types.ExecLogEntry{
			File: &types.File{
				Id:          1,
				Path:        "foo",
				LogPathHash: SAMPLE_LOG_PATH_HASH,
				Hash:        SAMPLE_FILE_HASH,
			},
		})

		filePath.Run(types.ExecLogEntry{
			File: &types.File{
				Id:          2,
				Path:        "bar",
				LogPathHash: SAMPLE_LOG_PATH_HASH,
				Hash:        SAMPLE_FILE_HASH,
			},
		})

		filePath.Run(types.ExecLogEntry{
			File: &types.File{
				Id:          3,
				Path:        "other",
				LogPathHash: SAMPLE_LOG_PATH_HASH,
				Hash:        SAMPLE_FILE_HASH,
			},
		})

		storeInputSets := StoreInputSets(normalizer)

		id := int64(4)
		storeInputSets.Run(types.ExecLogEntry{
			InputSet: &types.InputSet{
				Id:              id,
				FileIds:         []int64{1},
				NestedInputSets: []int64{},
				LogPathHash:     SAMPLE_LOG_PATH_HASH,
			},
		})

		id2 := int64(5)
		storeInputSets.Run(types.ExecLogEntry{
			InputSet: &types.InputSet{
				Id:              id2,
				FileNames:       []string{"baz"},
				FileIds:         []int64{2},
				NestedInputSets: []int64{4},
				LogPathHash:     SAMPLE_LOG_PATH_HASH,
			},
		})

		populate_files := PopulateFiles(normalizer)

		inputSetId := int64(4)
		outputFileIds := []int64{3}
		outputSetIds := []int64{5}
		spawn := &types.IncomingSpawn{
			Id:            5,
			LogPathHash:   SAMPLE_LOG_PATH_HASH,
			InputSetId:    inputSetId,
			OutputFileIds: outputFileIds,
			OutputSetIds:  outputSetIds,
		}

		result := populate_files.Run(types.ExecLogEntry{
			Spawn: spawn,
		})

		g.Expect(len(result.Spawn.InputFiles)).To(Equal(1))
		g.Expect(result.Spawn.InputFiles[0]).To(Equal("foo"))
		g.Expect(len(result.Spawn.OutputFiles)).To(Equal(4))
		g.Expect(result.Spawn.OutputFiles[0]).To(Equal("bar"))
		g.Expect(result.Spawn.OutputFiles[1]).To(Equal("baz"))
		g.Expect(result.Spawn.OutputFiles[2]).To(Equal("foo"))
		g.Expect(result.Spawn.OutputFiles[3]).To(Equal("other"))
	})
}
