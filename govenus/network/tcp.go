package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
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
		err = tcpSend(conn, packet.Data)
	} else {
		var addr *net.TCPAddr
		if addr, err = net.ResolveTCPAddr(packet.Address.Network(), packet.Address.String()); err != nil {
			return
		}
		if conn, err = net.DialTCP(packet.Address.Network(), nil, addr); err != nil {
			return
		}
		if err = tcpSend(conn, packet.Data); err == nil {
			go tcp.handleConnection(conn)
		} else {
			conn.Close()
		}
	}
	return
}

func tcpSend(writer io.Writer, data []byte) error {
	if err := binary.Write(writer, binary.LittleEndian, uint64(len(data))); err != nil {
		return err
	}
	_, err := writer.Write(data)
	return err
}

func tcpRead(reader io.Reader) (buffer []byte, err error) {
	var size uint64
	if err = binary.Read(reader, binary.LittleEndian, &size); err != nil {
		return nil, err
	}
	buffer = make([]byte, size)
	var sz int
	sz, err = reader.Read(buffer)
	return buffer[:sz], err
}

func (tcp *TcpChannel) handleConnection(conn *net.TCPConn) {
	tcp.connectionsRW.Lock()
	tcp.connections[conn.RemoteAddr().String()] = conn
	tcp.connectionsRW.Unlock()
	for {
		if buffer, err := tcpRead(conn); err == nil {
			fmt.Println("Got package of size", len(buffer), "from", conn.RemoteAddr())
			tcp.emitter <- comm.BinaryPacket[net.Addr]{
				Data:    buffer,
				Address: conn.RemoteAddr(),
			}
		} else {
			break
		}
	}
	tcp.connectionsRW.Lock()
	delete(tcp.connections, conn.RemoteAddr().String())
	tcp.connectionsRW.Unlock()
}
