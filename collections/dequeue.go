package collections

import (
	"iter"
)

type dequeueItem[T Comparable] struct {
	value T
	next  *dequeueItem[T]
	prev  *dequeueItem[T]
}

// Dequeue is a double-ended queue supporting insertion and
// removal from both ends.
type Dequeue[T Comparable] struct {
	first *dequeueItem[T]
	last  *dequeueItem[T]
	len   int
}

// Len returns the number of elements in the dequeue.
func (d *Dequeue[T]) Len() int {
	return d.len
}

func (d *Dequeue[T]) Less(i, j int) bool {
	left := d.at(i)
	right := d.at(j)
	return left.value.Lt(right.value)
}

func (d *Dequeue[T]) Swap(i, j int) {
	left := d.at(i)
	right := d.at(j)
	ll := left.prev
	rr := right.next
	left.next = rr
	right.prev = ll
	left.prev = right
	right.next = left
}

// Push appends an element to the right end.
func (d *Dequeue[T]) Push(item T) {
	if d.last == nil {
		dit := &dequeueItem[T]{
			value: item,
		}
		d.first = dit
		d.last = dit
	} else {
		dit := dequeueItem[T]{
			value: item,
			prev:  d.last,
		}
		d.last.next = &dit
		d.last = &dit
	}
	d.len++
}

// PushLeft prepends an element to the left end.
func (d *Dequeue[T]) PushLeft(item T) {
	if d.first == nil {
		dit := &dequeueItem[T]{
			value: item,
		}
		d.first = dit
		d.last = dit
	} else {
		dit := dequeueItem[T]{
			value: item,
			next:  d.first,
		}
		d.first.prev = &dit
		d.first = &dit
	}
	d.len++
}

// Pop removes and returns the element from the right end.
func (d *Dequeue[T]) Pop() T {
	last := d.last
	if last != nil {
		d.last = last.prev
		if d.last != nil {
			d.last.next = nil
		} else {
			d.first = nil
		}
		d.len--
		return last.value
	} else {
		var val T
		return val
	}
}

// PopLeft removes and returns the element from the left end.
func (d *Dequeue[T]) PopLeft() T {
	first := d.first
	if first != nil {
		d.first = first.next
		if d.first != nil {
			d.first.prev = nil
		} else {
			d.last = nil
		}
		d.len--
		return first.value
	} else {
		var val T
		return val
	}
}

// Pek returns the rightmost element without removing it.
func (d *Dequeue[T]) Pek() T {
	return d.last.value
}

// PekLeft returns the leftmost element without removing it.
func (d *Dequeue[T]) PekLeft() T {
	return d.first.value
}

// Iter returns a forward iterator over the dequeue.
func (d *Dequeue[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		current := d.first
		for current != nil {
			if !yield(current.value) {
				break
			}
			current = current.next
		}
	}
}

// IterRev returns a reverse iterator over the dequeue.
func (d *Dequeue[T]) IterRev() iter.Seq[T] {
	return func(yield func(T) bool) {
		current := d.last
		for current != nil {
			if !yield(current.value) {
				break
			}
			current = current.prev
		}
	}
}

// Iter2 returns an indexed forward iterator.
func (d *Dequeue[T]) Iter2() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		current := d.first
		idx := 0
		for current != nil {
			if !yield(idx, current.value) {
				break
			}
			current = current.next
			idx++
		}
	}
}

// Iter2Rev returns an indexed reverse iterator.
func (d *Dequeue[T]) Iter2Rev() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		current := d.last
		idx := d.len - 1
		for current != nil {
			if !yield(idx, current.value) {
				break
			}
			current = current.prev
			idx--
		}
	}
}

// At returns the element at index i and whether it exists.
func (d *Dequeue[T]) At(i int) (T, bool) {
	idx := 0
	current := d.first
	for idx < i && current != nil {
		current = current.next
	}
	if current == nil {
		var zero T
		return zero, false
	}
	return current.value, true
}

func (d *Dequeue[T]) at(i int) *dequeueItem[T] {
	idx := 0
	current := d.first
	for idx < i && current != nil {
		current = current.next
	}
	return current
}

// Cycle represents a circular doubly-linked list.
type Cycle[T Comparable] struct {
	first *dequeueItem[T]
	len   int
}

// Iter returns a forward infinite iterator over the cycle.
func (c *Cycle[V]) Iter() iter.Seq[V] {
	return func(yield func(V) bool) {
		current := c.first
		if current == nil {
			return
		}
		for {
			if !yield(current.value) {
				return
			}
			current = current.next
		}
	}
}

// IterRev returns a reverse infinite iterator over the cycle.
func (c *Cycle[V]) IterRev() iter.Seq[V] {
	return func(yield func(V) bool) {
		current := c.first
		if current == nil {
			return
		}
		for {
			if !yield(current.value) {
				return
			}
			current = current.prev
		}
	}
}

// Push inserts an element into the cycle.
func (c *Cycle[T]) Push(v T) {
	it := dequeueItem[T]{
		value: v,
	}
	if c.first == nil {
		it.next = &it
		it.prev = &it
		c.first = &it
	} else {
		first := c.first
		it.next = first
		it.prev = first.prev
		first.prev = &it
		c.first = &it
	}
	c.len++
}

// Pop removes and returns the current element from the cycle.
func (c *Cycle[T]) Pop() (v T) {
	if c.first == nil {
		return
	}
	it := c.first
	it.prev.next = it.next
	it.next.prev = it.prev
	c.len--
	return it.value
}
