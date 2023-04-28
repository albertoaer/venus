package govenus

import "github.com/albertoaer/venus/govenus/protocol"

type ClientService interface {
	GetId() protocol.ClientId
	Send(protocol.Message)
}

type MailBox interface {
	Notify(protocol.Message, ClientService)
}
