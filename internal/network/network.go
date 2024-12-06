package network

import (
	"log"
	"net"
	"os"
)

var networkLog = log.New(os.Stdout, "[Network] ", log.LstdFlags)

type Config struct {
	Address              string
	MaxConnections       int
	OnClientConnected    func(id int, conn *Conn)
	OnClientDisconnected func(id int, conn *Conn)
	OnDataReceived       func(id int, conn *Conn, bytes []byte)
}

type Network struct {
	config        *Config
	listen        net.Listener
	connectionIds []int
	connect       chan net.Conn
	disconnect    chan *Conn
}

func (net *Network) getAvailableConnectionId() int {
	if len(net.connectionIds) == 0 {
		return -1
	}

	id := net.connectionIds[0]

	net.connectionIds = net.connectionIds[1:]

	return id
}

func (network *Network) run() {
	defer func() {
		_ = network.listen.Close()

		close(network.connect)
		close(network.disconnect)

		networkLog.Print("Network subsystem has stopped")
	}()

	for {
		select {
		case conn := <-network.connect:
			connId := network.getAvailableConnectionId()
			if connId == -1 {
				_ = conn.Close()
				break
			}
			startConnection(network, connId, conn)

		case conn := <-network.disconnect:
			if conn.state == StateClosed {
				break
			}
			conn.state = StateClosed
			if conn.connId != -1 {
				network.config.OnClientDisconnected(conn.connId, conn)
				network.connectionIds = append(network.connectionIds, conn.connId)
			}
		}
	}
}

func Start(config Config) error {
	listen, err := net.Listen("tcp", config.Address)
	if err != nil {
		return err
	}

	network := Network{
		config:        &config,
		listen:        listen,
		connectionIds: make([]int, config.MaxConnections),
		connect:       make(chan net.Conn),
		disconnect:    make(chan *Conn),
	}

	for i := 0; i < config.MaxConnections; i++ {
		network.connectionIds[i] = i
	}

	networkLog.Println("Network subsystem has started on", config.Address)

	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Print(err)
				break
			}

			network.connect <- conn
		}
	}()

	go network.run()

	return nil
}
