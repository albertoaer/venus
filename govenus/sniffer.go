package govenus

import (
	"fmt"

	"github.com/albertoaer/venus/govenus/protocol"
)

type sniffer struct {
	action func(protocol.Message)
}

func NewSniffer() MailBox {
	return &sniffer{action: func(m protocol.Message) {
		fmt.Println(m)
	}}
}

func NewCustomSniffer(action func(protocol.Message)) MailBox {
	return &sniffer{action: action}
}

func (sniffer *sniffer) Notify(msg protocol.Message, _ protocol.Client, _ protocol.Sender) {
	sniffer.action(msg)
}
