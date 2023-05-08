package protocol

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
}

type Sender interface {
	// Returns a possible error and wether has ended or not
	Send(Message) (done bool, err error)
}

type ChannelEvent struct {
	Message
	Sender
}

type MessageChannel interface {
	Emitter() <-chan ChannelEvent
	Start() error
}

type Client interface {
	GetId() ClientId
	GotMessage(Message, Sender) error
	Send(Message) error
}
