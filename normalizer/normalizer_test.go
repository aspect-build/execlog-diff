package normalizer

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/converter"
	"github.com/aspect-build/execlog-diff/test-data/helpers"
	"github.com/aspect-build/execlog-diff/types"
	"github.com/aspect-build/execlog-diff/diagnostic"
	"gotest.tools/v3/golden"

	. "github.com/onsi/gomega"
)

func TestNormalizer(t *testing.T) {
	t.Run("Normalizer stops when told to by the context class", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)
		normalizer := New(ctx)

		execlogChan := make(chan types.ExecLogEntry, 100)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		execlogChan <- types.ExecLogEntry{
			File: &types.File{
				Id:          1,
				Path:        "foo",
				LogPathHash: "someLogPathHash",
				Hash:        SAMPLE_FILE_HASH,
			},
		}

		execlogChan <- types.ExecLogEntry{
			InputSet: &types.InputSet{
				Id:              2,
				FileIds:         []int64{1},
				NestedInputSets: []int64{},
				LogPathHash:     "someLogPathHash",
			},
		}

		execlogChan <- types.ExecLogEntry{
			File: &types.File{
				Id:          3,
				Path:        "bar",
				LogPathHash: "someLogPathHash",
				Hash:        "sadgefwyrtasgehwtrghwetyrgwergtewrg",
			},
		}

		execlogChan <- types.ExecLogEntry{
			InputSet: &types.InputSet{
				Id:              4,
				FileIds:         []int64{1},
				NestedInputSets: []int64{},
				LogPathHash:     "someLogPathHash",
			},
		}

		execlogChan <- types.ExecLogEntry{
			Spawn: &types.IncomingSpawn{
				Id:            5,
				LogPathHash:   "someLogPathHash",
				Args:          make([]string, 0),
				Env:           make(map[string]string),
				Label:         "//some/label",
				Hash:          "",
				OutputFiles:   make([]string, 0),
				InputSetId:    2,
				OutputFileIds: make([]int64, 0),
				OutputSetIds:  []int64{4},
			},
		}

		time.Sleep(time.Millisecond * 200)
		ctx.Error(&diagnostic.Diagnostic{Message: "Stop!"})

		execlogChan <- types.ExecLogEntry{
			Spawn: &types.IncomingSpawn{
				Id:            6,
				LogPathHash:   "someLogPathHash",
				Args:          make([]string, 0),
				Env:           make(map[string]string),
				Label:         "//some/label",
				Hash:          "",
				OutputFiles:   make([]string, 0),
				InputSetId:    2,
				OutputFileIds: make([]int64, 0),
				OutputSetIds:  []int64{4},
			},
		}

		execlogChan <- types.ExecLogEntry{
			Spawn: &types.IncomingSpawn{
				Id:            7,
				LogPathHash:   "someLogPathHash",
				Args:          make([]string, 0),
				Env:           make(map[string]string),
				Label:         "//some/label",
				Hash:          "",
				OutputFiles:   make([]string, 0),
				InputSetId:    2,
				OutputFileIds: make([]int64, 0),
				OutputSetIds:  []int64{4},
			},
		}

		execlogChan <- types.ExecLogEntry{
			Spawn: &types.IncomingSpawn{
				Id:            8,
				LogPathHash:   "someLogPathHash",
				Args:          make([]string, 0),
				Env:           make(map[string]string),
				Label:         "//some/label",
				Hash:          "",
				OutputFiles:   make([]string, 0),
				InputSetId:    2,
				OutputFileIds: make([]int64, 0),
				OutputSetIds:  []int64{4},
			},
		}

		time.Sleep(time.Millisecond * 1000)

		normalizerWait.Wait()

		g.Expect(len(spawnExecLogChan)).To(Equal(1))
	})

	t.Run("Spawns are created", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)
		normalizer := New(ctx)

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		execlogChan <- types.ExecLogEntry{
			File: &types.File{
				Id:          1,
				Path:        "foo",
				LogPathHash: SAMPLE_LOG_PATH_HASH,
				Hash:        SAMPLE_FILE_HASH,
			},
		}

		execlogChan <- types.ExecLogEntry{
			File: &types.File{
				Id:          2,
				Path:        "bar",
				LogPathHash: SAMPLE_LOG_PATH_HASH,
				Hash:        SAMPLE_FILE_HASH,
			},
		}

		execlogChan <- types.ExecLogEntry{
			File: &types.File{
				Id:          3,
				Path:        "other",
				LogPathHash: SAMPLE_LOG_PATH_HASH,
				Hash:        "asdfsadf",
			},
		}

		execlogChan <- types.ExecLogEntry{
			InputSet: &types.InputSet{
				Id:              4,
				FileIds:         []int64{1},
				NestedInputSets: []int64{},
				LogPathHash:     SAMPLE_LOG_PATH_HASH,
			},
		}

		execlogChan <- types.ExecLogEntry{
			InputSet: &types.InputSet{
				Id:              5,
				FileNames:       []string{"baz"},
				FileIds:         []int64{2},
				NestedInputSets: []int64{4},
				LogPathHash:     SAMPLE_LOG_PATH_HASH,
			},
		}

		logPath := "someLogPath"
		hash := "9ffd5df9b67abe4ef6fa321c96b2da38ab9f4aa17fe1aa8d5a91c4b37855425e352e48c099ec0bcf046a143da49b2ee6"
		args := []string{"some", "command"}
		env := make(map[string]string)
		env["key1"] = "val1"
		env["key2"] = "val2"
		label := "//foo/bar:baz"
		outputFiles := []string{"someOther"}
		inputSetId := int64(4)
		outputFileIds := []int64{3}
		outputSetIds := []int64{5}

		execlogChan <- types.ExecLogEntry{
			Spawn: &types.IncomingSpawn{
				Id:             6,
				LogPath:        logPath,
				LogPathHash:    SAMPLE_LOG_PATH_HASH,
				Args:           args,
				Env:            env,
				Label:          label,
				CacheHit:       true,
				Remoteable:     true,
				Cacheable:      true,
				RemoteCachable: true,
				Hash:           hash,
				OutputFiles:    outputFiles,
				InputSetId:     inputSetId,
				OutputFileIds:  outputFileIds,
				OutputSetIds:   outputSetIds,
			},
		}
		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		spawnList := helpers.CreateListFromSpawnChan(spawnExecLogChan, t)

		g.Expect(len(spawnList)).To(Equal(1))
		g.Expect(spawnList[0].LogPath).To(Equal(logPath))
		g.Expect(spawnList[0].LogPathHash).To(Equal(SAMPLE_LOG_PATH_HASH))
		g.Expect(spawnList[0].Args).To(Equal(args))
		g.Expect(spawnList[0].EnvHash).To(Equal("732b1733d27de74f98c454331b3663ba073d802784812f3f2a38fda7d4ed16d38728c65a1f8a0fb22bf4e3403dc951d8"))
		g.Expect(spawnList[0].Label).To(Equal(label))
		g.Expect(spawnList[0].CacheHit).To(BeTrue())
		g.Expect(spawnList[0].Remoteable).To(BeTrue())
		g.Expect(spawnList[0].Cacheable).To(BeTrue())
		g.Expect(spawnList[0].RemoteCachable).To(BeTrue())
		g.Expect(spawnList[0].Hash).To(Equal(hash))
		g.Expect(len(spawnList[0].InputFiles)).To(Equal(1))
		g.Expect(spawnList[0].InputFiles[0]).To(Equal("foo"))
		g.Expect(len(spawnList[0].OutputFiles)).To(Equal(5))
		g.Expect(spawnList[0].OutputFiles[0]).To(Equal("bar"))
		g.Expect(spawnList[0].OutputFiles[1]).To(Equal("baz"))
		g.Expect(spawnList[0].OutputFiles[2]).To(Equal("foo"))
		g.Expect(spawnList[0].OutputFiles[3]).To(Equal("other"))
		g.Expect(spawnList[0].OutputFiles[4]).To(Equal("someOther"))
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})

	t.Run("ListLabels mode creates simplified spawns with minimal fields", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.ListLabels)

		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)
		normalizer := New(ctx)

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		label := "//foo/bar:baz"

		execlogChan <- types.ExecLogEntry{
			Spawn: &types.IncomingSpawn{
				Id:             6,
				LogPath:        "someLogPath",
				LogPathHash:    SAMPLE_LOG_PATH_HASH,
				Args:           []string{"some", "command"},
				Env:            make(map[string]string),
				Label:          label,
				CacheHit:       true,
				Remoteable:     true,
				Cacheable:      true,
				RemoteCachable: true,
				Hash:           "",
				OutputFiles:    []string{},
				InputFiles:     []string{},
				InputSetId:     4,
				OutputFileIds:  []int64{},
				OutputSetIds:   []int64{},
			},
		}
		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		spawnList := helpers.CreateListFromSpawnChan(spawnExecLogChan, t)

		g.Expect(len(spawnList)).To(Equal(1))

		// ListLabels mode should only include minimal fields
		g.Expect(spawnList[0].Label).To(Equal(label))
		g.Expect(spawnList[0].CacheHit).To(BeTrue())
		g.Expect(spawnList[0].Remoteable).To(BeTrue())
		g.Expect(spawnList[0].Cacheable).To(BeTrue())
		g.Expect(spawnList[0].RemoteCachable).To(BeTrue())
		g.Expect(spawnList[0].Hash).To(Equal("e67c542ddadd760f3908a6493e94e46f250a7cdabde4b34235cbe36cd43d9ad7ea47148a74de43a2364729057c99e65e"))

		// Fields that should be empty/zero in ListLabels mode
		g.Expect(spawnList[0].LogPath).To(Equal(""))
		g.Expect(spawnList[0].LogPathHash).To(Equal(""))
		g.Expect(len(spawnList[0].Args)).To(Equal(0))
		g.Expect(spawnList[0].EnvHash).To(Equal(""))
		g.Expect(len(spawnList[0].InputFiles)).To(Equal(0))
		g.Expect(len(spawnList[0].OutputFiles)).To(Equal(0))
		g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
	})

	t.Run("golden files", func(t *testing.T) {
		files := map[string]string{
			"golden_files/changed_directory.txt": "../test-data/changed_directory_7.3.0_1.compact.log",
			"golden_files/changed_file.txt":      "../test-data/changed_file_7.3.0_1.compact.log",
		}

		// Run this 1000 times to confirm the result is static
		for i := 0; i < 1000; i++ {
			for golden_filepath, execlog_filepath := range files {
				g := NewGomegaWithT(t)

				var stderr strings.Builder
				ctx := context.New(&stderr, context.Diff)

				execLogChan := make(chan types.ExecLogEntry, 100)
				spawnChan := make(chan types.Spawn, 100)
				var compactWait sync.WaitGroup
				compactWait.Add(1)
				go converter.Convert(execlog_filepath, execLogChan, &compactWait, ctx)

				var normalizerWait sync.WaitGroup
				normalizerWait.Add(1)
				normalizer := New(ctx)
				go normalizer.Normalize(execLogChan, spawnChan, &normalizerWait)

				compactWait.Wait()
				close(execLogChan)
				normalizerWait.Wait()
				close(spawnChan)

				g.Expect(stderr.String()).To(Equal(""))
				g.Expect(ctx.Continue()).To(BeTrue())
				g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))

				output := helpers.CreateGoldenFileFromSpawnChan(spawnChan, t)
				golden.Assert(t, output, golden_filepath)
			}
		}
	})

	t.Run("lookupId", func(t *testing.T) {
		g := NewGomegaWithT(t)

		// int and int64 literals must produce the same MapKey (int is normalized to int64)
		g.Expect(lookupId("foo", int(123))).To(Equal(MapKey{LogPathHash: "foo", Key: int64(123)}))
		g.Expect(lookupId("foo", int64(123))).To(Equal(MapKey{LogPathHash: "foo", Key: int64(123)}))
		g.Expect(lookupId("foo", string("bar"))).To(Equal(MapKey{LogPathHash: "foo", Key: string("bar")}))
	})

	t.Run("nestedLookupId", func(t *testing.T) {
		g := NewGomegaWithT(t)

		g.Expect(nestedLookupId("foo", int(123))).To(Equal("nested:foo:123"))
		g.Expect(nestedLookupId("foo", int64(123))).To(Equal("nested:foo:123"))
		g.Expect(nestedLookupId("foo", string("bar"))).To(Equal("nested:foo:bar"))
	})

	t.Run("getInputsFromInputSet", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)
		normalizer := New(ctx)

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		execlogChan <- types.ExecLogEntry{
			File: &types.File{
				Id:          1,
				Path:        "foo",
				LogPathHash: SAMPLE_LOG_PATH_HASH,
				Hash:        SAMPLE_FILE_HASH,
			},
		}

		execlogChan <- types.ExecLogEntry{
			File: &types.File{
				Id:          2,
				Path:        "bar",
				LogPathHash: SAMPLE_LOG_PATH_HASH,
				Hash:        SAMPLE_FILE_HASH,
			},
		}

		execlogChan <- types.ExecLogEntry{
			File: &types.File{
				Id:          3,
				Path:        "other",
				LogPathHash: SAMPLE_LOG_PATH_HASH,
				Hash:        "asdfsadf",
			},
		}

		execlogChan <- types.ExecLogEntry{
			InputSet: &types.InputSet{
				Id:              4,
				FileIds:         []int64{1},
				NestedInputSets: []int64{},
				LogPathHash:     SAMPLE_LOG_PATH_HASH,
			},
		}

		execlogChan <- types.ExecLogEntry{
			InputSet: &types.InputSet{
				Id:              5,
				FileIds:         []int64{2},
				NestedInputSets: []int64{4},
				LogPathHash:     SAMPLE_LOG_PATH_HASH,
			},
		}

		execlogChan <- types.ExecLogEntry{
			InputSet: &types.InputSet{
				Id:              6,
				FileNames:       []string{"baz"},
				FileIds:         []int64{2, 3},
				NestedInputSets: []int64{4, 5},
				LogPathHash:     SAMPLE_LOG_PATH_HASH,
			},
		}

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		g.Expect(len(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, 4)))).To(Equal(1))
		g.Expect(len(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, 5)))).To(Equal(2))
		g.Expect(len(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, 6)))).To(Equal(4))

		// make sure files are always returned in the same order
		for i := 0; i < 1000; i++ {
			g.Expect(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, 4))[0]).To(Equal("foo"))

			g.Expect(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, 5))[0]).To(Equal("bar"))
			g.Expect(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, 5))[1]).To(Equal("foo"))

			g.Expect(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, 6))[0]).To(Equal("bar"))
			g.Expect(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, 6))[1]).To(Equal("baz"))
			g.Expect(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, 6))[2]).To(Equal("foo"))
			g.Expect(normalizer.getInputsFromInputSet(lookupId(SAMPLE_LOG_PATH_HASH, 6))[3]).To(Equal("other"))
		}
	})
}

