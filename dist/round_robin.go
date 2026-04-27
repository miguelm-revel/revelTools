package dist

import "sync"

// RoundRobinPool distributes items from a fixed slice in round-robin order.
// All operations are safe for concurrent use.
type RoundRobinPool[T any] struct {
	pool    []T
	current int
	lock    sync.Mutex
}

func RoundRobinPoolFromConstructor[T any](constructor func() T, poolSize int) *RoundRobinPool[T] {
	pool := make([]T, poolSize)
	for i := 0; i < poolSize; i++ {
		pool[i] = constructor()
	}
	return &RoundRobinPool[T]{
		pool:    pool,
		current: 0,
		lock:    sync.Mutex{},
	}
}

// Next returns the next item in round-robin order, cycling back to the first item after the last.
func (r *RoundRobinPool[T]) Next() T {
	r.lock.Lock()
	defer r.lock.Unlock()
	item := r.pool[r.current]
	if r.current == len(r.pool)-1 {
		r.current = 0
	} else {
		r.current++
	}
	return item
}

// NewRoundRobinPool creates a RoundRobinPool backed by items, starting from the first element.
func NewRoundRobinPool[T any](items []T) *RoundRobinPool[T] {
	return &RoundRobinPool[T]{
		pool:    items,
		current: 0,
		lock:    sync.Mutex{},
	}
}
