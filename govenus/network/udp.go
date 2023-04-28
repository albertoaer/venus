package network

import (
	"errors"
	"net"
	"strings"

	"github.com/albertoaer/venus/govenus/protocol"
)

type UdpPackageChannel struct {
	port    int
	conn    *net.UDPConn
	emitter chan protocol.Packet
}

func NewUdpPackageChannel() *UdpPackageChannel {
	return &UdpPackageChannel{
		port:    DefaultPort,
		conn:    nil,
		emitter: make(chan protocol.Packet),
	}
}

func (udp *UdpPackageChannel) SetPort(port int) *UdpPackageChannel {
	udp.port = port
	return udp
}

func (udp *UdpPackageChannel) Emitter() <-chan protocol.Packet {
	return udp.emitter
}

func (udp *UdpPackageChannel) Start() (err error) {
	udp.conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: udp.port})
	if err == nil {
		go func() {
			buffer := make([]byte, NetBufferSize)
			for {
				if size, addr, err := udp.conn.ReadFromUDP(buffer); err == nil {
					udp.emitter <- protocol.Packet{Data: buffer[:size], Address: addr, Channel: udp}
				}
			}
		}()
	}
	return
}

func (udp *UdpPackageChannel) Send(packet protocol.Packet) (err error) {
	if !strings.HasPrefix(packet.Address.Network(), "udp") {
		return errors.New("expecting udp address")
	}
	_, err = udp.conn.WriteTo(packet.Data, packet.Address)
	return
}
