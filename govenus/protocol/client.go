package protocol

import (
	"math/rand"
	"net"

	"github.com/oklog/ulid"

	"github.com/albertoaer/venus/govenus/utils"
)

func (msg Message) Priority() int64 {
	return msg.Timestamp
}

type activeCommunication struct {
	unordered     *utils.PriorityQueue[Message]
	provider      PacketProvider
	address       net.Addr
	lastTimestamp int64
}

func newActiveCommunication(provider PacketProvider, address net.Addr) *activeCommunication {
	return &activeCommunication{
		unordered:     utils.NewPriorityQueue[Message](),
		provider:      provider,
		address:       address,
		lastTimestamp: -1,
	}
}

type baseClient struct {
	serializer           MessageSerializer
	id                   ClientId
	packetCallback       func(Packet)
	messageCallback      func(Message)
	activeCommunications map[ClientId]*activeCommunication
}

func NewClient() Client {
	timestamp := ulid.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(int64(timestamp))), 0)
	id := ulid.MustNew(timestamp, entropy)
	return &baseClient{
		serializer:           &jsonSerializer{},
		id:                   ClientId(id.String()),
		packetCallback:       nil,
		messageCallback:      nil,
		activeCommunications: make(map[ClientId]*activeCommunication),
	}
}

func (client *baseClient) GetId() ClientId {
	return client.id
}

func (client *baseClient) SetPacketCallback(callback func(Packet)) {
	client.packetCallback = callback
}

func (client *baseClient) SetMessageCallback(callback func(Message)) {
	client.messageCallback = callback
}

func (client *baseClient) processMessage(msg Message, comm *activeCommunication) {
	switch msg.Type {
	case MESSAGE_TYPE_BEGIN:
	case MESSAGE_TYPE_PERFORM:
	case MESSAGE_TYPE_INFO:
	}
}

func (client *baseClient) ProcessPacket(packet Packet) {
	msg, err := client.serializer.Deserialize(packet.Data)
	if err != nil { // TODO: handle properly
		return
	}
	if msg.ReceiverClientId != client.id {
		return
	}
	if comm, exists := client.activeCommunications[msg.ReceiverClientId]; exists {
		client.processMessage(msg, comm)
	} else {
		comm = newActiveCommunication(packet.Provider, packet.Address)
		client.activeCommunications[msg.ReceiverClientId] = comm
		client.processMessage(msg, comm)
	}
}

func (client *baseClient) ProcessMessage(msg Message) {
	if msg.SenderClientId != client.id {
		return
	}
	if comm, exists := client.activeCommunications[msg.SenderClientId]; exists && client.packetCallback != nil {
		data, err := client.serializer.Serialize(msg)
		if err != nil { // TODO: handle properly
			return
		}
		client.packetCallback(Packet{
			Data:     data,
			Address:  comm.address,
			Provider: comm.provider,
		})
	}
}
