package protocol

import "time"

type messageBuilder struct {
	message Message
}

type MessageBuilder interface {
	SetSender(ClientId) MessageBuilder
	SetReceiver(ClientId) MessageBuilder
	SetTimestamp(int64) MessageBuilder
	SetVerb(Verb) MessageBuilder
	SetArgs([]string) MessageBuilder
	SetOptions(map[string]string) MessageBuilder
	SetPayload([]byte) MessageBuilder
	Build() Message
}

func NewMessageBuilder(sender ClientId) MessageBuilder {
	return &messageBuilder{
		message: Message{
			Args:      make([]string, 0),
			Options:   make(map[string]string, 0),
			Payload:   make([]byte, 0),
			Receiver:  nil,
			Sender:    sender,
			Timestamp: time.Now().UnixMilli(),
			Verb:      Ping,
		},
	}
}

func NewMessageBuilderFrom(original Message) MessageBuilder {
	return &messageBuilder{
		message: original,
	}
}

func (mb *messageBuilder) Build() Message {
	return mb.message
}

func (mb *messageBuilder) SetArgs(args []string) MessageBuilder {
	mb.message.Args = args
	return mb
}

func (mb *messageBuilder) SetOptions(options map[string]string) MessageBuilder {
	mb.message.Options = options
	return mb
}

func (mb *messageBuilder) SetPayload(payload []byte) MessageBuilder {
	mb.message.Payload = payload
	return mb
}

func (mb *messageBuilder) SetReceiver(receiver ClientId) MessageBuilder {
	mb.message.Receiver = &receiver
	return mb
}

func (mb *messageBuilder) SetSender(sender ClientId) MessageBuilder {
	mb.message.Sender = sender
	return mb
}

func (mb *messageBuilder) SetTimestamp(timestamp int64) MessageBuilder {
	mb.message.Timestamp = timestamp
	return mb
}

func (mb *messageBuilder) SetVerb(verb Verb) MessageBuilder {
	mb.message.Verb = verb
	return mb
}
