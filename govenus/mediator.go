package govenus

import (
	"net"

	"github.com/albertoaer/venus/govenus/protocol"
	"github.com/albertoaer/venus/govenus/utils"
)

type Mediator struct {
	mailboxes *utils.ConcurrentArray[MailBox]
	channel   protocol.PacketChannel
	client    protocol.Client
}

func NewMediator(channel protocol.PacketChannel) *Mediator {
	client := protocol.NewClient()
	mediator := &Mediator{
		mailboxes: utils.NewArray[MailBox](20),
		channel:   channel,
		client:    client,
	}
	client.SetMessageCallback(mediator.onMessage)
	client.SetPacketCallback(mediator.onPacket)
	return mediator
}

func (mediator *Mediator) onMessage(message protocol.Message) {
	mediator.mailboxes.ForEach(func(_ int, mb MailBox) {
		mb.Notify(message, mediator)
	})
}

func (mediator *Mediator) onPacket(packet protocol.Packet) {
	packet.Channel.Send(packet)
}

func (mediator *Mediator) Attach(mailbox MailBox) {
	mediator.mailboxes.Add(mailbox)
}

func (mediator *Mediator) Detach(mailbox MailBox) {
	mediator.mailboxes.Remove(mailbox)
}

func (mediator *Mediator) SetAddress(id protocol.ClientId, address net.Addr) {
	mediator.client.ForceAlias(id, address, mediator.channel)
}

func (mediator *Mediator) GetId() protocol.ClientId {
	return mediator.client.GetId()
}

func (mediator *Mediator) Send(message protocol.Message) {
	mediator.client.ProcessMessage(message)
}

func (mediator *Mediator) Start() (err error) {
	err = mediator.channel.Start()
	if err == nil {
		for {
			emitter := mediator.channel.Emitter()
			packet := <-emitter
			mediator.client.ProcessPacket(packet)
		}
	}
	return
}
