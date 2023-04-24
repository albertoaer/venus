package protocol

import (
	"fmt"
	"math/rand"
	"net"

	"github.com/oklog/ulid"

	"github.com/albertoaer/venus/govenus/utils"
)

type activeCommunication struct {
	unordered     *utils.PriorityQueue[Message]
	provider      PacketProvider
	address       net.Addr
	lastTimestamp int64
}

func newActiveCommunication(provider PacketProvider, address net.Addr) *activeCommunication {
	return &activeCommunication{
		unordered: utils.NewPriorityQueue(
			func(m1, m2 Message) bool {
				return m1.Timestamp() < m2.Timestamp()
			},
		),
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

func (client *baseClient) ProcessPacket(packet Packet) {
	msg, err := client.serializer.Deserialize(packet.Data)
	if err != nil { // TODO: handle properly
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	if msg.Receiver() != nil && *msg.Receiver() != client.id {
		return
	}
	if _, exists := client.activeCommunications[msg.Sender()]; exists {
		client.messageCallback(msg)
	} else {
		client.activeCommunications[msg.Sender()] = newActiveCommunication(packet.Provider, packet.Address)
		client.messageCallback(msg)
	}
}

func (client *baseClient) ProcessMessage(msg Message) {
	if msg.Sender() != client.id {
		return
	}
	if comm, exists := client.activeCommunications[msg.Sender()]; exists && client.packetCallback != nil {
		data, err := client.serializer.Serialize(msg)
		if err != nil { // TODO: handle properly
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
		client.packetCallback(Packet{
			Data:     data,
			Address:  comm.address,
			Provider: comm.provider,
		})
	}
}
