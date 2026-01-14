# reveltools
[![Ask DeepWiki](https://devin.ai/assets/askdeepwiki.png)](https://deepwiki.com/miguelm-revel/revelTools)

This module provides a small, self-contained set of **concurrency primitives** and **generic data structures** implemented in Go. It is designed to be lightweight, idiomatic, and easy to integrate into other projects without external dependencies.

The module is organized into two main packages:
-   **`sync`**: Simple synchronization helpers built on top of Go's standard library primitives.
-   **`collections`**: A suite of generic containers, iterators, and their thread-safe counterparts.

## Installation

Add the module to your project using `go get`:

```bash
go get github.com/miguelm-revel/revelTools
```

## Package `sync`

The `sync` package provides lightweight concurrency utilities.

### Barrier

`Barrier` is a "wait-for-all" synchronization primitive. A group of goroutines can use a `Barrier` to block until all of them have reached a certain point. Each goroutine registers its participation by calling `Lock()` and signals its arrival at the barrier by calling `Unlock()`.

#### API

- `NewBarrier() *Barrier`
- `(*Barrier).Lock()`
- `(*Barrier).Unlock()`

#### Example

```go
import "github.com/miguelm-revel/revelTools/sync"

// Create a new barrier
b := sync.NewBarrier()

for i := 0; i < 3; i++ {
    // Register a new participant
    b.Lock()
    go func() {
        // Signal arrival and wait for all others
        defer b.Unlock() 
        // ... do work ...
    }()
}
// Execution continues here only after all 3 goroutines have called Unlock()
```

### Semaphore

`Semaphore` is a classic counting semaphore implemented using a buffered channel. The capacity of the semaphore defines the maximum number of concurrent resource holders.

#### API

- `NewSemaphore(n int) Semaphore`
- `(Semaphore).Lock()`
- `(Semaphore).Unlock()`

#### Example

```go
import "github.com/miguelm-revel/revelTools/sync"

// Create a semaphore that allows up to 5 concurrent jobs
sem := sync.NewSemaphore(5)

for _, job := range jobs {
    sem.Lock() // Will block if 5 jobs are already running
    go func(j Job) {
        defer sem.Unlock() // Release the semaphore slot
        handle(j)
    }(job)
}
```

## Package `collections`

The `collections` package provides generic data structures, common interfaces, and concurrent-safe wrappers.

### Shared Interfaces

- `Iterable[V]`: Defines methods for creating forward and indexed iterators.
- `Comparable`: Defines a total ordering over values (`<`, `>`, `==`, etc.), required for ordered collections like `Heap`.
- `QueueLike[T]`: Represents a basic FIFO data structure.
- `StackLike[T]`: Represents a basic LIFO data structure.

### Set

`Set[T]` is an unordered collection of unique elements, implemented with a `map` for efficient O(1) operations.

#### Features

- Constant-time add (`Add`), delete (`Del`), and membership checks (`Has`).
- Set-theoretic operations like `Union` and `Intersection`.
- Forward and indexed iteration (`Iter`, `Iter2`).

#### Example

```go
import "github.com/miguelm-revel/revelTools/collections"

s := collections.NewSet([]int{1, 2, 3})
s.Add(4) // s is now {1, 2, 3, 4}

if s.Has(2) {
    s.Del(2) // s is now {1, 3, 4}
}

s2 := collections.NewSet([]int{3, 4, 5})
union := s.Union(s2) // {1, 3, 4, 5}

for v := range s.Iter() {
    fmt.Println(v)
}
```

### Stack

`Stack[T]` is a LIFO (Last-In, First-Out) data structure.

#### API & Example

```go
import "github.com/miguelm-revel/revelTools/collections"

st := collections.NewStack[string]()
st.Push("A")
st.Push("B")

top := st.Peek() // "B"
val := st.Pop()  // "B"
len := st.Len()  // 1
```

### Queue

`Queue[T]` is a FIFO (First-In, First-Out) data structure.

#### API & Example

```go
import "github.com/miguelm-revel/revelTools/collections"

q := collections.NewQueue[string]()
q.Enqueue("A")
q.Enqueue("B")

front := q.Dequeue() // "A"
len := q.Len()       // 1
```

### Heap and PriorityQueue

`PriorityQueue[T]` is a priority queue built on a generic heap implementation. It can be configured as either a `MinHeap` or `MaxHeap`. Elements must implement the `collections.Comparable` interface.

#### Features

- `MinHeap` or `MaxHeap` ordering.
- `Enqueue` and `Dequeue` operations that respect element priority.

#### Example

```go
import "github.com/miguelm-revel/revelTools/collections"

// Assuming MyInt implements collections.Comparable
pq := collections.NewPriorityQueue[MyInt](collections.MinHeap)

pq.Enqueue(MyInt(10))
pq.Enqueue(MyInt(5))
pq.Enqueue(MyInt(15))

next := pq.Dequeue() // Dequeues `MyInt(5)`
```

### Concurrent Collections

The package provides thread-safe, blocking wrappers for queues and stacks, ideal for producer-consumer patterns.

#### GoQueue (Thread-Safe Queue)

`GoQueue[T]` is a concurrent-safe, blocking, and optionally bounded queue. It wraps any type that implements `QueueLike`.

- **Features**:
    - `Enqueue`: Adds an item, blocking if the queue is full (if bounded).
    - `Dequeue`: Removes an item, blocking if the queue is empty.
    - `TryDequeue`: Removes an item without blocking.
    - `Close`: Shuts down the queue, unblocking all waiting goroutines.
    - `NewGoQueue(queue Like[T], buffer int)`: `buffer = 0` creates an unbounded queue.

##### Example

```go
// Create a bounded queue with a capacity of 10
q := collections.NewGoQueue(collections.NewQueue[string](), 10)

// Producer goroutine
go func() {
    q.Enqueue("some data")
}()

// Consumer goroutine
data, ok := q.Dequeue()
if ok {
    // process data
}
```

#### GoStack (Thread-Safe Stack)

`GoStack[T]` is a concurrent-safe, blocking, and optionally bounded stack. It wraps any type that implements `StackLike`.

- **Features**:
    - `Push`: Adds an item, blocking if the stack is full (if bounded).
    - `Pop`: Removes an item, blocking if the stack is empty.
    - `TryPop`: Removes an item without blocking.
    - `Close`: Shuts down the stack, unblocking all waiting goroutines.
    - `NewGoStack(stack StackLike[T], buffer int)`: `buffer = 0` creates an unbounded stack.

##### Example

```go
// Create an unbounded thread-safe stack
st := collections.NewGoStack(collections.NewStack[int](), 0)

// Producer goroutine
go func() {
    st.Push(100)
}()

// Consumer goroutine
val, ok := st.Pop()
if ok {
    // process val
}