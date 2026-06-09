package bazel8

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/aspect-build/execlog-diff/third_party/github.com/bazelbuild/bazel8/src/main/protobuf/spawn"
	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"
	"github.com/aspect-build/execlog-diff/diagnostic"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type processor struct {
	idDescriptor protoreflect.FieldDescriptor

	invocationDescriptor        protoreflect.FieldDescriptor
	fileDescriptor              protoreflect.FieldDescriptor
	directoryDescriptor         protoreflect.FieldDescriptor
	unresolvedSymlinkDescriptor protoreflect.FieldDescriptor
	inputSetDescriptor          protoreflect.FieldDescriptor
	spawnDescriptor             protoreflect.FieldDescriptor
	symlinkActionDescriptor     protoreflect.FieldDescriptor
	symlinkEntrySetDescriptor   protoreflect.FieldDescriptor
	runfilesTreeDescriptor      protoreflect.FieldDescriptor

	logPath     string
	logPathHash string

	protoMessageChan <-chan protoreflect.ProtoMessage
	execLogEntryChan chan<- types.ExecLogEntry
	waitGroup        *sync.WaitGroup
	ctx              *context.Context
}

func newProcessor(
	logPath string,
	protoMessageChan <-chan protoreflect.ProtoMessage,
	execLogEntryChan chan<- types.ExecLogEntry,
	waitGroup *sync.WaitGroup,
	ctx *context.Context,
) *processor {
	hasher := md5.New()
	hasher.Write([]byte(logPath))
	logPathHash := hex.EncodeToString(hasher.Sum(nil))

	execLogEntry := spawn.ExecLogEntry{}

	return &processor{
		idDescriptor: descriptor(&execLogEntry, "id"),

		invocationDescriptor:        descriptor(&execLogEntry, "invocation"),
		fileDescriptor:              descriptor(&execLogEntry, "file"),
		directoryDescriptor:         descriptor(&execLogEntry, "directory"),
		unresolvedSymlinkDescriptor: descriptor(&execLogEntry, "unresolved_symlink"),
		inputSetDescriptor:          descriptor(&execLogEntry, "input_set"),
		spawnDescriptor:             descriptor(&execLogEntry, "spawn"),
		symlinkActionDescriptor:     descriptor(&execLogEntry, "symlink_action"),
		symlinkEntrySetDescriptor:   descriptor(&execLogEntry, "symlink_entry_set"),
		runfilesTreeDescriptor:      descriptor(&execLogEntry, "runfiles_tree"),

		logPath:          logPath,
		logPathHash:      logPathHash,
		protoMessageChan: protoMessageChan,
		execLogEntryChan: execLogEntryChan,
		waitGroup:        waitGroup,
		ctx:              ctx,
	}
}

func (processor *processor) process() {
	for message := range processor.protoMessageChan {
		switch {
		case has(message, processor.invocationDescriptor):
			// Do nothing
			continue
		case has(message, processor.fileDescriptor):
			processor.processFile(message)
		case has(message, processor.directoryDescriptor):
			processor.processDirectory(message)
		case has(message, processor.unresolvedSymlinkDescriptor):
			processor.processUnresolvedSymlink(message)
		case has(message, processor.inputSetDescriptor):
			processor.processInputSet(message)
		case has(message, processor.spawnDescriptor):
			processor.processSpawn(message)
		case has(message, processor.symlinkActionDescriptor):
			processor.processSymlinkAction(message)
		case has(message, processor.symlinkEntrySetDescriptor):
			processor.processSymlinkEntrySet(message)
		case has(message, processor.runfilesTreeDescriptor):
			processor.processRunfilesTree(message)
		default:
			processor.ctx.Warning(&diagnostic.Diagnostic{
				Message: fmt.Sprintf("Unknown message: %s", message),
			})
		}
	}

	processor.waitGroup.Done()
}

func (processor *processor) getId(message protoreflect.ProtoMessage) int64 {
	return int64(message.ProtoReflect().Get(processor.idDescriptor).Uint())
}

func (processor *processor) processFile(message protoreflect.ProtoMessage) {
	file := message.ProtoReflect().Get(processor.fileDescriptor)

	path := get(file, "path")
	digest := get(file, "digest")
	hash := get(digest, "hash")

	sizeBytes := get(digest, "size_bytes")

	processor.execLogEntryChan <- types.ExecLogEntry{
		File: &types.File{
			Id:          processor.getId(message),
			Path:        path.String(),
			Hash:        hash.String(),
			LogPathHash: processor.logPathHash,
			SizeBytes:   sizeBytes.Int(),
		},
	}
}

