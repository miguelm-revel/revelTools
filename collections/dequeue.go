package collections

import (
	"iter"
)

type dequeueItem[T any] struct {
	value T
	next  *dequeueItem[T]
	prev  *dequeueItem[T]
}

type Dequeue[T any] struct {
	first *dequeueItem[T]
	last  *dequeueItem[T]
	len   int
}

func (d *Dequeue[T]) Len() int {
	return d.len
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
