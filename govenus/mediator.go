package govenus

import (
	"fmt"
	"sync"

	"github.com/albertoaer/venus/govenus/protocol"
	"github.com/albertoaer/venus/govenus/utils"
)

type Mediator struct {
	mailboxes      *utils.ConcurrentArray[MailBox]
	client         protocol.Client
	channels       *utils.ConcurrentArray[protocol.PacketChannel]
	receptionMutex sync.Mutex
}

func NewMediator(client protocol.Client) *Mediator {
	mediator := &Mediator{
		mailboxes:      utils.NewArray[MailBox](20),
		client:         client,
		channels:       utils.NewArray[protocol.PacketChannel](1),
		receptionMutex: sync.Mutex{},
	}
	return mediator
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

func (mediator *Mediator) GetChannel(index int) protocol.PacketChannel {
	return mediator.channels.Get(index)
}

func (mediator *Mediator) ChannelsCount() int {
	return mediator.channels.Length()
}

func (mediator *Mediator) StartChannel(channel protocol.PacketChannel) (err error) {
	if err = channel.Start(); err == nil {
		mediator.channels.Add(channel)
		go func() {
			for {
				emitter := channel.Emitter()
				packet := <-emitter
				mediator.receptionMutex.Lock()
				msg, err := mediator.client.ProcessPacket(packet)
				if msg.Valid {
					mediator.mailboxes.ForEach(func(_ int, mb MailBox) {
						mb.Notify(msg, mediator)
					})
				} else if err != nil {
					fmt.Println("Message error: " + err.Error())
				}
				mediator.receptionMutex.Unlock()
			}
		}()
	}
	return
}
