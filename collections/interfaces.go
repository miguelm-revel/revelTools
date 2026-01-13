package collections

import "iter"

// Iterable represents a collection that supports indexed iteration.
type Iterable[V any] interface {
	Iter() iter.Seq[V]
	Iter2() iter.Seq2[int, V]
}

// Comparable defines a total ordering over values.
type Comparable interface {
	Eq(Comparable) bool
	Neq(Comparable) bool
	Gt(Comparable) bool
	Gte(Comparable) bool
	Lt(Comparable) bool
	Lte(Comparable) bool
}

// QueueLike represents a FIFO data structure.
type QueueLike[T any] interface {
	Enqueue(T)
	Dequeue() T
	Len() int
}

// StackLike represents a LIFO data structure.
type StackLike[T any] interface {
	Push(T)
	Pop() T
	Pek() T
	Len() int
}
