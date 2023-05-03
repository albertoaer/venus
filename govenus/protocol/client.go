package protocol

import (
	"fmt"
	"net"
)

type baseClient struct {
	serializer  MessageSerializer
	id          ClientId
	peerManager *peerManager
}

func NewClient(id ClientId) Client {
	fmt.Println("Client created with id: " + id)
	return &baseClient{
		serializer:  &jsonSerializer{},
		id:          id,
		peerManager: newPeerManager(),
	}
}

func (client *baseClient) GetId() ClientId {
	return client.id
}

func (client *baseClient) spreadPacket(data []byte, msg Message) {
	if msg.Receiver != nil && *msg.Receiver == client.id {
		return
	}
	packets := client.peerManager.spreadDataByClient(msg.Receiver, msg.Sender, data, true)
	for _, packet := range packets {
		fmt.Println("Packet sent to " + packet.Address.String())
		packet.Channel.Send(packet)
	}
}

func (client *baseClient) ProcessPacket(packet Packet) (msg Message, err error) {
	if msg, err = client.serializer.Deserialize(packet.Data); err != nil {
		return
	}
	if msg.Sender == client.id {
		return
	}
	valid := client.peerManager.handlePeerData(msg.Sender, packet.Address, packet.Channel, msg.Timestamp)
	if !valid {
		return
	}
	client.spreadPacket(packet.Data, msg)
	return
}

func (client *baseClient) ProcessMessage(msg Message) error {
	packet, err := client.serializer.Serialize(msg)
	if err != nil {
		return err
	}
	client.spreadPacket(packet, msg)
	return nil
}

func (client *baseClient) ForceAlias(id ClientId, address net.Addr, channel PacketChannel) error {
	client.peerManager.addPeer(id, address, channel)
	return nil
}
