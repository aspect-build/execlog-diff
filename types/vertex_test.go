package types

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestVertex(t *testing.T) {
	t.Run("can save spawns to vertex", func(t *testing.T) {
		g := NewGomegaWithT(t)

		vertex := Vertex{}

		vertex.SetSpawn("foo", &Spawn{
			Label: "//bar",
		})

		vertex.SetSpawn("baz", &Spawn{
			Label: "//foo/bar/baz",
		})

		g.Expect(len(vertex.Spawns)).To(Equal(2))
		g.Expect(vertex.Spawns["foo"].Label).To(Equal("//bar"))
		g.Expect(vertex.Spawns["baz"].Label).To(Equal("//foo/bar/baz"))
	})

	t.Run("can get spawns from vertex", func(t *testing.T) {
		g := NewGomegaWithT(t)

		spawns := make(map[string]*Spawn)
		spawns["foo"] = &Spawn{
			Label: "//bar",
		}
		spawns["baz"] = &Spawn{
			Label: "//foo/bar/baz",
		}
		vertex := Vertex{
			Spawns: spawns,
		}

		spawn1 := vertex.GetSpawn("foo")
		g.Expect(spawn1).NotTo(BeNil())
		g.Expect(spawn1.Label).To(Equal("//bar"))

		spawn2 := vertex.GetSpawn("baz")
		g.Expect(spawn2).NotTo(BeNil())
		g.Expect(spawn2.Label).To(Equal("//foo/bar/baz"))
	})

	t.Run("returns nil if spawn map has not been created", func(t *testing.T) {
		g := NewGomegaWithT(t)

		vertex := Vertex{}

		spawn := vertex.GetSpawn("foo")

		g.Expect(spawn).To(BeNil())
	})

	t.Run("HasSpawnForLog: can confirm if vertex contains spawn", func(t *testing.T) {
		g := NewGomegaWithT(t)

		spawns := make(map[string]*Spawn)
		spawns["foo"] = &Spawn{
			Label: "//bar",
		}
		vertex := Vertex{
			Spawns: spawns,
		}

		g.Expect(vertex.HasSpawnForLog("foo")).To(BeTrue())
		g.Expect(vertex.HasSpawnForLog("bar")).To(BeFalse())
	})

	t.Run("HasSpawnForLog: returns false if map is nil", func(t *testing.T) {
		g := NewGomegaWithT(t)

		vertex := Vertex{}

		g.Expect(vertex.HasSpawnForLog("bar")).To(BeFalse())
	})

	t.Run("can get log path hashes", func(t *testing.T) {
		g := NewGomegaWithT(t)

		vertex := Vertex{}

		vertex.SetSpawn("foo", &Spawn{
			LogPathHash: "foo",
			Label:       "//foo",
		})

		vertex.SetSpawn("bar", &Spawn{
			LogPathHash: "bar",
			Label:       "//bar",
		})

		vertex.SetSpawn("baz", &Spawn{
			LogPathHash: "baz",
			Label:       "//baz",
		})

		g.Expect(len(vertex.GetLogPathHashes())).To(Equal(3))
		g.Expect(vertex.GetLogPathHashes()[0]).To(Equal("bar"))
		g.Expect(vertex.GetLogPathHashes()[1]).To(Equal("baz"))
		g.Expect(vertex.GetLogPathHashes()[2]).To(Equal("foo"))
	})
}
