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

func (p *Player) SendWelcome() {
	p.SendMessage("Type /help for help on commands.  Use arrow keys to move, hold down shift to run, and use ctrl to attack.", Cyan)

	if len(Motd) > 0 {
		p.SendMessage(fmt.Sprintf("MOTD: %s", Motd), BrightCyan)
	}

	p.SendPlayersOnline()
}

func (p *Player) SendPlayersOnline() {

	names := make([]string, 0, MaxPlayers)
	for i := 0; i < MaxPlayers; i++ {
		if i == p.Id {
			continue
		}

		if Players[i].IsPlaying() {
			names = append(names, Players[i].Char.Name)
		}
	}

	if len(names) == 0 {
		p.SendMessage("There are no other players online.", WhoColor)
		return
	}

	p.SendMessage(fmt.Sprintf("There are %d other players online: %s.", len(names), strings.Join(names, ", ")), WhoColor)
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

	writer.WriteInteger(SvAlert)
	writer.WriteString(message)

	p.Send(writer.Bytes())
	p.Disconnect()
}

func (p *Player) SendCharacters() {
	writer := packet.NewWriter()

	writer.WriteInteger(SvCharacters)

	for _, c := range p.Characters {
		writer.WriteLong(c.Sprite)
		writer.WriteString(c.Name)
		writer.WriteByte(byte(c.Level))
	}

	p.Send(writer.Bytes())
}

func (p *Player) SendLoginOk() {
	writer := packet.NewWriter()
	writer.WriteInteger(SvLoginOk)
	writer.WriteLong(p.Id)

	p.Send(writer.Bytes())
}

