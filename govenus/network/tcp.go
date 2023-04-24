package network

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/albertoaer/venus/govenus/protocol"
)

type TcpPackageProvider struct {
	port          int
	conn          *net.TCPListener
	emitter       chan protocol.Packet
	connections   map[string]*net.TCPConn
	connectionsRW sync.RWMutex
}

func NewTcpPackageProvider() *TcpPackageProvider {
	return &TcpPackageProvider{
		port:          DefaultPort,
		emitter:       make(chan protocol.Packet),
		connections:   make(map[string]*net.TCPConn),
		connectionsRW: sync.RWMutex{},
	}
}

func (tcp *TcpPackageProvider) SetPort(port int) *TcpPackageProvider {
	tcp.port = port
	return tcp
}

func (tcp *TcpPackageProvider) Emitter() <-chan protocol.Packet {
	return tcp.emitter
}

func (tcp *TcpPackageProvider) Start() (err error) {
	tcp.conn, err = net.ListenTCP("tcp", &net.TCPAddr{Port: tcp.port})
	if err == nil {
		fmt.Printf("Starting tcp server at port: %d\n", tcp.port)
		go func() {
			for {
				if conn, err := tcp.conn.AcceptTCP(); err == nil {
					go tcp.handleConnection(conn)
				}
			}
		}()
	}
	return
}

func (tcp *TcpPackageProvider) Send(packet protocol.Packet) (err error) {
	if !strings.HasPrefix(packet.Address.Network(), "tcp") {
		return errors.New("expecting tcp address")
	}
	tcp.connectionsRW.RLock()
	conn, exists := tcp.connections[packet.Address.String()]
	tcp.connectionsRW.RUnlock()
	if exists {
		_, err = conn.Write(packet.Data)
	} else {
		var addr *net.TCPAddr
		if addr, err = net.ResolveTCPAddr(packet.Address.Network(), packet.Address.String()); err != nil {
			return
		}
		if conn, err = net.DialTCP(packet.Address.Network(), nil, addr); err != nil {
			return
		}
		if _, err = conn.Write(packet.Data); err == nil {
			go tcp.handleConnection(conn)
		} else {
			conn.Close()
		}
	}
	return
}

func (tcp *TcpPackageProvider) handleConnection(conn *net.TCPConn) {
	tcp.connectionsRW.Lock()
	tcp.connections[conn.RemoteAddr().String()] = conn
	tcp.connectionsRW.Unlock()
	// TODO: maybe reduce the number of buffers
	buffer := make([]byte, NetBufferSize)
	for {
		size, err := conn.Read(buffer)
		if err != nil {
			tcp.connectionsRW.Lock()
			delete(tcp.connections, conn.RemoteAddr().String())
			tcp.connectionsRW.Unlock()
			break
		}
		fmt.Printf("Got package of size %d\n", size)
		tcp.emitter <- protocol.Packet{
			Data:     buffer[:size],
			Address:  conn.RemoteAddr(),
			Provider: tcp,
		}
	}
}
