package protocol

import "net"

type Packet struct {
	Data    []byte
	Address net.Addr
	Channel PacketChannel
}

type PacketChannel interface {
	Emitter() <-chan Packet
	Start() error
	Send(Packet) error
}

type ClientId string

type Verb string

type Message struct {
	Sender    ClientId
	Receiver  *ClientId // Optional
	Timestamp int64
	Verb      Verb
	Args      []string
	Options   map[string]string
	Payload   []byte
	Valid     bool // Not serialized, true for any incoming/created not default Message
}

type MessageSerializer interface {
	Deserialize([]byte) (Message, error)
	Serialize(Message) ([]byte, error)
}

type Client interface {
	GetId() ClientId
	ProcessPacket(Packet) (Message, error)
	ProcessMessage(Message) error
	ForceAlias(ClientId, net.Addr, PacketChannel) error
}
