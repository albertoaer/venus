package comm

type MessageSerializer interface {
	Deserialize([]byte) (Message, error)
	Serialize(Message) ([]byte, error)
}

func newSerializer() MessageSerializer {
	return &jsonSerializer{}
}
