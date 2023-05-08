package mailbox

import (
	"github.com/albertoaer/venus/govenus/protocol"
	"github.com/albertoaer/venus/govenus/runtime"
)

type MailEvent struct {
	Message protocol.Message
	Client  protocol.Client
	Sender  protocol.Sender
}

type MailContext = runtime.EventContext[MailEvent]

type MailTask = runtime.EventTask[MailEvent]

type RuntimeMailbox struct {
	runtime         runtime.Runtime
	responses       map[protocol.Verb]MailTask
	defaultResponse MailTask
}

func Mailboxed(runtime runtime.Runtime) *RuntimeMailbox {
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

func (rm *RuntimeMailbox) Notify(event protocol.ChannelEvent, client protocol.Client) {
	if event.Message.Receiver != nil && *event.Message.Receiver != client.GetId() {
		return
	}
	context := rm.runtime.NewContext()
	taskBuilder := runtime.NewEventTaskBuilder[MailEvent]()
	taskBuilder.SetEvent(MailEvent{
		Message: event.Message,
		Client:  client,
		Sender:  event.Sender,
	})
	if task, exists := rm.responses[event.Message.Verb]; exists {
		taskBuilder.SetTask(runtime.EventTask[MailEvent](task))
	} else if rm.defaultResponse != nil {
		taskBuilder.SetTask(runtime.EventTask[MailEvent](rm.defaultResponse))
	}
	rm.runtime.LaunchWith(taskBuilder.Build(), context)
}
