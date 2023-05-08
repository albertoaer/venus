package comm

type OpenableChannel[T any] interface {
	MessageChannel
	Open(T) (Sender, error)
}

type BinaryPacket[T any] struct {
	Data    []byte
	Address T
}

type BinaryChannel[T any] interface {
	Emitter() <-chan BinaryPacket[T]
	Send(BinaryPacket[T]) error
	Start() error
}

type binaryAdapter[T any] struct {
	adapted   BinaryChannel[T]
	serialier MessageSerializer
	emitter   chan ChannelEvent
	senders   map[any]Sender
}

func AdaptBinaryChannel[T any](adapted BinaryChannel[T]) OpenableChannel[T] {
	return &binaryAdapter[T]{
		adapted:   adapted,
		serialier: newSerializer(),
		emitter:   make(chan ChannelEvent),
		senders:   make(map[any]Sender),
	}
}

func (adapter *binaryAdapter[_]) Emitter() <-chan ChannelEvent {
	return adapter.emitter
}

func (adapter *binaryAdapter[T]) Start() error {
	go func() {
		emitter := adapter.adapted.Emitter()
		for packet := range emitter {
			if data, err := adapter.serialier.Deserialize(packet.Data); err == nil {
				sender := adapter.open(packet.Address)
				adapter.emitter <- struct {
					Message
					Sender
				}{
					data,
					sender,
				}
			}
		}
	}()
	return adapter.adapted.Start()
}

func (adapter *binaryAdapter[T]) Open(address T) (Sender, error) {
	return adapter.open(address), nil
}

func (adapter *binaryAdapter[T]) open(address T) Sender {
	if sender, exists := adapter.senders[address]; exists {
		return sender
	}
	sender := &binarySender[T]{
		channel:   adapter.adapted,
		serialier: adapter.serialier,
		address:   address,
	}
	adapter.senders[address] = sender
	return sender
}

type binarySender[T any] struct {
	channel   BinaryChannel[T]
	serialier MessageSerializer
	address   T
}

func (sender binarySender[T]) Send(message Message) (bool, error) {
	data, err := sender.serialier.Serialize(message)
	if err != nil {
		return false, err
	}
	sender.channel.Send(BinaryPacket[T]{
		Data:    data,
		Address: sender.address,
	})
	return false, nil
}