func (processor *processor) processInputSet(message protoreflect.ProtoMessage) {
	input_set := message.ProtoReflect().Get(processor.inputSetDescriptor)

	input_ids := get(input_set, "input_ids")
	transitive_set_ids := get(input_set, "transitive_set_ids")

	input_ids_actual := make([]int64, 0)
	for i := 0; i < input_ids.List().Len(); i++ {
		input_ids_actual = append(input_ids_actual, int64(input_ids.List().Get(i).Uint()))
	}

	transitive_set_ids_actual := make([]int64, 0)
	for i := 0; i < transitive_set_ids.List().Len(); i++ {
		transitive_set_ids_actual = append(transitive_set_ids_actual, int64(transitive_set_ids.List().Get(i).Uint()))
	}

	processor.execLogEntryChan <- types.ExecLogEntry{
		InputSet: &types.InputSet{
			Id:              processor.getId(message),
			FileNames:       make([]string, 0),
			FileIds:         input_ids_actual,
			NestedInputSets: transitive_set_ids_actual,
			LogPathHash:     processor.logPathHash,
		},
	}
}

func (processor *processor) processDirectory(message protoreflect.ProtoMessage) {
	directory := message.ProtoReflect().Get(processor.directoryDescriptor)

	directory_path := get(directory, "path")

	files := get(directory, "files")
	formatted_files := make([]string, files.List().Len())
	for i := 0; i < files.List().Len(); i++ {
		file := files.List().Get(i)

		path := get(file, "path")
		digest := get(file, "digest")
		hash := get(digest, "hash")

		sizeBytes := get(digest, "size_bytes")

		formatted_filepath := fmt.Sprintf("%s/%s", directory_path.String(), path.String())
		formatted_files[i] = formatted_filepath

		processor.execLogEntryChan <- types.ExecLogEntry{
			File: &types.File{
				Id:          0,
				Path:        fmt.Sprintf("%s/%s", directory_path.String(), path.String()),
				Hash:        hash.String(),
				LogPathHash: processor.logPathHash,
				SizeBytes:   sizeBytes.Int(),
			},
		}
	}

	processor.execLogEntryChan <- types.ExecLogEntry{
		InputSet: &types.InputSet{
			Id:              processor.getId(message),
			FileNames:       formatted_files,
			FileIds:         make([]int64, 0),
			NestedInputSets: make([]int64, 0),
			LogPathHash:     processor.logPathHash,
		},
	}
}

func (processor *processor) processRunfilesTree(message protoreflect.ProtoMessage) {
	// TODO: Create a test execlog containing this message type
	return
}

func (processor *processor) processUnresolvedSymlink(message protoreflect.ProtoMessage) {
	// TODO: Create a test execlog containing this message type
	return
}

func (processor *processor) processSymlinkAction(message protoreflect.ProtoMessage) {
	// TODO: Create a test execlog containing this message type
	return
}

func (processor *processor) processSymlinkEntrySet(message protoreflect.ProtoMessage) {
	// TODO: Create a test execlog containing this message type
	return
}

func (processor *processor) processSpawn(message protoreflect.ProtoMessage) {
	spawn := message.ProtoReflect().Get(processor.spawnDescriptor)

	arg_values := get(spawn, "args")
	args := make([]string, 0)
	for i := 0; i < arg_values.List().Len(); i++ {
		args = append(args, arg_values.List().Get(i).String())
	}

	env_values := get(spawn, "env_vars")
	env := make(map[string]string)
	for i := 0; i < env_values.List().Len(); i++ {
		env_value := env_values.List().Get(i)
		name := get(env_value, "name").String()
		value := get(env_value, "value").String()
		env[name] = value
	}

	input_set_id := get(spawn, "input_set_id")

	output_files_values := get(spawn, "outputs")
	output_file_ids := make([]int64, 0)
	output_set_ids := make([]int64, 0)
	for i := 0; i < output_files_values.List().Len(); i++ {
		output_files_value := output_files_values.List().Get(i)

		if output_files_value.Message().Has(output_files_value.Message().Descriptor().Fields().ByName("output_id")) {
			output_file_ids = append(output_file_ids, int64(get(output_files_value, "output_id").Uint()))
		}
	}

	label := get(spawn, "target_label")

	cache_hit := get(spawn, "cache_hit")
	remotable := get(spawn, "remotable")
	cacheable := get(spawn, "cacheable")
	remote_cacheable := get(spawn, "remote_cacheable")
	digest := get(spawn, "digest")
	hash := get(digest, "hash")

	processor.execLogEntryChan <- types.ExecLogEntry{
		Spawn: &types.IncomingSpawn{
			LogPath:        processor.logPath,
			LogPathHash:    processor.logPathHash,
			Id:             processor.getId(message),
			Args:           args,
			Env:            env,
			InputSetId:     int64(input_set_id.Uint()),
			OutputFileIds:  output_file_ids,
			OutputFiles:    make([]string, 0),
			OutputSetIds:   output_set_ids,
			Label:          label.String(),
			CacheHit:       cache_hit.Bool(),
			Remoteable:     remotable.Bool(),
			Cacheable:      cacheable.Bool(),
			RemoteCachable: remote_cacheable.Bool(),
			Hash:           hash.String(),
		},
	}
}
