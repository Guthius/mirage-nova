package main

import (
	"fmt"
	"log"

	"github.com/guthius/mirage-nova/net"
	"github.com/guthius/mirage-nova/server/color"
	"github.com/guthius/mirage-nova/server/config"
	"github.com/guthius/mirage-nova/server/data"
)

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
