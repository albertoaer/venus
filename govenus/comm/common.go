package comm

type Message struct {
	Sender    string
	Receiver  *string // Optional
	Timestamp int64
	Verb      string
	Args      []string
	Options   map[string]string
	Payload   []byte
	Distance  uint32
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

type MessageSerializer interface {
	Deserialize([]byte) (Message, error)
	Serialize(Message) ([]byte, error)
}

type Mailbox interface {
	Notify(ChannelEvent, Client)
}

type Client interface {
	GetId() string
	Send(Message) error
	Attach(Mailbox)
	Detach(Mailbox)
	StartChannel(MessageChannel) error
}
