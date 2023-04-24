package protocol

import "net"

type Packet struct {
	Data     []byte
	Address  net.Addr
	Provider PacketProvider
}

type PacketProvider interface {
	Emitter() <-chan Packet
	Start() error
	Send(Packet) error
}

type MessageType int8

type ClientId string

type Verb string

type Message interface {
	Sender() ClientId
	Receiver() *ClientId // Optional
	Timestamp() int64
	PreviousTimestamp() *int64 // Optional
	Verb() Verb
	Type() MessageType
	Args() []string
	Options() map[string]string
	Payload() []byte
}

type MessageSerializer interface {
	Deserialize([]byte) (Message, error)
	Serialize(Message) ([]byte, error)
}

type Client interface {
	GetId() ClientId
	SetPacketCallback(func(Packet))
	SetMessageCallback(func(Message))
	ProcessPacket(Packet)
	ProcessMessage(Message)
}
