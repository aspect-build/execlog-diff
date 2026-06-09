package helpers

import (
	"testing"

	"github.com/aspect-build/execlog-diff/types"

	. "github.com/onsi/gomega"
)

func TestContainsFile(t *testing.T) {
	t.Run("returns true when ExecLogEntry contains the file", func(t *testing.T) {
		g := NewGomegaWithT(t)

		g.Expect(
			ContainsFile(
				[]types.ExecLogEntry{
					types.ExecLogEntry{
						File: &types.File{
							Path: "foobarbaz",
						},
					},
				},
				"foobarbaz",
			),
		).To(BeTrue())

		g.Expect(
			ContainsFile(
				[]types.ExecLogEntry{
					types.ExecLogEntry{
						InputSet: &types.InputSet{},
					},
					types.ExecLogEntry{
						Spawn: &types.IncomingSpawn{},
					},
					types.ExecLogEntry{
						File: &types.File{
							Path: "foo",
						},
					},
					types.ExecLogEntry{
						InputSet: &types.InputSet{},
					},
					types.ExecLogEntry{
						Spawn: &types.IncomingSpawn{},
					},
					types.ExecLogEntry{
						File: &types.File{
							Path: "bar",
						},
					},
					types.ExecLogEntry{
						InputSet: &types.InputSet{},
					},
					types.ExecLogEntry{
						Spawn: &types.IncomingSpawn{},
					},
					types.ExecLogEntry{
						File: &types.File{
							Path: "baz",
						},
					},
					types.ExecLogEntry{
						InputSet: &types.InputSet{},
					},
					types.ExecLogEntry{
						Spawn: &types.IncomingSpawn{},
					},
					types.ExecLogEntry{
						File: &types.File{
							Path: "foobarbaz",
						},
					},
					types.ExecLogEntry{
						InputSet: &types.InputSet{},
					},
					types.ExecLogEntry{
						Spawn: &types.IncomingSpawn{},
					},
				},
				"foobarbaz",
			),
		).To(BeTrue())
	})

	t.Run("returns false when ExecLogEntry does not contain the file", func(t *testing.T) {
		g := NewGomegaWithT(t)

		g.Expect(
			ContainsFile(
				[]types.ExecLogEntry{
					types.ExecLogEntry{
						InputSet: &types.InputSet{},
					},
					types.ExecLogEntry{
						Spawn: &types.IncomingSpawn{},
					},
					types.ExecLogEntry{
						File: &types.File{
							Path: "foo",
						},
					},
					types.ExecLogEntry{
						InputSet: &types.InputSet{},
					},
					types.ExecLogEntry{
						Spawn: &types.IncomingSpawn{},
					},
					types.ExecLogEntry{
						File: &types.File{
							Path: "bar",
						},
					},
					types.ExecLogEntry{
						InputSet: &types.InputSet{},
					},
					types.ExecLogEntry{
						Spawn: &types.IncomingSpawn{},
					},
					types.ExecLogEntry{
						File: &types.File{
							Path: "baz",
						},
					},
					types.ExecLogEntry{
						InputSet: &types.InputSet{},
					},
					types.ExecLogEntry{
						Spawn: &types.IncomingSpawn{},
					},
					types.ExecLogEntry{
						File: &types.File{
							Path: "foobarbaz",
						},
					},
					types.ExecLogEntry{
						InputSet: &types.InputSet{},
					},
					types.ExecLogEntry{
						Spawn: &types.IncomingSpawn{},
					},
				},
				"foobar",
			),
		).To(BeFalse())
	})
}
