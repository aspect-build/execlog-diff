package main

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"fmt"
)

type Streams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

var DefaultStreams = Streams{
	Stdin:  os.Stdin,
	Stdout: os.Stdout,
	Stderr: os.Stderr,
}

func NewDefaultCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Short: "CLI to analyze bazel execution logs.",
		Long:  "CLI to analyze bazel execution logs.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Please select a run mode")
		},
	}

	cmd.AddCommand(NewDefaultExplainCmd())

	return cmd
}

func NewCmd() *cobra.Command {
	return NewDefaultCmd()
}

func main() {
	if err := NewDefaultCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
