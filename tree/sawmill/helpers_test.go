package sawmill

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestHelpers(t *testing.T) {
	t.Run("stringListFind: Returns the correct int when expected location is correct", func(t *testing.T) {
		g := NewGomegaWithT(t)

		list := []string{"foo", "bar", "baz"}

		g.Expect(stringListFind(list, "foo", 0)).To(Equal(0))
		g.Expect(stringListFind(list, "bar", 1)).To(Equal(1))
		g.Expect(stringListFind(list, "baz", 2)).To(Equal(2))
	})

	t.Run("stringListFind: Returns the correct int when expected location is wrong", func(t *testing.T) {
		g := NewGomegaWithT(t)

		list := []string{"foo", "bar", "baz"}

		g.Expect(stringListFind(list, "foo", -1)).To(Equal(0))
		g.Expect(stringListFind(list, "foo", 1)).To(Equal(0))
		g.Expect(stringListFind(list, "foo", 2)).To(Equal(0))

		g.Expect(stringListFind(list, "bar", -1)).To(Equal(1))
		g.Expect(stringListFind(list, "bar", 0)).To(Equal(1))
		g.Expect(stringListFind(list, "bar", 2)).To(Equal(1))

		g.Expect(stringListFind(list, "baz", -1)).To(Equal(2))
		g.Expect(stringListFind(list, "baz", 0)).To(Equal(2))
		g.Expect(stringListFind(list, "baz", 1)).To(Equal(2))
	})

	t.Run("stringListFind: Returns -1 when not found", func(t *testing.T) {
		g := NewGomegaWithT(t)

		list := []string{"foo", "bar", "baz"}

		g.Expect(stringListFind(list, "foobarbaz", -1)).To(Equal(-1))
		g.Expect(stringListFind(list, "foobarbaz", 0)).To(Equal(-1))
		g.Expect(stringListFind(list, "foobarbaz", 1)).To(Equal(-1))
		g.Expect(stringListFind(list, "foobarbaz", 2)).To(Equal(-1))
	})
}
