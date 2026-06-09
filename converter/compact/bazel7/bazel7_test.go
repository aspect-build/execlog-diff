package bazel7

import (
	"strings"
	"sync"
	"testing"

	"gotest.tools/v3/golden"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/test-data/helpers"
	"github.com/aspect-build/execlog-diff/types"
	. "github.com/onsi/gomega"
)

func TestBazel7(t *testing.T) {
	t.Run("source files files are processed", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		filepath := "../../../test-data/changed_file_7.3.0_1.compact.log"

		execLogChan := make(chan types.ExecLogEntry, 100)
		var compactWait sync.WaitGroup
		compactWait.Add(1)
		go Convert(filepath, execLogChan, &compactWait, ctx)
		compactWait.Wait()

		g.Expect(ctx.Continue()).To(BeTrue())
		diagnostics := ctx.GetDiagnostics()

		g.Expect(len(execLogChan)).To(Equal(15))
		g.Expect(stderr.String()).To(Equal(""))
		g.Expect(len(diagnostics)).To(Equal(0))

		close(execLogChan)

		execLogList := helpers.CreateListFromExecLogChan(execLogChan, t)

		g.Expect(helpers.ContainsFile(execLogList, "some_file.txt")).To(BeTrue())
	})

	t.Run("generated files are processed", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		filepath := "../../../test-data/changed_file_7.3.0_1.compact.log"

		execLogChan := make(chan types.ExecLogEntry, 100)
		var compactWait sync.WaitGroup
		compactWait.Add(1)
		go Convert(filepath, execLogChan, &compactWait, ctx)
		compactWait.Wait()

		g.Expect(ctx.Continue()).To(BeTrue())
		diagnostics := ctx.GetDiagnostics()

		g.Expect(len(execLogChan)).To(Equal(15))
		g.Expect(stderr.String()).To(Equal(""))
		g.Expect(len(diagnostics)).To(Equal(0))

		close(execLogChan)

		execLogList := helpers.CreateListFromExecLogChan(execLogChan, t)

		g.Expect(helpers.ContainsFile(execLogList, "bazel-out/darwin_arm64-fastbuild/bin/foo/foo.txt")).To(BeTrue())
		g.Expect(helpers.ContainsFile(execLogList, "bazel-out/darwin_arm64-fastbuild/bin/bar/bar.txt")).To(BeTrue())
		g.Expect(helpers.ContainsFile(execLogList, "bazel-out/darwin_arm64-fastbuild/bin/baz/baz.txt")).To(BeTrue())
	})

	t.Run("directories are turned into files & input sets", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		filepath := "../../../test-data/changed_directory_7.3.0_1.compact.log"

		execLogChan := make(chan types.ExecLogEntry, 100)
		var compactWait sync.WaitGroup
		compactWait.Add(1)
		go Convert(filepath, execLogChan, &compactWait, ctx)
		compactWait.Wait()

		g.Expect(ctx.Continue()).To(BeTrue())
		diagnostics := ctx.GetDiagnostics()

		g.Expect(len(execLogChan)).To(Equal(27))
		g.Expect(stderr.String()).To(Equal(""))
		g.Expect(len(diagnostics)).To(Equal(0))

		close(execLogChan)

		execLogList := helpers.CreateListFromExecLogChan(execLogChan, t)

		g.Expect(helpers.ContainsFile(execLogList, "bazel-out/darwin_arm64-fastbuild/bin/some_custom_dir/foo.txt")).To(BeTrue())
		g.Expect(helpers.ContainsFile(execLogList, "bazel-out/darwin_arm64-fastbuild/bin/some_custom_dir/bar.txt")).To(BeTrue())
		g.Expect(helpers.ContainsFile(execLogList, "bazel-out/darwin_arm64-fastbuild/bin/some_custom_dir/baz.txt")).To(BeTrue())
		g.Expect(helpers.ContainsFile(execLogList, "bazel-out/darwin_arm64-fastbuild/bin/some_custom_dir/baz.txt")).To(BeTrue())
		g.Expect(helpers.ContainsInputSet(execLogList, []string{
			"bazel-out/darwin_arm64-fastbuild/bin/some_custom_dir/foo.txt",
			"bazel-out/darwin_arm64-fastbuild/bin/some_custom_dir/bar.txt",
			"bazel-out/darwin_arm64-fastbuild/bin/some_custom_dir/baz.txt",
		})).To(BeTrue())
	})

	t.Run("spawns are processed", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		filepath := "../../../test-data/changed_file_7.3.0_1.compact.log"

		execLogChan := make(chan types.ExecLogEntry, 100)
		var compactWait sync.WaitGroup
		compactWait.Add(1)
		go Convert(filepath, execLogChan, &compactWait, ctx)
		compactWait.Wait()

		g.Expect(ctx.Continue()).To(BeTrue())
		diagnostics := ctx.GetDiagnostics()

		g.Expect(len(execLogChan)).To(Equal(15))
		g.Expect(stderr.String()).To(Equal(""))
		g.Expect(len(diagnostics)).To(Equal(0))

		close(execLogChan)

		execLogList := helpers.CreateListFromExecLogChan(execLogChan, t)

		g.Expect(helpers.ContainsSpawn(execLogList, "//foo:foo")).To(BeTrue())
		g.Expect(helpers.ContainsSpawn(execLogList, "//bar:bar")).To(BeTrue())
		g.Expect(helpers.ContainsSpawn(execLogList, "//baz:baz")).To(BeTrue())
	})

	t.Run("create golden files", func(t *testing.T) {
		files := map[string]string{
			"golden_files/changed_directory.txt": "../../../test-data/changed_directory_7.3.0_1.compact.log",
			"golden_files/changed_file.txt":      "../../../test-data/changed_file_7.3.0_1.compact.log",
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
