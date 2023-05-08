package network

import (
	"net/http"
	"strconv"

	"github.com/albertoaer/venus/govenus/protocol"
)

type HttpChannel struct {
	port    int
	emitter chan protocol.ChannelEvent
	handler http.Handler
}

func newHttpChannel(handler http.Handler) *HttpChannel {
	channel := &HttpChannel{
		port:    80,
		emitter: make(chan protocol.ChannelEvent),
		handler: handler,
	}
	channel.handler = handler
	return channel
}

func (httpChannel *HttpChannel) SetPort(port int) *HttpChannel {
	httpChannel.port = port
	return httpChannel
}

func (httpChannel *HttpChannel) Emitter() <-chan protocol.ChannelEvent {
	return httpChannel.emitter
}

func (httpChannel *HttpChannel) Start() error {
	server := http.Server{
		Addr:    ":" + strconv.Itoa(httpChannel.port),
		Handler: httpChannel.handler,
	}
	go server.ListenAndServe()
	return nil
}

type genericMessageSender struct {
	messageReceiver chan protocol.Message
	errorReceiver   chan error
}

func (gms *genericMessageSender) Send(message protocol.Message) (bool, error) {
	gms.messageReceiver <- message
	return true, <-gms.errorReceiver
}
