package comm

import (
	"fmt"
	"sync"

	"github.com/albertoaer/venus/govenus/utils"
)

type registeredSender struct {
	distance uint32
	sender   Sender
}

func registeredSenderFromEvent(event ChannelEvent) *registeredSender {
	return &registeredSender{
		distance: event.Message.Distance,
		sender:   event.Sender,
	}
}

func (rs *registeredSender) tryReplaceWithEvent(event ChannelEvent) {
	if event.Message.Distance <= rs.distance {
		rs.distance = event.Message.Distance
		rs.sender = event.Sender
	}
}

func (rs *registeredSender) getSender() Sender {
	return rs.sender
}

type baseClient struct {
	id           string
	endpoint     bool
	senders      map[string]*registeredSender
	sendersMutex sync.RWMutex
	mailboxes    *utils.ConcurrentArray[Mailbox]
}

func NewClient(id string) Client {
	fmt.Println("Client created with id: " + id)
	return &baseClient{
		id:        id,
		endpoint:  true,
		senders:   make(map[string]*registeredSender),
		mailboxes: utils.NewArray[Mailbox](20),
	}
}

func NewRouterClient(id string) Client {
	fmt.Println("Router Client created with id: " + id)
	return &baseClient{
		id:        id,
		endpoint:  false,
		senders:   make(map[string]*registeredSender),
		mailboxes: utils.NewArray[Mailbox](20),
	}
}

func (client *baseClient) GetId() string {
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
			sender.getSender().Send(message)
			return
		}
	}
	if allowBroadcast {
		for id, sender := range client.senders {
			if id != message.Sender {
				sender.getSender().Send(message)
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
	if registered, exists := client.senders[event.Message.Sender]; !exists {
		client.senders[event.Message.Sender] = registeredSenderFromEvent(event)
	} else {
		registered.tryReplaceWithEvent(event)
	}
	event.Message.Distance += 1
	client.spreadMessage(event.Message, !client.endpoint)
}

func (client *baseClient) StartChannel(channel MessageChannel) (err error) {
	if err = channel.Start(); err == nil {
		go func() {
			emitter := channel.Emitter()
			for {
				received := <-emitter
				if received.Message.Sender == client.id {
					continue
				}
				client.onEvent(received)
				client.mailboxes.ForEach(func(_ int, mb Mailbox) {
					mb.Notify(received.Message, client)
				})
			}
		}()
	}
	return
}
