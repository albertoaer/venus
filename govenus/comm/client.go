package comm

import (
	"fmt"
	"sync"

	"github.com/albertoaer/venus/govenus/utils"
)

type baseClient struct {
	id           ClientId
	endpoint     bool
	senders      map[ClientId]Sender // TODO: allow multiple senders to prevent client replacement
	sendersMutex sync.RWMutex
	mailboxes    *utils.ConcurrentArray[Mailbox]
}

func NewClient(id ClientId) Client {
	fmt.Println("Client created with id: " + id)
	return &baseClient{
		id:        id,
		endpoint:  true,
		senders:   make(map[ClientId]Sender),
		mailboxes: utils.NewArray[Mailbox](20),
	}
}

func NewRouterClient(id ClientId) Client {
	fmt.Println("Router Client created with id: " + id)
	return &baseClient{
		id:        id,
		endpoint:  false,
		senders:   make(map[ClientId]Sender),
		mailboxes: utils.NewArray[Mailbox](20),
	}
}

func (client *baseClient) GetId() ClientId {
	return client.id
}

func (client *baseClient) Attach(mailbox Mailbox) {
	client.mailboxes.Add(mailbox)
}

func (client *baseClient) Detach(mailbox Mailbox) {
	client.mailboxes.Remove(mailbox)
}

func (client *baseClient) spreadMessage(message Message, allowBroadcast bool) {
	if message.Receiver != nil {
		if *message.Receiver == client.id {
			return
		}
		if sender, exists := client.senders[*message.Receiver]; exists {
			sender.Send(message)
			return
		}
	}
	if allowBroadcast {
		for id, sender := range client.senders {
			if id != message.Sender {
				sender.Send(message)
			}
		}
	}
}

func (client *baseClient) Send(message Message) error {
	client.sendersMutex.RLock()
	defer client.sendersMutex.RUnlock()
	client.spreadMessage(message, true)
	return nil
}

func (client *baseClient) onEvent(event ChannelEvent) {
	client.sendersMutex.Lock()
	defer client.sendersMutex.Unlock()
	if event.Message.Sender == client.id {
		return
	}
	if _, exists := client.senders[event.Message.Sender]; !exists {
		client.senders[event.Message.Sender] = event.Sender
	}
	client.spreadMessage(event.Message, !client.endpoint)
}

func (client *baseClient) StartChannel(channel MessageChannel) (err error) {
	if err = channel.Start(); err == nil {
		go func() {
			emitter := channel.Emitter()
			for {
				received := <-emitter
				if err != nil {
					fmt.Println("Message error: " + err.Error())
				} else {
					client.mailboxes.ForEach(func(_ int, mb Mailbox) {
						mb.Notify(received, client)
					})
				}
				client.onEvent(received)
			}
		}()
	}
	return
}
