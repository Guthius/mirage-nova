package main

import (
	"github.com/guthius/mirage-nova/net"
	"github.com/guthius/mirage-nova/server/config"
)

type Level struct {
	Id      int
	Data    *LevelData
	Players []*PlayerData
}

var Levels [config.MaxMaps]Level

func init() {
	for i := 0; i < len(Levels); i++ {
		Levels[i] = Level{
			Id:      i + 1,
			Data:    &LevelData{},
			Players: make([]*PlayerData, config.MaxPlayers),
		}
	}
}

// Send a packet with the specified bytes to all players on the level
func (level *Level) Send(bytes []byte) {
	for _, p := range level.Players {
		p.Send(bytes)
	}
}

// AddPlayer adds the player to the level
func (level *Level) AddPlayer(p *PlayerData) {
	for _, other := range level.Players {
		if other == p {
			return
		}
		playerData := getPlayerDataPacket(other)
		p.Send(playerData)
	}

	level.Players = append(level.Players, p)

	playerData := getPlayerDataPacket(p)
	for _, other := range level.Players {
		other.Send(playerData)
	}
}

// RemovePlayer removes the specified Player from the level
func (level *Level) RemovePlayer(p *PlayerData) {
	for i := 0; i < len(level.Players); i++ {
		if level.Players[i] != p {
			continue
		}
		level.Players = append(level.Players[:i], level.Players[i+1:]...)
		return
	}

	writer := net.NewWriter()
	writer.WriteInteger(SvLeft)
	writer.WriteLong(p.Id + 1)

	// Notify the remaining players that the specified player has left
	for _, other := range level.Players {
		other.Send(writer.Bytes())
	}
}

func getPlayerDataPacket(pl *PlayerData) []byte {
	return nil
}
