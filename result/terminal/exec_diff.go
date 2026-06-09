package terminal

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/aspect-build/execlog-diff/result/proto"
)

func (t *terminal) printExecDiff(execDiff *proto.ExecDiff) {
	sort.Slice(execDiff.Nodes, func(i, j int) bool {
		return execDiff.Nodes[i].TargetLabel < execDiff.Nodes[j].TargetLabel
	})
	for _, node := range execDiff.Nodes {
		t.printNode(node, execDiff.Nodes)
	}
}

func (t *terminal) printNode(node *proto.ActionGraphDiffNode, nodes []*proto.ActionGraphDiffNode) {
	t.printInColor(node.TargetLabel, GREEN)

	if node.RemoteCachable {
		t.printInColor(" (WARNING: Target is cacheable)  ", PURPLE)
	}

	rebuilds := TransitiveInfluenceCount(node.ActionHash, nodes)
	if rebuilds > 0 {
		t.whitespace(5)

		t.printInColor("(Caused ", RED)
		t.printInColor(strconv.Itoa(TransitiveInfluenceCount(node.ActionHash, nodes)), RED)
		t.printInColor(" Rebuilds)", RED)
	}

	t.printOutputChanges(node)

	t.printInputChanges(node)

	t.printChangedDependencies(node, nodes)

	t.newLine()
	t.newLine()
}

func (t *terminal) printOutputChanges(node *proto.ActionGraphDiffNode) {
	headerPrinted := false
	sort.Slice(node.Changes, func(i, j int) bool {
		return node.Changes[i].String() < node.Changes[j].String()
	})
	for _, change := range node.Changes {
		switch item := change.Item.(type) {
		case *proto.InputItem_OutputFileListChanges:
			if !headerPrinted {
				t.newLine()
				t.whitespace(4)

				t.print("Changed Output Files")

				headerPrinted = true
			}

			t.newLine()
			t.whitespace(8)
			t.print("-")
			t.whitespace(1)
			t.print(item.OutputFileListChanges.Path)

			t.whitespace(2)
			t.printInColor("(", CYAN)
			t.printBytes(item.OutputFileListChanges.SizeBytes)
			t.printInColor(")", CYAN)
		}
	}

	if headerPrinted {
		t.newLine()
	}
}

// TODO: make this better. Do one of the following things
// 1 - add a (src) / (dep) before each thing so its clear where it came from
// 2 - split out into 2 lists. 1 for src changes & 1 for changes
// 3 - Do not print dep changes (thos are already shown in printChangedDependencies)
func (t *terminal) printInputChanges(node *proto.ActionGraphDiffNode) {
	headerPrinted := false
	sort.Slice(node.Changes, func(i, j int) bool {
		return node.Changes[i].String() < node.Changes[j].String()
	})
	for _, change := range node.Changes {
		switch item := change.Item.(type) {
		case *proto.InputItem_InputFileListChanges:
			if !headerPrinted {
				t.newLine()
				t.whitespace(4)

				t.print("Changed Input Files")

				headerPrinted = true
			}

			t.newLine()
			t.whitespace(8)
			t.print("-")
			t.whitespace(1)
			t.print(item.InputFileListChanges.Path)

			t.whitespace(2)
			t.printInColor("(", CYAN)
			t.printBytes(item.InputFileListChanges.SizeBytes)
			t.printInColor(")", CYAN)
		}
	}

	if headerPrinted {
		t.newLine()
	}
}

func (t *terminal) printChangedDependencies(node *proto.ActionGraphDiffNode, nodes []*proto.ActionGraphDiffNode) {
	headerPrinted := false
	sort.Slice(node.Children, func(i, j int) bool {
		return GetTarget(node.Children[i].Child, nodes) < GetTarget(node.Children[j].Child, nodes)
	})
	for _, child := range node.Children {
		if !headerPrinted {
			t.newLine()
			t.whitespace(4)

			t.print("Changed Dependencies")

			headerPrinted = true
		}

		t.newLine()
		t.whitespace(8)

		// TODO: handle "" case
		t.printInColor(GetTarget(child.Child, nodes), YELLOW)

		for _, file := range child.ChangedFiles {
			switch item := file.Item.(type) {
			case *proto.InputItem_InputFileListChanges:

				t.newLine()
				t.whitespace(12)
				t.print("-")
				t.whitespace(1)
				t.print(item.InputFileListChanges.Path)
				t.newLine()
			}
		}
	}
}

func (t *terminal) printBytes(b int64) {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case b >= GB:
		t.printInColor(fmt.Sprintf("%.2f GB", float64(b)/float64(GB)), CYAN)
	case b >= MB:
		t.printInColor(fmt.Sprintf("%.2f MB", float64(b)/float64(MB)), CYAN)
	case b >= KB:
		t.printInColor(fmt.Sprintf("%.2f KB", float64(b)/float64(KB)), CYAN)
	default:
		t.printInColor(fmt.Sprintf("%d Bytes", b), CYAN)
	}
}
