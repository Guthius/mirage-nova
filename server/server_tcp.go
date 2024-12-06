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

func SendDataToAll(bytes []byte) {
	for _, player := range Players {
		player.Send(bytes)
	}
}

func SendGlobalMessage(message string, color Color) {
	writer := packet.NewWriter()

	writer.WriteInteger(SGlobalMsg)
	writer.WriteString(message)
	writer.WriteByte(byte(color))

	SendDataToAll(writer.Bytes())
}

func (p *Player) ReportHack(message string) {
	log.Printf("[%d] Terminating connection with %s (%s)\n", p.Id, p.Connection.RemoteAddr(), message)

	if p.IsPlaying() {
		SendGlobalMessage(fmt.Sprintf("%s/%s has been booted (%s)", p.Account.Name, p.Char.Name, message), White)
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
		writer.WriteLong(class.GetMaxVital(database.VitalHP, class.Stats.Strength))
		writer.WriteLong(class.GetMaxVital(database.VitalMP, class.Stats.Magic))
		writer.WriteLong(class.GetMaxVital(database.VitalSP, class.Stats.Speed))
		writer.WriteByte(byte(class.Stats.Strength))
		writer.WriteByte(byte(class.Stats.Defense))
		writer.WriteByte(byte(class.Stats.Speed))
		writer.WriteByte(byte(class.Stats.Magic))
	}

	p.Send(writer.Bytes())
}

func (p *Player) SendMaxes() {
	writer := packet.NewWriter()

	writer.WriteInteger(SSendMaxes)
	writer.WriteInteger(MaxPlayers)
	writer.WriteInteger(database.MaxItems)
	writer.WriteInteger(database.MaxNpcs)
	writer.WriteInteger(database.MaxShops)
	writer.WriteInteger(database.MaxSpells)
	writer.WriteInteger(database.MaxMaps)

	p.Send(writer.Bytes())
}

func (p *Player) SendMapRevs() {
	//     Dim I As Long
	//     Dim Buffer As clsBuffer

	//     Set Buffer = New clsBuffer

	//     Buffer.PreAllocate (MAX_MAPS * 4) + 2
	//     Buffer.WriteInteger SMapRevs
	//     For I = 1 To MAX_MAPS
	//         Buffer.WriteLong Map(I).Revision
	//     Next

	// Call SendDataTo(Index, Buffer.ToArray())
}
