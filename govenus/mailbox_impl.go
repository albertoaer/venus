package govenus

import "github.com/albertoaer/venus/govenus/protocol"

type RuntimeMailbox[T any] struct {
	runtime         Runtime[T]
	responses       map[protocol.Verb]Task[T]
	defaultResponse Task[T]
}

func Mailboxed[T any](runtime Runtime[T]) *RuntimeMailbox[T] {
	return &RuntimeMailbox[T]{
		runtime:         runtime,
		responses:       make(map[protocol.Verb]Task[T]),
		defaultResponse: nil,
	}
}

func (rm *RuntimeMailbox[T]) On(verb protocol.Verb, task Task[T]) {
	rm.responses[verb] = task
}

func (rm *RuntimeMailbox[T]) OnDefault(task Task[T]) {
	rm.defaultResponse = task
}

func (rm *RuntimeMailbox[T]) Notify(message protocol.Message) {
	if message.Type != protocol.MESSAGE_TYPE_PERFORM {
		return
	}
	context := rm.runtime.InitializeContextBuilder()
	context.SetMessage(message)
	if task, exists := rm.responses[message.Verb]; exists {
		rm.runtime.LaunchWith(task, context)
	} else if rm.defaultResponse != nil {
		rm.runtime.LaunchWith(rm.defaultResponse, context)
	}
}
