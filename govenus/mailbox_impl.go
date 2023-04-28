package govenus

import "github.com/albertoaer/venus/govenus/protocol"

type MailEvent struct {
	Message protocol.Message
	Client  ClientService
}

type MailContext = EventContext[MailEvent]

type MailTask = EventTask[MailEvent]

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

func (rm *RuntimeMailbox) Notify(message protocol.Message, client ClientService) {
	if message.Type() != protocol.MESSAGE_TYPE_PERFORM {
		return
	}
	context := rm.runtime.NewContext()
	taskBuilder := NewEventTaskBuilder[MailEvent]()
	taskBuilder.SetEvent(MailEvent{
		Message: message,
		Client:  client,
	})
	if task, exists := rm.responses[message.Verb()]; exists {
		taskBuilder.SetTask(EventTask[MailEvent](task))
	} else if rm.defaultResponse != nil {
		taskBuilder.SetTask(EventTask[MailEvent](rm.defaultResponse))
	}
	rm.runtime.LaunchWith(taskBuilder.Build(), context)
}
