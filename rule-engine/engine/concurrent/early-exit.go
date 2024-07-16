// Under the Apache-2.0 License
package concurrent

import (
	"context"
	"fmt"
)

type EarlyExitWorker[T any, R any] interface {
	// Perform runs the execution work.
	//
	// Arguments:
	//   v: value to process
	//   otherWorkers: when new work is found, send it here for handling in other goroutines.
	// Returns:
	//   R - early-exit result, to use if
	//   bool - is true.
	Perform(v T, otherWorkers chan<- T) (*R, bool)
}

func RunEarlyExit[T any, R any](
	operation EarlyExitWorker[T, R],
	ctx context.Context,
	initial []T,
) <-chan *R {
	inp := make(chan T)

	// Add some buffer space on the 'done' channel to keep from blocking.
	count := NewCounter()
	count.AddCloser(func() {
		close(inp)
	})

	ctx, cancel := context.WithCancelCause(ctx)

	handler := func(val T) {
		defer func() {
			count.Decr()

			// Note: errors due to pushing into closed channel "inp" should
			// be ignored.
			if rec := recover(); rec != nil {
				if err, ok := rec.(error); ok {
					cancel(err)
				} else {
					cancel(fmt.Errorf("unrecoverable error: %v", rec))
				}
				count.End(nil)
			}
		}()
		if count.IsFinished() {
			// early exit, before calling perform.
			return
		}
		r, d := operation.Perform(val, inp)
		if d {
			count.End(r)
		}
	}

	done := make(chan *R, 1)
	go func() {
		defer count.End(nil)

		for {
			select {
			case r := <-inp:
				// inp: the values ready to work.
				count.Incr()
				go handler(r)

			case r := <-count.Wait():
				// count.Zero: the work is all done with no exit-early.
				if r == nil {
					done <- nil
				} else {
					done <- r.(*R)
				}
				return

			case <-ctx.Done():
				// ctx.Done: something has cancelled the goroutine stuff.
				done <- nil
				count.End(nil)
				return
			}
		}
	}()

	count.Add(len(initial))
	for _, v := range initial {
		go handler(v)
	}

	return done
}
