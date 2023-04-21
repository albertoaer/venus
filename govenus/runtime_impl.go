package govenus

import (
	"sync"

	"github.com/albertoaer/venus/govenus/utils"
)

type funcPromise[T any] struct {
	task     Task[T]
	done     bool
	context  Context[T]
	prepared bool
	mutex    sync.RWMutex
}

func createFuncPromise[T any](task Task[T], context Context[T], prepared bool) *funcPromise[T] {
	return &funcPromise[T]{
		task:     task,
		done:     false,
		context:  context,
		prepared: prepared,
	}
}

func (p *funcPromise[T]) IsDone() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.done
}

func (p *funcPromise[T]) OnDone(task Task[T]) Promise[T] {
	return p.OnDoneWith(task, p.context.Runtime().InitializeContextBuilder())
}

func (p *funcPromise[T]) OnDoneWith(task Task[T], contextBuilder ContextBuilder[T]) Promise[T] {
	contextBuilder.AddAvailabilityCondition(
		func() bool {
			return p.IsDone()
		},
	)
	return p.context.Runtime().LaunchWith(task, contextBuilder)
}

func (p *funcPromise[T]) runOnce() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.done {
		p.done = p.task(p.context)
	}
}

type queueRuntime[T any] struct {
	queue      utils.ConcurrentQueue[*funcPromise[T]]
	state      T
	on         bool
	onMutex    sync.Mutex
	startMutex sync.Mutex
}

func NewDefaultRuntime[T any](initial T) Runtime[T] {
	return &queueRuntime[T]{
		queue:      utils.NewConcurrentQueue[*funcPromise[T]](),
		state:      initial,
		on:         true,
		onMutex:    sync.Mutex{},
		startMutex: sync.Mutex{},
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
	contextBuilder.SetRuntime(mqr) // Prevent unexpected behaviour
	context, err := contextBuilder.Build()
	if err != nil {
		panic(err)
	}
	promise := createFuncPromise(task, context, false)
	mqr.queue.Enqueue(promise)
	return promise
}

func (mqr *queueRuntime[T]) Start() {
	mqr.startMutex.Lock()
	mqr.onMutex.Lock()
	mqr.on = true
	mqr.onMutex.Unlock()
	for mqr.on {
		if !mqr.queue.Empty() {
			promise := mqr.queue.Dequeue()
			if promise.context.IsAvailable() {
				promise.runOnce()
				if !promise.IsDone() {
					mqr.queue.Enqueue(promise)
				}
			} else {
				mqr.queue.Enqueue(promise)
			}
		}
	}
}

func (mqr *queueRuntime[T]) Stop() {
	mqr.onMutex.Lock()
	mqr.on = false
	mqr.onMutex.Unlock()
	mqr.startMutex.TryLock()
	mqr.startMutex.Unlock()
}
