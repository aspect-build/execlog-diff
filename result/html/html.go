package html

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/aspect-build/execlog-diff/result/proto"
	"github.com/aspect-build/execlog-diff/diagnostic"
)

type html struct {
	// Save any errors hit during generation to return to the caller
	diagnostics map[string]*diagnostic.Diagnostic

	totalBytes          int64
	totalCacheableBytes int64
	outputPath          string
}

func New(outputPath string) *html {
	return &html{
		diagnostics: make(map[string]*diagnostic.Diagnostic),
		outputPath:  outputPath,
	}
}

func (h *html) error(d *diagnostic.Diagnostic) {
	d.Severity = diagnostic.Severity_ERROR

	hasher := md5.New()
	hasher.Write([]byte(d.Title))
	hasher.Write([]byte(d.Message))
	diagnosticHash := hex.EncodeToString(hasher.Sum(nil))
	h.diagnostics[diagnosticHash] = d
}

func (h *html) Print(result *proto.Result) []*diagnostic.Diagnostic {
	switch result_type := result.Result.(type) {
	case *proto.Result_ExecDiff:
		h.printExecDiff(result.GetExecDiff())
	default:
		h.error(&diagnostic.Diagnostic{
			Title:   "Result type not defined",
			Message: fmt.Sprintf("Unable to print results of type: %s", reflect.TypeOf(result_type)),
		})
	}

	diagnostics := []*diagnostic.Diagnostic{}
	for _, diagnostic := range h.diagnostics {
		diagnostics = append(diagnostics, diagnostic)
	}
	return diagnostics
}
