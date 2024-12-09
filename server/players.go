package main

import (
	"strings"

	"github.com/guthius/mirage-nova/net"
	"github.com/guthius/mirage-nova/server/character"
	"github.com/guthius/mirage-nova/server/config"
	"github.com/guthius/mirage-nova/server/data"
	"github.com/guthius/mirage-nova/server/data/vitals"
	"github.com/guthius/mirage-nova/server/user"
)

const (
	BufferSize = 4096
)

type TargetType int

const (
	TargetNone TargetType = iota
	TargetPlayer
	TargetNpc
)

type PlayerData struct {
	Id            int
	Connection    *net.Conn
	Account       *user.Account
	Buffer        []byte
	CharacterList [config.MaxChars]character.Character
	Character     *character.Character
	TargetType    TargetType
	Target        int
	GettingLevel  bool
	Room          *Room
}

var players [config.MaxPlayers]PlayerData

// GetPlayer returns the player at the specified index.
func GetPlayer(index int) *PlayerData {
	if index < 0 || index >= config.MaxPlayers {
		return nil
	}
	return &players[index]
}

// GetPlayersInGame returns a slice that contains all players that are currently in game.
func GetPlayersInGame() []*PlayerData {
	result := make([]*PlayerData, 0, config.MaxPlayers)
	for i := 0; i < config.MaxPlayers; i++ {
		if players[i].IsPlaying() {
			result = append(result, &players[i])
		}
	}
	return result
}

func (p *PlayerData) Clear() {
	p.Connection = nil
	p.Account = nil
	p.Buffer = make([]byte, 0, BufferSize)
	p.Character = nil
	p.TargetType = TargetNone
	p.Target = -1
	p.GettingLevel = false
	p.Room = nil

	for i := 0; i < config.MaxChars; i++ {
		p.CharacterList[i].Clear()
	}
}

// Send the specified bytes to the player.
func (p *PlayerData) Send(bytes []byte) {
	if p == nil || p.Connection == nil {
		return
	}

	size := len(bytes)
	if size == 0 {
		return
	}

	sizeBytes := []byte{byte(size), byte(size >> 8)}

	p.Connection.Send(sizeBytes)
	p.Connection.Send(bytes)
}

// Disconnect closes the connection with the player.
func (p *PlayerData) Disconnect() {
	if p == nil || p.Connection == nil {
		return
	}
	p.Connection.Close()
}

// IsConnected returns true if the player is currently connected to the server; otherwise, returns false.
func (p *PlayerData) IsConnected() bool {
	return p.Connection != nil && p.Connection.State() == net.StateOpen
}

// IsLoggedIn returns true if the player is currently logged into the server; otherwise, returns false.
func (p *PlayerData) IsLoggedIn() bool {
	return p.IsConnected() && p.Account != nil
}

// IsPlaying returns true if the player is currently in game; otherwise, returns false.
func (p *PlayerData) IsPlaying() bool {
	return p.IsLoggedIn() && p.Character != nil
}

// GetMaxVital returns the maximum value of the specified vital type.
func (p *PlayerData) GetMaxVital(vital vitals.Type) int {
	if p.Character == nil {
		return 0
	}

	switch vital {
	case vitals.HP:
		return data.GetClass(p.Character.Class).GetMaxVital(vital, p.Character.Stats.Strength)
	case vitals.MP:
		return data.GetClass(p.Character.Class).GetMaxVital(vital, p.Character.Stats.Magic)
	case vitals.SP:
		return data.GetClass(p.Character.Class).GetMaxVital(vital, p.Character.Stats.Speed)
	}

	return 0
}

// GetVital returns the current value of the specified vital type.
func (p *PlayerData) GetVital(vital vitals.Type) int {
	if p.Character == nil {
		return 0
	}

	switch vital {
	case vitals.HP:
		return p.Character.Vitals.HP
	case vitals.MP:
		return p.Character.Vitals.MP
	case vitals.SP:
		return p.Character.Vitals.SP
	}

	return 0
}

// WarpTo moves the player to the specified room and position.
func (p *PlayerData) WarpTo(room *Room, x, y int) {
	p.Room = room
	p.Character.X = x
	p.Character.Y = y
	p.Room.AddPlayer(p)
}

// IsAccountLoggedIn returns true if there is a player logged in with the specified account name; otherwise, returns false.
func IsAccountLoggedIn(accountName string) bool {
	for _, p := range players {
		if p.IsLoggedIn() && strings.EqualFold(p.Account.Name, accountName) {
			return true
		}
	}
	return false
}
