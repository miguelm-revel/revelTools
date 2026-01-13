package collections

import (
	"iter"
)

type dequeueItem[T Comparable] struct {
	value T
	next  *dequeueItem[T]
	prev  *dequeueItem[T]
}

type Dequeue[T Comparable] struct {
	first *dequeueItem[T]
	last  *dequeueItem[T]
	len   int
}

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

func (d *Dequeue[T]) Pek() T {
	return d.last.value
}

func (d *Dequeue[T]) PekLeft() T {
	return d.first.value
}

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

type Cycle[T Comparable] struct {
	first *dequeueItem[T]
	len   int
}

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
