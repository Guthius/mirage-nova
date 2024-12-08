package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/guthius/mirage-nova/net"
	"github.com/guthius/mirage-nova/server/config"

	_ "github.com/guthius/mirage-nova/server/internal/logger"
)

var IsShuttingDown = false
var Motd = ""
var PlayersOnline = 0

func HandleClientConnected(id int, conn *net.Conn) {
	log.Printf("[%d] Client connected from %s\n", id, conn.RemoteAddr())

	player := GetPlayer(id)
	player.Connection = conn
	player.Id = id

	if IsBanned(conn.RemoteAddr()) {
		SendAlert(player, fmt.Sprintf("You have been banned from %s, and you are no longer able to play.", config.GameName))
	}
}

func HandleClientDisconnected(id int, conn *net.Conn) {
	log.Printf("[%d] Connection with %s has been terminated\n", id, conn.RemoteAddr())

	player := GetPlayer(id)
	if player.IsPlaying() {
		// TODO: Call LeftGame
	}

	player.Clear()
}

func HandleDataReceived(id int, _ *net.Conn, bytes []byte) {
	const headerSize = 2

	player := GetPlayer(id)

	player.Buffer = append(player.Buffer, bytes...)
	if len(player.Buffer) < headerSize {
		return
	}

	buf := player.Buffer
	off := 0

	// Handle all packets in the buffer
	for len(buf) >= headerSize {
		size := int(binary.LittleEndian.Uint16(buf))
		if len(buf) < size+headerSize {
			return
		}
		off += headerSize
		buf = buf[headerSize:]

		reader := net.NewReader(buf[:size])
		HandlePacket(player, reader)

		off += size
		buf = buf[size:]
	}

	// Move the bytes that are remaining to the front of the buffer
	bytesLeft := len(player.Buffer) - off
	if bytesLeft > 0 {
		copy(player.Buffer, player.Buffer[off:])
	}

	player.Buffer = player.Buffer[:bytesLeft]
}

func LoadMotd() {
	file, err := os.Open("motd.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Printf("error loading motd (%s)", err)
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("error loading motd (%s)", err)
		return
	}

	Motd = string(bytes)
}

func main() {
	networkConfig := net.Config{
		Address:              config.GameAddr,
		MaxConnections:       config.MaxPlayers,
		OnClientConnected:    HandleClientConnected,
		OnClientDisconnected: HandleClientDisconnected,
		OnDataReceived:       HandleDataReceived,
	}

	LoadMotd()

	err := net.Start(networkConfig)
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
