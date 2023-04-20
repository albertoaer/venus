package utils

import "container/heap"

type PriorityItem interface {
	Priority() int64
}

type priorityQueueArr[T PriorityItem] []T

func (pqa priorityQueueArr[T]) Len() int {
	return len(pqa)
}

func (pqa priorityQueueArr[T]) Less(i, j int) bool {
	return pqa[i].Priority() < pqa[i].Priority()
}

func (pqa priorityQueueArr[T]) Swap(i, j int) {
	pqa[i], pqa[j] = pqa[j], pqa[i]
}

func (pqa *priorityQueueArr[T]) Push(item any) {
	*pqa = append(*pqa, item.(T))
}

func (pqa *priorityQueueArr[T]) Pop() any {
	old := *pqa
	item := old[len(old)-1]
	*pqa = old[0 : len(old)-1]
	return item
}

type PriorityQueue[T PriorityItem] struct {
	nodes priorityQueueArr[T]
}

func NewPriorityQueue[T PriorityItem]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		nodes: make([]T, 0),
	}
}

func (pq *PriorityQueue[T]) Push(item T) {
	heap.Push(&pq.nodes, item)
}

func (pq *PriorityQueue[T]) Pop() T {
	return heap.Pop(&pq.nodes).(T)
}

func (pq *PriorityQueue[T]) Top() T {
	return pq.nodes[0]
}

func (pq *PriorityQueue[T]) Clear() {
	pq.nodes = make([]T, 0)
}
