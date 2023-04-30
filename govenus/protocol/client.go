package protocol

import (
	"errors"
	"net"
)

type knownHost struct {
	channel PacketChannel
	address net.Addr
}

func newKnownHost(channel PacketChannel, address net.Addr) *knownHost {
	return &knownHost{
		channel: channel,
		address: address,
	}
}

type baseClient struct {
	serializer      MessageSerializer
	id              ClientId
	messageCallback func(Message)
	knownHosts      map[ClientId]*knownHost
}

func NewClient(id ClientId) Client {
	return &baseClient{
		serializer:      &jsonSerializer{},
		id:              id,
		messageCallback: nil,
		knownHosts:      make(map[ClientId]*knownHost),
	}
}

func (client *baseClient) GetId() ClientId {
	return client.id
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
	} else if comm, exists := client.knownHosts[*msg.Receiver]; exists {
		data, err := client.serializer.Serialize(msg)
		if err != nil {
			return err
		}
		comm.channel.Send(Packet{
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
