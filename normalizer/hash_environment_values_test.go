package normalizer

import (
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"

	. "github.com/onsi/gomega"
)

func TestCredentials(t *testing.T) {
	t.Run("Files are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		credentials := HashEnvironmentValues(normalizer)

		filesAreNotUpdated(t, credentials, ctx)
	})

	t.Run("Input Sets are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		credentials := HashEnvironmentValues(normalizer)

		inputSetsAreNotUpdated(t, credentials, ctx)
	})

	t.Run("Only spawn environment variables are updated", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		credentials := HashEnvironmentValues(normalizer)

		id := int64(0)
		logPath := "someLogPath"
		hash := SAMPLE_FILE_HASH
		args := []string{"some", "command"}
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

		result := credentials.Run(types.ExecLogEntry{
			Spawn: spawn,
		})

		g.Expect(result.Spawn.Id).To(Equal(id))
		g.Expect(result.Spawn.LogPath).To(Equal(logPath))
		g.Expect(result.Spawn.LogPathHash).To(Equal(SAMPLE_LOG_PATH_HASH))
		g.Expect(result.Spawn.Args).To(Equal(args))
		g.Expect(len(result.Spawn.Env)).To(Equal(2))
		g.Expect(result.Spawn.Env["key1"]).To(Equal("7d9d8e4825eeb4202b5dfe3dca79b52284b973e53c00150735a5ad7979ce09cfe80e4175641a3e61796e783734b70ec9"))
		g.Expect(result.Spawn.Env["key2"]).To(Equal("e12d3a46f80d47ee06a014ce4d84bce5722e871518315d469a69f846b738eba0bdc934fb25f076980228e96878ebb150"))
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
}
