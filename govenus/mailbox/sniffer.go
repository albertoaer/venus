package mailbox

import (
	"fmt"

	"github.com/albertoaer/venus/govenus/protocol"
)

type sniffer struct {
	action func(protocol.Message)
}

func NewSniffer() protocol.Mailbox {
	return &sniffer{action: func(m protocol.Message) {
		fmt.Println(m)
	}}
}

func NewCustomSniffer(action func(protocol.Message)) protocol.Mailbox {
	return &sniffer{action: action}
}

func (sniffer *sniffer) Notify(event protocol.ChannelEvent, _ protocol.Client) {
	sniffer.action(event.Message)
}
