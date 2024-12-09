package main

import (
	"encoding/binary"
	"unicode/utf16"

	"github.com/guthius/mirage-nova/net"
	"github.com/guthius/mirage-nova/server/color"
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
			LevelCache: buildLevelCache(i+1, levelData),
			TempTiles:  make([]TempTile, len(levelData.Tiles)),
			Players:    make([]*PlayerData, config.MaxPlayers),
			DoorTimer:  0,
		}

		rooms[i].resetTempTiles()
	}
}

func (r *Room) resetTempTiles() {
	for i := 0; i < len(r.TempTiles); i++ {
		r.TempTiles[i].DoorOpen = false
	}
}

// stringToUtf16 converts a string to a byte array of UTF-16 characters.
func stringToUtf16(s string, maxLen int) []byte {
	bytes := make([]byte, maxLen*2)

	codes := utf16.Encode([]rune(s))
	codesLen := len(codes)
	if maxLen > codesLen {
		maxLen = codesLen
	}

	for i := 0; i < maxLen; i++ {
		binary.LittleEndian.PutUint16(bytes[i*2:], codes[i])
	}

	return bytes
}

// buildLevelCache creates a byte array of the specified level data that can be sent to players.
func buildLevelCache(id int, l *data.LevelData) []byte {
	writer := net.NewWriter()

	writer.WriteInteger(SvLevelData)
	writer.WriteLong(id)
	writer.Write(stringToUtf16(l.Name, config.NameLength))
	writer.WriteLong(l.Revision)
	writer.WriteByte(byte(l.Type))
	writer.WriteInteger(l.TileSet)
	writer.WriteInteger(l.Up + 1)
	writer.WriteInteger(l.Down + 1)
	writer.WriteInteger(l.Left + 1)
	writer.WriteInteger(l.Right + 1)
	writer.WriteByte(byte(l.Music))
	writer.WriteInteger(l.BootMap + 1)
	writer.WriteByte(byte(l.BootX))
	writer.WriteByte(byte(l.BootY))
	writer.WriteByte(byte(l.Shop + 1))

	for i := 0; i < len(l.Tiles); i++ {
		for j := 0; j < len(l.Tiles[i].Num); j++ {
			writer.WriteInteger(l.Tiles[i].Num[j])
		}

		writer.WriteByte(byte(l.Tiles[i].Type))
		writer.WriteInteger(l.Tiles[i].Data1)
		writer.WriteInteger(l.Tiles[i].Data2)
		writer.WriteInteger(l.Tiles[i].Data3)
	}

	for i := 0; i < config.MaxMapNpcs; i++ {
		writer.WriteByte(byte(l.Npcs[i] + 1))
	}

	return writer.Bytes()
}

// Send a packet with the specified bytes to all players on the level.
func (level *Room) Send(bytes []byte) {
	for _, p := range level.Players {
		p.Send(bytes)
	}
}

// Contains returns true if the specified player is in the level; otherwise, returns false.
func (r *Room) Contains(p *PlayerData) bool {
	for _, player := range r.Players {
		if player == p {
			return true
		}
	}
	return false
}

// AddPlayer adds the player to the level
func (r *Room) AddPlayer(p *PlayerData) {
	// If the player is already in the room, return
	if p.Room == r {
		return
	}

	// If the player is already in a room, remove them from that room
	if p.Room != nil {
		p.Room.RemovePlayer(p)
	}

	// Send the player data of all players in the room to the new player
	for _, o := range r.Players {
		playerData := getPlayerDataPacket(o)
		p.Send(playerData)
	}

	r.Players = append(r.Players, p)

	p.TargetType = TargetNone
	p.Target = -1
	p.GettingLevel = true
	p.Room = r

	// If there is a shop in the room, say hello to the player
	shop := data.GetShop(r.LevelData.Shop)
	if shop != nil {
		if shop.JoinSay != "" {
			SendMessage(p, shop.Name+" says, '"+shop.JoinSay+"'", color.SayColor)
		}
	}

	// Send the player data to all players in the room including the new player
	playerData := getPlayerDataPacket(p)
	for _, o := range r.Players {
		o.Send(playerData)
	}

	r.sendDoorDataTo(p)

	// Request the player to check if they have the correct revision of the level
	SendCheckForLevel(p, r.Id)
}

func (r *Room) sendDoorDataTo(player *PlayerData) {
	width := r.LevelData.Width
	height := r.LevelData.Height
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			tid := y*width + x
			tile := &r.TempTiles[tid]

			if tile.DoorOpen {
				writer := net.NewWriter()
				writer.WriteInteger(SDoor)
				writer.WriteLong(x)
				writer.WriteLong(y)

				player.Send(writer.Bytes())
			}
		}
	}
}

// RemovePlayer removes the specified Player from the level
func (r *Room) RemovePlayer(p *PlayerData) {
	// Remove the player from the list of players in the room
	for i := 0; i < len(r.Players); i++ {
		if r.Players[i] != p {
			continue
		}
		r.Players = append(r.Players[:i], r.Players[i+1:]...)
		break
	}

	// If there is a shop in the room, say goodbye to the player
	shop := data.GetShop(r.LevelData.Shop)
	if shop != nil {
		if shop.LeaveSay != "" {
			SendMessage(p, shop.Name+" says, '"+shop.LeaveSay+"'", color.SayColor)
		}
	}

	writer := net.NewWriter()
	writer.WriteInteger(SvLeft)
	writer.WriteLong(p.Id + 1)

	// Notify the remaining players that the specified player has left
	for _, other := range r.Players {
		other.Send(writer.Bytes())
	}
}

func getPlayerDataPacket(p *PlayerData) []byte {
	char := p.Character
	if char == nil {
		return []byte{}
	}

	pk := 0
	if char.PK {
		pk = 1
	}

	writer := net.NewWriter()

	writer.WriteInteger(SPlayerData)
	writer.WriteLong(p.Id + 1)
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
