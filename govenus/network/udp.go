package network

import (
	"errors"
	"net"
	"strings"

	"github.com/albertoaer/venus/govenus/protocol"
)

type UdpPackageProvider struct {
	port    int
	conn    *net.UDPConn
	emitter chan protocol.Packet
}

func NewUdpPackageProvider() *UdpPackageProvider {
	return &UdpPackageProvider{
		port:    DefaultPort,
		conn:    nil,
		emitter: make(chan protocol.Packet),
	}
}

func (udp *UdpPackageProvider) SetPort(port int) *UdpPackageProvider {
	udp.port = port
	return udp
}

func (udp *UdpPackageProvider) Emitter() <-chan protocol.Packet {
	return udp.emitter
}

func (udp *UdpPackageProvider) Start() (err error) {
	udp.conn, err = net.ListenUDP("udp", &net.UDPAddr{Port: udp.port})
	if err == nil {
		go func() {
			buffer := make([]byte, NetBufferSize)
			for {
				if size, addr, err := udp.conn.ReadFromUDP(buffer); err == nil {
					udp.emitter <- protocol.Packet{Data: buffer[:size], Address: addr, Provider: udp}
				}
			}
		}()
	}
	return
}

func (udp *UdpPackageProvider) Send(packet protocol.Packet) (err error) {
	if !strings.HasPrefix(packet.Address.Network(), "udp") {
		return errors.New("expecting udp address")
	}
	_, err = udp.conn.WriteTo(packet.Data, packet.Address)
	return
}
