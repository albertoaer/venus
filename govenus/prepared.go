package govenus

import (
	"net"

	"github.com/albertoaer/venus/govenus/network"
	"github.com/albertoaer/venus/govenus/protocol"
	"github.com/albertoaer/venus/govenus/utils"
)

func SayHi[K any](client protocol.Client, address K, channel protocol.OpenableChannel[K]) error {
	sender, err := channel.Open(address)
	if err != nil {
		return err
	}
	messageBuilder := protocol.NewMessageBuilder(client.GetId())
	messageBuilder.SetVerb(protocol.Hi)
	sender.Send(messageBuilder.Build())
	return nil
}

func SetupTcpClient(port int) (protocol.Client, protocol.OpenableChannel[net.Addr]) {
	tcpChannel := network.NewTcpChannel().SetPort(port).AsMessageChannel()
	client := protocol.NewClient(protocol.ClientId(utils.NewUlidIdGenerator().NextId()))
	if err := client.StartChannel(tcpChannel); err != nil {
		panic(err)
	}
	return client, tcpChannel
}
