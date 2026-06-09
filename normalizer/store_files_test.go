package normalizer

import (
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"

	. "github.com/onsi/gomega"
)

func TestStoreFiles(t *testing.T) {
	t.Run("Files are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		storeFiles := StoreFiles(normalizer)

		filesAreNotUpdated(t, storeFiles, ctx)
	})

	t.Run("Spawns are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		storeFiles := StoreFiles(normalizer)

		spawnsAreNotUpdated(t, storeFiles, ctx)
	})

	t.Run("Input Sets are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		storeFiles := StoreFiles(normalizer)

		inputSetsAreNotUpdated(t, storeFiles, ctx)
	})

	t.Run("File paths are stored", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		storeFiles := StoreFiles(normalizer)

		id := int64(7)
		hash := SAMPLE_FILE_HASH
		path := NORMALIZER_BAZEL_PATH + "/foo.go"
		file := &types.File{
			Id:          id,
			Path:        path,
			LogPathHash: SAMPLE_LOG_PATH_HASH,
			Hash:        hash,
		}

		storeFiles.Run(types.ExecLogEntry{
			File: file,
		})

		g.Expect(len(normalizer.file_id_to_path_map)).To(Equal(1))
		g.Expect(normalizer.file_id_to_path_map[lookupId(SAMPLE_LOG_PATH_HASH, id)]).To(Equal(path))
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})

	t.Run("File hashes are stored", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		storeFiles := StoreFiles(normalizer)

		id := int64(7)
		hash := SAMPLE_FILE_HASH
		path := NORMALIZER_BAZEL_PATH + "/foo.go"
		file := &types.File{
			Id:          id,
			Path:        path,
			LogPathHash: SAMPLE_LOG_PATH_HASH,
			Hash:        hash,
		}

		storeFiles.Run(types.ExecLogEntry{
			File: file,
		})

		g.Expect(normalizer.GetFileHash(SAMPLE_LOG_PATH_HASH, path)).To(Equal(hash))
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})
}
