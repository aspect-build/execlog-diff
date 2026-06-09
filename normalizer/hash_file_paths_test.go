package normalizer

import (
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"

	. "github.com/onsi/gomega"
)

const SAMPLE_FILE_HASH = "jlhkuqrghfekjwhtyweurkqhjlrtuwrie3hqtfgrtsdfs"
const SAMPLE_FILE_HASH_2 = "0c4db600e4e612d86ad35b2e79c9f067"

func TestHashFilePath(t *testing.T) {
	t.Run("Does not replace if the filename is shorter", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		filePath := HashFilePaths(normalizer)

		id := int64(0)
		hash := SAMPLE_FILE_HASH
		path := "foo.bar"
		file := &types.File{
			Id:          id,
			Path:        path,
			LogPathHash: SAMPLE_LOG_PATH_HASH,
			Hash:        hash,
		}
		result := filePath.Run(types.ExecLogEntry{
			File: file,
		})

		g.Expect(result.File.Path).To(Equal(path))
		g.Expect(result.File.Id).To(Equal(id))
		g.Expect(result.File.Hash).To(Equal(hash))
		g.Expect(result.File.LogPathHash).To(Equal(SAMPLE_LOG_PATH_HASH))
		g.Expect(result.InputSet).To(BeNil())
		g.Expect(result.Spawn).To(BeNil())
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})

	t.Run("Only the file path is updated on a File", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		filePath := HashFilePaths(normalizer)

		id := int64(0)
		hash := SAMPLE_FILE_HASH
		path := NORMALIZER_BAZEL_PATH + "/foo.go"
		file := &types.File{
			Id:          id,
			Path:        path,
			LogPathHash: SAMPLE_LOG_PATH_HASH,
			Hash:        hash,
		}
		result := filePath.Run(types.ExecLogEntry{
			File: file,
		})

		g.Expect(result.File.Path).To(Equal(SAMPLE_FILE_HASH_2))
		g.Expect(result.File.Id).To(Equal(id))
		g.Expect(result.File.Hash).To(Equal(hash))
		g.Expect(result.File.LogPathHash).To(Equal(SAMPLE_LOG_PATH_HASH))
		g.Expect(result.InputSet).To(BeNil())
		g.Expect(result.Spawn).To(BeNil())
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})

	t.Run("Spawns are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		hashFilePaths := HashFilePaths(normalizer)

		spawnsAreNotUpdated(t, hashFilePaths, ctx)
	})

	t.Run("Input Sets are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		filePath := HashFilePaths(normalizer)

		inputSetsAreNotUpdated(t, filePath, ctx)
	})

	t.Run("File path hashes are stored", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		filePath := HashFilePaths(normalizer)

		id := int64(1)
		hash := SAMPLE_FILE_HASH
		path := NORMALIZER_BAZEL_PATH + "/foo.go"
		file := &types.File{
			Id:          id,
			Path:        path,
			LogPathHash: SAMPLE_LOG_PATH_HASH,
			Hash:        hash,
		}
		result := filePath.Run(types.ExecLogEntry{
			File: file,
		})

		g.Expect(result.File.Path).To(Equal(SAMPLE_FILE_HASH_2))
		g.Expect(normalizer.GetFilePathFromFilePathHash(SAMPLE_FILE_HASH_2)).To(Equal(path))
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})

	t.Run("Duplicate file paths are stored only once", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		filePath := HashFilePaths(normalizer)

		hash := SAMPLE_FILE_HASH
		path := NORMALIZER_BAZEL_PATH + "/foo.go"

		file1 := &types.File{
			Id:          0,
			Path:        path,
			LogPathHash: SAMPLE_LOG_PATH_HASH,
			Hash:        hash,
		}
		result := filePath.Run(types.ExecLogEntry{
			File: file1,
		})

		g.Expect(result.File.Path).To(Equal(SAMPLE_FILE_HASH_2))
		g.Expect(normalizer.GetFilePathFromFilePathHash(SAMPLE_FILE_HASH_2)).To(Equal(path))

		file2 := &types.File{
			Id:          0,
			Path:        path,
			LogPathHash: SAMPLE_LOG_PATH_HASH,
			Hash:        hash,
		}
		result = filePath.Run(types.ExecLogEntry{
			File: file2,
		})

		g.Expect(result.File.Path).To(Equal(SAMPLE_FILE_HASH_2))
		g.Expect(normalizer.GetFilePathFromFilePathHash(SAMPLE_FILE_HASH_2)).To(Equal(path))

		counter := 0
		normalizer.hash_to_file_path_map.Range(func(key, value any) bool {
			counter++
			return true
		})
		g.Expect(counter).To(Equal(1))
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})

	t.Run("Handles hash clashes", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		otherPath := NORMALIZER_BAZEL_PATH + "/other_file.go"
		normalizer.hash_to_file_path_map.Store(SAMPLE_FILE_HASH_2, otherPath)

		filePath := HashFilePaths(normalizer)

		hash := SAMPLE_FILE_HASH
		path := NORMALIZER_BAZEL_PATH + "/foo.go"

		file := &types.File{
			Id:          0,
			Path:        path,
			LogPathHash: SAMPLE_LOG_PATH_HASH,
			Hash:        hash,
		}
		result := filePath.Run(types.ExecLogEntry{
			File: file,
		})

		g.Expect(result.File.Path).To(Equal(SAMPLE_FILE_HASH_2 + "_"))
		g.Expect(normalizer.GetFilePathFromFilePathHash(SAMPLE_FILE_HASH_2)).To(Equal(otherPath))
		g.Expect(normalizer.GetFilePathFromFilePathHash(SAMPLE_FILE_HASH_2 + "_")).To(Equal(path))
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})
}
