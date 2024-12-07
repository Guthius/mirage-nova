package main

import (
	"mirage/internal/database"
	"mirage/internal/packet"
)

type Level struct {
	Data    *database.Map
	Players []*Player
}

// Send a packet with the specified bytes to all players on the level
func (level *Level) Send(bytes []byte) {
	for i := 0; i < len(level.Players); i++ {
		player := level.Players[i]
		player.Send(bytes)
	}
}

// AddPlayer adds the specified Player to the level
func (level *Level) AddPlayer(player *Player) {
	for _, other := range level.Players {
		if other == player {
			return
		}

		playerData := getPlayerDataPacket(other)
		player.Send(playerData)
	}

	level.Players = append(level.Players, player)

	playerData := getPlayerDataPacket(player)
	for _, other := range level.Players {
		other.Send(playerData)
	}
}

// RemovePlayer removes the specified Player from the level
func (level *Level) RemovePlayer(player *Player) {
	for i := 0; i < len(level.Players); i++ {
		if level.Players[i] != player {
			continue
		}
		level.Players = append(level.Players[:i], level.Players[i+1:]...)
		return
	}

	writer := packet.NewWriter()
	writer.WriteInteger(SvLeft)
	writer.WriteLong(player.Id + 1)

	// Notify the remaining players that the specified player has left
	for _, other := range level.Players {
		other.Send(writer.Bytes())
	}
}

func getPlayerDataPacket(player *Player) []byte {
	return nil
}
