package bazel8

import (
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/diagnostic"
	. "github.com/onsi/gomega"
)

func TestInvocationIsBazel8(t *testing.T) {
	t.Run("errors opening files are surfaced", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		filepath := "./some_nonexistant_file_123456.log"

		result := InvocationIsBazel8(filepath, ctx)

		diagnostics := ctx.GetDiagnostics()

		g.Expect(stderr.String()).To(HavePrefix("ERROR: Failed to open file"))
		g.Expect(stderr.String()).To(ContainSubstring(filepath))
		g.Expect(stderr.String()).To(ContainSubstring("no such file or directory"))
		g.Expect(ctx.Continue()).To(BeFalse())
		g.Expect(len(diagnostics)).To(Equal(1))
		g.Expect(diagnostics[0].Severity).To(Equal(diagnostic.Severity_ERROR))
		g.Expect(result).To(BeFalse())
	})

	t.Run("decompression errors are surfaced", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		filepath := "../../../test-data/changed_file_7.3.0_1.binary.log"
		result := InvocationIsBazel8(filepath, ctx)

		diagnostics := ctx.GetDiagnostics()

		g.Expect(stderr.String()).To(HavePrefix("ERROR: Failed to decompress"))
		g.Expect(stderr.String()).To(ContainSubstring(filepath))
		g.Expect(ctx.Continue()).To(BeFalse())
		g.Expect(len(diagnostics)).To(Equal(1))
		g.Expect(diagnostics[0].Severity).To(Equal(diagnostic.Severity_ERROR))
		g.Expect(result).To(BeFalse())
	})

	t.Run("bazel 8 execlogs return true", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		result := InvocationIsBazel8("../../../test-data/changed_file_8.0.0_1.compact.log", ctx)

		g.Expect(ctx.Continue()).To(BeTrue())
		diagnostics := ctx.GetDiagnostics()

		g.Expect(stderr.String()).To(Equal(""))
		g.Expect(len(diagnostics)).To(Equal(0))
		g.Expect(result).To(BeTrue())
	})

	t.Run("bazel 7 execlogs return false", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		result := InvocationIsBazel8("../../../test-data/changed_file_7.3.0_1.compact.log", ctx)

		g.Expect(ctx.Continue()).To(BeTrue())
		diagnostics := ctx.GetDiagnostics()

		g.Expect(stderr.String()).To(Equal(""))
		g.Expect(len(diagnostics)).To(Equal(0))
		g.Expect(result).To(BeFalse())
	})
}
