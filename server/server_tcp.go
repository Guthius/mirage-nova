package main

import (
	"mirage/internal/database"
	"mirage/internal/packet"
)

func IsAccountLoggedIn(accountName string) bool {
	for _, player := range Players {
		if player.IsLoggedIn() && player.Account.Name == accountName {
			return true
		}
	}
	return false
}

func (p *Player) SendAlert(message string) {
	if len(message) == 0 {
		return
	}

	writer := packet.NewWriter()

	writer.WriteInteger(SAlertMsg)
	writer.WriteString(message)

	p.Send(writer.Bytes())
	p.Disconnect()
}

func (p *Player) SendNewCharClasses() {
	writer := packet.NewWriter()

	numberOfClasses := len(database.Classes)

	writer.WriteInteger(SNewCharClasses)
	writer.WriteLong(numberOfClasses)

	for _, class := range database.Classes {
		writer.WriteString(class.Name)
		writer.WriteLong(class.GetMaxVital(database.VitalHP))
		writer.WriteLong(class.GetMaxVital(database.VitalMP))
		writer.WriteLong(class.GetMaxVital(database.VitalSP))
		writer.WriteLong(class.Stats.Strength)
		writer.WriteLong(class.Stats.Defense)
		writer.WriteLong(class.Stats.Speed)
		writer.WriteLong(class.Stats.Magic)
	}

	p.Send(writer.Bytes())
}
