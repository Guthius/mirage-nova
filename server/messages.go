package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/guthius/mirage-nova/net"
	"github.com/guthius/mirage-nova/server/color"
	"github.com/guthius/mirage-nova/server/config"
	"github.com/guthius/mirage-nova/server/data"
	"github.com/guthius/mirage-nova/server/data/vitals"
)

func ReportHack(p *PlayerData, reason string) {
	log.Printf("[%d] Terminating connection with %s (%s)\n", p.Id, p.Connection.RemoteAddr(), reason)

	if p.IsPlaying() {
		SendGlobalMessage(fmt.Sprintf("%s has been booted", p.Character.Name), color.White)
	}

	SendAlert(p, fmt.Sprintf("You have lost your connection with %s", config.GameName))
}

func SendDataToAll(bytes []byte) {
	for _, p := range players {
		p.Send(bytes)
	}
}

func SendWelcome(player *PlayerData) {
	SendMessage(player, "Type /help for help on commands. Use arrow keys to move, hold down shift to run, and use ctrl to attack.", color.Cyan)

	if len(Motd) > 0 {
		SendMessage(player, fmt.Sprintf("MOTD: %s", Motd), color.BrightCyan)
	}

	SendPlayersOnline(player)
}

func SendGlobalMessage(message string, color color.Color) {
	writer := net.NewWriter()

	writer.WriteInteger(SvGlobalMessage)
	writer.WriteString(message)
	writer.WriteByte(byte(color))

	SendDataToAll(writer.Bytes())
}

func SendPlayersOnline(player *PlayerData) {
	// Get a slice with all the in game players.
	playing := GetPlayersInGame()
	if len(playing) == 0 {
		SendMessage(player, "There are no other players online.", color.WhoColor)
		return
	}

	// Get the names of all the in game players
	names := make([]string, 0, config.MaxPlayers)
	for _, p := range playing {
		names = append(names, p.Character.Name)
	}

	// Create a comma separated list of the names
	nameList := strings.Join(names, ", ")

	// Send the list to the player
	SendMessage(player, fmt.Sprintf("There are %d other players online: %s.", len(names), nameList), color.WhoColor)
}

func SendAlert(player *PlayerData, message string) {
	if len(message) == 0 {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvAlert)
	writer.WriteString(message)

	player.Send(writer.Bytes())
	player.Disconnect()
}

func SendCharacters(player *PlayerData) {
	writer := net.NewWriter()

	writer.WriteInteger(SvCharacters)

	for _, c := range player.CharacterList {
		writer.WriteLong(c.Sprite)
		writer.WriteString(c.Name)
		writer.WriteByte(byte(c.Level))
	}

	player.Send(writer.Bytes())
}

func SendLoginOk(player *PlayerData) {
	writer := net.NewWriter()
	writer.WriteInteger(SvLoginOk)
	writer.WriteLong(player.Id + 1)

	player.Send(writer.Bytes())
}

func SendNewCharClasses(player *PlayerData) {
	writer := net.NewWriter()

	numberOfClasses := data.GetClassCount()

	writer.WriteInteger(SvNewCharClasses)
	writer.WriteByte(byte(numberOfClasses))

	for i := 0; i < numberOfClasses; i++ {
		class := data.GetClass(i)
		if class == nil {
			continue
		}

		writer.WriteString(class.Name)
		writer.WriteLong(class.Sprite)
		writer.WriteLong(class.GetMaxVital(vitals.HP, class.Stats.Strength))
		writer.WriteLong(class.GetMaxVital(vitals.MP, class.Stats.Magic))
		writer.WriteLong(class.GetMaxVital(vitals.SP, class.Stats.Speed))
		writer.WriteByte(byte(class.Stats.Strength))
		writer.WriteByte(byte(class.Stats.Defense))
		writer.WriteByte(byte(class.Stats.Speed))
		writer.WriteByte(byte(class.Stats.Magic))
	}

	player.Send(writer.Bytes())
}

