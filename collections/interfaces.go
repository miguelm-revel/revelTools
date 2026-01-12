package collections

import "iter"

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

type Lengthy interface {
	Len() int
}
