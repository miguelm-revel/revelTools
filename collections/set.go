package collections

import (
	"bytes"
	"encoding/json"
	"iter"
)

// Set represents an unordered collection of unique elements.
type Set[T comparable] map[T]struct{}

func (s *Set[T]) UnmarshalJSON(bts []byte) error {
	if bytes.Equal(bts, []byte("null")) {
		*s = nil
		return nil
	}
	raw := make([]T, 0)
	if err := json.Unmarshal(bts, &raw); err != nil {
		return err
	}
	*s = make(Set[T])
	for _, it := range raw {
		s.Add(it)
	}
	return nil
}

func (s *Set[T]) MarshalJSON() ([]byte, error) {
	if s != nil {
		sl := make([]T, s.Len())
		for i, element := range s.Iter2() {
			sl[i] = element
		}
		return json.Marshal(sl)
	} else {
		return []byte("null"), nil
	}
}

func NewSet[T comparable](sub []T) Set[T] {
	newSet := make(Set[T])
	for _, v := range sub {
		newSet.Add(v)
	}
	return newSet
}

// Add inserts a value into the set.
func (s *Set[T]) Add(v T) {
	(*s)[v] = struct{}{}
}

// Has reports whether the value exists in the set.
func (s *Set[T]) Has(v T) bool {
	_, ok := (*s)[v]
	return ok
}

// Del removes a value from the set.
func (s *Set[T]) Del(v T) {
	delete(*s, v)
}

// Union returns a new set containing all elements from both sets.
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

// Intersection returns a new set containing only elements
// present in both sets.
func (s *Set[T]) Intersection(set Set[T]) Set[T] {
	result := make(Set[T])
	for it := range s.Iter() {
		if set.Has(it) {
			result.Add(it)
		}
	}
	return result
}

// Iter returns a forward iterator over the set.
func (s *Set[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range *s {
			if !yield(v) {
				return
			}
		}
	}
}

// Iter2 returns an indexed iterator over the set.
func (s *Set[T]) Iter2() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		idx := 0
		for v := range *s {
			if !yield(idx, v) {
				return
			}
			idx++
		}
	}
}

// Len returns the number of elements in the set.
func (s *Set[T]) Len() int {
	return len(*s)
}
