package main

import (
	"fmt"
	"io"
	"log"
	"mirage/internal/database"
	_ "mirage/internal/logger"
	"mirage/internal/network"
	"os"
	"time"
)

var IsShuttingDown = false
var Motd = ""
var PlayersOnline = 0

func HandleClientConnected(id int, conn *network.Conn) {
	log.Printf("[%d] Client connected from %s\n", id, conn.RemoteAddr())

	player := GetPlayer(id)
	player.Connection = conn
	player.Id = id

	if IsBanned(conn.RemoteAddr()) {
		player.SendAlert(fmt.Sprintf("You have been banned from %s, and you are no longer able to play.", GameName))
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

	networkConfig := network.Config{
		Address:              GameAddr,
		MaxConnections:       MaxPlayers,
		OnClientConnected:    HandleClientConnected,
		OnClientDisconnected: HandleClientDisconnected,
		OnDataReceived:       HandleDataReceived,
	}

	LoadMotd()

	err := network.Start(networkConfig)
	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
