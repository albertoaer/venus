package network

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/albertoaer/venus/govenus/comm"
	"github.com/albertoaer/venus/govenus/utils"
)

type httpAppHandler struct {
	channel     *HttpChannel
	idGenerator utils.IdGenerator
}

func NewHttpAppChannel() *HttpChannel {
	transport := &httpAppHandler{
		idGenerator: utils.NewUlidIdGenerator(),
	}
	channel := newHttpChannel(transport)
	transport.channel = channel
	return channel
}

func (handler *httpAppHandler) prepareMessage(request *http.Request) (msg comm.Message, err error) {
	msg.Sender = comm.ClientId(handler.idGenerator.NextId())
	msg.Timestamp = time.Now().Unix()
	msg.Verb = comm.Verb(request.Method)
	msg.Args = []string{request.RequestURI}
	msg.Payload, err = io.ReadAll(request.Body)
	if err != nil {
		return
	}
	msg.Options = make(map[string]string, len(request.Header))
	for k, v := range request.Header {
		msg.Options[k] = strings.Join(v, ", ")
	}
	return
}

func (handler *httpAppHandler) awaitAndReply(writer http.ResponseWriter, gms genericMessageSender) {
	message := <-gms.messageReceiver
	writer.Header().Set("content-type", "text/plain")
	for k, v := range message.Options {
		writer.Header().Set(k, v)
	}
	if status, err := strconv.Atoi(string(message.Verb)); err == nil {
		writer.WriteHeader(status)
	} else {
		writer.WriteHeader(200)
	}
	if _, err := writer.Write(message.Payload); err != nil {
		fmt.Println("Got error", err)
		gms.errorReceiver <- err
	} else {
		gms.errorReceiver <- nil
	}
}

func (handler *httpAppHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	message, err := handler.prepareMessage(request)
	if err != nil {
		return
	}

	gms := genericMessageSender{
		messageReceiver: make(chan comm.Message),
		errorReceiver:   make(chan error),
	}

	handler.channel.emitter <- comm.ChannelEvent{
		Message: message,
		Sender:  &gms,
	}

	handler.awaitAndReply(writer, gms)
}
