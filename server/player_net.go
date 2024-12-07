package main

import (
	"fmt"
	"strings"

	"github.com/guthius/mirage-nova/net"
	"github.com/guthius/mirage-nova/server/color"
	"github.com/guthius/mirage-nova/server/config"
	"github.com/guthius/mirage-nova/server/data"
	"github.com/guthius/mirage-nova/server/data/vitals"
)

func SendPlayersOnline(p *PlayerData) {
	// Get a slice with all the in game players.
	playing := GetPlaying()
	if len(playing) == 0 {
		SendMessage(p, "There are no other players online.", color.WhoColor)
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
	SendMessage(p, fmt.Sprintf("There are %d other players online: %s.", len(names), nameList), color.WhoColor)
}

func SendAlert(p *PlayerData, message string) {
	if len(message) == 0 {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvAlert)
	writer.WriteString(message)

	p.Send(writer.Bytes())
	p.Disconnect()
}

func SendCharacters(p *PlayerData) {
	writer := net.NewWriter()

	writer.WriteInteger(SvCharacters)

	for _, c := range p.CharacterList {
		writer.WriteLong(c.Sprite)
		writer.WriteString(c.Name)
		writer.WriteByte(byte(c.Level))
	}

	p.Send(writer.Bytes())
}

func SendLoginOk(p *PlayerData) {
	writer := net.NewWriter()
	writer.WriteInteger(SvLoginOk)
	writer.WriteLong(p.Id + 1)

	p.Send(writer.Bytes())
}

func SendNewCharClasses(p *PlayerData) {
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

	p.Send(writer.Bytes())
}

func SendClasses(p *PlayerData) {
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

	p.Send(writer.Bytes())
}

func SendInGame(p *PlayerData) {
	writer := net.NewWriter()

	writer.WriteInteger(SvInGame)

	p.Send(writer.Bytes())
}

func SendInventory(p *PlayerData) {
	if p.Character == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvPlayerInventory)

	for i := 0; i < config.MaxInventory; i++ {
		writer.WriteLong(p.Character.Inv[i].Item + 1)
		writer.WriteLong(p.Character.Inv[i].Value)
		writer.WriteLong(p.Character.Inv[i].Dur)
	}

	p.Send(writer.Bytes())
}

func SendEquipment(p *PlayerData) {
	if p.Character == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvPlayerEquipment)
	writer.WriteByte(byte(p.Character.Equipment.Weapon + 1))
	writer.WriteByte(byte(p.Character.Equipment.Armor + 1))
	writer.WriteByte(byte(p.Character.Equipment.Helmet + 1))
	writer.WriteByte(byte(p.Character.Equipment.Shield + 1))

	p.Send(writer.Bytes())
}

func SendVital(p *PlayerData, vital vitals.Type) {
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

	writer.WriteLong(p.GetMaxVital(vital))
	writer.WriteLong(p.GetVital(vital))

	p.Send(writer.Bytes())
}

func SendStats(p *PlayerData) {
	writer := net.NewWriter()

	writer.WriteInteger(SvPlayerStats)
	writer.WriteLong(p.Character.Stats.Strength)
	writer.WriteLong(p.Character.Stats.Defense)
	writer.WriteLong(p.Character.Stats.Speed)
	writer.WriteLong(p.Character.Stats.Magic)

	p.Send(writer.Bytes())
}

func SendCheckForMap(p *PlayerData, levelId int) {
	levelData := data.GetLevel(levelId)
	if levelData == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvCheckForMap)
	writer.WriteLong(levelId + 1)
	writer.WriteLong(levelData.Revision)

	p.Send(writer.Bytes())
}

func SendMessage(p *PlayerData, message string, color color.Color) {
	writer := net.NewWriter()

	writer.WriteInteger(SvPlayerMessage)
	writer.WriteString(message)
	writer.WriteByte(byte(color))

	p.Send(writer.Bytes())
}

func SendItems(p *PlayerData) {
	for i := 0; i < config.MaxItems; i++ {
		item := data.GetItem(i)
		if item == nil || len(item.Name) == 0 {
			continue
		}

		SendUpdateItem(p, i)
	}
}

func SendUpdateItem(p *PlayerData, itemId int) {
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

	p.Send(writer.Bytes())
}

func SendNpcs(p *PlayerData) {
	for i := 0; i < config.MaxNpcs; i++ {
		npcData := data.GetNpc(i)
		if npcData == nil || len(npcData.Name) == 0 {
			continue
		}
		SendUpdateNpc(p, i)
	}
}

func SendUpdateNpc(p *PlayerData, npcId int) {
	npcData := data.GetNpc(npcId)
	if npcData == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvUpdateNpc)
	writer.WriteLong(npcId + 1)
	writer.WriteString(npcData.Name)
	writer.WriteInteger(npcData.Sprite)

	p.Send(writer.Bytes())
}

func SendShops(p *PlayerData) {
	for i := 0; i < config.MaxShops; i++ {
		shop := data.GetShop(i)
		if shop == nil || len(shop.Name) == 0 {
			continue
		}

		SendUpdateShop(p, i)
	}
}

func SendUpdateShop(p *PlayerData, shopId int) {
	shop := data.GetShop(shopId)
	if shop == nil {
		return
	}

	writer := net.NewWriter()

	writer.WriteInteger(SvUpdateShop)
	writer.WriteLong(shopId + 1)
	writer.WriteString(shop.Name)

	p.Send(writer.Bytes())
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

func SendUpdateSpell(p *PlayerData, spellId int) {
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

	p.Send(writer.Bytes())
}

func SendLimits(p *PlayerData) {
	writer := net.NewWriter()

	writer.WriteInteger(SvLimits)
	writer.WriteInteger(config.MaxPlayers)
	writer.WriteInteger(config.MaxItems)
	writer.WriteInteger(config.MaxNpcs)
	writer.WriteInteger(config.MaxShops)
	writer.WriteInteger(config.MaxSpells)
	writer.WriteInteger(config.MaxMaps)

	p.Send(writer.Bytes())
}

func SendMapRevisions(p *PlayerData) {
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

	p.Send(writer.Bytes())
}
