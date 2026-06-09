package converter

import (
	"strings"
	"sync"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/test-data/helpers"
	"github.com/aspect-build/execlog-diff/types"
	"github.com/aspect-build/execlog-diff/diagnostic"
	"gotest.tools/v3/golden"

	. "github.com/onsi/gomega"
)

func TestConverter(t *testing.T) {
	t.Run("errors opening files are surfaced", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		filepath := "./some_nonexistant_file_123456.log"

		execLogChan := make(chan types.ExecLogEntry)
		var compactWait sync.WaitGroup
		compactWait.Add(1)
		go Convert(filepath, execLogChan, &compactWait, ctx)
		compactWait.Wait()

		diagnostics := ctx.GetDiagnostics()

		g.Expect(stderr.String()).To(HavePrefix("ERROR: Failed to open file"))
		g.Expect(stderr.String()).To(ContainSubstring(filepath))
		g.Expect(stderr.String()).To(ContainSubstring("no such file or directory"))
		g.Expect(ctx.Continue()).To(BeFalse())
		g.Expect(len(diagnostics)).To(Equal(1))
		g.Expect(diagnostics[0].Severity).To(Equal(diagnostic.Severity_ERROR))
		g.Expect(len(execLogChan)).To(Equal(0))

		close(execLogChan)
	})

	t.Run("it flags empty files as invalid", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		filepath := "../test-data/empty.execlog.binary.log"
		execLogChan := make(chan types.ExecLogEntry)
		var compactWait sync.WaitGroup
		compactWait.Add(1)
		go Convert(filepath, execLogChan, &compactWait, ctx)
		compactWait.Wait()

		diagnostics := ctx.GetDiagnostics()

		g.Expect(stderr.String()).To(HavePrefix("ERROR: Provided file is empty"))
		g.Expect(stderr.String()).To(ContainSubstring(filepath))
		g.Expect(ctx.Continue()).To(BeFalse())
		g.Expect(len(diagnostics)).To(Equal(1))
		g.Expect(diagnostics[0].Severity).To(Equal(diagnostic.Severity_ERROR))
		g.Expect(len(execLogChan)).To(Equal(0))

		close(execLogChan)
	})

	t.Run("it picks up JSON files", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		filepath := "../test-data/changed_file_7.3.0_1.json"
		execLogChan := make(chan types.ExecLogEntry)
		var compactWait sync.WaitGroup
		compactWait.Add(1)
		go Convert(filepath, execLogChan, &compactWait, ctx)
		compactWait.Wait()

		diagnostics := ctx.GetDiagnostics()

		g.Expect(stderr.String()).To(ContainSubstring("ERROR: JSON execlogs are not yet supported"))
		g.Expect(ctx.Continue()).To(BeFalse())
		g.Expect(len(diagnostics)).To(Equal(1))
		g.Expect(diagnostics[0].Severity).To(Equal(diagnostic.Severity_ERROR))
		g.Expect(len(execLogChan)).To(Equal(0))

		close(execLogChan)
	})

	t.Run("it picks up BINARY files", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		filepath := "../test-data/changed_file_7.3.0_1.binary.log"
		execLogChan := make(chan types.ExecLogEntry)
		var compactWait sync.WaitGroup
		compactWait.Add(1)
		go Convert(filepath, execLogChan, &compactWait, ctx)
		compactWait.Wait()

		diagnostics := ctx.GetDiagnostics()

		g.Expect(stderr.String()).To(ContainSubstring("ERROR: BINARY execlogs are not yet supported"))
		g.Expect(ctx.Continue()).To(BeFalse())
		g.Expect(len(diagnostics)).To(Equal(1))
		g.Expect(diagnostics[0].Severity).To(Equal(diagnostic.Severity_ERROR))
		g.Expect(len(execLogChan)).To(Equal(0))

		close(execLogChan)
	})

	t.Run("it picks up COMPACT files", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		filepath := "../test-data/changed_file_7.3.0_1.compact.log"
		execLogChan := make(chan types.ExecLogEntry, 100)
		var compactWait sync.WaitGroup
		compactWait.Add(1)
		go Convert(filepath, execLogChan, &compactWait, ctx)
		compactWait.Wait()

		g.Expect(ctx.Continue()).To(BeTrue())
		diagnostics := ctx.GetDiagnostics()

		g.Expect(stderr.String()).To(Equal(""))
		g.Expect(len(execLogChan)).To(Equal(15))
		g.Expect(len(diagnostics)).To(Equal(0))

		close(execLogChan)
	})

	t.Run("create golden files", func(t *testing.T) {
		files := map[string]string{
			"golden_files/changed_directory.txt": "../test-data/changed_directory_7.3.0_1.compact.log",
			"golden_files/changed_file.txt":      "../test-data/changed_file_7.3.0_1.compact.log",
		}

		for golden_filepath, execlog_filepath := range files {
			g := NewGomegaWithT(t)

			var stderr strings.Builder
			ctx := context.New(&stderr, context.Diff)

			execLogChan := make(chan types.ExecLogEntry, 100)
			var compactWait sync.WaitGroup
			compactWait.Add(1)
			go Convert(execlog_filepath, execLogChan, &compactWait, ctx)
			compactWait.Wait()

			g.Expect(ctx.Continue()).To(BeTrue())
			diagnostics := ctx.GetDiagnostics()

			g.Expect(stderr.String()).To(Equal(""))
			g.Expect(len(diagnostics)).To(Equal(0))

			close(execLogChan)

			output := helpers.CreateGoldenFileFromExecLogChan(execLogChan, t)
			golden.Assert(t, output, golden_filepath)
		}
	})
}
