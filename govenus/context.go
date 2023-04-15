package govenus

type simpleContext[T any] struct {
	runtime Runtime[T]
}

func NewMockContext[T any](runtime Runtime[T]) Context[T] {
	return &simpleContext[T]{
		runtime: runtime,
	}
}

func (sc *simpleContext[T]) Runtime() Runtime[T] {
	return sc.runtime
}
