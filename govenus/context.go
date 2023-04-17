package govenus

type Context[T any] interface {
	Runtime() Runtime[T]
}

type ContextBuilder[T any] interface {
	Context[T]
	Build() Context[T]
	SetRuntime(Runtime[T]) ContextBuilder[T]
}

type simpleContext[T any] struct {
	runtime Runtime[T]
}

func NewContextBuilder[T any]() ContextBuilder[T] {
	return &simpleContext[T]{}
}

func (sc *simpleContext[T]) Build() Context[T] {
	return sc
}

func (sc *simpleContext[T]) Runtime() Runtime[T] {
	return sc.runtime
}

func (sc *simpleContext[T]) SetRuntime(runtime Runtime[T]) ContextBuilder[T] {
	sc.runtime = runtime
	return sc
}
