package utils

import "container/heap"

type priorityQueueArr[T any] struct {
	arr        []T
	comparator func(T, T) bool
}

func (pqa priorityQueueArr[T]) Len() int {
	return len(pqa.arr)
}

func (pqa priorityQueueArr[T]) Less(i, j int) bool {
	return pqa.comparator(pqa.arr[i], pqa.arr[i])
}

func (pqa priorityQueueArr[T]) Swap(i, j int) {
	pqa.arr[i], pqa.arr[j] = pqa.arr[j], pqa.arr[i]
}

func (pqa *priorityQueueArr[T]) Push(item any) {
	pqa.arr = append(pqa.arr, item.(T))
}

func (pqa *priorityQueueArr[T]) Pop() any {
	old := pqa.arr
	item := old[len(old)-1]
	pqa.arr = old[0 : len(old)-1]
	return item
}

type PriorityQueue[T any] struct {
	nodes priorityQueueArr[T]
}

func NewPriorityQueue[T any](comparator func(T, T) bool) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		nodes: priorityQueueArr[T]{
			arr:        make([]T, 0),
			comparator: comparator,
		},
	}
}

func (pq *PriorityQueue[T]) Push(item T) {
	heap.Push(&pq.nodes, item)
}

func (pq *PriorityQueue[T]) Pop() T {
	return heap.Pop(&pq.nodes).(T)
}

func (pq *PriorityQueue[T]) Top() T {
	return pq.nodes.arr[0]
}

func (pq *PriorityQueue[T]) Clear() {
	pq.nodes.arr = make([]T, 0)
}
