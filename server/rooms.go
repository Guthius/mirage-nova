package main

import (
	"github.com/guthius/mirage-nova/net"
	"github.com/guthius/mirage-nova/server/config"
	"github.com/guthius/mirage-nova/server/data"
)

type TempTile struct {
	DoorOpen bool
}

type Room struct {
	Id         int
	LevelData  *data.LevelData
	LevelCache []byte
	TempTiles  []TempTile
	Players    []*PlayerData
	DoorTimer  int64
}

var rooms [config.MaxMaps]Room

func init() {
	for i := 0; i < len(rooms); i++ {
		levelData := data.GetLevel(i)

		rooms[i] = Room{
			Id:         i + 1,
			LevelData:  levelData,
			LevelCache: buildLevelCache(levelData),
			TempTiles:  make([]TempTile, len(levelData.Tiles)),
			Players:    make([]*PlayerData, config.MaxPlayers),
			DoorTimer:  0,
		}
	}
}

// buildLevelCache creates a byte array of the specified level data that can be sent to players.
func buildLevelCache(_ *data.LevelData) []byte {
	// TODO: Implement me
	return nil
}

// Send a packet with the specified bytes to all players on the level
func (level *Room) Send(bytes []byte) {
	for _, p := range level.Players {
		p.Send(bytes)
	}
}

// AddPlayer adds the player to the level
func (level *Room) AddPlayer(p *PlayerData) {
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
func (level *Room) RemovePlayer(p *PlayerData) {
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

func getPlayerDataPacket(_ *PlayerData) []byte {
	return nil
}
