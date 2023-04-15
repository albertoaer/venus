package govenus

type Task[T any] interface {
	Perform(Context[T])
	Done() bool
}

type Promise[T any] interface {
	IsDone() bool
	OnDone(Task[T]) Promise[T]
}

type Runtime[T any] interface {
	State() T
	SetState(T)
	Launch(Task[T]) Promise[T]
	Start()
	Stop()
}

type Context[T any] interface {
	Runtime() Runtime[T]
}
