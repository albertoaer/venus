package utils

import "sync"

type concurrentQueueNode[T any] struct {
	value T
	next  *concurrentQueueNode[T]
}

type ConcurrentQueue[T any] struct {
	start     *concurrentQueueNode[T]
	end       *concurrentQueueNode[T]
	queueLock sync.RWMutex
}

func NewConcurrentQueue[T any]() ConcurrentQueue[T] {
	return ConcurrentQueue[T]{
		start:     nil,
		end:       nil,
		queueLock: sync.RWMutex{},
	}
}

func (cq *ConcurrentQueue[T]) Empty() bool {
	return cq.start == nil
}

func (cq *ConcurrentQueue[T]) Enqueue(item T) {
	cq.queueLock.Lock()
	defer cq.queueLock.Unlock()
	node := &concurrentQueueNode[T]{value: item, next: nil}
	if cq.end != nil {
		cq.end.next = node
	}
	cq.end = node
	if cq.start == nil {
		cq.start = cq.end
	}
}

func (cq *ConcurrentQueue[T]) Dequeue() T {
	cq.queueLock.Lock()
	defer cq.queueLock.Unlock()
	node := cq.start
	if node != nil {
		cq.start = cq.start.next
		if cq.start == nil {
			cq.end = nil
		}
		return node.value
	}
	return *new(T)
}
