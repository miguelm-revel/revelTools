package collections

import "sync"

// GoQueue is a concurrent-safe, optionally bounded FIFO queue.
// It wraps any Queuer implementation and synchronizes access with a mutex and condition variables.
// A buffer of 0 creates an unbounded queue.
type GoQueue[T any] struct {
	queue    Queuer[T]
	mutex    *sync.Mutex
	nonEmpty *sync.Cond
	nonFull  *sync.Cond
	buffer   int
	closed   bool
}

// Enqueue adds t to the queue. Blocks if the queue is bounded and full.
// Returns immediately without enqueuing if the queue is closed.
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

// Dequeue removes and returns the front element, blocking until one is available.
// Returns (zero, false) if the queue is closed and empty.
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

// TryDequeue removes and returns the front element without blocking.
// Returns (zero, false) if the queue is empty.
func (a *GoQueue[T]) TryDequeue() (t T, ok bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.queue.Len() != 0 {
		defer a.nonFull.Signal()
		t, ok = a.queue.Dequeue(), true
	}
	return
}

// Len returns the current number of elements in the queue.
func (a *GoQueue[T]) Len() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.queue.Len()
}

// Close shuts down the queue and unblocks all goroutines waiting on Enqueue or Dequeue.
func (a *GoQueue[T]) Close() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.closed = true
	a.nonEmpty.Broadcast()
	a.nonFull.Broadcast()
}

// NewGoQueue wraps queue in a thread-safe GoQueue with the given capacity.
// A buffer of 0 creates an unbounded queue.
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

// GoStack is a concurrent-safe, optionally bounded LIFO stack.
// It wraps any Stacker implementation and synchronizes access with a mutex and condition variables.
// A buffer of 0 creates an unbounded stack.
type GoStack[T any] struct {
	stack    Stacker[T]
	mutex    *sync.RWMutex
	nonEmpty *sync.Cond
	nonFull  *sync.Cond
	buffer   int
	closed   bool
}

// Len returns the current number of elements in the stack.
func (a *GoStack[T]) Len() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.stack.Len()
}

// Push adds t to the top of the stack. Blocks if the stack is bounded and full.
// Returns immediately without pushing if the stack is closed.
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

// Pop removes and returns the top element, blocking until one is available.
// Returns (zero, false) if the stack is closed and empty.
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

// TryPop removes and returns the top element without blocking.
// Returns (zero, false) if the stack is empty.
func (a *GoStack[T]) TryPop() (t T, ok bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if a.stack.Len() != 0 {
		defer a.nonFull.Signal()
		t, ok = a.stack.Pop(), true
	}
	return
}

// Peek returns the top element without removing it.
func (a *GoStack[T]) Peek() T {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.stack.Peek()
}

// Close shuts down the stack and unblocks all goroutines waiting on Push or Pop.
func (a *GoStack[T]) Close() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.closed = true
	a.nonEmpty.Broadcast()
	a.nonFull.Broadcast()
}

// NewGoStack wraps stack in a thread-safe GoStack with the given capacity.
// A buffer of 0 creates an unbounded stack.
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
