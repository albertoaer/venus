package network

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/albertoaer/venus/govenus/comm"
)

type TcpChannel struct {
	port          int
	conn          *net.TCPListener
	emitter       chan comm.BinaryPacket[net.Addr]
	connections   map[string]*net.TCPConn
	connectionsRW sync.RWMutex
}

func NewTcpChannel() *TcpChannel {
	return &TcpChannel{
		port:          DefaultPort,
		emitter:       make(chan comm.BinaryPacket[net.Addr]),
		connections:   make(map[string]*net.TCPConn),
		connectionsRW: sync.RWMutex{},
	}
}

func (tcp *TcpChannel) AsMessageChannel() comm.OpenableChannel[net.Addr] {
	return comm.AdaptBinaryChannel[net.Addr](tcp)
}

func (tcp *TcpChannel) SetPort(port int) *TcpChannel {
	tcp.port = port
	return tcp
}

func (tcp *TcpChannel) Emitter() <-chan comm.BinaryPacket[net.Addr] {
	return tcp.emitter
}

func (tcp *TcpChannel) Start() (err error) {
	tcp.conn, err = net.ListenTCP("tcp", &net.TCPAddr{Port: tcp.port})
	if err == nil {
		fmt.Printf("Starting tcp server at port: %d\n", tcp.port)
		go func() {
			for {
				conn, err := tcp.conn.AcceptTCP()
				if err != nil {
					break
				}
				go tcp.handleConnection(conn)
			}
		}()
	}
	return
}

func (tcp *TcpChannel) Send(packet comm.BinaryPacket[net.Addr]) (err error) {
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

func (tcp *TcpChannel) handleConnection(conn *net.TCPConn) {
	tcp.connectionsRW.Lock()
	tcp.connections[conn.RemoteAddr().String()] = conn
	tcp.connectionsRW.Unlock()
	// TODO: maybe reduce the number of buffers
	buffer := make([]byte, NetBufferSize)
	for {
		size, err := conn.Read(buffer)
		if err != nil {
			break
		}
		fmt.Println("Got package of size", size, "from", conn.RemoteAddr())
		tcp.emitter <- comm.BinaryPacket[net.Addr]{
			Data:    buffer[:size],
			Address: conn.RemoteAddr(),
		}
	}
	tcp.connectionsRW.Lock()
	delete(tcp.connections, conn.RemoteAddr().String())
	tcp.connectionsRW.Unlock()
}
