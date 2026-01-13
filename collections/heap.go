package collections

import "sort"

// HeapType defines the ordering strategy of a Heap.
type HeapType uint

const (
	// MinHeap orders elements so that the smallest element has the highest priority.
	MinHeap HeapType = iota

	// MaxHeap orders elements so that the largest element has the highest priority.
	MaxHeap
)

// Heap is a generic heap data structure parameterized by a Comparable type.
//
// The ordering of elements is determined by the HeapType (MinHeap or MaxHeap).
type Heap[T Comparable] struct {
	heap     []T
	heapType HeapType
}

// Len returns the number of elements in the heap.
func (h *Heap[T]) Len() int {
	return len(h.heap)
}

// Less reports whether the element with index i should sort before
// the element with index j, according to the heap type.
func (h *Heap[T]) Less(i, j int) bool {
	if h.heapType == MaxHeap {
		return h.heap[i].Lt(h.heap[j])
	}
	return h.heap[i].Gt(h.heap[j])
}

// Swap swaps the elements with indexes i and j.
func (h *Heap[T]) Swap(i, j int) {
	h.heap[i], h.heap[j] = h.heap[j], h.heap[i]
}

// Push inserts a new element into the heap.
func (h *Heap[T]) Push(t T) {
	h.heap = append(h.heap, t)
	sort.Sort(h)
}

// Pop removes and returns the top-priority element from the heap.
func (h *Heap[T]) Pop() T {
	old := h.heap
	n := len(old)
	x := old[n-1]
	h.heap = old[0 : n-1]
	return x
}

// NewHeap creates and returns a new Heap with the given HeapType.
func NewHeap[T Comparable](heapType HeapType) *Heap[T] {
	return &Heap[T]{
		heap:     make([]T, 0),
		heapType: heapType,
	}
}

// PriorityQueue is a queue-like abstraction backed by a Heap.
type PriorityQueue[T Comparable] struct {
	heap *Heap[T]
}

// Enqueue inserts an element into the priority queue.
func (p *PriorityQueue[T]) Enqueue(t T) {
	p.heap.Push(t)
}

// Dequeue removes and returns the highest-priority element.
func (p *PriorityQueue[T]) Dequeue() T {
	return p.heap.Pop()
}

// NewPriorityQueue creates a new PriorityQueue using the given HeapType.
func NewPriorityQueue[T Comparable](heapType HeapType) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		heap: NewHeap[T](heapType),
	}
}
