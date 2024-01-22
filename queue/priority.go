package queue

import (
	"container/heap"
)

// Item represents an item in the priority queue.
type Item[T any] struct {
	Value    T   // 实际的值
	Priority int // 优先级
}

// PriorityQueue implements a min heap based on the priority field of Item.
type PriorityQueue[T any] []*Item[T]

// Len 返回堆中的元素数量
func (pq PriorityQueue[T]) Len() int { return len(pq) }

// Less 指定两个元素之间的比较规则，这里是按照优先级升序排列
func (pq PriorityQueue[T]) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}

// Swap 交换两个元素的位置
func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

// Push 将元素推入堆中
func (pq *PriorityQueue[T]) Push(x interface{}) {
	item := x.(*Item[T])
	*pq = append(*pq, item)
}

// Pop 从堆中弹出最小元素
func (pq *PriorityQueue[T]) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// FixedSizeHeap implements a fixed-size heap with a priority queue.
type FixedSizeHeap[T any] struct {
	maxSize int
	heap    PriorityQueue[T]
}

// NewFixedSizeHeap creates a new fixed-size heap with the specified max size.
func NewFixedSizeHeap[T any](maxSize int) *FixedSizeHeap[T] {
	return &FixedSizeHeap[T]{
		maxSize: maxSize,
		heap:    make(PriorityQueue[T], 0),
	}
}

// Push adds an item with priority to the heap.
func (h *FixedSizeHeap[T]) Push(item *Item[T]) {
	if len(h.heap) < h.maxSize {
		heap.Push(&h.heap, item)
	} else {
		if item.Priority > h.heap[0].Priority {
			heap.Pop(&h.heap)
			heap.Push(&h.heap, item)
		}
	}
}

// Pop removes and returns the item with the highest priority from the heap.
func (h *FixedSizeHeap[T]) Pop() *Item[T] {
	if !h.IsEmpty() {
		return heap.Pop(&h.heap).(*Item[T])
	}
	return nil
}

// IsEmpty returns true if the heap is empty.
func (h *FixedSizeHeap[T]) IsEmpty() bool {
	return len(h.heap) == 0
}

// Len returns the number of items in the heap.
func (h *FixedSizeHeap[T]) Len() int {
	return len(h.heap)
}
