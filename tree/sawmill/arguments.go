package sawmill

import (
	"github.com/aspect-build/execlog-diff/result/proto"
	"github.com/aspect-build/execlog-diff/types"
)

func Arguments(sawmill *SawMill) differ {
	return &arguments_differ{
		sawmill: sawmill,
	}
}

type arguments_differ struct {
	sawmill *SawMill
}

func (arguments_differ *arguments_differ) Run(spawn1 *types.Spawn, spawn2 *types.Spawn) []*proto.InputItem {
	diffResults := make([]*proto.InputItem, 0)

	for i1, item := range spawn1.Args {
		i2 := stringListFind(spawn2.Args, item, i1)
		if i1 == i2 {
			// Arguments are the same in spawn 1 & 2
			continue
		} else if i2 == -1 {
			// Argument is present in spawn 1 but not spawn 2
			diffResults = append(diffResults, &proto.InputItem{
				Item: &proto.InputItem_ArgListChanges{
					ArgListChanges: &proto.ArgChange{
						OldArg:      item,
						OldArgIndex: int64(i1),
						NewArg:      "",
						NewArgIndex: -1,
					},
				},
				Reason: proto.Reason_REMOVED,
			})
		} else {
			// Argument is present in spawn 1 & spawn 2 but is not in the same location
			diffResults = append(diffResults, &proto.InputItem{
				Item: &proto.InputItem_ArgListChanges{
					ArgListChanges: &proto.ArgChange{
						OldArg:      item,
						OldArgIndex: int64(i1),
						NewArg:      item,
						NewArgIndex: int64(i2),
					},
				},
				Reason: proto.Reason_CHANGED,
			})
		}
	}

	for i2, item := range spawn2.Args {
		i1 := stringListFind(spawn1.Args, item, i2)
		if i1 == -1 {
			// Argument is present in spawn 2 but not spawn 1
			diffResults = append(diffResults, &proto.InputItem{
				Item: &proto.InputItem_ArgListChanges{
					ArgListChanges: &proto.ArgChange{
						OldArg:      "",
						OldArgIndex: -1,
						NewArg:      item,
						NewArgIndex: int64(i2),
					},
				},
				Reason: proto.Reason_ADDED,
			})
		}
	}

	return diffResults
}
