package protocol

import "time"

type standardMessage struct {
	args              []string
	options           map[string]string
	payload           []byte
	previousTimestamp *int64
	receiver          *ClientId
	sender            ClientId
	timestamp         int64
	type_             MessageType
	verb              Verb
}

func (sm *standardMessage) Args() []string {
	return sm.args
}

func (sm *standardMessage) Options() map[string]string {
	return sm.options
}

func (sm *standardMessage) Payload() []byte {
	return sm.payload
}

func (sm *standardMessage) PreviousTimestamp() *int64 {
	return sm.previousTimestamp
}

func (sm *standardMessage) Receiver() *ClientId {
	return sm.receiver
}

func (sm *standardMessage) Sender() ClientId {
	return sm.sender
}

func (sm *standardMessage) Timestamp() int64 {
	return sm.timestamp
}

func (sm *standardMessage) Type() MessageType {
	return sm.type_
}

func (sm *standardMessage) Verb() Verb {
	return sm.verb
}

type messageBuilder struct {
	args              []string
	options           map[string]string
	payload           []byte
	previousTimestamp *int64
	receiver          *ClientId
	sender            ClientId
	timestamp         int64
	type_             MessageType
	verb              Verb
}

type MessageBuilder interface {
	SetSender(ClientId) MessageBuilder
	SetReceiver(ClientId) MessageBuilder
	SetTimestamp(int64) MessageBuilder
	SetPreviousTimestamp(int64) MessageBuilder
	SetVerb(Verb) MessageBuilder
	SetType(MessageType) MessageBuilder
	SetArgs([]string) MessageBuilder
	SetOptions(map[string]string) MessageBuilder
	SetPayload([]byte) MessageBuilder
	Build() Message
}

func NewMessageBuilder() MessageBuilder {
	return &messageBuilder{
		args:              make([]string, 0),
		options:           make(map[string]string, 0),
		payload:           make([]byte, 0),
		previousTimestamp: nil,
		receiver:          nil,
		sender:            *new(ClientId),
		timestamp:         time.Now().UnixMilli(),
		type_:             MESSAGE_TYPE_PERFORM,
		verb:              Ping,
	}
}

func (mb *messageBuilder) Build() Message {
	return &standardMessage{
		args:              mb.args,
		options:           mb.options,
		payload:           mb.payload,
		previousTimestamp: mb.previousTimestamp,
		receiver:          mb.receiver,
		sender:            mb.sender,
		timestamp:         mb.timestamp,
		type_:             mb.type_,
		verb:              mb.verb,
	}
}

func (mb *messageBuilder) SetArgs(args []string) MessageBuilder {
	mb.args = args
	return mb
}

func (mb *messageBuilder) SetOptions(options map[string]string) MessageBuilder {
	mb.options = options
	return mb
}

func (mb *messageBuilder) SetPayload(payload []byte) MessageBuilder {
	mb.payload = payload
	return mb
}

func (mb *messageBuilder) SetPreviousTimestamp(previousTimestamp int64) MessageBuilder {
	mb.previousTimestamp = &previousTimestamp
	return mb
}

func (mb *messageBuilder) SetReceiver(receiver ClientId) MessageBuilder {
	mb.receiver = &receiver
	return mb
}

func (mb *messageBuilder) SetSender(sender ClientId) MessageBuilder {
	mb.sender = sender
	return mb
}

func (mb *messageBuilder) SetTimestamp(timestamp int64) MessageBuilder {
	mb.timestamp = timestamp
	return mb
}

func (mb *messageBuilder) SetType(messageType MessageType) MessageBuilder {
	mb.type_ = messageType
	return mb
}

func (mb *messageBuilder) SetVerb(verb Verb) MessageBuilder {
	mb.verb = verb
	return mb
}
