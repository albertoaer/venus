package utils

import "sync"

type ConcurrentArray[T any] struct {
	elements []T
	rwMutex  sync.RWMutex
}

func NewArray[T any](initialCapacity int) *ConcurrentArray[T] {
	return &ConcurrentArray[T]{
		elements: make([]T, 0, initialCapacity),
		rwMutex:  sync.RWMutex{},
	}
}

func (arr *ConcurrentArray[T]) Length() int {
	arr.rwMutex.RLock()
	defer arr.rwMutex.RUnlock()
	return len(arr.elements)
}

func (arr *ConcurrentArray[T]) Add(item T) {
	arr.rwMutex.Lock()
	defer arr.rwMutex.Unlock()
	arr.elements = append(arr.elements, item)
}

func (arr *ConcurrentArray[T]) Remove(item T) {
	arr.rwMutex.Lock()
	defer arr.rwMutex.Unlock()
	for idx, element := range arr.elements {
		if interface{}(element) == interface{}(item) {
			arr.elements = append(arr.elements[:idx], arr.elements[idx+1:]...)
		}
	}
}

func (arr *ConcurrentArray[T]) Get(index int) T {
	arr.rwMutex.RLock()
	defer arr.rwMutex.RUnlock()
	return arr.elements[index]
}

func (arr *ConcurrentArray[T]) ForEach(op func(int, T)) {
	arr.rwMutex.RLock()
	defer arr.rwMutex.RUnlock()
	for idx, element := range arr.elements {
		op(idx, element)
	}
}
