package network

import (
	"errors"
	"net"
	"strings"

	"github.com/albertoaer/venus/govenus/comm"
)

type UdpChannel struct {
	port    int
	conn    *net.UDPConn
	emitter chan comm.BinaryPacket[net.Addr]
}

func NewUdpChannel() *UdpChannel {
	return &UdpChannel{
		port:    DefaultPort,
		conn:    nil,
		emitter: make(chan comm.BinaryPacket[net.Addr]),
	}
}

func (udp *UdpChannel) AsMessageChannel() comm.OpenableChannel[net.Addr] {
	return comm.AdaptBinaryChannel[net.Addr](udp)
}

func (udp *UdpChannel) SetPort(port int) *UdpChannel {
	udp.port = port
	return udp
}

func (udp *UdpChannel) Emitter() <-chan comm.BinaryPacket[net.Addr] {
	return udp.emitter
}

func (udp *UdpChannel) Start() (err error) {
	udp.conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: udp.port})
	if err == nil {
		go func() {
			buffer := make([]byte, NetBufferSize)
			for {
				if size, addr, err := udp.conn.ReadFromUDP(buffer); err == nil {
					udp.emitter <- comm.BinaryPacket[net.Addr]{
						Data:    buffer[:size],
						Address: addr,
					}
				}
			}
		}()
	}
	return
}

func (udp *UdpChannel) Send(packet comm.BinaryPacket[net.Addr]) (err error) {
	if !strings.HasPrefix(packet.Address.Network(), "udp") {
		return errors.New("expecting udp address")
	}
	_, err = udp.conn.WriteTo(packet.Data, packet.Address)
	return
}
