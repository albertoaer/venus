package govenus

import (
	"sync"

	"github.com/albertoaer/venus/govenus/utils"
)

type funcPromise[T any] struct {
	task           Task[T]
	onDone         *funcPromise[T]
	done           bool
	contextBuilder ContextBuilder[T]
}

func createFuncPromise[T any](task Task[T], contextBuilder ContextBuilder[T]) *funcPromise[T] {
	return &funcPromise[T]{
		task:           task,
		onDone:         nil,
		done:           false,
		contextBuilder: contextBuilder,
	}
}

func (p *funcPromise[T]) IsDone() bool {
	return p.done
}

func (p *funcPromise[T]) OnDone(task Task[T]) Promise[T] {
	promise := createFuncPromise(task, p.contextBuilder)
	p.onDone = promise
	return promise
}

func (p *funcPromise[T]) runOnce() {
	// TODO: avoid use the context without building it
	p.task.Perform(p.contextBuilder.Build())
	p.done = p.task.Done()
}

type queueRuntime[T any] struct {
	queue     utils.ConcurrentQueue[*funcPromise[T]]
	state     T
	on        bool
	onLocker  sync.Mutex
	startLock sync.Mutex
}

func NewDefaultRuntime[T any](initial T) Runtime[T] {
	return &queueRuntime[T]{
		queue:     utils.NewConcurrentQueue[*funcPromise[T]](),
		state:     initial,
		on:        true,
		onLocker:  sync.Mutex{},
		startLock: sync.Mutex{},
	}
}

func (mqr *queueRuntime[T]) State() T {
	return mqr.state
}

func (mqr *queueRuntime[T]) SetState(state T) {
	mqr.state = state
}

func (mqr *queueRuntime[T]) InitializeContextBuilder() ContextBuilder[T] {
	return NewContextBuilder[T]().SetRuntime(mqr)
}

func (mqr *queueRuntime[T]) Launch(task Task[T]) Promise[T] {
	return mqr.LaunchWith(task, mqr.InitializeContextBuilder())
}

func (mqr *queueRuntime[T]) LaunchWith(task Task[T], contextBuilder ContextBuilder[T]) Promise[T] {
	promise := createFuncPromise(task, contextBuilder)
	mqr.queue.Enqueue(promise)
	return promise
}

func (mqr *queueRuntime[T]) Start() {
	mqr.startLock.Lock()
	mqr.onLocker.Lock()
	mqr.on = true
	mqr.onLocker.Unlock()
	for mqr.on {
		if !mqr.queue.Empty() {
			promise := mqr.queue.Dequeue()
			promise.runOnce()
			if promise.IsDone() { // Only react to it finalization if it actually ends
				if promise.onDone != nil {
					mqr.queue.Enqueue(promise.onDone)
				}
			} else {
				mqr.queue.Enqueue(promise)
			}
		}
	}
}

func (mqr *queueRuntime[T]) Stop() {
	mqr.onLocker.Lock()
	mqr.on = false
	mqr.onLocker.Unlock()
	mqr.startLock.TryLock()
	mqr.startLock.Unlock()
}
