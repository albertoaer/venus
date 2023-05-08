package govenus

import (
	"net"

	"github.com/albertoaer/venus/govenus/comm"
	"github.com/albertoaer/venus/govenus/network"
	"github.com/albertoaer/venus/govenus/utils"
)

func SayHi[K any](client comm.Client, address K, channel comm.OpenableChannel[K]) error {
	sender, err := channel.Open(address)
	if err != nil {
		return err
	}
	messageBuilder := comm.NewMessageBuilder(client.GetId())
	messageBuilder.SetVerb(comm.Hi)
	_, err = sender.Send(messageBuilder.Build())
	return err
}

func SetupTcpClient(port int) (comm.Client, comm.OpenableChannel[net.Addr]) {
	tcpChannel := network.NewTcpChannel().SetPort(port).AsMessageChannel()
	client := comm.NewClient(comm.ClientId(utils.NewUlidIdGenerator().NextId()))
	if err := client.StartChannel(tcpChannel); err != nil {
		panic(err)
	}
	return client, tcpChannel
}
