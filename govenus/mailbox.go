package govenus

import (
	"github.com/albertoaer/venus/govenus/comm"
)

type MailEvent struct {
	Message comm.Message
	Client  comm.Client
}

type MailContext = EventContext[MailEvent]

type MailTask = EventTask[MailEvent]

type RuntimeMailbox struct {
	runtime         Runtime
	responses       map[string]MailTask
	defaultResponse MailTask
}

func Mailboxed(runtime Runtime) *RuntimeMailbox {
	return &RuntimeMailbox{
		runtime:         runtime,
		responses:       make(map[string]MailTask),
		defaultResponse: nil,
	}
}

func (rm *RuntimeMailbox) On(verb string, task MailTask) {
	rm.responses[verb] = task
}

func (rm *RuntimeMailbox) OnDefault(task MailTask) {
	rm.defaultResponse = task
}

func (rm *RuntimeMailbox) Notify(message comm.Message, client comm.Client) {
	if message.Receiver != nil && *message.Receiver != client.GetId() {
		return
	}
	context := rm.runtime.NewContext()
	taskBuilder := NewEventTaskBuilder[MailEvent]()
	taskBuilder.SetEvent(MailEvent{
		Message: message,
		Client:  client,
	})
	if task, exists := rm.responses[message.Verb]; exists {
		taskBuilder.SetTask(EventTask[MailEvent](task))
	} else if rm.defaultResponse != nil {
		taskBuilder.SetTask(EventTask[MailEvent](rm.defaultResponse))
	}
	rm.runtime.LaunchWith(taskBuilder.Build(), context)
}
