package terminal

import (
	"fmt"
	"strings"

	"github.com/aspect-build/execlog-diff/result/proto"
)

func (t *terminal) newLine() {
	fmt.Fprintln(t.stdout)
}

func (t *terminal) whitespace(n int) {
	fmt.Fprint(t.stdout, strings.Repeat(" ", n))
}

func (t *terminal) print(s string) {
	fmt.Fprint(t.stdout, s)
}

func (t *terminal) addToPrintBuffer(s string) {

}

// How many times is this node depended on directly or transitively
func TransitiveInfluenceCount(nodeActionHash string, nodes []*proto.ActionGraphDiffNode) int {
	count := 0
	for _, node := range nodes {
		if node.ActionHash != nodeActionHash {
			for _, child := range node.Children {
				if child.Child == nodeActionHash {
					count++
					count = count + TransitiveInfluenceCount(node.ActionHash, nodes)
				}
			}
		}
	}

	return count
}

func GetTarget(nodeActionHash string, nodes []*proto.ActionGraphDiffNode) string {
	for _, node := range nodes {
		if node.ActionHash == nodeActionHash {
			return node.TargetLabel
		}
	}

	return ""
}
