package collections

// Stack is a LIFO data structure backed by a Dequeue.
type Stack[T Comparable] struct {
	dequeue *Dequeue[T]
}

// NewStack creates and returns an empty Stack.
func NewStack[T Comparable]() Stack[T] {
	return Stack[T]{
		dequeue: &Dequeue[T]{},
	}
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(v T) {
	s.dequeue.Push(v)
}

// Pop removes and returns the top element of the stack.
func (s *Stack[T]) Pop() T {
	return s.dequeue.Pop()
}

// Pek returns the top element of the stack without removing it.
func (s *Stack[T]) Pek() T {
	return s.dequeue.Pek()
}

// Queue is a FIFO data structure backed by a Dequeue.
type Queue[T Comparable] struct {
	dequeue *Dequeue[T]
}

// NewQueue creates and returns an empty Queue.
func NewQueue[T Comparable]() Queue[T] {
	return Queue[T]{
		dequeue: &Dequeue[T]{},
	}
}

// Enqueue adds an element to the end of the queue.
func (q *Queue[T]) Enqueue(v T) {
	q.dequeue.Push(v)
}

// Dequeue removes and returns the front element of the queue.
func (q *Queue[T]) Dequeue() T {
	return q.dequeue.PopLeft()
}
