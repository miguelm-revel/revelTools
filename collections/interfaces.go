package collections

import "iter"

type IterableSimple[V any] interface {
	Iter() iter.Seq[V]
	IterRev() iter.Seq[V]
}

type Iterable[V any] interface {
	Iter() iter.Seq[V]
	Iter2() iter.Seq2[int, V]
}

type RevIterable[V any] interface {
	IterRev() iter.Seq[V]
	Iter2Rev() iter.Seq2[int, V]
}

type Indexable[T any] interface {
	At(int) (T, bool)
}

type Comparable interface {
	Eq(Comparable) bool
	Neq(Comparable) bool
	Gt(Comparable) bool
	Gte(Comparable) bool
	Lt(Comparable) bool
	Lte(Comparable) bool
}

type QueueLike[T any] interface {
	Enqueue(T)
	Dequeue() T
}

type StackLike[T any] interface {
	Push(T)
	Pop() T
}
