package normalizer

import (
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"

	. "github.com/onsi/gomega"
)

func TestHashEnv(t *testing.T) {
	t.Run("Files are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		credentials := HashEnv(normalizer)

		filesAreNotUpdated(t, credentials, ctx)
	})

	t.Run("Input Sets are not updated", func(t *testing.T) {
		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		credentials := HashEnv(normalizer)

		inputSetsAreNotUpdated(t, credentials, ctx)
	})

	t.Run("Only Spawn environment hash is updated", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		credentials := HashEnv(normalizer)

		id := int64(0)
		logPath := "someLogPath"
		hash := SAMPLE_FILE_HASH
		args := []string{"some", "command"}
		env := make(map[string]string)
		env["key1"] = "val1"
		env["key2"] = "val2"
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
			EnvHash:        "",
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
		g.Expect(result.Spawn.Env["key1"]).To(Equal("val1"))
		g.Expect(result.Spawn.Env["key2"]).To(Equal("val2"))
		g.Expect(result.Spawn.EnvHash).To(Equal("ab293b0819bf8a23bfac3cd0dfdc1b2546253825b998657598b2d812eadfcea6310fd6320590fc55905e2719043c3d85"))
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

	t.Run("2 Spawns with the same env share a map entry", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		credentials := HashEnv(normalizer)

		env := make(map[string]string)
		env["key1"] = "val1"
		env["key2"] = "val2"
		env["key3"] = "val3"
		result1 := credentials.Run(types.ExecLogEntry{
			Spawn: &types.IncomingSpawn{
				Id:      1,
				Env:     env,
				EnvHash: "",
			},
		})

		env2 := make(map[string]string)
		env2["key1"] = "val1"
		env2["key2"] = "val2"
		env2["key3"] = "val3"
		result2 := credentials.Run(types.ExecLogEntry{
			Spawn: &types.IncomingSpawn{
				Id:      2,
				Env:     env2,
				EnvHash: "",
			},
		})

		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))

		g.Expect(result1.Spawn.Id).To(Equal(int64(1)))
		g.Expect(len(result1.Spawn.Env)).To(Equal(3))
		g.Expect(result1.Spawn.Env["key1"]).To(Equal("val1"))
		g.Expect(result1.Spawn.Env["key2"]).To(Equal("val2"))
		g.Expect(result1.Spawn.Env["key3"]).To(Equal("val3"))
		g.Expect(result1.Spawn.EnvHash).To(Equal("fcc5a850cce5e7f1a9d88a2c6c848941eb2f8efe94592d6dc8d2aaf47ff75b5c73ac23ccaf1434d8028d42a23e547cf2"))
		g.Expect(result1.InputSet).To(BeNil())
		g.Expect(result1.File).To(BeNil())

		g.Expect(result2.Spawn.Id).To(Equal(int64(2)))
		g.Expect(len(result2.Spawn.Env)).To(Equal(3))
		g.Expect(result2.Spawn.Env["key1"]).To(Equal("val1"))
		g.Expect(result2.Spawn.Env["key2"]).To(Equal("val2"))
		g.Expect(result2.Spawn.Env["key3"]).To(Equal("val3"))
		g.Expect(result2.Spawn.EnvHash).To(Equal("fcc5a850cce5e7f1a9d88a2c6c848941eb2f8efe94592d6dc8d2aaf47ff75b5c73ac23ccaf1434d8028d42a23e547cf2"))
		g.Expect(result2.InputSet).To(BeNil())
		g.Expect(result2.File).To(BeNil())

		counter := 0
		normalizer.hash_to_env_map.Range(func(key, value any) bool {
			counter++
			return true
		})
		g.Expect(counter).To(Equal(1))

		storedEnv := normalizer.getEnv("fcc5a850cce5e7f1a9d88a2c6c848941eb2f8efe94592d6dc8d2aaf47ff75b5c73ac23ccaf1434d8028d42a23e547cf2")
		g.Expect(len(storedEnv)).To(Equal(len(env)))
		for key, val := range storedEnv {
			g.Expect(env[key]).To(Equal(val))
		}
	})

	t.Run("2 Spawns with different envs do not share the same map entry", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := New(ctx)

		credentials := HashEnv(normalizer)

		env := make(map[string]string)
		env["key1"] = "val1"
		env["key2"] = "val2"
		env["key3"] = "val3"
		result1 := credentials.Run(types.ExecLogEntry{
			Spawn: &types.IncomingSpawn{
				Id:      1,
				Env:     env,
				EnvHash: "",
			},
		})

		env2 := make(map[string]string)
		env2["key1"] = "val1"
		env2["key2"] = "val2"
		env2["key3"] = "otherVal3"
		env2["key4"] = "val4"
		result2 := credentials.Run(types.ExecLogEntry{
			Spawn: &types.IncomingSpawn{
				Id:      2,
				Env:     env2,
				EnvHash: "",
			},
		})

		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))

		g.Expect(result1.Spawn.Id).To(Equal(int64(1)))
		g.Expect(len(result1.Spawn.Env)).To(Equal(3))
		g.Expect(result1.Spawn.Env["key1"]).To(Equal("val1"))
		g.Expect(result1.Spawn.Env["key2"]).To(Equal("val2"))
		g.Expect(result1.Spawn.Env["key3"]).To(Equal("val3"))
		g.Expect(result1.Spawn.EnvHash).To(Equal("fcc5a850cce5e7f1a9d88a2c6c848941eb2f8efe94592d6dc8d2aaf47ff75b5c73ac23ccaf1434d8028d42a23e547cf2"))
		g.Expect(result1.InputSet).To(BeNil())
		g.Expect(result1.File).To(BeNil())

		g.Expect(result2.Spawn.Id).To(Equal(int64(2)))
		g.Expect(len(result2.Spawn.Env)).To(Equal(4))
		g.Expect(result2.Spawn.Env["key1"]).To(Equal("val1"))
		g.Expect(result2.Spawn.Env["key2"]).To(Equal("val2"))
		g.Expect(result2.Spawn.Env["key3"]).To(Equal("otherVal3"))
		g.Expect(result2.Spawn.Env["key4"]).To(Equal("val4"))
		g.Expect(result2.Spawn.EnvHash).To(Equal("8db44cb91a8a6b23496e1d9aafffe79aa9970095c30971b77dcf8f25d9556f6a818b3ca016063fbdb3f1baea93bfcfad"))
		g.Expect(result2.InputSet).To(BeNil())
		g.Expect(result2.File).To(BeNil())

		counter := 0
		normalizer.hash_to_env_map.Range(func(key, value any) bool {
			counter++
			return true
		})
		g.Expect(counter).To(Equal(2))

		storedEnv := normalizer.getEnv("fcc5a850cce5e7f1a9d88a2c6c848941eb2f8efe94592d6dc8d2aaf47ff75b5c73ac23ccaf1434d8028d42a23e547cf2")
		g.Expect(len(storedEnv)).To(Equal(len(env)))
		for key, val := range storedEnv {
			g.Expect(env[key]).To(Equal(val))
		}

		storedEnv2 := normalizer.getEnv("8db44cb91a8a6b23496e1d9aafffe79aa9970095c30971b77dcf8f25d9556f6a818b3ca016063fbdb3f1baea93bfcfad")
		g.Expect(len(storedEnv2)).To(Equal(len(env2)))
		for key, val := range storedEnv2 {
			g.Expect(env2[key]).To(Equal(val))
		}
	})
}
