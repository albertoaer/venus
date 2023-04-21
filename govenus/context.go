package govenus

import (
	"errors"

	"github.com/albertoaer/venus/govenus/protocol"
)

type Context[T any] interface {
	Runtime() Runtime[T]        // The runtime where the context is being run
	Message() *protocol.Message // The message which the context depends on
	IsAvailable() bool          // Determines when the context is available to run
}

type ContextBuilder[T any] interface {
	Build() (Context[T], error)
	SetRuntime(Runtime[T]) ContextBuilder[T]
	SetMessage(protocol.Message) ContextBuilder[T]
	AddAvailabilityCondition(func() bool) ContextBuilder[T]
}

type simpleContext[T any] struct {
	runtime           Runtime[T]
	message           *protocol.Message
	checkAvailability func() bool
}

func NewContextBuilder[T any]() ContextBuilder[T] {
	return &simpleContext[T]{
		runtime:           nil,
		message:           nil,
		checkAvailability: func() bool { return true },
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

func (sc *simpleContext[T]) SetMessage(message protocol.Message) ContextBuilder[T] {
	sc.message = &message
	return sc
}

func (sc *simpleContext[T]) AddAvailabilityCondition(condition func() bool) ContextBuilder[T] {
	prev := sc.checkAvailability
	sc.checkAvailability = func() bool { return condition() && prev() }
	return sc
}

func (sc *simpleContext[T]) IsAvailable() bool {
	return sc.checkAvailability()
}
