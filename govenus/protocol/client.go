package protocol

import (
	"errors"
	"net"

	"github.com/albertoaer/venus/govenus/utils"
)

type knownHost struct {
	unordered     *utils.PriorityQueue[Message]
	channel       PacketChannel
	address       net.Addr
	lastTimestamp int64
}

func newKnownHost(channel PacketChannel, address net.Addr) *knownHost {
	return &knownHost{
		unordered: utils.NewPriorityQueue(
			func(m1, m2 Message) bool {
				return m1.Timestamp < m2.Timestamp
			},
		),
		channel:       channel,
		address:       address,
		lastTimestamp: -1,
	}
}

type baseClient struct {
	serializer      MessageSerializer
	id              ClientId
	packetCallback  func(Packet)
	messageCallback func(Message)
	knownHosts      map[ClientId]*knownHost
}

func NewClient(id ClientId) Client {
	return &baseClient{
		serializer:      &jsonSerializer{},
		id:              id,
		packetCallback:  nil,
		messageCallback: nil,
		knownHosts:      make(map[ClientId]*knownHost),
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

func (client *baseClient) ProcessPacket(packet Packet) error {
	msg, err := client.serializer.Deserialize(packet.Data)
	if err != nil {
		return err
	}
	if _, exists := client.knownHosts[msg.Sender]; exists {
		client.messageCallback(msg)
	} else {
		client.knownHosts[msg.Sender] = newKnownHost(packet.Channel, packet.Address)
		client.messageCallback(msg)
	}
	return nil
}

func (client *baseClient) ProcessMessage(msg Message) error {
	if msg.Receiver == nil {
		// TODO: Notify all hosts
		return nil
	} else if comm, exists := client.knownHosts[*msg.Receiver]; exists && client.packetCallback != nil {
		data, err := client.serializer.Serialize(msg)
		if err != nil {
			return err
		}
		client.packetCallback(Packet{
			Data:    data,
			Address: comm.address,
			Channel: comm.channel,
		})
		return nil
	} else {
		return errors.New("client not found")
	}
}

func (client *baseClient) ForceAlias(id ClientId, address net.Addr, channel PacketChannel) error {
	client.knownHosts[id] = newKnownHost(channel, address)
	return nil
}
