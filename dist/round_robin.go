package dist

type RoundRobinPool[T any] struct {
	pool    []T
	current int
}

func (r *RoundRobinPool[T]) Next() T {
	item := r.pool[r.current]
	if r.current == len(r.pool)-1 {
		r.current = 0
	} else {
		r.current++
	}
	return item
}

func NewRoundRobinPool[T any](items []T) *RoundRobinPool[T] {
	return &RoundRobinPool[T]{
		pool:    items,
		current: 0,
	}
}
