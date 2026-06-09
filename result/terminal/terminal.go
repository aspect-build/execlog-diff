package terminal

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/aspect-build/execlog-diff/result/proto"
	"github.com/aspect-build/execlog-diff/diagnostic"
)

type terminal struct {
	// Used for logging
	stdout io.Writer

	// Used for logging
	stderr io.Writer

	// Save any errors hit during printing to return to the caller
	diagnostics map[string]*diagnostic.Diagnostic
}

func New(stdout io.Writer, stderr io.Writer) *terminal {
	return &terminal{
		stdout:      stdout,
		stderr:      stderr,
		diagnostics: make(map[string]*diagnostic.Diagnostic),
	}
}

func NewDefault() *terminal {
	return New(os.Stdout, os.Stderr)
}

func (t *terminal) error(d *diagnostic.Diagnostic) {
	d.Severity = diagnostic.Severity_ERROR

	hasher := md5.New()
	hasher.Write([]byte(d.Title))
	hasher.Write([]byte(d.Message))
	diagnosticHash := hex.EncodeToString(hasher.Sum(nil))
	t.diagnostics[diagnosticHash] = d
}

func (t *terminal) Print(result *proto.Result) []*diagnostic.Diagnostic {
	switch result_type := result.Result.(type) {
	case *proto.Result_ExecDiff:
		t.printExecDiff(result.GetExecDiff())
	case *proto.Result_LabelQuery:
		t.printListLabels(result.GetLabelQuery())
	default:
		t.error(&diagnostic.Diagnostic{
			Title:   "Result type not defined",
			Message: fmt.Sprintf("Unable to print results of type: %s", reflect.TypeOf(result_type)),
		})
	}

	diagnostics := []*diagnostic.Diagnostic{}
	for _, diagnostic := range t.diagnostics {
		diagnostics = append(diagnostics, diagnostic)
	}
	return diagnostics
}
