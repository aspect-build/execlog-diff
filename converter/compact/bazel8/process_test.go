package bazel8

import (
	"strings"
	"sync"
	"testing"

	"github.com/aspect-build/execlog-diff/third_party/github.com/bazelbuild/bazel8/src/main/protobuf/spawn"
	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/test-data/helpers"
	"github.com/aspect-build/execlog-diff/types"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestProcess(t *testing.T) {
	t.Run("unpacks file messages correctly", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		execLogChan := make(chan types.ExecLogEntry, 10)
		protoMessageChan := make(chan protoreflect.ProtoMessage, 10)
		var wait sync.WaitGroup

		id := 7
		logPath := "some_log_path"
		filePath := "some_file_path"
		hash := "adsafaefgewfhdyerytgfedasgtgewg"

		execLogEntry := &spawn.ExecLogEntry{
			Id: uint32(id),
			Type: &spawn.ExecLogEntry_File_{
				File: &spawn.ExecLogEntry_File{
					Path: filePath,
					Digest: &spawn.Digest{
						Hash:      hash,
						SizeBytes: 1234,
					},
				},
			},
		}

		processor := newProcessor(logPath, protoMessageChan, execLogChan, &wait, ctx)
		processor.processFile(execLogEntry)

		g.Expect(ctx.Continue()).To(BeTrue())
		diagnostics := ctx.GetDiagnostics()

		g.Expect(len(diagnostics)).To(Equal(0))
		g.Expect(len(execLogChan)).To(Equal(1))
		g.Expect(stderr.String()).To(Equal(""))

		close(execLogChan)
		close(protoMessageChan)

		execLogList := helpers.CreateListFromExecLogChan(execLogChan, t)

		g.Expect(execLogList[0].InputSet).To(BeNil())
		g.Expect(execLogList[0].Spawn).To(BeNil())
		g.Expect(execLogList[0].File.Id).To(Equal(int64(id)))
		g.Expect(execLogList[0].File.Path).To(Equal(filePath))
		g.Expect(execLogList[0].File.LogPathHash).To(Equal("c2c87929d047a92999dc11b393cd2bda"))
		g.Expect(execLogList[0].File.Hash).To(Equal(hash))
	})

	t.Run("unpacks directory messages correctly", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		execLogChan := make(chan types.ExecLogEntry, 10)
		protoMessageChan := make(chan protoreflect.ProtoMessage, 10)
		var wait sync.WaitGroup

		id := 7
		logPath := "some_log_path"
		dirPath := "some_file_path"
		file1Path := "filepath1"
		file2Path := "filepath2"
		file3Path := "filepath3"
		hash1 := "adsafaefgewfhdyerytgfedasgtgewg"
		hash2 := "sdafdgdefsgsdfghretawdsfewrtgee"
		hash3 := "ehtjenthegrfdsghrewtrtiuytreere"

		execLogEntry := &spawn.ExecLogEntry{
			Id: uint32(id),
			Type: &spawn.ExecLogEntry_Directory_{
				Directory: &spawn.ExecLogEntry_Directory{
					Path: dirPath,
					Files: []*spawn.ExecLogEntry_File{
						&spawn.ExecLogEntry_File{
							Path: file1Path,
							Digest: &spawn.Digest{
								Hash:      hash1,
								SizeBytes: 1234,
							},
						},
						&spawn.ExecLogEntry_File{
							Path: file2Path,
							Digest: &spawn.Digest{
								Hash:      hash2,
								SizeBytes: 1234,
							},
						},
						&spawn.ExecLogEntry_File{
							Path: file3Path,
							Digest: &spawn.Digest{
								Hash:      hash3,
								SizeBytes: 1234,
							},
						},
					},
				},
			},
		}

		processor := newProcessor(logPath, protoMessageChan, execLogChan, &wait, ctx)
		processor.processDirectory(execLogEntry)

		g.Expect(ctx.Continue()).To(BeTrue())
		diagnostics := ctx.GetDiagnostics()

		g.Expect(len(diagnostics)).To(Equal(0))
		g.Expect(len(execLogChan)).To(Equal(4))
		g.Expect(stderr.String()).To(Equal(""))

		close(execLogChan)
		close(protoMessageChan)

		execLogList := helpers.CreateListFromExecLogChan(execLogChan, t)

		g.Expect(execLogList[0].InputSet).To(BeNil())
		g.Expect(execLogList[0].Spawn).To(BeNil())
		g.Expect(execLogList[0].File.Id).To(Equal(int64(0)))
		g.Expect(execLogList[0].File.Path).To(Equal(dirPath + "/" + file1Path))
		g.Expect(execLogList[0].File.LogPathHash).To(Equal("c2c87929d047a92999dc11b393cd2bda"))
		g.Expect(execLogList[0].File.Hash).To(Equal(hash1))

		g.Expect(execLogList[1].InputSet).To(BeNil())
		g.Expect(execLogList[1].Spawn).To(BeNil())
		g.Expect(execLogList[1].File.Id).To(Equal(int64(0)))
		g.Expect(execLogList[1].File.Path).To(Equal(dirPath + "/" + file2Path))
		g.Expect(execLogList[1].File.LogPathHash).To(Equal("c2c87929d047a92999dc11b393cd2bda"))
		g.Expect(execLogList[1].File.Hash).To(Equal(hash2))

		g.Expect(execLogList[2].InputSet).To(BeNil())
		g.Expect(execLogList[2].Spawn).To(BeNil())
		g.Expect(execLogList[2].File.Id).To(Equal(int64(0)))
		g.Expect(execLogList[2].File.Path).To(Equal(dirPath + "/" + file3Path))
		g.Expect(execLogList[2].File.LogPathHash).To(Equal("c2c87929d047a92999dc11b393cd2bda"))
		g.Expect(execLogList[2].File.Hash).To(Equal(hash3))

		g.Expect(execLogList[3].File).To(BeNil())
		g.Expect(execLogList[3].Spawn).To(BeNil())
		g.Expect(execLogList[3].InputSet.Id).To(Equal(int64(7)))
		g.Expect(len(execLogList[3].InputSet.FileIds)).To(Equal(0))
		g.Expect(execLogList[3].InputSet.LogPathHash).To(Equal("c2c87929d047a92999dc11b393cd2bda"))
		g.Expect(len(execLogList[3].InputSet.FileNames)).To(Equal(3))
		g.Expect(execLogList[3].InputSet.FileNames[0]).To(Equal(dirPath + "/" + file1Path))
		g.Expect(execLogList[3].InputSet.FileNames[1]).To(Equal(dirPath + "/" + file2Path))
		g.Expect(execLogList[3].InputSet.FileNames[2]).To(Equal(dirPath + "/" + file3Path))
		g.Expect(len(execLogList[3].InputSet.NestedInputSets)).To(Equal(0))
	})

	t.Run("unpacks spawn messages correctly", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		execLogChan := make(chan types.ExecLogEntry, 10)
		protoMessageChan := make(chan protoreflect.ProtoMessage, 10)
		var wait sync.WaitGroup

		id := 7
		logPath := "some_log_path"
		arg1 := "some_arg_1"
		arg2 := "some_arg_2"
		arg3 := "some_arg_3"
		envKey1 := "some_key_1"
		envVal1 := "some_val_1"
		envKey2 := "some_key_2"
		envVal2 := "some_val_2"
		envKey3 := "some_key_3"
		envVal3 := "some_val_3"
		inputSetId := 23
		toolSetId := 32
		output_id_1 := 123
		output_id_2 := 321
		label := "//foo/bar:baz"
		mnemonic := "My Awesome Action"
		exitCode := 0
		status := "some status"
		runner := "disk cache"
		hash := "dsfwertgeradfwert"
		timeout := 0

		execLogEntry := &spawn.ExecLogEntry{
			Id: uint32(id),
			Type: &spawn.ExecLogEntry_Spawn_{
				Spawn: &spawn.ExecLogEntry_Spawn{
					Args: []string{
						arg1,
						arg2,
						arg3,
					},
					EnvVars: []*spawn.EnvironmentVariable{
						&spawn.EnvironmentVariable{
							Name:  envKey1,
							Value: envVal1,
						},
						&spawn.EnvironmentVariable{
							Name:  envKey2,
							Value: envVal2,
						},
						&spawn.EnvironmentVariable{
							Name:  envKey3,
							Value: envVal3,
						},
					},
					Platform: &spawn.Platform{
						Properties: []*spawn.Platform_Property{
							&spawn.Platform_Property{},
						},
					},
					InputSetId: uint32(inputSetId),
					ToolSetId:  uint32(toolSetId),
					Outputs: []*spawn.ExecLogEntry_Output{
						&spawn.ExecLogEntry_Output{
							Type: &spawn.ExecLogEntry_Output_OutputId{
								OutputId: uint32(output_id_1),
							},
						},
						&spawn.ExecLogEntry_Output{
							Type: &spawn.ExecLogEntry_Output_OutputId{
								OutputId: uint32(output_id_2),
							},
						},
					},
					TargetLabel:     label,
					Mnemonic:        mnemonic,
					ExitCode:        int32(exitCode),
					Status:          status,
					Runner:          runner,
					CacheHit:        true,
					Remotable:       true,
					Cacheable:       true,
					RemoteCacheable: true,
					Digest: &spawn.Digest{
						Hash:             hash,
						SizeBytes:        1234,
						HashFunctionName: "sha512",
					},
					TimeoutMillis: int64(timeout),
					Metrics:       &spawn.SpawnMetrics{},
				},
			},
		}

		processor := newProcessor(logPath, protoMessageChan, execLogChan, &wait, ctx)
		processor.processSpawn(execLogEntry)

		g.Expect(ctx.Continue()).To(BeTrue())
		diagnostics := ctx.GetDiagnostics()

		g.Expect(len(diagnostics)).To(Equal(0))
		g.Expect(len(execLogChan)).To(Equal(1))
		g.Expect(stderr.String()).To(Equal(""))

		close(execLogChan)
		close(protoMessageChan)

		execLogList := helpers.CreateListFromExecLogChan(execLogChan, t)

		g.Expect(execLogList[0].InputSet).To(BeNil())
		g.Expect(execLogList[0].File).To(BeNil())
		g.Expect(len(execLogList[0].Spawn.Args)).To(Equal(3))
		g.Expect(execLogList[0].Spawn.Args[0]).To(Equal(arg1))
		g.Expect(execLogList[0].Spawn.Args[1]).To(Equal(arg2))
		g.Expect(execLogList[0].Spawn.Args[2]).To(Equal(arg3))
		g.Expect(execLogList[0].Spawn.CacheHit).To(BeTrue())
		g.Expect(execLogList[0].Spawn.Cacheable).To(BeTrue())
		g.Expect(len(execLogList[0].Spawn.Env)).To(Equal(3))
		g.Expect(execLogList[0].Spawn.Env[envKey1]).To(Equal(envVal1))
		g.Expect(execLogList[0].Spawn.Env[envKey2]).To(Equal(envVal2))
		g.Expect(execLogList[0].Spawn.Env[envKey3]).To(Equal(envVal3))
		g.Expect(execLogList[0].Spawn.Hash).To(Equal("dsfwertgeradfwert"))
		g.Expect(execLogList[0].Spawn.Id).To(Equal(int64(id)))
		g.Expect(execLogList[0].Spawn.InputSetId).To(Equal(int64(inputSetId)))
		g.Expect(execLogList[0].Spawn.Label).To(Equal(label))
		g.Expect(execLogList[0].Spawn.LogPath).To(Equal(logPath))
		g.Expect(execLogList[0].Spawn.LogPathHash).To(Equal("c2c87929d047a92999dc11b393cd2bda"))
		g.Expect(len(execLogList[0].Spawn.OutputFileIds)).To(Equal(2))
		g.Expect(execLogList[0].Spawn.OutputFileIds[0]).To(Equal(int64(output_id_1)))
		g.Expect(execLogList[0].Spawn.OutputFileIds[1]).To(Equal(int64(output_id_2)))
		g.Expect(execLogList[0].Spawn.RemoteCachable).To(BeTrue())
		g.Expect(execLogList[0].Spawn.Remoteable).To(BeTrue())
	})

	t.Run("unpacks input set messages correctly", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		execLogChan := make(chan types.ExecLogEntry, 10)
		protoMessageChan := make(chan protoreflect.ProtoMessage, 10)
		var wait sync.WaitGroup

		id := 7
		logPath := "some_log_path"
		fileId1 := 1
		fileId2 := 2
		fileId3 := 3
		directoryId1 := 4
		directoryId2 := 5
		transitiveSetId := 6

		execLogEntry := &spawn.ExecLogEntry{
			Id: uint32(id),
			Type: &spawn.ExecLogEntry_InputSet_{
				InputSet: &spawn.ExecLogEntry_InputSet{
					InputIds: []uint32{
						uint32(fileId1),
						uint32(fileId2),
						uint32(fileId3),
						uint32(directoryId1),
						uint32(directoryId2),
					},
					TransitiveSetIds: []uint32{
						uint32(transitiveSetId),
					},
				},
			},
		}

		processor := newProcessor(logPath, protoMessageChan, execLogChan, &wait, ctx)
		processor.processInputSet(execLogEntry)

		g.Expect(ctx.Continue()).To(BeTrue())
		diagnostics := ctx.GetDiagnostics()

		g.Expect(len(diagnostics)).To(Equal(0))
		g.Expect(len(execLogChan)).To(Equal(1))
		g.Expect(stderr.String()).To(Equal(""))

		close(execLogChan)
		close(protoMessageChan)

		execLogList := helpers.CreateListFromExecLogChan(execLogChan, t)

		g.Expect(execLogList[0].Spawn).To(BeNil())
		g.Expect(execLogList[0].File).To(BeNil())
		g.Expect(len(execLogList[0].InputSet.FileIds)).To(Equal(5))
		g.Expect(execLogList[0].InputSet.Id).To(Equal(int64(id)))
		g.Expect(execLogList[0].InputSet.FileIds[0]).To(Equal(int64(fileId1)))
		g.Expect(execLogList[0].InputSet.FileIds[1]).To(Equal(int64(fileId2)))
		g.Expect(execLogList[0].InputSet.FileIds[2]).To(Equal(int64(fileId3)))
		g.Expect(execLogList[0].InputSet.FileIds[3]).To(Equal(int64(directoryId1)))
		g.Expect(execLogList[0].InputSet.FileIds[4]).To(Equal(int64(directoryId2)))
		g.Expect(execLogList[0].InputSet.LogPathHash).To(Equal("c2c87929d047a92999dc11b393cd2bda"))
		g.Expect(len(execLogList[0].InputSet.NestedInputSets)).To(Equal(1))
		g.Expect(execLogList[0].InputSet.NestedInputSets[0]).To(Equal(int64(transitiveSetId)))
	})

	t.Run("process calls processInputSet", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		execLogChan := make(chan types.ExecLogEntry, 10)
		protoMessageChan := make(chan protoreflect.ProtoMessage, 10)
		var wait sync.WaitGroup
		wait.Add(1)

		processor := newProcessor("some_log_path", protoMessageChan, execLogChan, &wait, ctx)
		go processor.process()

		protoMessageChan <- &spawn.ExecLogEntry{
			Id: uint32(7),
			Type: &spawn.ExecLogEntry_InputSet_{
				InputSet: &spawn.ExecLogEntry_InputSet{},
			},
		}

		close(protoMessageChan)

		wait.Wait()

		close(execLogChan)

		g.Expect(ctx.Continue()).To(BeTrue())
		diagnostics := ctx.GetDiagnostics()

		g.Expect(len(diagnostics)).To(Equal(0))
		g.Expect(len(execLogChan)).To(Equal(1))
		g.Expect(stderr.String()).To(Equal(""))

		execLogList := helpers.CreateListFromExecLogChan(execLogChan, t)

		g.Expect(execLogList[0].Spawn).To(BeNil())
		g.Expect(execLogList[0].File).To(BeNil())
		g.Expect(execLogList[0].InputSet).NotTo(BeNil())
	})

	t.Run("process calls processSpawn", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		execLogChan := make(chan types.ExecLogEntry, 10)
		protoMessageChan := make(chan protoreflect.ProtoMessage, 10)
		var wait sync.WaitGroup
		wait.Add(1)

		processor := newProcessor("some_log_path", protoMessageChan, execLogChan, &wait, ctx)
		go processor.process()

		protoMessageChan <- &spawn.ExecLogEntry{
			Id: uint32(7),
			Type: &spawn.ExecLogEntry_Spawn_{
				Spawn: &spawn.ExecLogEntry_Spawn{},
			},
		}

		close(protoMessageChan)

		wait.Wait()

		close(execLogChan)

		g.Expect(ctx.Continue()).To(BeTrue())
		diagnostics := ctx.GetDiagnostics()

		g.Expect(len(diagnostics)).To(Equal(0))
		g.Expect(len(execLogChan)).To(Equal(1))
		g.Expect(stderr.String()).To(Equal(""))

		execLogList := helpers.CreateListFromExecLogChan(execLogChan, t)

		g.Expect(execLogList[0].Spawn).NotTo(BeNil())
		g.Expect(execLogList[0].File).To(BeNil())
		g.Expect(execLogList[0].InputSet).To(BeNil())
	})

	t.Run("process calls processFile", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		execLogChan := make(chan types.ExecLogEntry, 10)
		protoMessageChan := make(chan protoreflect.ProtoMessage, 10)
		var wait sync.WaitGroup
		wait.Add(1)

		processor := newProcessor("some_log_path", protoMessageChan, execLogChan, &wait, ctx)
		go processor.process()

		protoMessageChan <- &spawn.ExecLogEntry{
			Id: uint32(7),
			Type: &spawn.ExecLogEntry_File_{
				File: &spawn.ExecLogEntry_File{},
			},
		}

		close(protoMessageChan)

		wait.Wait()

		close(execLogChan)

		g.Expect(ctx.Continue()).To(BeTrue())
		diagnostics := ctx.GetDiagnostics()

		g.Expect(len(diagnostics)).To(Equal(0))
		g.Expect(len(execLogChan)).To(Equal(1))
		g.Expect(stderr.String()).To(Equal(""))

		execLogList := helpers.CreateListFromExecLogChan(execLogChan, t)

		g.Expect(execLogList[0].Spawn).To(BeNil())
		g.Expect(execLogList[0].File).NotTo(BeNil())
		g.Expect(execLogList[0].InputSet).To(BeNil())
	})

	t.Run("process calls processDirectory", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		execLogChan := make(chan types.ExecLogEntry, 10)
		protoMessageChan := make(chan protoreflect.ProtoMessage, 10)
		var wait sync.WaitGroup
		wait.Add(1)

		processor := newProcessor("some_log_path", protoMessageChan, execLogChan, &wait, ctx)
		go processor.process()

		protoMessageChan <- &spawn.ExecLogEntry{
			Id: uint32(7),
			Type: &spawn.ExecLogEntry_Directory_{
				Directory: &spawn.ExecLogEntry_Directory{
					Path: "some_dir_path",
					Files: []*spawn.ExecLogEntry_File{
						&spawn.ExecLogEntry_File{
							Path: "some_filepath_1",
							Digest: &spawn.Digest{
								Hash: "adfasdfsadf",
							},
						},
						&spawn.ExecLogEntry_File{
							Path: "some_filepath_2",
							Digest: &spawn.Digest{
								Hash: "sadfsdfsads",
							},
						},
						&spawn.ExecLogEntry_File{
							Path: "some_filepath_3",
							Digest: &spawn.Digest{
								Hash: "asdfewgterw",
							},
						},
					},
				},
			},
		}

		close(protoMessageChan)

		wait.Wait()

		close(execLogChan)

		g.Expect(ctx.Continue()).To(BeTrue())
		diagnostics := ctx.GetDiagnostics()

		g.Expect(len(diagnostics)).To(Equal(0))
		g.Expect(len(execLogChan)).To(Equal(4))
		g.Expect(stderr.String()).To(Equal(""))

		execLogList := helpers.CreateListFromExecLogChan(execLogChan, t)

		g.Expect(execLogList[0].Spawn).To(BeNil())
		g.Expect(execLogList[0].File).NotTo(BeNil())
		g.Expect(execLogList[0].InputSet).To(BeNil())

		g.Expect(execLogList[1].Spawn).To(BeNil())
		g.Expect(execLogList[1].File).NotTo(BeNil())
		g.Expect(execLogList[1].InputSet).To(BeNil())

		g.Expect(execLogList[2].Spawn).To(BeNil())
		g.Expect(execLogList[2].File).NotTo(BeNil())
		g.Expect(execLogList[2].InputSet).To(BeNil())

		g.Expect(execLogList[3].Spawn).To(BeNil())
		g.Expect(execLogList[3].File).To(BeNil())
		g.Expect(execLogList[3].InputSet).NotTo(BeNil())
	})
}
