package main

import (
	"encoding/binary"
	"fmt"
	"unicode/utf16"

	"github.com/guthius/mirage-nova/net"
	"github.com/guthius/mirage-nova/server/color"
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
			Players:    make([]*PlayerData, 0, config.MaxPlayers),
			DoorTimer:  0,
		}

		rooms[i].resetTempTiles()
	}

}

func (r *Room) resetTempTiles() {
	for i := 0; i < len(r.TempTiles); i++ {
		tile := &r.TempTiles[i]
		tile.Data = &r.LevelData.Tiles[i]
		tile.DoorOpen = false
		tile.DoorTimer = 0
	}
}

// stringToUtf16 converts a string to a byte array of UTF-16 characters.
func stringToUtf16(s string, maxLen int) []byte {
	const space uint16 = 0x20

	bytes := make([]byte, maxLen*2)

	codes := utf16.Encode([]rune(s))
	codesLen := len(codes)

	for i := 0; i < maxLen; i++ {
		if i < codesLen {
			binary.LittleEndian.PutUint16(bytes[i*2:], codes[i])
		} else {
			binary.LittleEndian.PutUint16(bytes[i*2:], space)
		}
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
	writer.WriteInteger(int(l.Type))
	writer.WriteInteger(l.TileSet)
	writer.WriteInteger(l.Up + 1)
	writer.WriteInteger(l.Down + 1)
	writer.WriteInteger(l.Left + 1)
	writer.WriteInteger(l.Right + 1)
	writer.WriteInteger(int(l.Music))
	writer.WriteInteger(l.BootMap + 1)
	writer.WriteByte(byte(l.BootX))
	writer.WriteByte(byte(l.BootY))
	writer.WriteInteger(int(l.Shop + 1))

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
func (r *Room) Send(bytes []byte) {
	for _, p := range r.Players {
		p.Send(bytes)
	}
}

// SendExclude sends a packet with the specified bytes to all players on the level except the specified player.
func (r *Room) SendExclude(bytes []byte, exclude *PlayerData) {
	for _, p := range r.Players {
		if p == exclude {
			continue
		}
		p.Send(bytes)
	}
}

// SendMessage sends a message to all players in the room.
func (r *Room) SendMessage(message string, color color.Color) {
	writer := net.NewWriter()

	writer.WriteInteger(SvRoomMessage)
	writer.WriteString(message)
	writer.WriteByte(byte(color))

	r.Send(writer.Bytes())
}

func (r *Room) SendPlayerData(p *PlayerData) {
	playerData := getPlayerDataPacket(p)
	for _, o := range r.Players {
		o.Send(playerData)
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

func (r *Room) AddPlayerAt(p *PlayerData, x int, y int) {
	p.Character.X = x
	p.Character.Y = y

	// If the player is already in the room just send the updated player data to all players in the room
	if p.Room == r {
		r.SendPlayerData(p)
		return
	}

	r.AddPlayer(p)

	TriggerTileEffect(p)
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
			SendMessage(p, fmt.Sprintf("%s says, '%s'", shop.Name, shop.JoinSay), color.SayColor)
		}
	}

	// Send the player data to all players in the room including the new player
	r.SendPlayerData(p)

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
			SendMessage(p, fmt.Sprintf("%s says, '%s'", shop.Name, shop.LeaveSay), color.SayColor)
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

	writer.WriteInteger(SvPlayerData)
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

// GetTile returns the tile at the specified position.
func (r *Room) GetTile(x int, y int) *TempTile {
	if !r.LevelData.Contains(x, y) {
		return nil
	}
	tid := y*r.LevelData.Width + x
	return &r.TempTiles[tid]
}
