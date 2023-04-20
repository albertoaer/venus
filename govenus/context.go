package govenus

import (
	"errors"

	"github.com/albertoaer/venus/govenus/protocol"
)

type Context[T any] interface {
	Runtime() Runtime[T]
	Message() *protocol.Message
}

type ContextBuilder[T any] interface {
	Build() (Context[T], error)
	SetRuntime(Runtime[T]) ContextBuilder[T]
	SetMessage(protocol.Message)
}

type simpleContext[T any] struct {
	runtime Runtime[T]
	message *protocol.Message
}

func NewContextBuilder[T any]() ContextBuilder[T] {
	return &simpleContext[T]{
		runtime: nil,
		message: nil,
	}
}

func (sc *simpleContext[T]) Build() (Context[T], error) {
	if sc.runtime == nil {
		return nil, errors.New("cannot build a context without runtime")
	}
	return sc, nil
}

func (sc *simpleContext[T]) Runtime() Runtime[T] {
	return sc.runtime
}

func (sc *simpleContext[T]) SetRuntime(runtime Runtime[T]) ContextBuilder[T] {
	sc.runtime = runtime
	return sc
}

func (sc *simpleContext[T]) Message() *protocol.Message {
	return sc.message
}

func (sc *simpleContext[T]) SetMessage(message protocol.Message) {
	sc.message = &message
}
