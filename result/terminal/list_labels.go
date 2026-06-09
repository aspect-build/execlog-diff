package terminal

import (
	"sort"

	"github.com/aspect-build/execlog-diff/result/proto"
)

func (t *terminal) printListLabels(listLabels *proto.ListLabelsQuery) {
	labels := listLabels.GetLabels()
	sort.Slice(labels, func(i, j int) bool {
		return labels[i].GetTargetLabel() < labels[j].GetTargetLabel()
	})
	for _, label := range labels {
		if label.GetCacheHit() {
			t.printInColor(label.GetTargetLabel(), GREEN)
		} else {
			t.printInColor(label.GetTargetLabel(), YELLOW)
		}
		t.newLine()
	}
}
