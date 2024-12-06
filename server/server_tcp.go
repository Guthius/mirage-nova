package main

import (
	"fmt"
	"log"
	"mirage/internal/database"
	"mirage/internal/packet"
	"strings"
)

func IsAccountLoggedIn(accountName string) bool {
	for _, player := range Players {
		if player.IsLoggedIn() && strings.EqualFold(player.Account.Name, accountName) {
			return true
		}
	}
	return false
}

func (p *Player) ReportHack(message string) {
	log.Printf("[%d] Terminating connection with %s (%s)\n", p.Id, p.Connection.RemoteAddr(), message)

	if p.IsPlaying() {
		// TODO:   Call GlobalMsg(GetPlayerLogin(Index) & "/" & GetPlayerName(Index) & " has been booted for (" & Reason & ")", White)
	}

	p.SendAlert(fmt.Sprintf("You have lost your connection with %s", GameName))
}

func (p *Player) SendAlert(message string) {
	if len(message) == 0 {
		return
	}

	writer := packet.NewWriter()

	writer.WriteInteger(SvpAlert)
	writer.WriteString(message)

	p.Send(writer.Bytes())
	p.Disconnect()
}

func (p *Player) SendCharacters() {
	writer := packet.NewWriter()

	writer.WriteInteger(SvpCharacters)

	for _, c := range p.Characters {
		writer.WriteLong(c.Sprite)
		writer.WriteString(c.Name)
		writer.WriteByte(byte(c.Level))
	}

	p.Send(writer.Bytes())
}

func (p *Player) SendLoginOk() {
	writer := packet.NewWriter()
	writer.WriteInteger(SvpLoginOk)

	p.Send(writer.Bytes())
}

func (p *Player) SendNewCharClasses() {
	writer := packet.NewWriter()

	numberOfClasses := len(database.Classes)

	writer.WriteInteger(SvpNewCharClasses)
	writer.WriteByte(byte(numberOfClasses))

	for _, class := range database.Classes {
		writer.WriteString(class.Name)
		writer.WriteLong(class.Sprite)
		writer.WriteLong(class.GetMaxVital(database.VitalHP))
		writer.WriteLong(class.GetMaxVital(database.VitalMP))
		writer.WriteLong(class.GetMaxVital(database.VitalSP))
		writer.WriteByte(byte(class.Stats.Strength))
		writer.WriteByte(byte(class.Stats.Defense))
		writer.WriteByte(byte(class.Stats.Speed))
		writer.WriteByte(byte(class.Stats.Magic))
	}

	p.Send(writer.Bytes())
}
