package collections

import "sync"

type GoQueue[T any] struct {
	queue  QueueLike[T]
	mutex  *sync.Mutex
	condF  *sync.Cond
	confR  *sync.Cond
	buffer int
}

func (a *GoQueue[T]) Enqueue(t T) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for a.buffer != 0 && a.queue.Len() == a.buffer {
		a.confR.Wait()
	}
	a.queue.Enqueue(t)
	a.condF.Broadcast()
}

func (a *GoQueue[T]) Dequeue() T {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	defer a.confR.Broadcast()
	for a.queue.Len() == 0 {
		a.condF.Wait()
	}
	return a.queue.Dequeue()
}

func (a *GoQueue[T]) Len() int {
	return a.queue.Len()
}

func NewGoQueue[T Comparable](queue QueueLike[T], buffer int) *GoQueue[T] {
	mutex := &sync.Mutex{}
	goQueue := &GoQueue[T]{
		queue:  queue,
		mutex:  mutex,
		condF:  sync.NewCond(mutex),
		buffer: buffer,
	}
	if buffer != 0 {
		goQueue.confR = sync.NewCond(mutex)
	}
	return goQueue
}

type GoStack[T any] struct {
	stack  StackLike[T]
	mutex  *sync.RWMutex
	condF  *sync.Cond
	confR  *sync.Cond
	buffer int
}

func (a *GoStack[T]) Len() int {
	return a.stack.Len()
}

func (a *GoStack[T]) Push(t T) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for a.buffer != 0 && a.stack.Len() == a.buffer {
		a.confR.Wait()
	}
	a.stack.Push(t)
	a.condF.Broadcast()
}

func (a *GoStack[T]) Pop() T {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	defer a.confR.Broadcast()
	for a.stack.Len() == 0 {
		a.condF.Wait()
	}
	return a.stack.Pop()
}

func (a *GoStack[T]) Pek() T {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.stack.Pek()
}

func NewGoStack[T any](stack StackLike[T], buffer int) *GoStack[T] {
	mutex := &sync.RWMutex{}
	goStack := &GoStack[T]{
		stack:  stack,
		mutex:  mutex,
		condF:  sync.NewCond(mutex),
		buffer: buffer,
	}
	if buffer != 0 {
		goStack.confR = sync.NewCond(mutex)
	}
	return goStack
}
