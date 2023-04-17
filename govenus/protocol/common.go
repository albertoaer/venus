package protocol

import "net"

type Verb string

type Package struct {
	Data    []byte
	Address net.Addr
}

type PackageProvider interface {
	Emitter() <-chan Package
	Start() error
	Send(Package) error
}

type Message struct {
	Verb    Verb
	Args    []string
	Content []byte
}

type ConversationId string

type ConversationParticipant interface {
	GetId() ConversationId
}

type MessageSerializer interface {
	Deserialize([]byte) Message
	Serialize(Message) []byte
}
