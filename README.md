# revelTools
[![Ask DeepWiki](https://devin.ai/assets/askdeepwiki.png)](https://deepwiki.com/miguelm-revel/revelTools)

A lightweight, self-contained Go module providing **concurrency primitives**, **generic data structures**, and **statistical distributions**. No external dependencies beyond the standard library.

## Installation

```bash
go get github.com/miguelm-revel/revelTools
```

## Packages

| Package | Description |
|---|---|
| `syncx` | Synchronization and concurrency utilities |
| `collections` | Generic containers, iterators, and thread-safe wrappers |
| `randx` | Probability distributions (sampling, PDF, CDF) |
| `dist` | Item distribution strategies |

---

## Package `syncx`

### Barrier

`Barrier` blocks until all registered goroutines have reached a common point. Each participant calls `Lock()` to register and `Unlock()` to signal arrival. All goroutines unblock once the last one calls `Unlock()`.

```go
import "github.com/miguelm-revel/revelTools/syncx"

b := syncx.NewBarrier()

for i := 0; i < 3; i++ {
    b.Lock()
    go func() {
        defer b.Unlock()
        // ... do work ...
    }()
}
// All 3 goroutines complete before execution continues here
```

### RateLimiter

`RateLimiter` is a token-bucket limiter. Tokens are pre-filled at creation and replenished at a steady interval by a background goroutine.

```go
import (
    "context"
    "github.com/miguelm-revel/revelTools/syncx"
    "time"
)

rl := syncx.NewRateLimiter(10, time.Second) // 10 requests per second
defer rl.Stop()

// Non-blocking check
if rl.Allow() {
    // handle request
}

// Blocking wait with context
if err := rl.Wait(ctx); err != nil {
    // context cancelled
}
```

### WorkerPool

`WorkerPool` runs a fixed number of worker goroutines that drain a job channel. Errors from jobs are collected and returned by `Wait`.

```go
import "github.com/miguelm-revel/revelTools/syncx"

pool := syncx.NewWorkerPool(4, 100) // 4 workers, queue capacity 100

pool.Submit(func() error {
    // ... do work ...
    return nil
})

pool.Close() // signal no more jobs
errs := pool.Wait() // blocks until all workers finish
```

---

## Package `collections`

### Shared Interfaces

| Interface | Description |
|---|---|
| `Iterable[V]` | Forward (`Iter`) and indexed (`Iter2`) iterators |
| `Comparable` | Total ordering: `Eq`, `Neq`, `Lt`, `Lte`, `Gt`, `Gte` |
| `Setter[T]` | Set membership: `Add`, `Has`, `Del` |
| `Queuer[T]` | FIFO: `Enqueue`, `Dequeue`, `Len` |
| `Stacker[T]` | LIFO: `Push`, `Pop`, `Peek`, `Len` |
| `IntoSlice[T]` | Materialize elements into a `[]T` slice |

### Set

`Set[T]` is an unordered collection of unique elements backed by a `map`. O(1) operations.

```go
import "github.com/miguelm-revel/revelTools/collections"

s := collections.NewSet([]int{1, 2, 3})
s.Add(4)
s.Del(2)
fmt.Println(s.Has(3)) // true

s2 := collections.NewSet([]int{3, 4, 5})
union := s.Union(s2)        // {1, 3, 4, 5}
inter := s.Intersection(s2) // {3, 4}

for v := range s.Iter() {
    fmt.Println(v)
}
```

`Set` implements `json.Marshaler` / `json.Unmarshaler` (serialized as a JSON array).

### Stack

`Stack[T]` is a LIFO structure backed by a doubly-linked list. `T` must implement `Comparable`.

```go
st := collections.NewStack[MyInt]()
st.Push(MyInt(1))
st.Push(MyInt(2))

top := st.Peek() // MyInt(2)
val := st.Pop()  // MyInt(2)
fmt.Println(st.Len()) // 1
```

### Queue

`Queue[T]` is a FIFO structure backed by a doubly-linked list. `T` must implement `Comparable`.

```go
q := collections.NewQueue[MyInt]()
q.Enqueue(MyInt(1))
q.Enqueue(MyInt(2))

front := q.Dequeue() // MyInt(1)
fmt.Println(q.Len()) // 1
```

### Heap and PriorityQueue

`PriorityQueue[T]` wraps a generic `Heap` and can be configured as `MinHeap` or `MaxHeap`. `T` must implement `Comparable`.

```go
pq := collections.NewPriorityQueue[MyInt](collections.MinHeap)
pq.Enqueue(MyInt(10))
pq.Enqueue(MyInt(5))
pq.Enqueue(MyInt(15))

next := pq.Dequeue() // MyInt(5)
```

### GoQueue (Thread-Safe Queue)

`GoQueue[T]` is a concurrent-safe, optionally bounded FIFO wrapper around any `Queuer`. A `buffer` of `0` is unbounded.

```go
q := collections.NewGoQueue(collections.NewQueue[MyInt](), 10)

// Producer
go func() { q.Enqueue(MyInt(42)) }()

// Consumer (blocks until item available)
val, ok := q.Dequeue()

// Non-blocking consumer
val, ok = q.TryDequeue()

q.Close() // unblocks all waiting goroutines
```

### GoStack (Thread-Safe Stack)

`GoStack[T]` is a concurrent-safe, optionally bounded LIFO wrapper around any `Stacker`. A `buffer` of `0` is unbounded.

```go
st := collections.NewGoStack(collections.NewStack[MyInt](), 0)

go func() { st.Push(MyInt(1)) }()

val, ok := st.Pop()     // blocks until item available
val, ok = st.TryPop()   // non-blocking

st.Close()
```

### BKTree

`BKTree` is a Burkhard-Keller tree for approximate string matching using edit distance. The `Fuzziness` field controls the maximum edit distance for `Has`.

```go
tree := collections.NewBKTree([]string{"hello", "world", "golang"})
tree.Fuzziness = 1

fmt.Println(tree.Has("helo"))  // true (1 edit away)
fmt.Println(tree.Has("xyz"))   // false

tree.Add("gopher")
tree.Del("world")
fmt.Println(tree.Len()) // 3
```

`BKTree` implements `json.Marshaler` / `json.Unmarshaler` (serialized as a JSON array). `Fuzziness` is reset to `2` on unmarshal.

### Zip

`Zip` and `ZipSlice` produce paired iterators from two sequences, stopping at the shorter one.

```go
import "github.com/miguelm-revel/revelTools/collections"

keys := []string{"a", "b", "c"}
vals := []int{1, 2, 3}

for k, v := range collections.ZipSlice(keys, vals) {
    fmt.Println(k, v) // "a" 1 / "b" 2 / "c" 3
}
```

---

## Package `randx`

All distributions implement the `Dist` interface:

```go
type Dist interface {
    Rand() float64       // draw a random sample
    PDF(x float64) float64 // probability density/mass at x
    CDF(x float64) float64 // cumulative probability P(X <= x)
}
```

### NormalDist

```go
import "github.com/miguelm-revel/revelTools/randx"

d := randx.NormalDist{Mu: 0, Sigma: 1}
sample := d.Rand()
fmt.Println(d.PDF(0))   // ≈ 0.3989
fmt.Println(d.CDF(1.96)) // ≈ 0.975
```

### ExpDist

```go
d := randx.ExpDist{Lambda: 2.0} // mean = 0.5
sample := d.Rand()
fmt.Println(d.CDF(1.0)) // ≈ 0.8647
```

### BinomDist

```go
d := randx.BinomDist{N: 10, P: 0.5}
sample := d.Rand()
fmt.Println(d.PDF(5))   // ≈ 0.2461
fmt.Println(d.CDF(5))   // ≈ 0.6230
```

### PoissonDist

```go
d := randx.PoissonDist{Lambda: 3.0}
sample := d.Rand()
fmt.Println(d.PDF(3))   // ≈ 0.2240
fmt.Println(d.CDF(3))   // ≈ 0.6472
```

### Chi2Dist

```go
d := randx.Chi2Dist{K: 4}
sample := d.Rand()
fmt.Println(d.PDF(2))   // ≈ 0.1839
fmt.Println(d.CDF(4))   // ≈ 0.5940
```

---

## Package `dist`

### RoundRobinPool

`RoundRobinPool[T]` cycles through a fixed slice of items in round-robin order. Safe for concurrent use.

```go
import "github.com/miguelm-revel/revelTools/dist"

pool := dist.NewRoundRobinPool([]string{"server-1", "server-2", "server-3"})

for i := 0; i < 6; i++ {
    fmt.Println(pool.Next())
    // server-1, server-2, server-3, server-1, server-2, server-3
}
```
