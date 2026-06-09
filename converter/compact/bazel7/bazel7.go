package bazel7

import (
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/types"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Convert(path string, execLogEntryChan chan types.ExecLogEntry, waitGroup *sync.WaitGroup, ctx *context.Context) {
	protoMessageChan := make(chan protoreflect.ProtoMessage)

	// Unmarshal proto
	var unmarshalWait sync.WaitGroup
	unmarshalWait.Add(1)

	go unmarshal(path, protoMessageChan, &unmarshalWait, ctx)

	// Process the unmarshalled proto
	var processWait sync.WaitGroup
	processWait.Add(1)

	processor := newProcessor(path, protoMessageChan, execLogEntryChan, &processWait, ctx)
	go processor.process()

	// Wait for goroutines to finish
	unmarshalWait.Wait()
	close(protoMessageChan)

	processWait.Wait()
	waitGroup.Done()
}
