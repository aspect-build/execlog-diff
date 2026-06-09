package compact

import (
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/converter/compact/bazel7"
	"github.com/aspect-build/execlog-diff/converter/compact/bazel8"
	"github.com/aspect-build/execlog-diff/converter/compact/bazel8_4_2"
	"github.com/aspect-build/execlog-diff/types"
)

func Convert(path string, execLogEntryChan chan types.ExecLogEntry, waitGroup *sync.WaitGroup, ctx *context.Context) {
	switch {
	case bazel8_4_2.InvocationIsBazel8_4_2(path, ctx):
		if !ctx.Continue() {
			waitGroup.Done()
			return
		}

		var bazel8_4_2Wait sync.WaitGroup
		bazel8_4_2Wait.Add(1)
		go bazel8_4_2.Convert(path, execLogEntryChan, &bazel8_4_2Wait, ctx)
		bazel8_4_2Wait.Wait()
	case bazel8.InvocationIsBazel8(path, ctx):
		if !ctx.Continue() {
			waitGroup.Done()
			return
		}

		var bazel8Wait sync.WaitGroup
		bazel8Wait.Add(1)
		go bazel8.Convert(path, execLogEntryChan, &bazel8Wait, ctx)
		bazel8Wait.Wait()
	default:
		if !ctx.Continue() {
			waitGroup.Done()
			return
		}

		var bazel7Wait sync.WaitGroup
		bazel7Wait.Add(1)
		go bazel7.Convert(path, execLogEntryChan, &bazel7Wait, ctx)
		bazel7Wait.Wait()
	}

	waitGroup.Done()
}
