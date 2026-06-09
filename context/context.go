/**
 * As execlog-diff is running its async go routines, they can log any diagnostic they have with this call and
 * check whether to continue running or which will update depending on whats going on in other go routines
 */
package context

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/aspect-build/execlog-diff/diagnostic"
)

type RunMode int

const (
	Unknown RunMode = iota
	Diff
	ListLabels
)

type Context struct {
	// Used for logging
	stderr io.Writer

	// Used to ensure all the diagnostics on the diagnosticChan are processed.
	diagnosticWaitGroup *sync.WaitGroup

	// This channel is used to allow non-error diagnostics to be passed to the context without blocking.
	// We do some logging and processing on the diagnostic that the main go routines dont need to wait for.
	diagnosticChan chan *diagnostic.Diagnostic

	// This map is used to store diagnostics until the context finishes running.
	// At which point we will return all the diagnostics as a list.
	// Storing as a Map over List to eliminate duplicates.
	diagnostics map[string]*diagnostic.Diagnostic

	// This boolean value will be passed to external go routines.
	// true -> those external go routines should continue running.
	// false -> those external go routines should stop running
	shouldContinue bool

	// This boolean is to track the internal state of our context.
	// false -> keep running dignostic handling
	// true -> we should be done running and do not accept any more diagnostics
	completed bool

	Mode RunMode
}

func New(
	stderr io.Writer,
	mode RunMode,
) *Context {
	var diagnosticWaitGroup sync.WaitGroup
	diagnosticWaitGroup.Add(1)

	ctx := &Context{
		stderr:              stderr,
		diagnosticWaitGroup: &diagnosticWaitGroup,
		diagnosticChan:      make(chan *diagnostic.Diagnostic, 10),
		diagnostics:         make(map[string]*diagnostic.Diagnostic),
		shouldContinue:      true,
		Mode:                mode,
	}

	go ctx.handleDiagnostics()

	return ctx
}

func NewQuery(mode RunMode) *Context {
	return New(os.Stderr, mode)
}

func NewDiff() *Context {
	return New(os.Stderr, Diff)
}

func NewDefault() *Context {
	return New(os.Stderr, Diff)
}

func (ctx *Context) Continue() bool {
	return ctx.shouldContinue
}

func (context *Context) handleDiagnostics() {
	for d := range context.diagnosticChan {
		fmt.Fprintf(context.stderr, "%s: ", d.Severity.String())
		if d.Title != "" {
			fmt.Fprintln(context.stderr, d.Title)
			fmt.Fprint(context.stderr, "  ")
		}
		fmt.Fprintln(context.stderr, d.Message)

		hasher := md5.New()
		hasher.Write([]byte(d.Title))
		hasher.Write([]byte(d.Message))
		diagnosticHash := hex.EncodeToString(hasher.Sum(nil))
		context.diagnostics[diagnosticHash] = d
	}

	context.diagnosticWaitGroup.Done()
}

func (context *Context) Diagnostic(d *diagnostic.Diagnostic) {
	if d.Severity == diagnostic.Severity_ERROR {
		context.shouldContinue = false
	}

	if !context.completed {
		context.diagnosticChan <- d
	}
}

func (context *Context) Warning(d *diagnostic.Diagnostic) {
	d.Severity = diagnostic.Severity_WARNING
	context.Diagnostic(d)
}

func (context *Context) Error(d *diagnostic.Diagnostic) {
	d.Severity = diagnostic.Severity_ERROR
	context.Diagnostic(d)
}

// GetDiagnostics should only be called once when execlog-diff is done running
func (context *Context) GetDiagnostics() []*diagnostic.Diagnostic {
	if !context.completed {
		context.completed = true
		context.shouldContinue = false

		// Intentionally destructive. If here, execlog-diff is done running
		close(context.diagnosticChan)
	}

	// Ensure that all diagnostics from the chan have been processed
	context.diagnosticWaitGroup.Wait()

	diagnostics := []*diagnostic.Diagnostic{}
	for _, diagnostic := range context.diagnostics {
		diagnostics = append(diagnostics, diagnostic)
	}
	return diagnostics
}