func SendClasses(player *PlayerData) {
	writer := net.NewWriter()

	numberOfClasses := data.GetClassCount()

	writer.WriteInteger(SvClasses)
	writer.WriteByte(byte(numberOfClasses))

	for i := 0; i < numberOfClasses; i++ {
		class := data.GetClass(i)
		if class == nil {
			continue
		}

		writer.WriteString(class.Name)
		writer.WriteLong(class.Sprite)
		writer.WriteLong(class.GetMaxVital(vitals.HP, class.Stats.Strength))
		writer.WriteLong(class.GetMaxVital(vitals.MP, class.Stats.Magic))
		writer.WriteLong(class.GetMaxVital(vitals.SP, class.Stats.Speed))
		writer.WriteByte(byte(class.Stats.Strength))
		writer.WriteByte(byte(class.Stats.Defense))
		writer.WriteByte(byte(class.Stats.Speed))
		writer.WriteByte(byte(class.Stats.Magic))
	}

	player.Send(writer.Bytes())
}

func SendInGame(player *PlayerData) {
	writer := net.NewWriter()

	writer.WriteInteger(SvInGame)

	player.Send(writer.Bytes())
}

func SendInventory(player *PlayerData) {
	if player.Character == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvPlayerInventory)

	for i := 0; i < config.MaxInventory; i++ {
		writer.WriteLong(player.Character.Inv[i].Item + 1)
		writer.WriteLong(player.Character.Inv[i].Value)
		writer.WriteLong(player.Character.Inv[i].Dur)
	}

	player.Send(writer.Bytes())
}

func SendEquipment(player *PlayerData) {
	if player.Character == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvPlayerEquipment)
	writer.WriteByte(byte(player.Character.Equipment.Weapon + 1))
	writer.WriteByte(byte(player.Character.Equipment.Armor + 1))
	writer.WriteByte(byte(player.Character.Equipment.Helmet + 1))
	writer.WriteByte(byte(player.Character.Equipment.Shield + 1))

	player.Send(writer.Bytes())
}

func SendVital(player *PlayerData, vital vitals.Type) {
	writer := net.NewWriter()

	switch vital {
	case vitals.HP:
		writer.WriteInteger(SvPlayerHP)
	case vitals.MP:
		writer.WriteInteger(SvPlayerMP)
	case vitals.SP:
		writer.WriteInteger(SvPlayerSP)
	default:
		return
	}

	writer.WriteLong(player.GetMaxVital(vital))
	writer.WriteLong(player.GetVital(vital))

	player.Send(writer.Bytes())
}

func SendStats(player *PlayerData) {
	writer := net.NewWriter()

	writer.WriteInteger(SvPlayerStats)
	writer.WriteLong(player.Character.Stats.Strength)
	writer.WriteLong(player.Character.Stats.Defense)
	writer.WriteLong(player.Character.Stats.Speed)
	writer.WriteLong(player.Character.Stats.Magic)

	player.Send(writer.Bytes())
}

func SendPlayerXY(player *PlayerData) {
	if player.Character == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvPlayerXY)
	writer.WriteLong(player.Character.X)
	writer.WriteLong(player.Character.Y)

	player.Send(writer.Bytes())
}

func SendCheckForLevel(player *PlayerData, levelId int) {
	levelData := data.GetLevel(levelId)
	if levelData == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvCheckForLevel)
	writer.WriteLong(levelId + 1)
	writer.WriteLong(levelData.Revision)

	player.Send(writer.Bytes())
}

func SendLevelData(player *PlayerData) {
	if player.Room == nil {
		return
	}
	player.Send(player.Room.LevelCache)
}

func SendMessage(player *PlayerData, message string, color color.Color) {
	writer := net.NewWriter()

	writer.WriteInteger(SvPlayerMessage)
	writer.WriteString(message)
	writer.WriteByte(byte(color))

	player.Send(writer.Bytes())
}

func SendItems(player *PlayerData) {
	for i := 0; i < config.MaxItems; i++ {
		item := data.GetItem(i)
		if item == nil || len(item.Name) == 0 {
			continue
		}

		SendUpdateItem(player, i)
	}
}

func SendUpdateItem(player *PlayerData, itemId int) {
	item := data.GetItem(itemId)
	if item == nil {
		return
	}

	writer := net.NewWriter()

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
	itemData := data.GetItem(itemId)
	if itemData == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvUpdateItem)
	writer.WriteLong(itemId + 1)
	writer.WriteString(itemData.Name)
	writer.WriteInteger(itemData.Pic)
	writer.WriteByte(byte(itemData.Type))
	writer.WriteInteger(itemData.Data1)
	writer.WriteInteger(itemData.Data2)
	writer.WriteInteger(itemData.Data3)

	SendDataToAll(writer.Bytes())
}

