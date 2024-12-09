package main

import (
	"fmt"

	"github.com/guthius/mirage-nova/net"
	"github.com/guthius/mirage-nova/server/color"
	"github.com/guthius/mirage-nova/server/compat"
	"github.com/guthius/mirage-nova/server/config"
	"github.com/guthius/mirage-nova/server/data"
)

type TempTile struct {
	Data      *data.Tile
	DoorOpen  bool
	DoorTimer int64
}

type Room struct {
	Id         int
	Level      *data.LevelData
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
			Level:      levelData,
			LevelCache: buildLevelCache(i+1, levelData),
			TempTiles:  make([]TempTile, len(levelData.Tiles)),
			Players:    make([]*PlayerData, 0, config.MaxPlayers),
			DoorTimer:  0,
		}

		rooms[i].resetTempTiles()
	}
}

func (room *Room) resetTempTiles() {
	for i := 0; i < len(room.TempTiles); i++ {
		tile := &room.TempTiles[i]
		tile.Data = &room.Level.Tiles[i]
		tile.DoorOpen = false
		tile.DoorTimer = 0
	}
}

// buildLevelCache creates a byte array of the specified level data that can be sent to players.
func buildLevelCache(id int, l *data.LevelData) []byte {
	writer := net.NewWriter()

	writer.WriteInteger(SvLevelData)
	writer.WriteLong(id)
	writer.Write(compat.StringToUtf16(l.Name, config.NameLength))
	writer.WriteLong(l.Revision)
	writer.WriteInteger(int(l.Type))
	writer.WriteInteger(l.TileSet)
	writer.WriteInteger(l.Up + 1)
	writer.WriteInteger(l.Down + 1)
	writer.WriteInteger(l.Left + 1)
	writer.WriteInteger(l.Right + 1)
	writer.WriteInteger(l.Music)
	writer.WriteInteger(l.BootMap + 1)
	writer.WriteByte(byte(l.BootX))
	writer.WriteByte(byte(l.BootY))
	writer.WriteInteger(l.Shop + 1)

	for i := 0; i < len(l.Tiles); i++ {
		for j := 0; j < len(l.Tiles[i].Num); j++ {
			writer.WriteInteger(l.Tiles[i].Num[j])
		}

		writer.WriteInteger(int(l.Tiles[i].Type))
		writer.WriteInteger(l.Tiles[i].Data1)
		writer.WriteInteger(l.Tiles[i].Data2)
		writer.WriteInteger(l.Tiles[i].Data3)
	}

	for i := 0; i < config.MaxMapNpcs; i++ {
		writer.WriteByte(byte(l.Npcs[i] + 1))
	}

	writer.WriteByte(0)
	writer.WriteByte(0)
	writer.WriteByte(0)

	return writer.Bytes()
}

// Send a packet with the specified bytes to all players on the level.
func (room *Room) Send(bytes []byte) {
	for _, p := range room.Players {
		p.Send(bytes)
	}
}

// SendExclude sends a packet with the specified bytes to all players on the level except the specified player.
func (room *Room) SendExclude(bytes []byte, exclude *PlayerData) {
	for _, p := range room.Players {
		if p == exclude {
			continue
		}
		p.Send(bytes)
	}
}

// SendMessage sends a message to all players in the room.
func (room *Room) SendMessage(message string, color color.Color) {
	writer := net.NewWriter()

	writer.WriteInteger(SvRoomMessage)
	writer.WriteString(message)
	writer.WriteByte(byte(color))

	room.Send(writer.Bytes())
}

// SendPlayerData sends the player data of the specified player to all players in the room.
func (room *Room) SendPlayerData(player *PlayerData) {
	playerData := getPlayerDataPacket(player)
	for _, p := range room.Players {
		p.Send(playerData)
	}
}

// Contains returns true if the specified player is in the level; otherwise, returns false.
func (room *Room) Contains(player *PlayerData) bool {
	for _, p := range room.Players {
		if player == p {
			return true
		}
	}
	return false
}

// AddPlayerAt adds the player to the level at the specified position.
func (room *Room) AddPlayerAt(player *PlayerData, x int, y int) {
	if player.Character == nil {
		return
	}

	player.Character.X = x
	player.Character.Y = y

	// If the player is already in the room just send the updated player data to all players in the room
	if player.Room == room {
		TriggerTileEffect(player)

		room.SendPlayerData(player)
		return
	}

	room.AddPlayer(player)
}

// AddPlayer adds the player to the level
func (room *Room) AddPlayer(player *PlayerData) {
	// If the player is already in the room, return
	if player.Room == room {
		return
	}

	// If the player is already in a room, remove them from that room
	if player.Room != nil {
		player.Room.RemovePlayer(player)
	}

	// Send the player data of all players in the room to the new player
	for _, p := range room.Players {
		playerData := getPlayerDataPacket(p)
		player.Send(playerData)
	}

	room.Players = append(room.Players, player)

	player.TargetType = TargetNone
	player.Target = -1
	player.GettingLevel = true
	player.Room = room

	// If there is a shop in the room, say hello to the player
	shop := data.GetShop(room.Level.Shop)
	if shop != nil {
		if shop.JoinSay != "" {
			SendMessage(player, fmt.Sprintf("%s says, '%s'", shop.Name, shop.JoinSay), color.SayColor)
		}
	}

	// Send the player data to all players in the room including the new player
	room.SendPlayerData(player)

	SendDoorData(player)
	SendCheckForLevel(player, room.Id)

	TriggerTileEffect(player)
}

// RemovePlayer removes the specified Player from the level
func (room *Room) RemovePlayer(player *PlayerData) {
	// Remove the player from the list of players in the room
	for i := 0; i < len(room.Players); i++ {
		if room.Players[i] != player {
			continue
		}
		room.Players = append(room.Players[:i], room.Players[i+1:]...)
		break
	}

	// If there is a shop in the room, say goodbye to the player
	shop := data.GetShop(room.Level.Shop)
	if shop != nil {
		if shop.LeaveSay != "" {
			SendMessage(player, fmt.Sprintf("%s says, '%s'", shop.Name, shop.LeaveSay), color.SayColor)
		}
	}

	writer := net.NewWriter()
	writer.WriteInteger(SvLeft)
	writer.WriteLong(player.Id + 1)

	// Notify the remaining players that the specified player has left
	for _, p := range room.Players {
		p.Send(writer.Bytes())
	}
}

func getPlayerDataPacket(player *PlayerData) []byte {
	char := player.Character
	if char == nil {
		return []byte{}
	}

	pk := 0
	if char.PK {
		pk = 1
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvPlayerData)
	writer.WriteLong(player.Id + 1)
	writer.WriteString(char.Name)
	writer.WriteLong(char.Sprite)
	writer.WriteLong(char.Room + 1)
	writer.WriteLong(char.X)
	writer.WriteLong(char.Y)
	writer.WriteString(char.Guild)
	writer.WriteLong(char.GuildAccess)
	writer.WriteLong(int(char.Dir))
	writer.WriteLong(int(char.Access))
	writer.WriteLong(pk)

	return writer.Bytes()
}

// GetTile returns the tile at the specified position.
func (room *Room) GetTile(x int, y int) *TempTile {
	if !room.Level.Contains(x, y) {
		return nil
	}
	tid := y*room.Level.Width + x
	return &room.TempTiles[tid]
}
