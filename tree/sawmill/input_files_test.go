package sawmill

import (
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/result/proto"
	"github.com/aspect-build/execlog-diff/test-data/helpers"
	"github.com/aspect-build/execlog-diff/types"
	. "github.com/onsi/gomega"
)

func TestInputFiles(t *testing.T) {
	t.Run("2 identical spawns returns nothing", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log1File("foo.txt", "abc")
		th.Log2File("foo.txt", "abc")

		th.Log1File("bar.txt", "abcd")
		th.Log2File("bar.txt", "abcd")

		th.Log1File("baz.txt", "abcde")
		th.Log2File("baz.txt", "abcde")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		spawn1 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog1Path(),
		}

		spawn2 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog2Path(),
		}

		result := InputFiles(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(0))
	})

	t.Run("changed input file gets flagged (1)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log1File("foo.txt", "abcd")
		th.Log2File("foo.txt", "abc")

		th.Log1File("bar.txt", "abcd")
		th.Log2File("bar.txt", "abcd")

		th.Log1File("baz.txt", "abcde")
		th.Log2File("baz.txt", "abcde")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		spawn1 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog1Path(),
		}

		spawn2 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog2Path(),
		}

		result := InputFiles(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(1))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[0].GetInputFileListChanges().Path).To(Equal("foo.txt"))
	})

	t.Run("changed input file gets flagged (2)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log1File("foo.txt", "abc")
		th.Log2File("foo.txt", "abc")

		th.Log1File("bar.txt", "abcd")
		th.Log2File("bar.txt", "abcdef")

		th.Log1File("baz.txt", "abcde")
		th.Log2File("baz.txt", "abcde")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		spawn1 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog1Path(),
		}

		spawn2 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog2Path(),
		}

		result := InputFiles(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(1))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[0].GetInputFileListChanges().Path).To(Equal("bar.txt"))
	})

	t.Run("changed input file gets flagged (3)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log1File("foo.txt", "abc")
		th.Log2File("foo.txt", "abc")

		th.Log1File("bar.txt", "abcd")
		th.Log2File("bar.txt", "abcd")

		th.Log1File("baz.txt", "abcde")
		th.Log2File("baz.txt", "abcdef")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		spawn1 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog1Path(),
		}

		spawn2 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog2Path(),
		}

		result := InputFiles(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(1))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[0].GetInputFileListChanges().Path).To(Equal("baz.txt"))
	})

	t.Run("changed input file gets flagged (4)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log1File("foo.txt", "abc")
		th.Log2File("foo.txt", "abcd")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		spawn1 := &types.Spawn{
			InputFiles:  []string{"foo.txt"},
			LogPathHash: th.GetLog1Path(),
		}

		spawn2 := &types.Spawn{
			InputFiles:  []string{"foo.txt"},
			LogPathHash: th.GetLog2Path(),
		}

		result := InputFiles(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(1))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[0].GetInputFileListChanges().Path).To(Equal("foo.txt"))
	})

	t.Run("changed input files gets flagged", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log1File("foo.txt", "abc")
		th.Log2File("foo.txt", "abcd")

		th.Log1File("bar.txt", "abcde")
		th.Log2File("bar.txt", "abcd")

		th.Log1File("baz.txt", "abcde")
		th.Log2File("baz.txt", "abcdef")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		spawn1 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog1Path(),
		}

		spawn2 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog2Path(),
		}

		result := InputFiles(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(3))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[0].GetInputFileListChanges().Path).To(Equal("foo.txt"))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[1].GetInputFileListChanges().Path).To(Equal("bar.txt"))

		g.Expect(result[2].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[2].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[2].GetInputFileListChanges().Path).To(Equal("baz.txt"))
	})

	t.Run("added input files gets flagged (1)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log2File("foo.txt", "abc")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		spawn1 := &types.Spawn{
			InputFiles:  []string{},
			LogPathHash: th.GetLog1Path(),
		}

		spawn2 := &types.Spawn{
			InputFiles:  []string{"foo.txt"},
			LogPathHash: th.GetLog2Path(),
		}

		result := InputFiles(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(1))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[0].GetInputFileListChanges().Path).To(Equal("foo.txt"))
	})

	t.Run("added input files gets flagged (2)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log2File("foo.txt", "abc")

		th.Log2File("bar.txt", "abcd")

		th.Log2File("baz.txt", "abcde")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		spawn1 := &types.Spawn{
			InputFiles:  []string{},
			LogPathHash: th.GetLog1Path(),
		}

		spawn2 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog2Path(),
		}

		result := InputFiles(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(3))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[0].GetInputFileListChanges().Path).To(Equal("foo.txt"))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[1].GetInputFileListChanges().Path).To(Equal("bar.txt"))

		g.Expect(result[2].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[2].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[2].GetInputFileListChanges().Path).To(Equal("baz.txt"))
	})

	t.Run("added input files gets flagged (3)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log1File("foo.txt", "abc")
		th.Log2File("foo.txt", "abc")

		th.Log2File("bar.txt", "abcd")

		th.Log2File("baz.txt", "abcde")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		spawn1 := &types.Spawn{
			InputFiles:  []string{"foo.txt"},
			LogPathHash: th.GetLog1Path(),
		}

		spawn2 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog2Path(),
		}

		result := InputFiles(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(2))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[0].GetInputFileListChanges().Path).To(Equal("bar.txt"))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[1].GetInputFileListChanges().Path).To(Equal("baz.txt"))
	})

	t.Run("removed input files gets flagged (1)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log1File("foo.txt", "abc")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		spawn1 := &types.Spawn{
			InputFiles:  []string{"foo.txt"},
			LogPathHash: th.GetLog1Path(),
		}

		spawn2 := &types.Spawn{
			InputFiles:  []string{},
			LogPathHash: th.GetLog2Path(),
		}

		result := InputFiles(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(1))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[0].GetInputFileListChanges().Path).To(Equal("foo.txt"))
	})

	t.Run("removed input files gets flagged (2)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log1File("foo.txt", "abc")
		th.Log1File("bar.txt", "abcd")
		th.Log1File("baz.txt", "abcde")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		spawn1 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog1Path(),
		}

		spawn2 := &types.Spawn{
			InputFiles:  []string{},
			LogPathHash: th.GetLog2Path(),
		}

		result := InputFiles(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(3))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[0].GetInputFileListChanges().Path).To(Equal("foo.txt"))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[1].GetInputFileListChanges().Path).To(Equal("bar.txt"))

		g.Expect(result[2].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[2].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[2].GetInputFileListChanges().Path).To(Equal("baz.txt"))
	})

	t.Run("removed input files gets flagged (3)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log1File("foo.txt", "abc")
		th.Log2File("foo.txt", "abc")

		th.Log1File("bar.txt", "abcd")
		th.Log1File("baz.txt", "abcde")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		spawn1 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt", "baz.txt"},
			LogPathHash: th.GetLog1Path(),
		}

		spawn2 := &types.Spawn{
			InputFiles:  []string{"foo.txt"},
			LogPathHash: th.GetLog2Path(),
		}

		result := InputFiles(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(2))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[0].GetInputFileListChanges().Path).To(Equal("bar.txt"))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[1].GetInputFileListChanges().Path).To(Equal("baz.txt"))
	})

	t.Run("files get added, removed and changed", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)
		var normalizerWait sync.WaitGroup
		normalizerWait.Add(1)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		execlogChan := make(chan types.ExecLogEntry)
		spawnExecLogChan := make(chan types.Spawn, 100)

		go normalizer.Normalize(execlogChan, spawnExecLogChan, &normalizerWait)

		th := helpers.NewPopulateHelper(execlogChan)

		th.Log1File("foo.txt", "abc")
		th.Log2File("foo.txt", "abcd")

		th.Log1File("bar.txt", "abcde")

		th.Log2File("baz.txt", "abcdef")

		close(execlogChan)

		normalizerWait.Wait()

		close(spawnExecLogChan)

		spawn1 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "bar.txt"},
			LogPathHash: th.GetLog1Path(),
		}

		spawn2 := &types.Spawn{
			InputFiles:  []string{"foo.txt", "baz.txt"},
			LogPathHash: th.GetLog2Path(),
		}

		result := InputFiles(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(3))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[0].GetInputFileListChanges().Path).To(Equal("foo.txt"))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[1].GetInputFileListChanges().Path).To(Equal("bar.txt"))

		g.Expect(result[2].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[2].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_InputFileListChanges{})))
		g.Expect(result[2].GetInputFileListChanges().Path).To(Equal("baz.txt"))
	})
}
