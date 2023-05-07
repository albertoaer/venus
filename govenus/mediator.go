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
	channels       *utils.ConcurrentArray[protocol.MessageChannel]
	receptionMutex sync.Mutex
}

func NewMediator(client protocol.Client) *Mediator {
	mediator := &Mediator{
		mailboxes:      utils.NewArray[MailBox](20),
		client:         client,
		channels:       utils.NewArray[protocol.MessageChannel](1),
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
	mediator.client.Send(message)
}

func (mediator *Mediator) GetChannel(index int) protocol.MessageChannel {
	return mediator.channels.Get(index)
}

func (mediator *Mediator) ChannelsCount() int {
	return mediator.channels.Length()
}

func (mediator *Mediator) StartChannel(channel protocol.MessageChannel) (err error) {
	if err = channel.Start(); err == nil {
		mediator.channels.Add(channel)
		go func() {
			for {
				emitter := channel.Emitter()
				received := <-emitter
				mediator.receptionMutex.Lock()
				err := mediator.client.GotMessage(received.Message, received.Sender)
				if err != nil {
					fmt.Println("Message error: " + err.Error())
				} else {
					mediator.mailboxes.ForEach(func(_ int, mb MailBox) {
						mb.Notify(received.Message, mediator.client, received.Sender)
					})
				}
				mediator.receptionMutex.Unlock()
			}
		}()
	}
	return
}
