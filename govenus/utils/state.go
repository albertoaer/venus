package utils

import "sync"

type State[T any] interface {
	State() T
	SetState(T)
}

type basicState[T any] struct {
	state T
	mutex sync.RWMutex
}

func NewState[T any](value T) State[T] {
	return &basicState[T]{
		state: value,
		mutex: sync.RWMutex{},
	}
}

func (bs *basicState[T]) SetState(state T) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	bs.state = state
}

func (bs *basicState[T]) State() T {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()
	return bs.state
}
