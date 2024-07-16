// Under the Apache-2.0 License
package concurrent

import "sync"

type Closer func()

// Counter acts a bit like a wait group, but with some channel capabilities.
//
// This tracks the count until it reaches zero.  It also allows for early end.
// When the counter reaches zero (not early end), it runs all the closers.
type Counter struct {
	mux         sync.RWMutex
	ready       bool     // not reached zero and not early exit
	reachedZero bool     // reached zero; does not care about early exit
	count       int      // counter
	total       int      // total incr calls
	zero        chan any // on-ready channel
	closers     []Closer
}

func NewCounter() *Counter {
	return &Counter{
		ready:       true,
		reachedZero: false,
		count:       0,
		total:       0,
		zero:        make(chan any, 1),
		closers:     make([]Closer, 0),
	}
}

// AddCloser adds a closer (or more) to the list of closers for this counter.
//
// The closers run when the count reaches zero.  If the counter has already reached
// zero, then the closers are called immediately.
func (c *Counter) AddCloser(cl ...Closer) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.reachedZero {
		for _, v := range cl {
			v()
		}
		return
	}

	c.closers = append(c.closers, cl...)
}

// IsFinished returns true when the channel has reached zero (after serving at least one bit of work).
func (c *Counter) IsFinished() bool {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return !c.ready
}

// Wait returns a channel that reads a boolean value when the counter reaches zero or if it ended early.
//
// The channel will only return a value once, then will close itself.
func (c *Counter) Wait() <-chan any {
	return c.zero
}

// IfFinished runs the function if the counter has finished.
func (c *Counter) IfFinished(f Closer) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if !c.ready {
		f()
	}
}

// Total returns the total number of times the incr was called.
func (c *Counter) Total() int {
	c.mux.RLock()
	defer c.mux.RUnlock()

	return c.total
}

// Incr increments the ready counter by one.
func (c *Counter) Incr() {
	c.Add(1)
}

func (c *Counter) Add(by int) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.total += by
	c.count += by
}

// Decr decrements the ready counter by one.
func (c *Counter) Decr() {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.count--
	if !c.reachedZero && c.count <= 0 {
		c.reachedZero = true
		c.onEnd(nil)
		for _, v := range c.closers {
			v()
		}
	}
}

func (c *Counter) onEnd(val any) {
	// Must be called from within a write lock.
	if c.ready {
		c.ready = false
		c.zero <- val
		close(c.zero)
	}
}

// End immediately terminates the counter.
func (c *Counter) End(val any) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.onEnd(val)
}
