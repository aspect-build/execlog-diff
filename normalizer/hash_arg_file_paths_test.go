package normalizer

import (
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"

	. "github.com/onsi/gomega"
)

const NORMALIZER_BAZEL_PATH = "bazel-out/darwin_arm64-fastbuild/bin/normalizer"
const SAMPLE_LOG_PATH_HASH = "dsafgwerfadsfgewg"

func TestArgFilePaths(t *testing.T) {
	t.Run("Spawn file paths are updated", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		filePath := HashFilePaths(normalizer)

		filePath.Run(types.ExecLogEntry{
			File: &types.File{
				Id:          0,
				Path:        NORMALIZER_BAZEL_PATH + "/foo.bar",
				LogPathHash: SAMPLE_LOG_PATH_HASH,
				Hash:        SAMPLE_FILE_HASH,
			},
		})

		filePath.Run(types.ExecLogEntry{
			File: &types.File{
				Id:          1,
				Path:        NORMALIZER_BAZEL_PATH + "/baz.txt",
				LogPathHash: SAMPLE_LOG_PATH_HASH,
				Hash:        SAMPLE_FILE_HASH,
			},
		})

		filePath.Run(types.ExecLogEntry{
			File: &types.File{
				Id:          1,
				Path:        NORMALIZER_BAZEL_PATH + "/foobarbaz.txt",
				LogPathHash: SAMPLE_LOG_PATH_HASH,
				Hash:        SAMPLE_FILE_HASH,
			},
		})

		hashAreFilePaths := HashArgFilePaths(normalizer)

		id := int64(0)
		logPath := "someLogPath"
		hash := SAMPLE_FILE_HASH
		args := []string{
			NORMALIZER_BAZEL_PATH + "/foobarbaz.txt=",
			NORMALIZER_BAZEL_PATH + "/foo.bar",
			"something=" + NORMALIZER_BAZEL_PATH + "/baz.txt",
			"=" + NORMALIZER_BAZEL_PATH + "/foobarbaz.txt=",
			"=" + NORMALIZER_BAZEL_PATH + "/foobarbaz.txt",
		}
		env := make(map[string]string)
		env["key1"] = "val1"
		env["key2"] = "val2"
		envHash := "wqedferagfewytgeawf"
		label := "//foo/bar:baz"
		inputFiles := []string{"foo"}
		outputFiles := []string{"bar"}
		inputSetId := int64(7)
		outputFileIds := []int64{1, 2, 3}
		outputSetIds := []int64{4, 5, 6}
		spawn := &types.IncomingSpawn{
			Id:             id,
			LogPath:        logPath,
			LogPathHash:    SAMPLE_LOG_PATH_HASH,
			Args:           args,
			Env:            env,
			EnvHash:        envHash,
			Label:          label,
			CacheHit:       true,
			Remoteable:     true,
			Cacheable:      true,
			RemoteCachable: true,
			Hash:           hash,
			InputFiles:     inputFiles,
			OutputFiles:    outputFiles,
			InputSetId:     inputSetId,
			OutputFileIds:  outputFileIds,
			OutputSetIds:   outputSetIds,
		}

		result := hashAreFilePaths.Run(types.ExecLogEntry{
			Spawn: spawn,
		})

		g.Expect(result.Spawn.Id).To(Equal(id))
		g.Expect(result.Spawn.LogPath).To(Equal(logPath))
		g.Expect(result.Spawn.LogPathHash).To(Equal(SAMPLE_LOG_PATH_HASH))
		g.Expect(len(result.Spawn.Args)).To(Equal(5))
		g.Expect(result.Spawn.Args[0]).To(Equal(NORMALIZER_BAZEL_PATH + "/foobarbaz.txt="))
		g.Expect(result.Spawn.Args[1]).To(Equal("5e72ac610ed9df49721cf2d92fc0b790"))
		g.Expect(result.Spawn.Args[2]).To(Equal("something=251a88f2dfa5ecfc14003eb9fcec769f"))
		g.Expect(result.Spawn.Args[3]).To(Equal("=" + NORMALIZER_BAZEL_PATH + "/foobarbaz.txt="))
		g.Expect(result.Spawn.Args[4]).To(Equal("=" + NORMALIZER_BAZEL_PATH + "/foobarbaz.txt"))
		g.Expect(result.Spawn.Env).To(Equal(env))
		g.Expect(result.Spawn.EnvHash).To(Equal(envHash))
		g.Expect(result.Spawn.Label).To(Equal(label))
		g.Expect(result.Spawn.CacheHit).To(BeTrue())
		g.Expect(result.Spawn.Remoteable).To(BeTrue())
		g.Expect(result.Spawn.Cacheable).To(BeTrue())
		g.Expect(result.Spawn.RemoteCachable).To(BeTrue())
		g.Expect(result.Spawn.Hash).To(Equal(hash))
		g.Expect(result.Spawn.InputFiles).To(Equal(inputFiles))
		g.Expect(result.Spawn.OutputFiles).To(Equal(outputFiles))
		g.Expect(result.Spawn.OutputFileIds).To(Equal(outputFileIds))
		g.Expect(result.Spawn.OutputSetIds).To(Equal(outputSetIds))
		g.Expect(result.InputSet).To(BeNil())
		g.Expect(result.File).To(BeNil())
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})

	t.Run("Spawn is not updated when path names are short", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		hashArgFilePaths := HashArgFilePaths(normalizer)

		hashArgFilePaths.Run(types.ExecLogEntry{
			File: &types.File{
				Id:          0,
				Path:        "foo.bar",
				LogPathHash: SAMPLE_LOG_PATH_HASH,
				Hash:        SAMPLE_FILE_HASH,
			},
		})

		hashArgFilePaths.Run(types.ExecLogEntry{
			File: &types.File{
				Id:          1,
				Path:        "baz.txt",
				LogPathHash: SAMPLE_LOG_PATH_HASH,
				Hash:        SAMPLE_FILE_HASH,
			},
		})

		hashArgFilePaths.Run(types.ExecLogEntry{
			File: &types.File{
				Id:          1,
				Path:        "foobarbaz.txt",
				LogPathHash: SAMPLE_LOG_PATH_HASH,
				Hash:        SAMPLE_FILE_HASH,
			},
		})

		argFilePaths := HashArgFilePaths(normalizer)

		id := int64(0)
		logPath := "someLogPath"
		hash := SAMPLE_FILE_HASH
		args := []string{"foobarbaz.txt=", "foo.bar", "something=baz.txt", "=foobarbaz.txt="}
		env := make(map[string]string)
		env["key1"] = "val1"
		env["key2"] = "val2"
		envHash := "wqedferagfewytgeawf"
		label := "//foo/bar:baz"
		inputFiles := []string{"foo"}
		outputFiles := []string{"bar"}
		inputSetId := int64(7)
		outputFileIds := []int64{1, 2, 3}
		outputSetIds := []int64{4, 5, 6}
		spawn := &types.IncomingSpawn{
			Id:             id,
			LogPath:        logPath,
			LogPathHash:    SAMPLE_LOG_PATH_HASH,
			Args:           args,
			Env:            env,
			EnvHash:        envHash,
			Label:          label,
			CacheHit:       true,
			Remoteable:     true,
			Cacheable:      true,
			RemoteCachable: true,
			Hash:           hash,
			InputFiles:     inputFiles,
			OutputFiles:    outputFiles,
			InputSetId:     inputSetId,
			OutputFileIds:  outputFileIds,
			OutputSetIds:   outputSetIds,
		}

		result := argFilePaths.Run(types.ExecLogEntry{
			Spawn: spawn,
		})

		g.Expect(result.Spawn.Id).To(Equal(id))
		g.Expect(result.Spawn.LogPath).To(Equal(logPath))
		g.Expect(result.Spawn.LogPathHash).To(Equal(SAMPLE_LOG_PATH_HASH))
		g.Expect(result.Spawn.Args).To(Equal(args))
		g.Expect(result.Spawn.Env).To(Equal(env))
		g.Expect(result.Spawn.EnvHash).To(Equal(envHash))
		g.Expect(result.Spawn.Label).To(Equal(label))
		g.Expect(result.Spawn.CacheHit).To(BeTrue())
		g.Expect(result.Spawn.Remoteable).To(BeTrue())
		g.Expect(result.Spawn.Cacheable).To(BeTrue())
		g.Expect(result.Spawn.RemoteCachable).To(BeTrue())
		g.Expect(result.Spawn.Hash).To(Equal(hash))
		g.Expect(result.Spawn.InputFiles).To(Equal(inputFiles))
		g.Expect(result.Spawn.OutputFiles).To(Equal(outputFiles))
		g.Expect(result.Spawn.OutputFileIds).To(Equal(outputFileIds))
		g.Expect(result.Spawn.OutputSetIds).To(Equal(outputSetIds))
		g.Expect(result.InputSet).To(BeNil())
		g.Expect(result.File).To(BeNil())
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})

	t.Run("Files are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		hashArgFilePaths := HashArgFilePaths(normalizer)

		filesAreNotUpdated(t, hashArgFilePaths, ctx)
	})

	t.Run("Input Sets are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		hashArgFilePaths := HashArgFilePaths(normalizer)

		inputSetsAreNotUpdated(t, hashArgFilePaths, ctx)
	})
}
