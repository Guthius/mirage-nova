package main

import (
	"fmt"
	"log"
	"mirage/internal/database"
	"mirage/internal/network"
	"os"
	"time"
)

var IsShuttingDown = false

func HandleClientConnected(id int, conn *network.Conn) {
	log.Printf("[%d] Client connected from %s\n", id, conn.RemoteAddr())

	player := GetPlayer(id)
	player.Connection = conn
	player.Id = id

	if IsBanned(conn.RemoteAddr()) {
		player.SendAlert(fmt.Sprintf("Your have been banned from %s, and you are no longer able to play.", GameName))
	}
}

func HandleClientDisconnected(id int, conn *network.Conn) {
	log.Printf("[%d] Connection with %s has been terminated\n", id, conn.RemoteAddr())

	player := GetPlayer(id)
	if player.IsPlaying() {
		// TODO: Call LeftGame
	}

	player.Clear()
}

func HandleDataReceived(id int, _ *network.Conn, bytes []byte) {
	HandleData(GetPlayer(id), bytes)
}

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.SetOutput(os.Stdout)
}

func main() {
	database.Create()

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
