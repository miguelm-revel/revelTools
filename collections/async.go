package collections

import "sync"

type GoQueue[T any] struct {
	queue QueueLike[T]
	mutex *sync.Mutex
	cond  *sync.Cond
}

func (a *GoQueue[T]) Enqueue(t T) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.queue.Enqueue(t)
	a.cond.Broadcast()
}

func (a *GoQueue[T]) Dequeue() T {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for a.queue.Len() == 0 {
		a.cond.Wait()
	}
	return a.queue.Dequeue()
}

func (a *GoQueue[T]) Len() int {
	return a.queue.Len()
}

func NewGoQueue[T Comparable](queue QueueLike[T]) *GoQueue[T] {
	mutex := &sync.Mutex{}
	return &GoQueue[T]{
		queue: queue,
		mutex: mutex,
		cond:  sync.NewCond(mutex),
	}
}

type GoStack[T any] struct {
	stack StackLike[T]
	mutex *sync.RWMutex
	cond  *sync.Cond
}

func (a *GoStack[T]) Len() int {
	return a.stack.Len()
}

func (a *GoStack[T]) Push(t T) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.stack.Push(t)
	a.cond.Broadcast()
}

func (a *GoStack[T]) Pop() T {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for a.stack.Len() == 0 {
		a.cond.Wait()
	}
	return a.stack.Pop()
}

func (a *GoStack[T]) Pek() T {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.stack.Pek()
}

func NewGoStack[T any](stack StackLike[T]) *GoStack[T] {
	mutex := &sync.RWMutex{}
	return &GoStack[T]{
		stack: stack,
		mutex: mutex,
		cond:  sync.NewCond(mutex),
	}
}
