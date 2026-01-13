package collections

import "sort"

type HeapType uint

const (
	MinHeap HeapType = iota
	MaxHeap
)

type Heap[T Comparable] struct {
	heap     []T
	heapType HeapType
}

func (h *Heap[T]) Len() int {
	return len(h.heap)
}

func (h *Heap[T]) Less(i, j int) bool {
	if h.heapType == MaxHeap {
		return h.heap[i].Lt(h.heap[j])
	} else {
		return h.heap[i].Gt(h.heap[j])
	}
}

func (h *Heap[T]) Swap(i, j int) {
	h.heap[i], h.heap[j] = h.heap[j], h.heap[i]
}

func (h *Heap[T]) Push(t T) {
	h.heap = append(h.heap, t)
	sort.Sort(h)
}

func (h *Heap[T]) Pop() T {
	old := h.heap
	n := len(old)
	x := old[n-1]
	h.heap = old[0 : n-1]
	return x
}

func NewHeap[T Comparable](heapType HeapType) *Heap[T] {
	return &Heap[T]{
		heap:     make([]T, 0),
		heapType: heapType,
	}
}

type PriorityQueue[T Comparable] struct {
	heap *Heap[T]
}

func (p *PriorityQueue[T]) Enqueue(t T) {
	p.heap.Push(t)
}

func (p *PriorityQueue[T]) Dequeue() T {
	return p.heap.Pop()
}

func NewPriorityQueue[T Comparable](heapType HeapType) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		heap: NewHeap[T](heapType),
	}
}
