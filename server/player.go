﻿package main

import "mirage/internal/network"

const (
	BufferSize = 4096
)

type Player struct {
	Connection *network.Conn
	Account    *Account
	Buffer     []byte
	InGame     bool
}

var Players [MaxPlayers]Player

func GetPlayer(index int) *Player {
	if index < 0 || index >= MaxPlayers {
		return nil
	}
	return &Players[index]
}

func (p *Player) Clear() {
	p.Connection = nil
	p.Account = nil
	p.Buffer = make([]byte, 0, BufferSize)
	p.InGame = false
}

func (p *Player) Send(bytes []byte) {
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

func (p *Player) Disconnect() {
	if p == nil || p.Connection == nil {
		return
	}
	p.Connection.Close()
}

func (p *Player) IsConnected() bool {
	return p.Connection != nil && p.Connection.State() == network.StateOpen
}

func (p *Player) IsLoggedIn() bool {
	return p.IsConnected() && p.Account != nil
}

func (p *Player) IsPlaying() bool {
	return p.IsLoggedIn() && p.InGame
}
