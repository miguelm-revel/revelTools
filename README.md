# Module Overview

This module provides a small, self-contained set of **concurrency primitives** and **generic data structures** implemented in Go. It is designed to be lightweight, idiomatic, and easy to integrate into other projects without external dependencies.

The module is organized into two main packages:

- **`sync`**: simple synchronization helpers built on top of Go primitives.
- **`collections`**: generic containers, iterators, and supporting interfaces.

---

## Installation

Add the module to your project:

```bash
go get github.com/miguelm-revel/revelTools
```

## Package sync

The ``sync`` package contains small concurrency utilities built on top of Go’s standard synchronization primitives.

### Barrier

Barrier is a “wait-for-all” synchronization primitive.
Each goroutine registers itself with Lock() and signals completion with Unlock().
The final Unlock() blocks until all registered goroutines have reached the barrier.

#### API

- NewBarrier() *Barrier
- (*Barrier).Lock()
- (*Barrier).Unlock()

````go
b := syncx.NewBarrier()

for i := 0; i < 3; i++ {
    b.Lock()
    go func() {
        defer b.Unlock()
        // do work
    }()
}
````

### Semaphore

Semaphore is a counting ``semaphore`` implemented using a buffered channel.
The channel capacity defines the maximum number of concurrent holders.

#### API

- NewSemaphore(n int) Semaphore
- (Semaphore).Lock()
- (Semaphore).Unlock()

````go
sem := syncx.NewSemaphore(5)

for _, job := range jobs {
    sem.Lock()
    go func(j Job) {
        defer sem.Unlock()
        handle(j)
    }(job)
}
````

## Package collections
The collections package provides generic containers and iterator-friendly abstractions.
### Shared Interfaces

Several small interfaces define common behaviors:
- Iteration
 - IterableSimple
 - Iterable
 - RevIterable
- Access
 - Indexable
- Ordering
 - Comparable
- Abstractions
 - QueueLike
 - StackLike
These interfaces allow containers to interoperate and be consumed generically.

### Set
``Set[T]`` is an unordered collection of unique elements, implemented as a map.
#### Features
- Constant-time add, delete, and membership checks
- Union and intersection operations
- Forward and indexed iteration

````go
s := collections.Set[int]{}
s.Add(1)
s.Add(2)

if s.Has(2) {
    s.Del(2)
}

for v := range s.Iter() {
    _ = v
}
````

### Dequeue (Double-Ended Queue)
``Dequeue[T]`` is a doubly-linked deque supporting insertion and removal at both ends.
#### Features
- Push and pop from both left and right
- Peek operations
- Safe indexed access
- Forward and reverse iteration (with optional indices)

````go
var d collections.Dequeue[MyComparable]

d.Push(a)
d.PushLeft(b)

right := d.Pop()
left := d.PopLeft()

for v := range d.IterRev() {
    _ = v
}
````
### Stack

``Stack[T]`` is a LIFO data structure implemented as a thin wrapper around ``Dequeue[T]``.

````go
st := collections.NewStack[MyComparable]()
st.Push(x)

top := st.Pek()
_ = st.Pop()
_ = top
````

### Queue

``Queue[T]`` is a FIFO data structure implemented as a thin wrapper around ``Dequeue[T]``.

````go
q := collections.NewQueue[MyComparable]()
q.Enqueue(x)

front := q.Dequeue()
_ = front
````

### Heap and PriorityQueue

``Heap[T]`` is a generic heap structure configurable as either a min-heap or max-heap.
``PriorityQueue[T]`` is a queue abstraction built on top of ``Heap[T]``.

#### Features
- Min-heap or max-heap ordering
- Push and pop operations
- Priority-based dequeueing

````go
pq := collections.NewPriorityQueue[MyComparable](collections.MinHeap)

pq.Enqueue(a)
pq.Enqueue(b)

next := pq.Dequeue()
_ = next
````

### Cycle

``Cycle[T]`` is a circular doubly-linked structure.

#### Features

- Circular push and pop
- Forward and reverse iteration
- Infinite iterators by design

````go
var c collections.Cycle[MyComparable]
c.Push(x)
c.Push(y)

v := c.Pop()
_ = v

// Iterators are infinite — consumer must stop iteration explicitly.
````