package protocol

import "fmt"

type baseClient struct {
	id       ClientId
	endpoint bool
	senders  map[ClientId]Sender // TODO: allow multiple senders to prevent client replacement
}

func NewClient(id ClientId) Client {
	fmt.Println("Client created with id: " + id)
	return &baseClient{
		id:       id,
		endpoint: true,
		senders:  make(map[ClientId]Sender),
	}
}

func NewRouterClient(id ClientId) Client {
	fmt.Println("Router Client created with id: " + id)
	return &baseClient{
		id:       id,
		endpoint: false,
		senders:  make(map[ClientId]Sender),
	}
}

func (client *baseClient) GetId() ClientId {
	return client.id
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

func (client *baseClient) GotMessage(message Message, sender Sender) error {
	if message.Sender != client.id {
		if _, exists := client.senders[message.Sender]; !exists {
			client.senders[message.Sender] = sender
		}
		client.spreadMessage(message, !client.endpoint)
	}
	return nil
}

func (client *baseClient) Send(message Message) error {
	client.spreadMessage(message, true)
	return nil
}
