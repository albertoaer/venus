package govenus

import "github.com/albertoaer/venus/govenus/protocol"

type MailBox interface {
	Notify(protocol.Message, protocol.Client, protocol.Sender)
}
