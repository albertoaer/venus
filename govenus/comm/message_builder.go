package comm

import "time"

type messageBuilder struct {
	message Message
}

type MessageBuilder interface {
	SetSender(string) MessageBuilder
	SetReceiver(string) MessageBuilder
	SetTimestamp(int64) MessageBuilder
	SetVerb(string) MessageBuilder
	SetArgs([]string) MessageBuilder
	SetOptions(map[string]string) MessageBuilder
	SetPayload([]byte) MessageBuilder
	SetDistance(uint32) MessageBuilder
	Build() Message
}

func NewMessageBuilder(sender string) MessageBuilder {
	return &messageBuilder{
		message: Message{
			Args:      make([]string, 0),
			Options:   make(map[string]string, 0),
			Payload:   make([]byte, 0),
			Receiver:  nil,
			Sender:    sender,
			Timestamp: time.Now().UnixMilli(),
			Verb:      Ping,
			Distance:  0,
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

func (mb *messageBuilder) SetDistance(distance uint32) MessageBuilder {
	mb.message.Distance = distance
	return mb
}

func (mb *messageBuilder) SetReceiver(receiver string) MessageBuilder {
	mb.message.Receiver = &receiver
	return mb
}

func (mb *messageBuilder) SetSender(sender string) MessageBuilder {
	mb.message.Sender = sender
	return mb
}

func (mb *messageBuilder) SetTimestamp(timestamp int64) MessageBuilder {
	mb.message.Timestamp = timestamp
	return mb
}

func (mb *messageBuilder) SetVerb(verb string) MessageBuilder {
	mb.message.Verb = verb
	return mb
}
