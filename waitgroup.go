package waitgroupext

import (
	"sync"
	"sync/atomic"
)

type WaitGroup struct {
	channel chan struct{}
	mutex   sync.Mutex
	counter int
	waiting int32
}

// Add adds delta, which may be negative, to the WaitGroup counter.
// If the counter becomes zero, all goroutines blocked on Wait are released.
// If the counter goes negative, Add panics.
//
// Note that calls with a positive delta that occur when the counter is zero
// must happen before a Wait. Calls with a negative delta, or calls with a
// positive delta that start when the counter is greater than zero, may happen
// at any time.
// Typically this means the calls to Add should execute before the statement
// creating the goroutine or other event to be waited for.
// If a WaitGroup is reused to wait for several independent sets of events,
// new Add calls must happen after all previous Wait calls have returned.
// See the WaitGroup example.
func (wg *WaitGroup) Add(delta int) {
	if delta == 0 {
		return
	}
	wg.mutex.Lock()
	defer wg.mutex.Unlock()
	if delta < 0 {
		wg.counter += delta
		if wg.counter == 0 {
			close(wg.channel)
		}
		if wg.counter < 0 {
			wg.counter = 0
			panic("waitgroupext: negative WaitGroup counter")
		}
	} else {
		if atomic.LoadInt32(&wg.waiting) != 0 {
			panic("waitgroupext: WaitGroup misuse: Add called concurrently with Wait")
		}
		if wg.counter == 0 {
			wg.channel = make(chan struct{})
		}
		wg.counter += delta
	}
}

// Done decrements the WaitGroup counter.
func (wg *WaitGroup) Done() {
	wg.Add(-1)
}

// Wait blocks until the WaitGroup counter is zero. Note
// that if Add() is performed with a positive parameter after
// WaitChannel() is called, the returned channel will not
// necessarily wait for such nwely added items.
func (wg *WaitGroup) Wait() {
	atomic.AddInt32(&wg.waiting, 1)
	defer atomic.AddInt32(&wg.waiting, -1)
	<-wg.WaitChannel()
}

// WaitChan returns a channel which is only readable when
// the waitgroup has reached zero entries. Note that if
// Add() is performed with a positive parameter after
// WaitChannel() is called, the returned channel will not
// necessarily wait for such nwely added items.
func (wg *WaitGroup) WaitChannel() <-chan struct{} {
	// There is some subtlety here about what happens if
	// Add() is called with a positive or negative
	// parameter after WaitChannel has returned a channel
	// but before the channel is waited on. If Add()
	// is called and the counter reaches zero, the channel
	// will be closed; this will stop any subsequent Wait()
	// which is fine. If Add() is subsequently called with
	// a positive index, the channel returned will not
	// be reopened (a new one will be created), which
	// means a select on the returned channel will not
	// block. Therefore this call is defined such that it
	// will not necessarily wait for Add() performed after
	// WaitChannel() was called. This is a slight relaxation
	// of the conventional WaitGroup semantics
	wg.mutex.Lock()
	defer wg.mutex.Unlock()
	if wg.counter == 0 && wg.channel == nil {
		wg.channel = make(chan struct{})
		close(wg.channel)
	}
	return wg.channel
}
