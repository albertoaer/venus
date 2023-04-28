package govenus

import (
	"github.com/albertoaer/venus/govenus/protocol"
	"github.com/albertoaer/venus/govenus/utils"
)

type Mediator struct {
	mailboxes *utils.ConcurrentArray[MailBox]
	provider  protocol.PacketProvider
	client    protocol.Client
}

func NewMediator(provider protocol.PacketProvider) *Mediator {
	client := protocol.NewClient()
	mediator := &Mediator{
		mailboxes: utils.NewArray[MailBox](20),
		provider:  provider,
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
	packet.Provider.Send(packet)
}

func (mediator *Mediator) Attach(mailbox MailBox) *Mediator {
	mediator.mailboxes.Add(mailbox)
	return mediator
}

func (mediator *Mediator) Detach(mailbox MailBox) *Mediator {
	mediator.mailboxes.Remove(mailbox)
	return mediator
}

func (mediator *Mediator) GetId() protocol.ClientId {
	return mediator.client.GetId()
}

func (mediator *Mediator) Send(message protocol.Message) {
	mediator.client.ProcessMessage(message)
}

func (mediator *Mediator) Start() (err error) {
	err = mediator.provider.Start()
	if err == nil {
		for {
			emitter := mediator.provider.Emitter()
			packet := <-emitter
			mediator.client.ProcessPacket(packet)
		}
	}
	return
}
