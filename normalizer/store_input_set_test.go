package normalizer

import (
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"

	. "github.com/onsi/gomega"
)

func TestStoreInputSets(t *testing.T) {
	t.Run("Files are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		storeInputSets := StoreInputSets(normalizer)

		filesAreNotUpdated(t, storeInputSets, ctx)
	})

	t.Run("Spawns are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		storeInputSets := StoreInputSets(normalizer)

		spawnsAreNotUpdated(t, storeInputSets, ctx)
	})

	t.Run("Input sets are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		storeInputSets := StoreInputSets(normalizer)

		inputSetsAreNotUpdated(t, storeInputSets, ctx)
	})

	t.Run("Input sets are stored and retrieved", func(t *testing.T) {
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

		storeInputSets := StoreInputSets(normalizer)

		id := int64(3)
		inputSet := &types.InputSet{
			Id:              id,
			FileNames:       []string{"baz"},
			FileIds:         []int64{1, 2},
			NestedInputSets: []int64{},
			LogPathHash:     SAMPLE_LOG_PATH_HASH,
		}

		storeInputSets.Run(types.ExecLogEntry{
			InputSet: inputSet,
		})

		g.Expect(len(normalizer.input_set_id_to_files_map)).To(Equal(1))
		g.Expect(len(normalizer.input_set_id_to_files_map[lookupId(SAMPLE_LOG_PATH_HASH, id)])).To(Equal(3))
		g.Expect(normalizer.input_set_id_to_files_map[lookupId(SAMPLE_LOG_PATH_HASH, id)]).To(ContainElement(ContainSubstring("foo")))
		g.Expect(normalizer.input_set_id_to_files_map[lookupId(SAMPLE_LOG_PATH_HASH, id)]).To(ContainElement(ContainSubstring("bar")))
		g.Expect(normalizer.input_set_id_to_files_map[lookupId(SAMPLE_LOG_PATH_HASH, id)]).To(ContainElement(ContainSubstring("baz")))
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))

		g.Expect(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, id))).To(ContainElement(ContainSubstring("foo")))
		g.Expect(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, id))).To(ContainElement(ContainSubstring("bar")))
		g.Expect(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, id))).To(ContainElement(ContainSubstring("baz")))
	})

	t.Run("Nested input sets are stored and retrieved", func(t *testing.T) {
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

		storeInputSets := StoreInputSets(normalizer)

		id := int64(3)
		storeInputSets.Run(types.ExecLogEntry{
			InputSet: &types.InputSet{
				Id:              id,
				FileIds:         []int64{1},
				NestedInputSets: []int64{},
				LogPathHash:     SAMPLE_LOG_PATH_HASH,
			},
		})

		id2 := int64(4)
		storeInputSets.Run(types.ExecLogEntry{
			InputSet: &types.InputSet{
				Id:              id2,
				FileNames:       []string{"baz"},
				FileIds:         []int64{2},
				NestedInputSets: []int64{3},
				LogPathHash:     SAMPLE_LOG_PATH_HASH,
			},
		})

		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))

		g.Expect(len(normalizer.input_set_id_to_files_map)).To(Equal(2))
		g.Expect(len(normalizer.input_set_id_to_files_map[lookupId(SAMPLE_LOG_PATH_HASH, id)])).To(Equal(1))
		g.Expect(normalizer.input_set_id_to_files_map[lookupId(SAMPLE_LOG_PATH_HASH, id)]).To(ContainElement(ContainSubstring("foo")))

		g.Expect(len(normalizer.input_set_id_to_files_map[lookupId(SAMPLE_LOG_PATH_HASH, id2)])).To(Equal(3))
		g.Expect(normalizer.input_set_id_to_files_map[lookupId(SAMPLE_LOG_PATH_HASH, id2)]).To(ContainElement(ContainSubstring("bar")))
		g.Expect(normalizer.input_set_id_to_files_map[lookupId(SAMPLE_LOG_PATH_HASH, id2)]).To(ContainElement(ContainSubstring("baz")))
		g.Expect(normalizer.input_set_id_to_files_map[lookupId(SAMPLE_LOG_PATH_HASH, id2)]).To(ContainElement(ContainSubstring(nestedLookupId(SAMPLE_LOG_PATH_HASH, 3))))

		g.Expect(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, id))).To(ContainElement(ContainSubstring("foo")))

		g.Expect(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, id2))).To(ContainElement(ContainSubstring("foo")))
		g.Expect(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, id2))).To(ContainElement(ContainSubstring("bar")))
		g.Expect(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, id2))).To(ContainElement(ContainSubstring("baz")))
		g.Expect(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, id2))).ToNot(ContainElement(ContainSubstring(nestedLookupId(SAMPLE_LOG_PATH_HASH, 3))))
	})
}