func filesAreNotUpdated(t *testing.T, normalizer normalizer, ctx *context.Context) {
	g := NewGomegaWithT(t)

	id := int64(0)
	hash := SAMPLE_FILE_HASH
	path := NORMALIZER_BAZEL_PATH + "/foo.go"
	file := &types.File{
		Id:          id,
		Path:        path,
		LogPathHash: SAMPLE_LOG_PATH_HASH,
		Hash:        hash,
	}
	result := normalizer.Run(types.ExecLogEntry{
		File: file,
	})

	g.Expect(result.File.Path).To(Equal(path))
	g.Expect(result.File.Id).To(Equal(id))
	g.Expect(result.File.Hash).To(Equal(hash))
	g.Expect(result.File.LogPathHash).To(Equal(SAMPLE_LOG_PATH_HASH))
	g.Expect(result.InputSet).To(BeNil())
	g.Expect(result.Spawn).To(BeNil())
	g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
}

func inputSetsAreNotUpdated(t *testing.T, normalizer normalizer, ctx *context.Context) {
	g := NewGomegaWithT(t)

	id := int64(0)
	fileNames := []string{"foo"}
	fileIds := []int64{1, 2, 3}
	nestedInputSets := []int64{4, 5, 6}
	inputSet := &types.InputSet{
		Id:              id,
		FileNames:       fileNames,
		FileIds:         fileIds,
		NestedInputSets: nestedInputSets,
		LogPathHash:     SAMPLE_LOG_PATH_HASH,
	}

	result := normalizer.Run(types.ExecLogEntry{
		InputSet: inputSet,
	})

	g.Expect(result.InputSet.Id).To(Equal(id))
	g.Expect(result.InputSet.LogPathHash).To(Equal(SAMPLE_LOG_PATH_HASH))
	g.Expect(result.InputSet.FileNames).To(Equal(fileNames))
	g.Expect(result.InputSet.FileIds).To(Equal(fileIds))
	g.Expect(result.InputSet.NestedInputSets).To(Equal(nestedInputSets))
	g.Expect(result.File).To(BeNil())
	g.Expect(result.Spawn).To(BeNil())
	g.Expect(len(ctx.GetDiagnostics())).To(Equal(0))
}

func spawnsAreNotUpdated(t *testing.T, normalizer normalizer, ctx *context.Context) {
	g := NewGomegaWithT(t)

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

	result := normalizer.Run(types.ExecLogEntry{
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
}
