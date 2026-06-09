package main

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	execlogdiff "github.com/aspect-build/execlog-diff"
	"github.com/aspect-build/execlog-diff/result/html"
	"github.com/aspect-build/execlog-diff/result/terminal"
	"github.com/aspect-build/execlog-diff/diagnostic"
)

type Explain struct {
	Streams
}

func NewDefaultExplainCmd() *cobra.Command {
	return NewExplainCmd(DefaultStreams)
}

func NewExplainCmd(streams Streams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "explain <execlog1> <execlog2>",
		Short: "Compare two Bazel execution logs to detect non-hermetic build issues.",
		Long: `OVERVIEW: Compare two Bazel execution logs to detect non-hermetic build issues
EXAMPLE: aspect explain ./compact-execlog.1.log ./compact-execlog.2.log --cache-misses
Detects build non-determinism by analyzing:
- Actions: Changed outputs, timing, or ordering
- Files: Unexpected file modifications
- Environment: Non-hermetic environment variables
- Arguments: Inconsistent build flags
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return New(streams).Run(cmd, args)
		},
	}

	cmd.Flags().Bool("cache-misses", false, "Find bazel actions that have changed and would not have been cached between runs")
	cmd.Flags().Bool("list", false, "List all labels")
	cmd.Flags().Bool("html", false, "Generate HTML output (only valid with --cache-misses)")
	cmd.Flags().String("html-output", "", "Path to save HTML output (implies --html, only valid with --cache-misses)")

	return cmd
}

func New(streams Streams) *Explain {
	return &Explain{
		Streams: streams,
	}
}

func (runner *Explain) Run(cmd *cobra.Command, args []string) error {
	cacheMisses, err := cmd.Flags().GetBool("cache-misses")
	if err != nil {
		return err
	}
	listLabels, err := cmd.Flags().GetBool("list")
	if err != nil {
		return err
	}
	htmlOutput, err := cmd.Flags().GetBool("html")
	if err != nil {
		return err
	}
	htmlOutputPath, err := cmd.Flags().GetString("html-output")
	if err != nil {
		return err
	}

	if htmlOutputPath != "" {
		htmlOutput = true
	}
	if htmlOutput && !cacheMisses {
		return errors.New("--html can only be used with --cache-misses")
	}

	switch {
	case cacheMisses:
		if len(args) != 2 {
			return errors.New("Expect 2 arguments when running with --cache-misses")
		}

		var log_1 = args[0]
		var log_2 = args[1]

		result := execlogdiff.DiffExecLogs(log_1, log_2)
		if len(result.RuntimeDiagnostics) > 0 {
			return diagnosticsToError(result.RuntimeDiagnostics)
		}

		var printDiagnostics []*diagnostic.Diagnostic
		if htmlOutput {
			outputPath := htmlOutputPath
			if outputPath == "" {
				outputPath = "execlog-diff-cache-misses.html"
			}
			htmlPrinter := html.New(outputPath)
			printDiagnostics = htmlPrinter.Print(result.Result)
		} else {
			term := terminal.NewDefault()
			printDiagnostics = term.Print(result.Result)
		}
		if len(printDiagnostics) > 0 {
			return diagnosticsToError(printDiagnostics)
		}
	case listLabels:
		if len(args) != 1 {
			return errors.New("Expect 1 argument when running with --list")
		}

		var log = args[0]

		result := execlogdiff.ListLabels(log)
		if len(result.RuntimeDiagnostics) > 0 {
			return diagnosticsToError(result.RuntimeDiagnostics)
		}

		term := terminal.NewDefault()
		printDiagnostics := term.Print(result.Result)
		if len(printDiagnostics) > 0 {
			return diagnosticsToError(printDiagnostics)
		}
	default:
		fmt.Fprintln(runner.Stdout, cmd.Commands())
		fmt.Fprintln(runner.Stdout, "Possible modes:")
		fmt.Fprintln(runner.Stdout, "   --cache-misses      Attempt to find why a bazel action changed between execlogs")
		fmt.Fprintln(runner.Stdout, "   --list              List all labels")
	}

	return nil
}

func diagnosticsToError(diagnostics []*diagnostic.Diagnostic) error {
	err := errors.New(diagnosticToString(diagnostics[0]))
	for _, diagnostic := range diagnostics[1:] {
		err = fmt.Errorf("%w; %s", err, diagnosticToString(diagnostic))
	}

	return err
}

func diagnosticToString(diagnostic *diagnostic.Diagnostic) string {
	result := ""
	if diagnostic.Title != "" {
		result = result + diagnostic.Title + ": "
	}

	result = result + diagnostic.Message

	return result
}
