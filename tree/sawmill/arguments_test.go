package sawmill

import (
	"reflect"
	"strings"
	"testing"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/normalizer"
	"github.com/aspect-build/execlog-diff/result/proto"
	"github.com/aspect-build/execlog-diff/types"
	. "github.com/onsi/gomega"
)

func TestArguments(t *testing.T) {
	t.Run("2 identical spawns returns nothing", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"foo", "bar", "baz"},
		}

		spawn2 := &types.Spawn{
			Args: []string{"foo", "bar", "baz"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(0))
	})

	t.Run("argument that has been moved gets updated indexes (1)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"foo", "bar", "baz"},
		}

		spawn2 := &types.Spawn{
			Args: []string{"foo", "baz", "bar"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(2))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal("bar"))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal("bar"))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(1)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(2)))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[1].GetArgListChanges().NewArg).To(Equal("baz"))
		g.Expect(result[1].GetArgListChanges().OldArg).To(Equal("baz"))
		g.Expect(result[1].GetArgListChanges().OldArgIndex).To(Equal(int64(2)))
		g.Expect(result[1].GetArgListChanges().NewArgIndex).To(Equal(int64(1)))
	})

	t.Run("argument that has been moved gets updated indexes (2)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"foo", "bar", "baz"},
		}

		spawn2 := &types.Spawn{
			Args: []string{"baz", "bar", "foo"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(2))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal("foo"))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal("foo"))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(0)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(2)))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[1].GetArgListChanges().NewArg).To(Equal("baz"))
		g.Expect(result[1].GetArgListChanges().OldArg).To(Equal("baz"))
		g.Expect(result[1].GetArgListChanges().OldArgIndex).To(Equal(int64(2)))
		g.Expect(result[1].GetArgListChanges().NewArgIndex).To(Equal(int64(0)))
	})

	t.Run("argument that has been moved gets updated indexes (3)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"foo", "bar"},
		}

		spawn2 := &types.Spawn{
			Args: []string{"bar", "foo"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(2))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal("foo"))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal("foo"))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(0)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(1)))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[1].GetArgListChanges().NewArg).To(Equal("bar"))
		g.Expect(result[1].GetArgListChanges().OldArg).To(Equal("bar"))
		g.Expect(result[1].GetArgListChanges().OldArgIndex).To(Equal(int64(1)))
		g.Expect(result[1].GetArgListChanges().NewArgIndex).To(Equal(int64(0)))
	})

	t.Run("removed argument gets specified (1)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"foo", "bar", "baz"},
		}

		spawn2 := &types.Spawn{
			Args: []string{"bar", "baz"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(3))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal(""))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal("foo"))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(0)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(-1)))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[1].GetArgListChanges().NewArg).To(Equal("bar"))
		g.Expect(result[1].GetArgListChanges().OldArg).To(Equal("bar"))
		g.Expect(result[1].GetArgListChanges().OldArgIndex).To(Equal(int64(1)))
		g.Expect(result[1].GetArgListChanges().NewArgIndex).To(Equal(int64(0)))

		g.Expect(result[2].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[2].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[2].GetArgListChanges().NewArg).To(Equal("baz"))
		g.Expect(result[2].GetArgListChanges().OldArg).To(Equal("baz"))
		g.Expect(result[2].GetArgListChanges().OldArgIndex).To(Equal(int64(2)))
		g.Expect(result[2].GetArgListChanges().NewArgIndex).To(Equal(int64(1)))
	})

	t.Run("removed argument gets specified (2)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"foo", "bar", "baz"},
		}

		spawn2 := &types.Spawn{
			Args: []string{"foo", "baz"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(2))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal(""))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal("bar"))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(1)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(-1)))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[1].GetArgListChanges().NewArg).To(Equal("baz"))
		g.Expect(result[1].GetArgListChanges().OldArg).To(Equal("baz"))
		g.Expect(result[1].GetArgListChanges().OldArgIndex).To(Equal(int64(2)))
		g.Expect(result[1].GetArgListChanges().NewArgIndex).To(Equal(int64(1)))
	})

	t.Run("removed argument gets specified (3)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"foo", "bar", "baz"},
		}

		spawn2 := &types.Spawn{
			Args: []string{"foo", "bar"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(1))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal(""))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal("baz"))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(2)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(-1)))
	})

	t.Run("removed argument gets specified (4)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"foo"},
		}

		spawn2 := &types.Spawn{
			Args: []string{},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(1))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal(""))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal("foo"))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(0)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(-1)))
	})

	t.Run("added argument gets specified (1)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"foo", "bar"},
		}

		spawn2 := &types.Spawn{
			Args: []string{"foo", "bar", "baz"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(1))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal("baz"))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal(""))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(-1)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(2)))
	})

	t.Run("added argument gets specified (2)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"foo", "baz"},
		}

		spawn2 := &types.Spawn{
			Args: []string{"foo", "bar", "baz"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(2))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal("baz"))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal("baz"))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(1)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(2)))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[1].GetArgListChanges().NewArg).To(Equal("bar"))
		g.Expect(result[1].GetArgListChanges().OldArg).To(Equal(""))
		g.Expect(result[1].GetArgListChanges().OldArgIndex).To(Equal(int64(-1)))
		g.Expect(result[1].GetArgListChanges().NewArgIndex).To(Equal(int64(1)))
	})

	t.Run("added argument gets specified (3)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"bar", "baz"},
		}

		spawn2 := &types.Spawn{
			Args: []string{"foo", "bar", "baz"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(3))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal("bar"))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal("bar"))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(0)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(1)))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[1].GetArgListChanges().NewArg).To(Equal("baz"))
		g.Expect(result[1].GetArgListChanges().OldArg).To(Equal("baz"))
		g.Expect(result[1].GetArgListChanges().OldArgIndex).To(Equal(int64(1)))
		g.Expect(result[1].GetArgListChanges().NewArgIndex).To(Equal(int64(2)))

		g.Expect(result[2].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[2].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[2].GetArgListChanges().NewArg).To(Equal("foo"))
		g.Expect(result[2].GetArgListChanges().OldArg).To(Equal(""))
		g.Expect(result[2].GetArgListChanges().OldArgIndex).To(Equal(int64(-1)))
		g.Expect(result[2].GetArgListChanges().NewArgIndex).To(Equal(int64(0)))
	})

	t.Run("added argument gets specified (4)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{},
		}

		spawn2 := &types.Spawn{
			Args: []string{"foo"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(1))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal("foo"))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal(""))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(-1)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(0)))
	})

	t.Run("changed argument gets specified (1)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"foo1", "bar", "baz"},
		}

		spawn2 := &types.Spawn{
			Args: []string{"foo2", "bar", "baz"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(2))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal(""))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal("foo1"))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(0)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(-1)))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[1].GetArgListChanges().NewArg).To(Equal("foo2"))
		g.Expect(result[1].GetArgListChanges().OldArg).To(Equal(""))
		g.Expect(result[1].GetArgListChanges().OldArgIndex).To(Equal(int64(-1)))
		g.Expect(result[1].GetArgListChanges().NewArgIndex).To(Equal(int64(0)))
	})

	t.Run("changed argument gets specified (2)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"foo", "bar1", "baz"},
		}

		spawn2 := &types.Spawn{
			Args: []string{"foo", "bar2", "baz"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(2))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal(""))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal("bar1"))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(1)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(-1)))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[1].GetArgListChanges().NewArg).To(Equal("bar2"))
		g.Expect(result[1].GetArgListChanges().OldArg).To(Equal(""))
		g.Expect(result[1].GetArgListChanges().OldArgIndex).To(Equal(int64(-1)))
		g.Expect(result[1].GetArgListChanges().NewArgIndex).To(Equal(int64(1)))
	})

	t.Run("changed argument gets specified (3)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"foo", "bar", "baz1"},
		}

		spawn2 := &types.Spawn{
			Args: []string{"foo", "bar", "baz2"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(2))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal(""))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal("baz1"))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(2)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(-1)))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[1].GetArgListChanges().NewArg).To(Equal("baz2"))
		g.Expect(result[1].GetArgListChanges().OldArg).To(Equal(""))
		g.Expect(result[1].GetArgListChanges().OldArgIndex).To(Equal(int64(-1)))
		g.Expect(result[1].GetArgListChanges().NewArgIndex).To(Equal(int64(2)))
	})

	t.Run("multiple changes (1)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"foo", "bar", "baz"},
		}

		spawn2 := &types.Spawn{
			Args: []string{"bar"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(3))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal(""))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal("foo"))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(0)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(-1)))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[1].GetArgListChanges().NewArg).To(Equal("bar"))
		g.Expect(result[1].GetArgListChanges().OldArg).To(Equal("bar"))
		g.Expect(result[1].GetArgListChanges().OldArgIndex).To(Equal(int64(1)))
		g.Expect(result[1].GetArgListChanges().NewArgIndex).To(Equal(int64(0)))

		g.Expect(result[2].Reason).To(Equal(proto.Reason_REMOVED))
		g.Expect(reflect.TypeOf(result[2].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[2].GetArgListChanges().NewArg).To(Equal(""))
		g.Expect(result[2].GetArgListChanges().OldArg).To(Equal("baz"))
		g.Expect(result[2].GetArgListChanges().OldArgIndex).To(Equal(int64(2)))
		g.Expect(result[2].GetArgListChanges().NewArgIndex).To(Equal(int64(-1)))
	})

	t.Run("multiple changes (2)", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)

		normalizer := normalizer.New(ctx)

		sawmill := &SawMill{
			normalizer: normalizer,
		}

		spawn1 := &types.Spawn{
			Args: []string{"bar"},
		}

		spawn2 := &types.Spawn{
			Args: []string{"foo", "bar", "baz"},
		}

		result := Arguments(sawmill).Run(spawn1, spawn2)

		g.Expect(len(result)).To(Equal(3))

		g.Expect(result[0].Reason).To(Equal(proto.Reason_CHANGED))
		g.Expect(reflect.TypeOf(result[0].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[0].GetArgListChanges().NewArg).To(Equal("bar"))
		g.Expect(result[0].GetArgListChanges().OldArg).To(Equal("bar"))
		g.Expect(result[0].GetArgListChanges().OldArgIndex).To(Equal(int64(0)))
		g.Expect(result[0].GetArgListChanges().NewArgIndex).To(Equal(int64(1)))

		g.Expect(result[1].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[1].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[1].GetArgListChanges().NewArg).To(Equal("foo"))
		g.Expect(result[1].GetArgListChanges().OldArg).To(Equal(""))
		g.Expect(result[1].GetArgListChanges().OldArgIndex).To(Equal(int64(-1)))
		g.Expect(result[1].GetArgListChanges().NewArgIndex).To(Equal(int64(0)))

		g.Expect(result[2].Reason).To(Equal(proto.Reason_ADDED))
		g.Expect(reflect.TypeOf(result[2].Item)).To(Equal(reflect.TypeOf(&proto.InputItem_ArgListChanges{})))
		g.Expect(result[2].GetArgListChanges().NewArg).To(Equal("baz"))
		g.Expect(result[2].GetArgListChanges().OldArg).To(Equal(""))
		g.Expect(result[2].GetArgListChanges().OldArgIndex).To(Equal(int64(-1)))
		g.Expect(result[2].GetArgListChanges().NewArgIndex).To(Equal(int64(2)))
	})
}