func SendNpcs(player *PlayerData) {
	for i := 0; i < config.MaxNpcs; i++ {
		npcData := data.GetNpc(i)
		if npcData == nil || len(npcData.Name) == 0 {
			continue
		}
		SendUpdateNpc(player, i)
	}
}

func SendUpdateNpc(player *PlayerData, npcId int) {
	npcData := data.GetNpc(npcId)
	if npcData == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvUpdateNpc)
	writer.WriteLong(npcId + 1)
	writer.WriteString(npcData.Name)
	writer.WriteInteger(npcData.Sprite)

	player.Send(writer.Bytes())
}

func SendUpdateNpcToAll(npcId int) {
	npcData := data.GetNpc(npcId)
	if npcData == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvUpdateNpc)
	writer.WriteLong(npcId + 1)
	writer.WriteString(npcData.Name)
	writer.WriteInteger(npcData.Sprite)

	SendDataToAll(writer.Bytes())
}

func SendShops(player *PlayerData) {
	for i := 0; i < config.MaxShops; i++ {
		shop := data.GetShop(i)
		if shop == nil || len(shop.Name) == 0 {
			continue
		}

		SendUpdateShop(player, i)
	}
}

func SendUpdateShop(player *PlayerData, shopId int) {
	shop := data.GetShop(shopId)
	if shop == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvUpdateShop)
	writer.WriteLong(shopId + 1)
	writer.WriteString(shop.Name)

	player.Send(writer.Bytes())
}

func SendUpdateShopToAll(shopId int) {
	shop := data.GetShop(shopId)
	if shop == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvUpdateShop)
	writer.WriteLong(shopId + 1)
	writer.WriteString(shop.Name)

	SendDataToAll(writer.Bytes())
}

func SendSpells(p *PlayerData) {
	for i := 0; i < config.MaxSpells; i++ {
		spell := data.GetSpell(i)
		if spell == nil || len(spell.Name) == 0 {
			continue
		}

		SendUpdateSpell(p, i)
	}
}

func SendUpdateSpell(player *PlayerData, spellId int) {
	spell := data.GetSpell(spellId)
	if spell == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvUpdateSpell)
	writer.WriteLong(spellId + 1)
	writer.WriteString(spell.Name)
	writer.WriteInteger(spell.MPReq)
	writer.WriteInteger(spell.Pic)

	player.Send(writer.Bytes())
}

func SendUpdateSpellToAll(spellId int) {
	spell := data.GetSpell(spellId)
	if spell == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvUpdateSpell)
	writer.WriteLong(spellId + 1)
	writer.WriteString(spell.Name)
	writer.WriteInteger(spell.MPReq)
	writer.WriteInteger(spell.Pic)

	SendDataToAll(writer.Bytes())
}

func SendDoorData(player *PlayerData) {
	if player.Room == nil {
		return
	}

	room := player.Room
	width := room.Level.Width
	height := room.Level.Height

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			tid := y*width + x
			tile := &room.TempTiles[tid]

			if tile.DoorOpen {
				writer := net.NewWriter()
				writer.WriteInteger(SvDoor)
				writer.WriteLong(x)
				writer.WriteLong(y)

				player.Send(writer.Bytes())
			}
		}
	}
}

func SendLimits(player *PlayerData) {
	writer := net.NewWriter()

	writer.WriteInteger(SvLimits)
	writer.WriteInteger(config.MaxPlayers)
	writer.WriteInteger(config.MaxItems)
	writer.WriteInteger(config.MaxNpcs)
	writer.WriteInteger(config.MaxShops)
	writer.WriteInteger(config.MaxSpells)
	writer.WriteInteger(config.MaxMaps)

	player.Send(writer.Bytes())
}

func SendMapRevisions(player *PlayerData) {
	writer := net.NewWriter()

	writer.WriteInteger(SvMapRevisions)
	for i := 0; i < config.MaxMaps; i++ {
		levelData := data.GetLevel(i)
		if levelData != nil {
			writer.WriteLong(levelData.Revision)
		} else {
			writer.WriteLong(0)
		}
	}

	player.Send(writer.Bytes())
}
