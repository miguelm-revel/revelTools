package syncx

import "sync"

// Barrier is a synchronization primitive that allows multiple goroutines
// to block until all of them have reached a certain point.
//
// Internally, it wraps a syncx.WaitGroup and provides Lock/Unlock semantics
// to increment the number of participants and wait for all of them to finish.
type Barrier struct{ *sync.WaitGroup }

// Lock registers a new participant in the barrier.
//
// Each call to Lock increments the internal counter and indicates that
// the calling goroutine must later call Unlock to signal completion.
func (b *Barrier) Lock() {
	b.Add(1)
}

// Unlock signals that the calling goroutine has reached the barrier
// and waits until all registered participants have also called Unlock.
//
// This method decrements the internal counter and blocks until the
// counter reaches zero.
func (b *Barrier) Unlock() {
	b.Done()
	b.Wait()
}

// NewBarrier creates and returns a new Barrier instance with
// an initialized internal WaitGroup.
func NewBarrier() *Barrier {
	return &Barrier{&sync.WaitGroup{}}
}
