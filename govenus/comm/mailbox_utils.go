package comm

import "fmt"

type sniffer struct {
	action func(Message)
}

func NewSniffer() Mailbox {
	return &sniffer{action: func(m Message) {
		fmt.Println(m)
	}}
}

func NewCustomSniffer(action func(Message)) Mailbox {
	return &sniffer{action: action}
}

func (sniffer *sniffer) Notify(event ChannelEvent, _ Client) {
	sniffer.action(event.Message)
}
