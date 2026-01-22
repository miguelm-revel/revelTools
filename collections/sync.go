package collections

import "sync"

type GoQueue[T any] struct {
	queue    Queuer[T]
	mutex    *sync.Mutex
	nonEmpty *sync.Cond
	nonFull  *sync.Cond
	buffer   int
	closed   bool
}

func (a *GoQueue[T]) Enqueue(t T) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.closed {
		return
	}
	for a.buffer != 0 && a.queue.Len() == a.buffer && !a.closed {
		a.nonFull.Wait()
	}
	if a.closed {
		return
	}
	a.queue.Enqueue(t)
	a.nonEmpty.Signal()
}

func (a *GoQueue[T]) Dequeue() (t T, ok bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for a.queue.Len() == 0 && !a.closed {
		a.nonEmpty.Wait()
	}
	if a.closed && a.queue.Len() == 0 {
		return
	}
	t = a.queue.Dequeue()
	a.nonFull.Signal()
	return t, true
}

func (a *GoQueue[T]) TryDequeue() (t T, ok bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.queue.Len() != 0 {
		defer a.nonFull.Signal()
		t, ok = a.queue.Dequeue(), true
	}
	return
}

func (a *GoQueue[T]) Len() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.queue.Len()
}

func (a *GoQueue[T]) Close() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.closed = true
	a.nonEmpty.Broadcast()
	a.nonFull.Broadcast()
}

func NewGoQueue[T any](queue Queuer[T], buffer int) *GoQueue[T] {
	mutex := &sync.Mutex{}
	return &GoQueue[T]{
		queue:    queue,
		mutex:    mutex,
		nonEmpty: sync.NewCond(mutex),
		nonFull:  sync.NewCond(mutex),
		buffer:   buffer,
	}
}

type GoStack[T any] struct {
	stack    Stacker[T]
	mutex    *sync.RWMutex
	nonEmpty *sync.Cond
	nonFull  *sync.Cond
	buffer   int
	closed   bool
}

func (a *GoStack[T]) Len() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.stack.Len()
}

func (a *GoStack[T]) Push(t T) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.closed {
		return
	}
	for a.buffer != 0 && a.stack.Len() == a.buffer && !a.closed {
		a.nonFull.Wait()
	}
	if a.closed {
		return
	}
	a.stack.Push(t)
	a.nonEmpty.Signal()
}

func (a *GoStack[T]) Pop() (t T, ok bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for a.stack.Len() == 0 && !a.closed {
		a.nonEmpty.Wait()
	}
	if a.closed && a.stack.Len() == 0 {
		return
	}
	t = a.stack.Pop()
	a.nonFull.Signal()
	return t, true
}

func (a *GoStack[T]) TryPop() (t T, ok bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.stack.Len() != 0 {
		defer a.nonFull.Signal()
		t, ok = a.stack.Pop(), true
	}
	return
}

func (a *GoStack[T]) Peek() T {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.stack.Peek()
}

func (a *GoStack[T]) Close() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.closed = true
	a.nonEmpty.Broadcast()
	a.nonFull.Broadcast()
}

func NewGoStack[T any](stack Stacker[T], buffer int) *GoStack[T] {
	mutex := &sync.RWMutex{}
	return &GoStack[T]{
		stack:    stack,
		mutex:    mutex,
		nonEmpty: sync.NewCond(mutex),
		buffer:   buffer,
		nonFull:  sync.NewCond(mutex),
	}
}
