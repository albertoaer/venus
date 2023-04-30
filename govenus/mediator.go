package govenus

import (
	"github.com/albertoaer/venus/govenus/protocol"
	"github.com/albertoaer/venus/govenus/utils"
)

type Mediator struct {
	mailboxes *utils.ConcurrentArray[MailBox]
	client    protocol.Client
}

func NewMediator(client protocol.Client) *Mediator {
	mediator := &Mediator{
		mailboxes: utils.NewArray[MailBox](20),
		client:    client,
	}
	client.SetMessageCallback(mediator.onMessage)
	return mediator
}

func (mediator *Mediator) onMessage(message protocol.Message) {
	mediator.mailboxes.ForEach(func(_ int, mb MailBox) {
		mb.Notify(message, mediator)
	})
}

func (mediator *Mediator) GetClient() protocol.Client {
	return mediator.client
}

func (mediator *Mediator) Attach(mailbox MailBox) {
	mediator.mailboxes.Add(mailbox)
}

func (mediator *Mediator) Detach(mailbox MailBox) {
	mediator.mailboxes.Remove(mailbox)
}

func (mediator *Mediator) GetId() protocol.ClientId {
	return mediator.client.GetId()
}

func (mediator *Mediator) Send(message protocol.Message) {
	mediator.client.ProcessMessage(message)
}

func (mediator *Mediator) StartChannel(channel protocol.PacketChannel) (err error) {
	if err = channel.Start(); err == nil {
		go func() {
			for {
				emitter := channel.Emitter()
				packet := <-emitter
				mediator.client.ProcessPacket(packet)
			}
		}()
	}
	return
}
