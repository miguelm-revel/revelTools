package collections

import (
	"container/list"
	"iter"
)

// Stack is a LIFO data structure backed by a Dequeue.
type Stack[T Comparable] struct {
	dequeue *list.List
}

func (s *Stack[V]) Iter() iter.Seq[V] {
	return func(yield func(V) bool) {
		for s.dequeue.Len() > 0 {
			el := s.Pop()
			if !yield(el) {
				return
			}
		}
	}
}

func (s *Stack[V]) Iter2() iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		idx := 0
		for s.dequeue.Len() > 0 {
			el := s.Pop()
			if !yield(idx, el) {
				return
			}
			idx++
		}
	}
}

// NewStack creates and returns an empty Stack.
func NewStack[T Comparable]() *Stack[T] {
	dequeue := list.New()
	return &Stack[T]{
		dequeue: dequeue,
	}
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(v T) {
	s.dequeue.PushBack(v)
}

// Pop removes and returns the top element of the stack.
func (s *Stack[T]) Pop() T {
	el := s.dequeue.Back()
	s.dequeue.Remove(el)
	return el.Value.(T)
}

// Pek returns the top element of the stack without removing it.
func (s *Stack[T]) Pek() T {
	return s.dequeue.Back().Value.(T)
}

func (s *Stack[T]) Len() int {
	return s.dequeue.Len()
}

// Queue is a FIFO data structure backed by a Dequeue.
type Queue[T Comparable] struct {
	dequeue *list.List
}

// NewQueue creates and returns an empty Queue.
func NewQueue[T Comparable]() *Queue[T] {
	dequeue := list.New()
	return &Queue[T]{
		dequeue: dequeue,
	}
}

// Enqueue adds an element to the end of the queue.
func (c *Queue[T]) Enqueue(v T) {
	c.dequeue.PushBack(v)
}

// Dequeue removes and returns the front element of the queue.
func (c *Queue[T]) Dequeue() T {
	el := c.dequeue.Front()
	c.dequeue.Remove(el)
	return el.Value.(T)
}

func (c *Queue[T]) Len() int {
	return c.dequeue.Len()
}

func (c *Queue[V]) Iter() iter.Seq[V] {
	return func(yield func(V) bool) {
		for c.dequeue.Len() > 0 {
			el := c.Dequeue()
			if !yield(el) {
				return
			}
		}
	}
}

func (c *Queue[V]) Iter2() iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		idx := 0
		for c.dequeue.Len() > 0 {
			el := c.Dequeue()
			if !yield(idx, el) {
				return
			}
			idx++
		}
	}
}
