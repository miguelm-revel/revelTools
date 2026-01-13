package collections

import "iter"

// IterableSimple represents a collection that can be iterated
// forwards and backwards.
type IterableSimple[V any] interface {
	Iter() iter.Seq[V]
	IterRev() iter.Seq[V]
}

// Iterable represents a collection that supports indexed iteration.
type Iterable[V any] interface {
	Iter() iter.Seq[V]
	Iter2() iter.Seq2[int, V]
}

// RevIterable represents a collection that supports reverse indexed iteration.
type RevIterable[V any] interface {
	IterRev() iter.Seq[V]
	Iter2Rev() iter.Seq2[int, V]
}

// Indexable represents a collection that allows safe indexed access.
type Indexable[T any] interface {
	At(int) (T, bool)
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
}

// StackLike represents a LIFO data structure.
type StackLike[T any] interface {
	Push(T)
	Pop() T
}
