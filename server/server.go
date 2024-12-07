package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/guthius/mirage-nova/internal/database"
	_ "github.com/guthius/mirage-nova/internal/logger"
	"github.com/guthius/mirage-nova/net"
	"github.com/guthius/mirage-nova/server/config"
)

var IsShuttingDown = false
var Motd = ""
var PlayersOnline = 0

func HandleClientConnected(id int, conn *net.Conn) {
	log.Printf("[%d] Client connected from %s\n", id, conn.RemoteAddr())

	p := Get(id)
	p.Connection = conn
	p.Id = id

	if IsBanned(conn.RemoteAddr()) {
		SendAlert(p, fmt.Sprintf("You have been banned from %s, and you are no longer able to play.", config.GameName))
	}
}

func HandleClientDisconnected(id int, conn *net.Conn) {
	log.Printf("[%d] Connection with %s has been terminated\n", id, conn.RemoteAddr())

	pl := Get(id)
	if pl.IsPlaying() {
		// TODO: Call LeftGame
	}

	pl.Clear()
}

func HandleDataReceived(id int, _ *net.Conn, bytes []byte) {
	HandleData(Get(id), bytes)
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
	database.Create()

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
