package execlogdiff

import (
	"sync"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/converter"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/result/proto"
	"github.com/aspect-build/execlog-diff/tree"
	"github.com/aspect-build/execlog-diff/types"
)

type execLogDiff struct {
	ctx *context.Context

	// This chan will contain execlog entries from both the first and second execlog
	execlogChan chan types.ExecLogEntry

	// This chan will contain the spawn execlog entries from "execlogChan" after they have been processed
	spawnChan chan types.Spawn

	normalizer *normalizer.Normalizer
	diffTree   tree.Tree
	queryTree  tree.Tree
}

func DiffExecLogs(execlog_path_1 string, execlog_path_2 string) *proto.ExecLogDiffResult {
	ctx := context.NewDiff()
	n := normalizer.New(ctx)

	eld := &execLogDiff{
		ctx:         ctx,
		execlogChan: make(chan types.ExecLogEntry),
		spawnChan:   make(chan types.Spawn),

		normalizer: n,
		diffTree:   tree.NewDiffTree(ctx, n),
	}

	return eld.diffExecLogs(execlog_path_1, execlog_path_2)
}

func ListLabels(execlog_path string) *proto.ExecLogDiffResult {
	return QueryExecLog(execlog_path, context.ListLabels)
}

func QueryExecLog(execlog_path string, mode context.RunMode) *proto.ExecLogDiffResult {
	ctx := context.NewQuery(mode)
	n := normalizer.New(ctx)

	eld := &execLogDiff{
		ctx:         ctx,
		execlogChan: make(chan types.ExecLogEntry),
		spawnChan:   make(chan types.Spawn),

		normalizer: n,
		queryTree:  tree.NewQueryTree(ctx, n),
	}

	return eld.queryExecLog(execlog_path)
}

func (f *execLogDiff) diffExecLogs(execlog_path_1 string, execlog_path_2 string) *proto.ExecLogDiffResult {
	// This wait group is used to ensure that we have finished converting the 2 incoming execlogs.
	// When "Done" all from BOTH execlogs should be on "execlogChan".
	// It will not be "add"ed to anywhere else.
	var execlogWait sync.WaitGroup
	execlogWait.Add(2)

	go converter.Convert(execlog_path_1, f.execlogChan, &execlogWait, f.ctx)
	go converter.Convert(execlog_path_2, f.execlogChan, &execlogWait, f.ctx)

	// This wait group will only be used to wait on the "normalizer.Normalize" function completing.
	var normalizerWait sync.WaitGroup

	// normalizer.Normalize() is currently only called once in this context.
	// TODO: Start multiple normalizer.Normalize() go routines to increase throughput
	// This wait group will not be "add"ed to anywhere else.
	normalizerWait.Add(1)
	go f.normalizer.Normalize(f.execlogChan, f.spawnChan, &normalizerWait)

	var diffTreeWait sync.WaitGroup
	diffTreeWait.Add(1)
	go f.diffTree.Grow(f.spawnChan, &diffTreeWait)

	execlogWait.Wait()

	close(f.execlogChan)

	normalizerWait.Wait()
	close(f.spawnChan)

	diffTreeWait.Wait()

	return &proto.ExecLogDiffResult{
		RuntimeDiagnostics: f.ctx.GetDiagnostics(),
		Result:             f.diffTree.CreateResult(),
	}
}

func (f *execLogDiff) queryExecLog(execlog_path string) *proto.ExecLogDiffResult {
	var execlogWait sync.WaitGroup
	execlogWait.Add(1)

	go converter.Convert(execlog_path, f.execlogChan, &execlogWait, f.ctx)

	var normalizerWait sync.WaitGroup
	normalizerWait.Add(1)
	go f.normalizer.Normalize(f.execlogChan, f.spawnChan, &normalizerWait)

	var queryTreeWait sync.WaitGroup
	queryTreeWait.Add(1)
	go f.queryTree.Grow(f.spawnChan, &queryTreeWait)

	execlogWait.Wait()

	close(f.execlogChan)

	normalizerWait.Wait()
	close(f.spawnChan)

	queryTreeWait.Wait()

	return &proto.ExecLogDiffResult{
		RuntimeDiagnostics: f.ctx.GetDiagnostics(),
		Result:             f.queryTree.CreateResult(),
	}
}
