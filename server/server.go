package main

import (
	"fmt"
	"log"
	"mirage/internal/network"
	"os"
	"time"
)

var ServerLog = log.New(os.Stdout, "[Server] ", log.LstdFlags)
var PlayerLog = log.New(os.Stdout, "[Player] ", log.LstdFlags)
var IsShuttingDown = false

func HandleClientConnected(id int, conn *network.Conn) {
	ServerLog.Printf("[%d] Client connected from %s\n", id, conn.RemoteAddr())

	player := GetPlayer(id)
	player.Connection = conn

	if IsBanned(conn.RemoteAddr()) {
		player.SendAlert(fmt.Sprintf("Your have been banned from %s, and you are no longer able to play.", GameName))
	}
}

func HandleClientDisconnected(id int, conn *network.Conn) {
	ServerLog.Printf("[%d] Connection with %s has been terminated\n", id, conn.RemoteAddr())

	player := GetPlayer(id)
	if player.IsPlaying() {
		// TODO: Call LeftGame
	}

	player.Clear()
}

func HandleDataReceived(id int, _ *network.Conn, bytes []byte) {
	HandleData(GetPlayer(id), bytes)
}

func main() {
	networkConfig := network.Config{
		Address:              GameAddr,
		MaxConnections:       MaxPlayers,
		OnClientConnected:    HandleClientConnected,
		OnClientDisconnected: HandleClientDisconnected,
		OnDataReceived:       HandleDataReceived,
	}

	err := network.Start(networkConfig)
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
