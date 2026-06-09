package context

import (
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/diagnostic"
	. "github.com/onsi/gomega"
)

func TestContext(t *testing.T) {
	t.Run("Prints a title and message diagnostic and returns it when requested", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := New(&stderr, Diff)

		message := "some warning"
		title := "Some Title"
		ctx.Warning(&diagnostic.Diagnostic{
			Title:   title,
			Message: message,
		})

		g.Expect(ctx.Continue()).To(BeTrue())

		diagnostics := ctx.GetDiagnostics()

		g.Expect(stderr.String()).To(HavePrefix("WARNING:"))
		g.Expect(stderr.String()).To(ContainSubstring(title))
		g.Expect(stderr.String()).To(ContainSubstring(message))
		g.Expect(len(diagnostics)).To(Equal(1))
		g.Expect(diagnostics[0].Title).To(Equal(title))
		g.Expect(diagnostics[0].Message).To(Equal(message))
		g.Expect(ctx.Continue()).To(BeFalse())
	})

	t.Run("Prints a message to stderr and returns it when requested", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := New(&stderr, Diff)

		message := "some warning"
		ctx.Warning(&diagnostic.Diagnostic{
			Message: message,
		})

		g.Expect(ctx.Continue()).To(BeTrue())

		diagnostics := ctx.GetDiagnostics()

		g.Expect(stderr.String()).To(HavePrefix("WARNING:"))
		g.Expect(stderr.String()).To(ContainSubstring(message))
		g.Expect(len(diagnostics)).To(Equal(1))
		g.Expect(diagnostics[0].Message).To(Equal(message))
		g.Expect(ctx.Continue()).To(BeFalse())
	})

	t.Run("Prints multiple messages to stderr and returns them when requested", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := New(&stderr, Diff)

		message1 := "some warning 1"
		ctx.Warning(&diagnostic.Diagnostic{
			Message: message1,
		})

		message2 := "some warning 2"
		ctx.Warning(&diagnostic.Diagnostic{
			Message: message2,
		})

		message3 := "some warning 3"
		ctx.Warning(&diagnostic.Diagnostic{
			Message: message3,
		})

		g.Expect(ctx.Continue()).To(BeTrue())

		diagnostics := ctx.GetDiagnostics()

		g.Expect(stderr.String()).To(HavePrefix("WARNING:"))
		g.Expect(stderr.String()).To(ContainSubstring(message1))
		g.Expect(stderr.String()).To(ContainSubstring(message2))
		g.Expect(stderr.String()).To(ContainSubstring(message3))
		g.Expect(len(diagnostics)).To(Equal(3))
		g.Expect(diagnostics[0].Severity).To(Equal(diagnostic.Severity_WARNING))
		g.Expect(diagnostics[1].Severity).To(Equal(diagnostic.Severity_WARNING))
		g.Expect(diagnostics[2].Severity).To(Equal(diagnostic.Severity_WARNING))
		g.Expect(ctx.Continue()).To(BeFalse())
	})

	t.Run("Prints error and sets continue to false", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := New(&stderr, Diff)

		message := "some warning"
		ctx.Error(&diagnostic.Diagnostic{
			Message: message,
		})

		g.Expect(ctx.Continue()).To(BeFalse())

		diagnostics := ctx.GetDiagnostics()

		g.Expect(stderr.String()).To(HavePrefix("ERROR:"))
		g.Expect(stderr.String()).To(ContainSubstring(message))
		g.Expect(len(diagnostics)).To(Equal(1))
		g.Expect(diagnostics[0].Message).To(Equal(message))
		g.Expect(ctx.Continue()).To(BeFalse())
	})

	t.Run("Prints error with title and sets continue to false", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := New(&stderr, Diff)

		message := "some warning"
		title := "Some Title"
		ctx.Error(&diagnostic.Diagnostic{
			Title:   title,
			Message: message,
		})

		g.Expect(ctx.Continue()).To(BeFalse())

		diagnostics := ctx.GetDiagnostics()

		g.Expect(stderr.String()).To(HavePrefix("ERROR:"))
		g.Expect(stderr.String()).To(ContainSubstring(title))
		g.Expect(stderr.String()).To(ContainSubstring(message))
		g.Expect(len(diagnostics)).To(Equal(1))
		g.Expect(diagnostics[0].Message).To(Equal(message))
		g.Expect(ctx.Continue()).To(BeFalse())
	})

	t.Run("Does not throw any errors if diagnostics passed after destructive channel close", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := New(&stderr, Diff)

		message := "some warning"
		title := "Some Title"
		ctx.Error(&diagnostic.Diagnostic{
			Title:   title,
			Message: message,
		})

		g.Expect(ctx.Continue()).To(BeFalse())

		diagnostics := ctx.GetDiagnostics()

		g.Expect(stderr.String()).To(HavePrefix("ERROR:"))
		g.Expect(stderr.String()).To(ContainSubstring(title))
		g.Expect(stderr.String()).To(ContainSubstring(message))
		g.Expect(len(diagnostics)).To(Equal(1))
		g.Expect(diagnostics[0].Message).To(Equal(message))
		g.Expect(ctx.Continue()).To(BeFalse())

		diagnostics = ctx.GetDiagnostics()
		g.Expect(len(diagnostics)).To(Equal(1))
		g.Expect(ctx.Continue()).To(BeFalse())

		ctx.Diagnostic(&diagnostic.Diagnostic{
			Title:   "foo1",
			Message: "bar1",
		})

		ctx.Warning(&diagnostic.Diagnostic{
			Title:   "foo2",
			Message: "bar2",
		})

		ctx.Error(&diagnostic.Diagnostic{
			Title:   "foo3",
			Message: "bar3",
		})

		diagnostics = ctx.GetDiagnostics()
		g.Expect(len(diagnostics)).To(Equal(1))
		g.Expect(ctx.Continue()).To(BeFalse())
	})
}
