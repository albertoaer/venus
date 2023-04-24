package govenus

import "github.com/albertoaer/venus/govenus/protocol"

type MailContext = EventContext[protocol.Message]

type MailTask = EventTask[protocol.Message]

type RuntimeMailbox struct {
	runtime         Runtime
	responses       map[protocol.Verb]MailTask
	defaultResponse MailTask
}

func Mailboxed(runtime Runtime) *RuntimeMailbox {
	return &RuntimeMailbox{
		runtime:         runtime,
		responses:       make(map[protocol.Verb]MailTask),
		defaultResponse: nil,
	}
}

func (rm *RuntimeMailbox) On(verb protocol.Verb, task MailTask) {
	rm.responses[verb] = task
}

func (rm *RuntimeMailbox) OnDefault(task MailTask) {
	rm.defaultResponse = task
}

func (rm *RuntimeMailbox) Notify(message protocol.Message) {
	if message.Type() != protocol.MESSAGE_TYPE_PERFORM {
		return
	}
	context := rm.runtime.NewContext()
	taskBuilder := NewEventTaskBuilder[protocol.Message]()
	taskBuilder.SetEvent(message)
	if task, exists := rm.responses[message.Verb()]; exists {
		taskBuilder.SetTask(EventTask[protocol.Message](task))
	} else if rm.defaultResponse != nil {
		taskBuilder.SetTask(EventTask[protocol.Message](rm.defaultResponse))
	}
	rm.runtime.LaunchWith(taskBuilder.Build(), context)
}
