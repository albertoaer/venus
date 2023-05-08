package runtime

import (
	"sync"

	"github.com/albertoaer/venus/govenus/utils"
)

type funcPromise struct {
	task     Task
	done     bool
	context  RuntimeContext
	prepared bool
	mutex    sync.RWMutex
}

func createFuncPromise(task Task, context RuntimeContext, prepared bool) *funcPromise {
	return &funcPromise{
		task:     task,
		done:     false,
		context:  context,
		prepared: prepared,
	}
}

func (p *funcPromise) IsDone() bool {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.done
}

func (p *funcPromise) OnDone(task Task) Promise {
	return p.OnDoneWith(task, p.context.Runtime().NewContext())
}

func (p *funcPromise) OnDoneWith(task Task, contextBuilder RuntimeContextBuilder) Promise {
	contextBuilder.AddAvailabilityCondition(
		func() bool {
			return p.IsDone()
		},
	)
	return p.context.Runtime().LaunchWith(task, contextBuilder)
}

func (p *funcPromise) runOnce() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.done {
		p.done = p.task(p.context)
	}
}

type singleRoutineRuntime struct {
	queue      utils.ConcurrentQueue[*funcPromise]
	on         bool
	onMutex    sync.Mutex
	startMutex sync.Mutex
}

func NewSRRuntime() Runtime {
	return &singleRoutineRuntime{
		queue:      utils.NewConcurrentQueue[*funcPromise](),
		on:         true,
		onMutex:    sync.Mutex{},
		startMutex: sync.Mutex{},
	}
}

func (mqr *singleRoutineRuntime) NewContext() RuntimeContextBuilder {
	return NewContextBuilder().SetRuntime(mqr)
}

func (mqr *singleRoutineRuntime) Launch(task Task) Promise {
	return mqr.LaunchWith(task, mqr.NewContext())
}

func (mqr *singleRoutineRuntime) LaunchWith(task Task, contextBuilder RuntimeContextBuilder) Promise {
	contextBuilder.SetRuntime(mqr) // Prevent unexpected behaviour
	context, err := contextBuilder.Build()
	if err != nil {
		panic(err)
	}
	promise := createFuncPromise(task, context, false)
	mqr.queue.Enqueue(promise)
	return promise
}

func (mqr *singleRoutineRuntime) Start() {
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

func (mqr *singleRoutineRuntime) Stop() {
	mqr.onMutex.Lock()
	mqr.on = false
	mqr.onMutex.Unlock()
	mqr.startMutex.TryLock()
	mqr.startMutex.Unlock()
}
