package collections

import "iter"

type Set[T comparable] map[T]struct{}

func (s *Set[T]) Add(v T) {
	(*s)[v] = struct{}{}
}

func (s *Set[T]) Has(v T) bool {
	_, ok := (*s)[v]
	return ok
}

func (s *Set[T]) Del(v T) {
	delete(*s, v)
}

func (s *Set[T]) Union(set Set[T]) Set[T] {
	result := make(Set[T])
	for it := range s.Iter() {
		result.Add(it)
	}
	for it := range set.Iter() {
		result.Add(it)
	}
	return result
}

func (s *Set[T]) Intersection(set Set[T]) Set[T] {
	result := make(Set[T])
	for it := range s.Iter() {
		if set.Has(it) {
			result.Add(it)
		}
	}
	return result
}

func (s *Set[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v, _ := range *s {
			if !yield(v) {
				return
			}
		}
	}
}

func (s *Set[T]) Iter2() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		idx := 0
		for v, _ := range *s {
			if !yield(idx, v) {
				return
			}
			idx++
		}
	}
}

func (s *Set[T]) Len() int {
	return len(*s)
}
