// Under the Apache-2.0 License
package concurrent_test

import (
	"sync"
	"testing"

	"github.com/groboclown/qazaar-testing/rule-engine/engine/concurrent"
)

func Test_EarlyExit(t *testing.T) {
	t.Run("queue-read-close", func(t *testing.T) {
		q, r := concurrent.NewEarlyExit[string](0)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			if q.Queue("a") {
				t.Error("Queue reported a stop before a stop happened.")
			}
		}()
		go func() {
			defer wg.Done()
			v := <-r.Reader()
			if v != "a" {
				t.Errorf("expected to read value 'a', but found '%s'", v)
			}
		}()
		wg.Wait()

		if q.Stopped() {
			t.Error("reported stopped before calling stop")
		}
		r.Stop()
		if !q.Stopped() {
			t.Error("reported not-stopped after calling stop")
		}
		// Ensure stop is idempotent
		r.Stop()
		if !q.Stopped() {
			t.Error("reported not-stopped after calling stop twice")
		}
	})

	t.Run("close-then-queue", func(t *testing.T) {
		q, r := concurrent.NewEarlyExit[string](0)
		r.Stop()
		if !q.Stopped() {
			t.Error("reported not-stopped after calling stop")
		}

		if !q.Queue("a") {
			t.Error("Queue reported not-stopped after calling stop.")
		}

		select {
		case _, ok := <-r.Reader():
			if ok {
				t.Error("Reader did not report a closed queue.")
			}
		default:
			t.Error("Reader did not report a closed queue.")
		}
	})

	t.Run("partial-queue", func(t *testing.T) {
		q, r := concurrent.NewEarlyExit[int](0)
		values := []int{1, 2, 3, 4, 5, 6, 7, 8}

		go func() {
			c, ok := q.QueueAll(values...)
			if !ok {
				t.Error("QueueAll reported not-closed scenario")
			}
			if c != 4 {
				t.Errorf("expected 4 writes, QueueAll reported %d writes", c)
			}
		}()

		for v := range r.Reader() {
			if v == 4 {
				r.Stop()
			}
			if v > 4 {
				t.Errorf("Read returned values past 4 (%d)", v)
			}
		}
	})
}
