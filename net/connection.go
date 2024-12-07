package net

import (
	tcp "net"
)

type ConnState int

const (
	StateClosed ConnState = iota
	StateClosing
	StateOpen
)

type Conn struct {
	network    *Network
	config     *Config
	connId     int
	conn       tcp.Conn
	state      ConnState
	send       chan []byte
	remoteAddr string
}

func getRemoteAddr(conn tcp.Conn) string {
	hostAndPort := conn.RemoteAddr().String()

	host, _, err := tcp.SplitHostPort(hostAndPort)
	if err != nil {
		return "unknown"
	}

	return host
}

func startConnection(network *Network, connId int, conn tcp.Conn) {
	connection := &Conn{
		network:    network,
		config:     network.config,
		connId:     connId,
		conn:       conn,
		state:      StateOpen,
		send:       make(chan []byte),
		remoteAddr: getRemoteAddr(conn),
	}

	go connection.doReceive()
	go connection.doSend()

	connection.config.OnClientConnected(connId, connection)
}

func (conn *Conn) doReceive() {
	defer func() {
		_ = conn.conn.Close()
		conn.network.disconnect <- conn
	}()

	buf := make([]byte, 4096)

	for {
		bytesReceived, err := conn.conn.Read(buf)
		if err != nil {
			return
		}

		conn.config.OnDataReceived(conn.connId, conn, buf[:bytesReceived])
	}
}

func (conn *Conn) doSend() {
	defer func() {
		_ = conn.conn.Close()
	}()

	for packet := range conn.send {
		for len(packet) > 0 {
			sent, e := conn.conn.Write(packet)
			if e != nil {
				return
			}

			packet = packet[sent:]
		}
	}
}

func (conn *Conn) Send(bytes []byte) {
	if conn.state != StateOpen {
		return
	}
	conn.send <- bytes
}

func (conn *Conn) Id() int { return conn.connId }

func (conn *Conn) State() ConnState { return conn.state }

func (conn *Conn) RemoteAddr() string {
	return conn.remoteAddr
}

func (conn *Conn) Close() {
	if conn.state == StateClosed {
		return
	}
	close(conn.send)
	conn.state = StateClosing
}
