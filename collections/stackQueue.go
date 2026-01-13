package collections

type Stack[T Comparable] struct {
	dequeue *Dequeue[T]
}

func NewStack[T Comparable]() Stack[T] {
	return Stack[T]{
		dequeue: &Dequeue[T]{},
	}
}

func (s *Stack[T]) Push(v T) {
	s.dequeue.Push(v)
}

func (s *Stack[T]) Pop() T {
	return s.dequeue.Pop()
}

func (s *Stack[T]) Pek() T {
	return s.dequeue.Pek()
}

type Queue[T Comparable] struct {
	dequeue *Dequeue[T]
}

func NewQueue[T Comparable]() Queue[T] {
	return Queue[T]{
		dequeue: &Dequeue[T]{},
	}
}

func (q *Queue[T]) Enqueue(v T) {
	q.dequeue.Push(v)
}

func (q *Queue[T]) Dequeue() T {
	return q.dequeue.PopLeft()
}
