// Under the Apache-2.0 License
package concurrent

import (
	"sync/atomic"
)

type EarlyExitQueue[T any] interface {
	// Queue adds one item to the queue, and returns `true` if it could not add it due to a closed queue.
	//
	// This will not return until the queue has sufficient space to load the value.
	Queue(v T) bool

	// QueueAll adds all the items to the queue, and returns the count of items added and `true` if the queue closed before it could finish adding all the items.
	//
	// This will not return until the queue has sufficient space to load the value.
	QueueAll(v ...T) (int, bool)

	// Stopped checks if the reader stopped the queue.
	Stopped() bool

	// Finished marks the queue process as complete.
	//
	// Finished may safely be called multiples times, and safely called when the queue is already stopped.
	Finished()
}

type EarlyExitReader[T any] interface {
	// Reader returns the channel that reads from the queue.
	Reader() <-chan T

	// Stop prohibits the queue from writing any more values.
	//
	// Stop may safely be called multiple times.
	Stop()
}

type earlyExitGenerator[T any] struct {
	queue chan T

	// stopped state may only move from false to true.
	stopped atomic.Bool
}

// NewEarlyExit creates the EarlyExitQueue and matching EarlyExitReader.
//
// `size` defines the number of items the queue can store before blocking writes.
func NewEarlyExit[T any](size int) (EarlyExitQueue[T], EarlyExitReader[T]) {
	ret := &earlyExitGenerator[T]{
		queue: make(chan T, size),
	}
	return ret, ret
}

func (e *earlyExitGenerator[T]) Queue(v T) (b bool) {
	b = false
	defer func() {
		if err := recover(); err != nil {
			b = true
		}
	}()
	e.queue <- v
	return
}

func (e *earlyExitGenerator[T]) QueueAll(v ...T) (c int, b bool) {
	b = false
	c = 0
	defer func() {
		if err := recover(); err != nil {
			b = true
		}
	}()
	for _, x := range v {
		e.queue <- x
		// Only count after a successful queue operation.
		c++
	}
	return
}

func (e *earlyExitGenerator[T]) Reader() <-chan T {
	return e.queue
}

func (e *earlyExitGenerator[T]) Stop() {
	if old := e.stopped.Swap(true); !old {
		// Was not previously stopped
		close(e.queue)
	}
}

func (e *earlyExitGenerator[T]) Finished() {
	e.Stop()
}

func (e *earlyExitGenerator[T]) Stopped() bool {
	return e.stopped.Load()
}
