package converter

import (
	"bufio"
	"os"
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/converter/compact"
	"github.com/aspect-build/execlog-diff/converter/compact/decompress"
	"github.com/aspect-build/execlog-diff/types"
	"github.com/aspect-build/execlog-diff/diagnostic"
)

/**
 * This function is meant to be called as a go routine.
 * It should take an incoming ExecLog of any kind and call the appropriate function
 * to convert it into an internal format (that should be the same across all execlog types).
 */
func Convert(path string, execLogEntryChan chan types.ExecLogEntry, waitGroup *sync.WaitGroup, ctx *context.Context) {
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

	scanner := bufio.NewScanner(file)
	scan_result := scanner.Scan()
	if !scan_result {
		err := scanner.Err()
		if err == nil {
			ctx.Error(&diagnostic.Diagnostic{
				Title:   "Provided file is empty",
				Message: "File: " + path,
			})

			waitGroup.Done()
			return
		} else {
			ctx.Error(&diagnostic.Diagnostic{
				Title:   "Failed to scan file",
				Message: err.Error(),
			})

			waitGroup.Done()
			return
		}
	}

	if scanner.Text() == "" {
		ctx.Error(&diagnostic.Diagnostic{
			Title:   "Provided file is empty",
			Message: "File: " + path,
		})

		waitGroup.Done()
		return
	} else if scanner.Text() == "{" {
		// Bazel pretty prints the JSON execlogs so the first line will always be a solitary "{" in this case
		/********************************
		 * JSON ExecLog
		 ********************************/
		ctx.Error(&diagnostic.Diagnostic{
			Message: "JSON execlogs are not yet supported",
		})

		waitGroup.Done()
		return
	}

	if ok, err := decompress.CanDecompress(path); !ok {
		/********************************
		 * BINARY ExecLog
		 ********************************/
		ctx.Error(&diagnostic.Diagnostic{
			Message: "BINARY execlogs are not yet supported. " + err.Error(),
		})

		waitGroup.Done()
		return
	} else {
		/********************************
		 * COMPACT ExecLog
		 ********************************/
		var compactWait sync.WaitGroup
		compactWait.Add(1)
		go compact.Convert(path, execLogEntryChan, &compactWait, ctx)
		compactWait.Wait()
	}

	waitGroup.Done()
}
