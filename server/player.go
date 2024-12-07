package main

import (
	"mirage/internal/database"
	"mirage/internal/network"
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

type Player struct {
	Id            int
	Connection    *network.Conn
	Account       *database.Account
	Buffer        []byte
	CharacterList [database.MaxChars]database.Character
	Character     *database.Character
	TargetType    TargetType
	Target        int
	GettingMap    bool
}

var Players [MaxPlayers]Player

func GetPlayer(index int) *Player {
	if index < 0 || index >= MaxPlayers {
		return nil
	}
	return &Players[index]
}

func (player *Player) Clear() {
	player.Connection = nil
	player.Account = nil
	player.Buffer = make([]byte, 0, BufferSize)
	player.Character = nil
	player.TargetType = TargetNone
	player.Target = -1
	player.GettingMap = false

	for i := 0; i < database.MaxChars; i++ {
		player.CharacterList[i].Clear()
	}
}

func (player *Player) Send(bytes []byte) {
	if player == nil || player.Connection == nil {
		return
	}

	size := len(bytes)
	if size == 0 {
		return
	}

	sizeBytes := []byte{byte(size), byte(size >> 8)}

	player.Connection.Send(sizeBytes)
	player.Connection.Send(bytes)
}

func (player *Player) Disconnect() {
	if player == nil || player.Connection == nil {
		return
	}
	player.Connection.Close()
}

func (player *Player) IsConnected() bool {
	return player.Connection != nil && player.Connection.State() == network.StateOpen
}

func (player *Player) IsLoggedIn() bool {
	return player.IsConnected() && player.Account != nil
}

func (player *Player) IsPlaying() bool {
	return player.IsLoggedIn() && player.Character != nil
}

func (player *Player) GetMaxVital(vital database.VitalType) int {
	if player.Character == nil {
		return 0
	}

	switch vital {
	case database.VitalHP:
		return database.Classes[player.Character.Class].GetMaxVital(vital, player.Character.Stats.Strength)
	case database.VitalMP:
		return database.Classes[player.Character.Class].GetMaxVital(vital, player.Character.Stats.Magic)
	case database.VitalSP:
		return database.Classes[player.Character.Class].GetMaxVital(vital, player.Character.Stats.Speed)
	}

	return 0
}

func (player *Player) GetVital(vital database.VitalType) int {
	if player.Character == nil {
		return 0
	}

	switch vital {
	case database.VitalHP:
		return player.Character.Vitals.HP
	case database.VitalMP:
		return player.Character.Vitals.MP
	case database.VitalSP:
		return player.Character.Vitals.SP
	}

	return 0
}
