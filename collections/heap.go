package collections

import "container/heap"

// HeapType defines the ordering strategy of a rHeap.
type HeapType uint

const (
	// MinHeap orders elements so that the smallest element has the highest priority.
	MinHeap HeapType = iota

	// MaxHeap orders elements so that the largest element has the highest priority.
	MaxHeap
)

// rHeap is a generic heap data structure parameterized by a Comparable type.
//
// The ordering of elements is determined by the HeapType (MinHeap or MaxHeap).
type rHeap[T Comparable] struct {
	heap     []T
	heapType HeapType
}

// Len returns the number of elements in the heap.
func (h *rHeap[T]) Len() int {
	return len(h.heap)
}

// Less reports whether the element with index i should sort before
// the element with index j, according to the heap type.
func (h *rHeap[T]) Less(i, j int) bool {
	if h.heapType == MaxHeap {
		return h.heap[i].Lt(h.heap[j])
	}
	return h.heap[i].Gt(h.heap[j])
}

// Swap swaps the elements with indexes i and j.
func (h *rHeap[T]) Swap(i, j int) {
	h.heap[i], h.heap[j] = h.heap[j], h.heap[i]
}

// Push inserts a new element into the heap.
func (h *rHeap[T]) Push(t any) {
	typedT := t.(T)
	h.heap = append(h.heap, typedT)
}

// Pop removes and returns the top-priority element from the heap.
func (h *rHeap[T]) Pop() any {
	old := h.heap
	n := len(old)
	x := old[n-1]
	h.heap = old[0 : n-1]
	return x
}

// newHeap creates and returns a new rHeap with the given HeapType.
func newHeap[T Comparable](heapType HeapType) *rHeap[T] {
	return &rHeap[T]{
		heap:     make([]T, 0),
		heapType: heapType,
	}
}

// PriorityQueue is a queue-like abstraction backed by a rHeap.
type PriorityQueue[T Comparable] struct {
	heap *rHeap[T]
}

// Enqueue inserts an element into the priority queue.
func (p *PriorityQueue[T]) Enqueue(t T) {
	heap.Push(p.heap, t)
}

// Dequeue removes and returns the highest-priority element.
func (p *PriorityQueue[T]) Dequeue() T {
	t := heap.Pop(p.heap).(T)
	return t
}

// NewPriorityQueue creates a new PriorityQueue using the given HeapType.
func NewPriorityQueue[T Comparable](heapType HeapType) *PriorityQueue[T] {
	h := newHeap[T](heapType)
	heap.Init(h)
	return &PriorityQueue[T]{
		heap: h,
	}
}
