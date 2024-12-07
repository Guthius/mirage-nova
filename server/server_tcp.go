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

func (player *Player) SendWelcome() {
	player.SendMessage("Type /help for help on commands. Use arrow keys to move, hold down shift to run, and use ctrl to attack.", Cyan)

	if len(Motd) > 0 {
		player.SendMessage(fmt.Sprintf("MOTD: %s", Motd), BrightCyan)
	}

	player.SendPlayersOnline()
}

func (player *Player) SendPlayersOnline() {

	names := make([]string, 0, MaxPlayers)
	for i := 0; i < MaxPlayers; i++ {
		if i == player.Id {
			continue
		}

		if Players[i].IsPlaying() {
			names = append(names, Players[i].Character.Name)
		}
	}

	if len(names) == 0 {
		player.SendMessage("There are no other players online.", WhoColor)
		return
	}

	player.SendMessage(fmt.Sprintf("There are %d other players online: %s.", len(names), strings.Join(names, ", ")), WhoColor)
}

func (player *Player) ReportHack(reason string) {
	log.Printf("[%d] Terminating connection with %s (%s)\n", player.Id, player.Connection.RemoteAddr(), reason)

	if player.IsPlaying() {
		SendGlobalMessage(fmt.Sprintf("%s has been booted", player.Character.Name), White)
	}

	player.SendAlert(fmt.Sprintf("You have lost your connection with %s", GameName))
}

func (player *Player) SendAlert(message string) {
	if len(message) == 0 {
		return
	}

	writer := packet.NewWriter()

	writer.WriteInteger(SvAlert)
	writer.WriteString(message)

	player.Send(writer.Bytes())
	player.Disconnect()
}

func (player *Player) SendCharacters() {
	writer := packet.NewWriter()

	writer.WriteInteger(SvCharacters)

	for _, c := range player.CharacterList {
		writer.WriteLong(c.Sprite)
		writer.WriteString(c.Name)
		writer.WriteByte(byte(c.Level))
	}

	player.Send(writer.Bytes())
}

func (player *Player) SendLoginOk() {
	writer := packet.NewWriter()
	writer.WriteInteger(SvLoginOk)
	writer.WriteLong(player.Id + 1)

	player.Send(writer.Bytes())
}

func (player *Player) SendNewCharClasses() {
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

	player.Send(writer.Bytes())
}

func (player *Player) SendClasses() {
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

	player.Send(writer.Bytes())
}

func (player *Player) SendInGame() {
	writer := packet.NewWriter()

	writer.WriteInteger(SvInGame)

	player.Send(writer.Bytes())
}

func (player *Player) SendInventory() {
	if player.Character == nil {
		return
	}

	writer := packet.NewWriter()

	writer.WriteInteger(SvPlayerInventory)

	for i := 0; i < database.MaxInventory; i++ {
		writer.WriteLong(player.Character.Inv[i].Item + 1)
		writer.WriteLong(player.Character.Inv[i].Value)
		writer.WriteLong(player.Character.Inv[i].Dur)
	}

	player.Send(writer.Bytes())
}

func (player *Player) SendEquipment() {
	if player.Character == nil {
		return
	}

	writer := packet.NewWriter()

	writer.WriteInteger(SvPlayerEquipment)
	writer.WriteByte(byte(player.Character.Equipment.Weapon + 1))
	writer.WriteByte(byte(player.Character.Equipment.Armor + 1))
	writer.WriteByte(byte(player.Character.Equipment.Helmet + 1))
	writer.WriteByte(byte(player.Character.Equipment.Shield + 1))

	player.Send(writer.Bytes())
}

func (player *Player) SendVital(vital database.VitalType) {
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

	writer.WriteLong(player.GetMaxVital(vital))
	writer.WriteLong(player.GetVital(vital))

	player.Send(writer.Bytes())
}

func (player *Player) SendStats() {
	writer := packet.NewWriter()

	writer.WriteInteger(SvPlayerStats)
	writer.WriteLong(player.Character.Stats.Strength)
	writer.WriteLong(player.Character.Stats.Defense)
	writer.WriteLong(player.Character.Stats.Speed)
	writer.WriteLong(player.Character.Stats.Magic)

	player.Send(writer.Bytes())
}

func (player *Player) SendCheckForMap(mapId int) {
	mapData := database.GetMap(mapId)
	if mapData == nil {
		return
	}

	writer := packet.NewWriter()

	writer.WriteInteger(SvCheckForMap)
	writer.WriteLong(mapId + 1)
	writer.WriteLong(mapData.Revision)

	player.Send(writer.Bytes())
}

func SendGlobalMessage(message string, color Color) {
	writer := packet.NewWriter()

	writer.WriteInteger(SvGlobalMessage)
	writer.WriteString(message)
	writer.WriteByte(byte(color))

	SendDataToAll(writer.Bytes())
}

func (player *Player) SendMessage(message string, color Color) {
	writer := packet.NewWriter()

	writer.WriteInteger(SvPlayerMessage)
	writer.WriteString(message)
	writer.WriteByte(byte(color))

	player.Send(writer.Bytes())
}

func (player *Player) SendItems() {
	for i := 0; i < database.MaxItems; i++ {
		if len(database.Items[i].Name) == 0 {
			continue
		}

		player.SendUpdateItem(i)
	}
}

func (player *Player) SendUpdateItem(itemId int) {
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

	player.Send(writer.Bytes())
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

func (player *Player) SendNpcs() {
	for i := 0; i < database.MaxNpcs; i++ {
		if len(database.Npcs[i].Name) == 0 {
			continue
		}

		player.SendUpdateNpc(i)
	}
}

func (player *Player) SendUpdateNpc(npcId int) {
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

	player.Send(writer.Bytes())
}

func (player *Player) SendUpdateNpcToAll(npcId int) {
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

func (player *Player) SendShops() {
	for i := 0; i < database.MaxShops; i++ {
		if len(database.Shops[i].Name) == 0 {
			continue
		}

		player.SendUpdateShop(i)
	}
}

func (player *Player) SendUpdateShop(shopId int) {
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

	player.Send(writer.Bytes())
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

func (player *Player) SendSpells() {
	for i := 0; i < database.MaxSpells; i++ {
		if len(database.Spells[i].Name) == 0 {
			continue
		}

		player.SendUpdateSpell(i)
	}
}

func (player *Player) SendUpdateSpell(spellId int) {
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

	player.Send(writer.Bytes())
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

func (player *Player) SendLimits() {
	writer := packet.NewWriter()

	writer.WriteInteger(SvLimits)
	writer.WriteInteger(MaxPlayers)
	writer.WriteInteger(database.MaxItems)
	writer.WriteInteger(database.MaxNpcs)
	writer.WriteInteger(database.MaxShops)
	writer.WriteInteger(database.MaxSpells)
	writer.WriteInteger(database.MaxMaps)

	player.Send(writer.Bytes())
}

func (player *Player) SendMapRevisions() {
	writer := packet.NewWriter()

	writer.WriteInteger(SvMapRevisions)
	for i := 0; i < database.MaxMaps; i++ {
		writer.WriteLong(database.Maps[i].Revision)
	}

	player.Send(writer.Bytes())
}
