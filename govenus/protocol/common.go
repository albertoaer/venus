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

type MessageResolutionMethod int8

type MessageType int8

type ClientId string

type Verb string

type Message struct {
	SenderClientId    ClientId                `json:"sender"`
	ReceiverClientId  ClientId                `json:"receiver"`
	ResolutionMethod  MessageResolutionMethod `json:"resolution"`
	Timestamp         int64                   `json:"timestamp"`
	PreviousTimestamp int64                   `json:"previousTimestamp"`
	Type              MessageType             `json:"type"`
	Verb              Verb                    `json:"verb"`
	Payload           []byte                  `json:"payload"`
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