func (p *Player) SendNewCharClasses() {
	writer := packet.NewWriter()

	numberOfClasses := len(database.Classes)

	writer.WriteInteger(SvNewCharClasses)
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

func (p *Player) SendClasses() {
	writer := packet.NewWriter()

	numberOfClasses := len(database.Classes)

	writer.WriteInteger(SvClasses)
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

func (p *Player) SendInGame() {
	writer := packet.NewWriter()

	writer.WriteInteger(SvInGame)

	p.Send(writer.Bytes())
}

func (p *Player) SendInventory() {
	if p.Char == nil {
		return
	}

	writer := packet.NewWriter()

	writer.WriteInteger(SvPlayerInventory)

	for i := 0; i < database.MaxInventory; i++ {
		writer.WriteLong(p.Char.Inv[i].Item + 1)
		writer.WriteLong(p.Char.Inv[i].Value)
		writer.WriteLong(p.Char.Inv[i].Dur)
	}

	p.Send(writer.Bytes())
}

func (p *Player) SendEquipment() {
	if p.Char == nil {
		return
	}

	writer := packet.NewWriter()

	writer.WriteInteger(SvPlayerEquipment)
	writer.WriteByte(byte(p.Char.Equipment.Weapon + 1))
	writer.WriteByte(byte(p.Char.Equipment.Armor + 1))
	writer.WriteByte(byte(p.Char.Equipment.Helmet + 1))
	writer.WriteByte(byte(p.Char.Equipment.Shield + 1))

	p.Send(writer.Bytes())
}

func (p *Player) SendVital(vital database.VitalType) {
	writer := packet.NewWriter()

	switch vital {
	case database.VitalHP:
		writer.WriteInteger(SvPlayerHP)
	case database.VitalMP:
		writer.WriteInteger(SvPlayerMP)
	case database.VitalSP:
		writer.WriteInteger(SvPlayerSP)
	default:
		return
	}

	writer.WriteLong(p.GetMaxVital(vital))
	writer.WriteLong(p.GetVital(vital))

	p.Send(writer.Bytes())
}

func (p *Player) SendStats() {
	writer := packet.NewWriter()

	writer.WriteInteger(SvPlayerStats)
	writer.WriteLong(p.Char.Stats.Strength)
	writer.WriteLong(p.Char.Stats.Defense)
	writer.WriteLong(p.Char.Stats.Speed)
	writer.WriteLong(p.Char.Stats.Magic)

	p.Send(writer.Bytes())
}

func (p *Player) SendCheckForMap(mapId int) {
	mapData := database.GetMap(mapId)
	if mapData == nil {
		return
	}

	writer := packet.NewWriter()

	writer.WriteInteger(SvCheckForMap)
	writer.WriteLong(mapId + 1)
	writer.WriteLong(mapData.Revision)

	p.Send(writer.Bytes())
}

func SendGlobalMessage(message string, color Color) {
	writer := packet.NewWriter()

	writer.WriteInteger(SvGlobalMessage)
	writer.WriteString(message)
	writer.WriteByte(byte(color))

	SendDataToAll(writer.Bytes())
}

func (p *Player) SendMessage(message string, color Color) {
	writer := packet.NewWriter()

	writer.WriteInteger(SvPlayerMessage)
	writer.WriteString(message)
	writer.WriteByte(byte(color))

	p.Send(writer.Bytes())
}

func (p *Player) SendItems() {
	for i := 0; i < database.MaxItems; i++ {
		if len(database.Items[i].Name) == 0 {
			continue
		}

		p.SendUpdateItem(i)
	}
}

func (p *Player) SendUpdateItem(itemId int) {
	if itemId < 0 || itemId >= database.MaxItems {
		return
	}

	item := database.GetItem(itemId)
	if item == nil {
		return
	}

	writer := packet.NewWriter()

	writer.WriteInteger(SvUpdateItem)
	writer.WriteLong(itemId + 1)
	writer.WriteString(item.Name)
	writer.WriteInteger(item.Pic)
	writer.WriteByte(byte(item.Type))
	writer.WriteInteger(item.Data1)
	writer.WriteInteger(item.Data2)
	writer.WriteInteger(item.Data3)

	p.Send(writer.Bytes())
}

func SendUpdateItemToAll(itemId int) {
	if itemId < 0 || itemId >= database.MaxItems {
		return
	}

	item := database.GetItem(itemId)
	if item == nil {
		return
	}

	writer := packet.NewWriter()

	writer.WriteInteger(SvUpdateItem)
	writer.WriteLong(itemId + 1)
	writer.WriteString(item.Name)
	writer.WriteInteger(item.Pic)
	writer.WriteByte(byte(item.Type))
	writer.WriteInteger(item.Data1)
	writer.WriteInteger(item.Data2)
	writer.WriteInteger(item.Data3)

	SendDataToAll(writer.Bytes())
}

func (p *Player) SendNpcs() {
	for i := 0; i < database.MaxNpcs; i++ {
		if len(database.Npcs[i].Name) == 0 {
			continue
		}

		p.SendUpdateNpc(i)
	}
}

func (p *Player) SendUpdateNpc(npcId int) {
	if npcId < 0 || npcId >= database.MaxNpcs {
		return
	}

	npc := database.GetNpc(npcId)
	if npc == nil {
		return
	}

	writer := packet.NewWriter()

	writer.WriteInteger(SvUpdateNpc)
	writer.WriteLong(npcId + 1)
	writer.WriteString(npc.Name)
	writer.WriteInteger(npc.Sprite)

	p.Send(writer.Bytes())
}

func (p *Player) SendUpdateNpcToAll(npcId int) {
	if npcId < 0 || npcId >= database.MaxNpcs {
		return
	}

	npc := &database.Npcs[npcId]

	writer := packet.NewWriter()

	writer.WriteInteger(SvUpdateNpc)
	writer.WriteLong(npcId + 1)
	writer.WriteString(npc.Name)
	writer.WriteInteger(npc.Sprite)

	SendDataToAll(writer.Bytes())
}

func (p *Player) SendShops() {
	for i := 0; i < database.MaxShops; i++ {
		if len(database.Shops[i].Name) == 0 {
			continue
		}

		p.SendUpdateShop(i)
	}
}

func (p *Player) SendUpdateShop(shopId int) {
	if shopId < 0 || shopId >= database.MaxShops {
		return
	}

	shop := database.GetShop(shopId)
	if shop == nil {
		return
	}

	writer := packet.NewWriter()

	writer.WriteInteger(SvUpdateShop)
	writer.WriteLong(shopId + 1)
	writer.WriteString(shop.Name)

	p.Send(writer.Bytes())
}

func SendUpdateShopToAll(shopId int) {
	if shopId < 0 || shopId >= database.MaxShops {
		return
	}

	shop := database.GetShop(shopId)
	if shop == nil {
		return
	}

	writer := packet.NewWriter()

	writer.WriteInteger(SvUpdateShop)
	writer.WriteLong(shopId + 1)
	writer.WriteString(shop.Name)

	SendDataToAll(writer.Bytes())
}

func (p *Player) SendSpells() {
	for i := 0; i < database.MaxSpells; i++ {
		if len(database.Spells[i].Name) == 0 {
			continue
		}

		p.SendUpdateSpell(i)
	}
}

func (p *Player) SendUpdateSpell(spellId int) {
	if spellId < 0 || spellId >= database.MaxSpells {
		return
	}

	spell := &database.Spells[spellId]

	writer := packet.NewWriter()

	writer.WriteInteger(SvUpdateSpell)
	writer.WriteLong(spellId + 1)
	writer.WriteString(spell.Name)
	writer.WriteInteger(spell.MPReq)
	writer.WriteInteger(spell.Pic)

	p.Send(writer.Bytes())
}

func SendUpdateSpellToAll(spellId int) {
	if spellId < 0 || spellId >= database.MaxSpells {
		return
	}

	spell := &database.Spells[spellId]

	writer := packet.NewWriter()

	writer.WriteInteger(SvUpdateSpell)
	writer.WriteLong(spellId + 1)
	writer.WriteString(spell.Name)
	writer.WriteInteger(spell.MPReq)
	writer.WriteInteger(spell.Pic)

	SendDataToAll(writer.Bytes())
}

func (p *Player) SendLimits() {
	writer := packet.NewWriter()

	writer.WriteInteger(SvLimits)
	writer.WriteInteger(MaxPlayers)
	writer.WriteInteger(database.MaxItems)
	writer.WriteInteger(database.MaxNpcs)
	writer.WriteInteger(database.MaxShops)
	writer.WriteInteger(database.MaxSpells)
	writer.WriteInteger(database.MaxMaps)

	p.Send(writer.Bytes())
}

func (p *Player) SendMapRevisions() {
	writer := packet.NewWriter()

	writer.WriteInteger(SvMapRevisions)
	for i := 0; i < database.MaxMaps; i++ {
		writer.WriteLong(database.Maps[i].Revision)
	}

	p.Send(writer.Bytes())
}
