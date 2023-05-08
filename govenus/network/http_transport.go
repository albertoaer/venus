package network

import (
	"io"
	"net/http"

	"github.com/albertoaer/venus/govenus/comm"
)

type httpTransportHandler struct {
	channel    *HttpChannel
	serializer comm.MessageSerializer
}

func NewHttpTransportChannel() *HttpChannel {
	transport := &httpTransportHandler{
		serializer: comm.NewJsonSerializer(),
	}
	channel := newHttpChannel(transport)
	transport.channel = channel
	return channel
}

func (handler *httpTransportHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	if request.Method != http.MethodPost {
		writer.WriteHeader(400)
		writer.Write([]byte("{\"error\":\"only post allowed\"}"))
		return
	}
	data, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte("{\"error\":\"body message error\"}"))
		return
	}
	message, err := handler.serializer.Deserialize(data)
	if err != nil {
		writer.WriteHeader(400)
		writer.Write([]byte("{\"error\":\"body message error\"}"))
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
	reply := <-gms.messageReceiver
	response, err := handler.serializer.Serialize(reply)
	if err != nil {
		gms.errorReceiver <- err
		return
	}
	if _, err := writer.Write(response); err != nil {
		gms.errorReceiver <- err
		return
	}
	gms.errorReceiver <- nil
}
