package govenus

type OneShotTask[T any] struct {
	Action func(Context[T])
}

func NewOneShotTask[T any](action func(Context[T])) Task[T] {
	return &OneShotTask[T]{action}
}

func (task *OneShotTask[T]) Perform(ctx Context[T]) {
	task.Action(ctx)
}

func (task *OneShotTask[T]) Done() bool {
	return true
}

type CyclicTask[T any] struct {
	Action func(Context[T])
}

func NewCyclicTask[T any](action func(Context[T])) Task[T] {
	return &CyclicTask[T]{action}
}

func (task *CyclicTask[T]) Perform(ctx Context[T]) {
	task.Action(ctx)
}

func (task *CyclicTask[T]) Done() bool {
	return false
}

type ConditionalTask[T any] struct {
	Action    func(Context[T])
	Condition func() bool
}

func NewConditionalTask[T any](action func(Context[T]), condition func() bool) Task[T] {
	return &ConditionalTask[T]{action, condition}
}

func (task *ConditionalTask[T]) Perform(ctx Context[T]) {
	task.Action(ctx)
}

func (task *ConditionalTask[T]) Done() bool {
	return task.Condition()
}
