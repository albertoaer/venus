package govenus

import (
	"github.com/albertoaer/venus/govenus/comm"
)

type MailEvent struct {
	Message comm.Message
	Client  comm.Client
	Sender  comm.Sender
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

func (rm *RuntimeMailbox) Notify(event comm.ChannelEvent, client comm.Client) {
	if event.Message.Receiver != nil && *event.Message.Receiver != client.GetId() {
		return
	}
	context := rm.runtime.NewContext()
	taskBuilder := NewEventTaskBuilder[MailEvent]()
	taskBuilder.SetEvent(MailEvent{
		Message: event.Message,
		Client:  client,
		Sender:  event.Sender,
	})
	if task, exists := rm.responses[event.Message.Verb]; exists {
		taskBuilder.SetTask(EventTask[MailEvent](task))
	} else if rm.defaultResponse != nil {
		taskBuilder.SetTask(EventTask[MailEvent](rm.defaultResponse))
	}
	rm.runtime.LaunchWith(taskBuilder.Build(), context)
}
