package index

import "github.com/miguelm-revel/revelTools/collections"

type Boolean[T comparable] struct {
	t collections.Set[T]
	f collections.Set[T]
}

func NewBoolean[T comparable]() *Boolean[T] {
	return &Boolean[T]{
		t: make(collections.Set[T]),
		f: make(collections.Set[T]),
	}
}

func (b *Boolean[T]) Insert(v bool, value T) {
	if v {
		b.t.Add(value)
	} else {
		b.f.Add(value)
	}
}

func (b *Boolean[T]) Search(v bool) collections.Set[T] {
	if v {
		return copySet(b.t)
	} else {
		return copySet(b.f)
	}
}

func copySet[T comparable](src collections.Set[T]) (tar collections.Set[T]) {
	tar = make(collections.Set[T])
	for el := range src.Iter() {
		tar.Add(el)
	}
	return tar
}
