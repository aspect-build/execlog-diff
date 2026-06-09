package bazel8

import (
	"bytes"
	"os"

	"github.com/aspect-build/execlog-diff/third_party/github.com/bazelbuild/bazel8/src/main/protobuf/spawn"
	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/converter/compact/decompress"
	"github.com/aspect-build/execlog-diff/diagnostic"
	"google.golang.org/protobuf/encoding/protodelim"
)

func InvocationIsBazel8(path string, ctx *context.Context) bool {
	file, err := os.Open(path)
	if err != nil {
		ctx.Error(&diagnostic.Diagnostic{
			Title:   "Failed to open file",
			Message: err.Error(),
		})

		return false
	}
	defer file.Close()

	var b bytes.Buffer
	err = decompress.Decompress(file, &b)
	if err != nil {
		ctx.Error(&diagnostic.Diagnostic{
			Title:   "Failed to decompress " + path,
			Message: err.Error(),
		})

		return false
	}

	execLogEntry := spawn.ExecLogEntry{}
	message := execLogEntry.ProtoReflect().Interface()
	err = protodelim.UnmarshalFrom(&b, message)

	if err != nil {
		ctx.Error(&diagnostic.Diagnostic{
			Title:   "Failed to unmarshal proto",
			Message: err.Error(),
		})

		return false
	}

	invocation_descriptor := descriptor(&execLogEntry, "invocation")
	invocation := message.ProtoReflect().Get(invocation_descriptor)
	workspace_runfiles_directory_descriptor := descriptor(&invocation, "workspace_runfiles_directory")

	return invocation.Message().Interface().ProtoReflect().Has(workspace_runfiles_directory_descriptor)
}
