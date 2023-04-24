package govenus

type EventContext[T any] interface {
	RuntimeContext
	Event() T
}

type eventContextWrapper[T any] struct {
	context RuntimeContext
	event   T
}

func (cw *eventContextWrapper[T]) Runtime() Runtime {
	return cw.context.Runtime()
}

func (cw *eventContextWrapper[T]) IsAvailable() bool {
	return cw.context.IsAvailable()
}

func (cw *eventContextWrapper[T]) Event() T {
	return cw.event
}

type EventTask[T any] func(EventContext[T]) bool

type EventTaskBuilder[T any] interface {
	SetTask(EventTask[T]) EventTaskBuilder[T]
	SetEvent(Event T) EventTaskBuilder[T]
	Build() Task
}

type eventTaskBuilder[T any] struct {
	task  EventTask[T]
	event T
}

func NewEventTaskBuilder[T any]() EventTaskBuilder[T] {
	return &eventTaskBuilder[T]{
		task: func(sc EventContext[T]) bool { return true },
	}
}

func (tb *eventTaskBuilder[T]) SetTask(behave EventTask[T]) EventTaskBuilder[T] {
	tb.task = behave
	return tb
}

func (tb *eventTaskBuilder[T]) SetEvent(event T) EventTaskBuilder[T] {
	tb.event = event
	return tb
}

func (tb *eventTaskBuilder[T]) Build() Task {
	return func(context RuntimeContext) bool {
		return tb.task(&eventContextWrapper[T]{context: context, event: tb.event})
	}
}
