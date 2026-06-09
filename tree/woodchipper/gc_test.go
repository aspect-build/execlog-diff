package woodchipper

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aspect-build/execlog-diff/context"
	"github.com/aspect-build/execlog-diff/diagnostic"
	. "github.com/onsi/gomega"
)

func TestGarbageCollection(t *testing.T) {
	t.Run("garbage collector stops when told to by the context", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)
		gcCounter := 0
		garbageCollector := &gc{
			ctx:     ctx,
			spawner: make(chan struct{}, 5),
			wait:    &sync.WaitGroup{},
			runGarbageCollector: func() {
				time.Sleep(10 * time.Millisecond)
				gcCounter++
			},
		}

		garbageCollector.EnableCollection()

		garbageCollector.Collect()

		time.Sleep(100 * time.Millisecond)

		garbageCollector.Collect()

		time.Sleep(100 * time.Millisecond)

		ctx.Error(&diagnostic.Diagnostic{
			Message: "Stop!",
		})

		time.Sleep(100 * time.Millisecond)

		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()

		time.Sleep(100 * time.Millisecond)

		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()

		time.Sleep(100 * time.Millisecond)

		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()

		garbageCollector.Wait()

		g.Expect(gcCounter).To(Equal(2))
	})

	t.Run("only 1 garbage collection request gets cached", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)
		gcCounter := 0
		garbageCollector := &gc{
			ctx:     ctx,
			spawner: make(chan struct{}, 5),
			wait:    &sync.WaitGroup{},
			runGarbageCollector: func() {
				time.Sleep(100 * time.Millisecond)
				gcCounter++
			},
		}

		garbageCollector.EnableCollection()

		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()

		time.Sleep(10 * time.Millisecond)

		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()

		time.Sleep(10 * time.Millisecond)

		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()

		time.Sleep(10 * time.Millisecond)

		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()
		garbageCollector.Collect()

		garbageCollector.Wait()

		g.Expect(gcCounter).To(Equal(2))
	})

	t.Run("garbage collection gets started when not currently running", func(t *testing.T) {
		g := NewGomegaWithT(t)

		var stderr strings.Builder
		ctx := context.New(&stderr, context.Diff)
		gcCounter := 0
		garbageCollector := &gc{
			ctx:     ctx,
			spawner: make(chan struct{}, 5),
			wait:    &sync.WaitGroup{},
			runGarbageCollector: func() {
				time.Sleep(10 * time.Millisecond)
				gcCounter++
			},
		}

		garbageCollector.EnableCollection()

		garbageCollector.Collect()

		time.Sleep(100 * time.Millisecond)

		garbageCollector.Collect()

		time.Sleep(100 * time.Millisecond)

		garbageCollector.Collect()

		time.Sleep(100 * time.Millisecond)

		garbageCollector.Collect()

		garbageCollector.Wait()

		g.Expect(gcCounter).To(Equal(4))
	})
}
