package bazel8_4_2

import (
	"bytes"
	go_errors "errors"
	"io"
	"os"
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/converter/compact/decompress"
	"github.com/aspect-build/execlog-diff/third_party/github.com/bazelbuild/bazel8_4_2/src/main/protobuf/spawn"
	"github.com/aspect-build/execlog-diff/diagnostic"
	"google.golang.org/protobuf/encoding/protodelim"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func unmarshal(path string, protoMessageChan chan protoreflect.ProtoMessage, waitGroup *sync.WaitGroup, ctx *context.Context) {
	file, err := os.Open(path)
	if err != nil {
		ctx.Error(&diagnostic.Diagnostic{
			Title:   "Failed to open file",
			Message: err.Error(),
		})

		waitGroup.Done()
		return
	}
	defer file.Close()

	var buffer bytes.Buffer
	err = decompress.Decompress(file, &buffer)
	if err != nil {
		ctx.Error(&diagnostic.Diagnostic{
			Title:   "Failed to decompress " + path,
			Message: err.Error(),
		})

		waitGroup.Done()
		return
	}

	// Create unmarshal options with 100MB size limit to handle large exec log entries
	unmarshalOpts := protodelim.UnmarshalOptions{
		MaxSize: 100 * 1024 * 1024, // 100 MB
	}

	for {
		if !ctx.Continue() {
			break
		}

		execLogEntry := spawn.ExecLogEntry{}
		message := execLogEntry.ProtoReflect().Interface()

		err := unmarshalOpts.UnmarshalFrom(&buffer, message)
		if go_errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			ctx.Warning(&diagnostic.Diagnostic{
				Title:   "Proto Unmarshal Error",
				Message: err.Error(),
			})
			continue
		}

		protoMessageChan <- message
	}

	waitGroup.Done()
}
