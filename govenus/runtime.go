package govenus

type Task[T any] func(Context[T]) bool

type Promise[T any] interface {
	IsDone() bool
	OnDone(Task[T]) Promise[T]
	OnDoneWith(Task[T], ContextBuilder[T]) Promise[T]
}

type Runtime[T any] interface {
	State() T
	SetState(T)
	InitializeContextBuilder() ContextBuilder[T]
	Launch(Task[T]) Promise[T]
	LaunchWith(Task[T], ContextBuilder[T]) Promise[T]
	Start()
	Stop()
}
