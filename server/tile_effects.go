package main

import (
	"github.com/guthius/mirage-nova/net"
	"github.com/guthius/mirage-nova/server/color"
	"github.com/guthius/mirage-nova/server/data"
	"github.com/guthius/mirage-nova/server/data/vitals"
	"github.com/guthius/mirage-nova/server/utils"
)

// TriggerTileEffect triggers the effect of the tile the player is standing on.
func TriggerTileEffect(player *PlayerData) {
	if player.Room == nil || player.Character == nil {
		return
	}

	x := player.Character.X
	y := player.Character.Y

	tile := player.Room.GetTile(x, y)
	if tile == nil {
		return
	}

	switch tile.Data.Type {
	case data.TileTypeWarp:
		MovePlayerToRoom(player, tile.Data.Data1, tile.Data.Data2, tile.Data.Data3)

	case data.TileTypeKeyOpen:
		doorX := tile.Data.Data1
		doorY := tile.Data.Data2
		if !tile.DoorOpen {
			tile.DoorOpen = true
			tile.DoorTimer = utils.GetTickCount()

			buffer := net.NewWriter()
			buffer.WriteInteger(SMapKey)
			buffer.WriteLong(doorX)
			buffer.WriteLong(doorY)
			buffer.WriteLong(1)

			player.Room.Send(buffer.Bytes())
			player.Room.SendMessage("A door has been unlocked.", color.White)
		}

	case data.TileTypeHeal:
		player.Character.Vitals.HP = player.GetMaxVital(vitals.HP)
		SendVital(player, vitals.HP)
		SendMessage(player, "You feel odd as a strange glow eminated from you and your a lifted into the air. Bright orbs of light travel around you. You are miraculously healed!", color.BrightGreen)

	case data.TileTypeKill:
		TileKillPlayer(player)

	case data.TileTypeSprite:
		sprite := tile.Data.Data1
		player.Character.Sprite = sprite
		player.Room.SendPlayerData(player)
	}
}

func TileKillPlayer(player *PlayerData) {
	/*
		Dim nodamage As Boolean
		nodamage = False
		'Check the Armor Slot
		'If GetPlayerArmorSlot(Index) > 0 Then
		'    If GetPlayerInvItemNum(Index, GetPlayerArmorSlot(Index)) = Map(GetPlayerMap(Index)).Tile(GetPlayerX(Index), GetPlayerY(Index)).Data2 Then
		'        nodamage = True
		'    End If
		'End If
		'Check the Helmet Slot
		'If GetPlayerHelmetSlot(Index) > 0 Then
		'    If GetPlayerInvItemNum(Index, GetPlayerHelmetSlot(Index)) = Map(GetPlayerMap(Index)).Tile(GetPlayerX(Index), GetPlayerY(Index)).Data2 Then
		'        nodamage = True
		'    End If
		'End If
		'Check the Shield Slot
		'If GetPlayerShieldSlot(Index) > 0 Then
		'    If GetPlayerInvItemNum(Index, GetPlayerShieldSlot(Index)) = Map(GetPlayerMap(Index)).Tile(GetPlayerX(Index), GetPlayerY(Index)).Data2 Then
		'        nodamage = True
		'    End If
		'End If
		'Check the Weapon Slot
		'If GetPlayerWeaponSlot(Index) > 0 Then
		'    If GetPlayerInvItemNum(Index, GetPlayerWeaponSlot(Index)) = Map(GetPlayerMap(Index)).Tile(GetPlayerX(Index), GetPlayerY(Index)).Data2 Then
		'        nodamage = True
		'    End If
		'End If
		' Do Nothing
		If nodamage = False Then
		    ' Check to see if the sucker is going to die!
		    If GetPlayerVital(Index, Vitals.HP) > Trim$(Map(GetPlayerMap(Index)).Tile(GetPlayerX(Index), GetPlayerY(Index)).Data1) Then
		        Call SetPlayerVital(Index, Vitals.HP, GetPlayerVital(Index, Vitals.HP) - Trim$(Map(GetPlayerMap(Index)).Tile(GetPlayerX(Index), GetPlayerY(Index)).Data1))
		        Call SendVital(Index, Vitals.HP)
		        Call PlayerMsg(Index, "You've taken " & Trim$(Map(GetPlayerMap(Index)).Tile(GetPlayerX(Index), GetPlayerY(Index)).Data1) & " damage!", BrightRed)
		    ElseIf GetPlayerVital(Index, Vitals.HP) <= Trim$(Map(GetPlayerMap(Index)).Tile(GetPlayerX(Index), GetPlayerY(Index)).Data1) Then
		        Call PlayerMsg(Index, "You've taken " & Trim$(Map(GetPlayerMap(Index)).Tile(GetPlayerX(Index), GetPlayerY(Index)).Data1) & " damage, which has killed you!", BrightRed)
		        Call GlobalMsg("The player " & GetPlayerName(Index) & " has died!", BrightRed)
		        ' Warp player away
		        If Map(GetPlayerMap(Index)).BootMap > 0 And Map(GetPlayerMap(Index)).BootX > 0 And Map(GetPlayerMap(Index)).BootY > 0 Then
		            Call PlayerWarp(Index, Map(GetPlayerMap(Index)).BootMap, Map(GetPlayerMap(Index)).BootX, Map(GetPlayerMap(Index)).BootY)
		            Moved = YES
		        Else
		            Call PlayerWarp(Index, START_MAP, START_X, START_Y)
		            Moved = YES
		        End If
		        ' Restore vitals
		        Call SetPlayerVital(Index, Vitals.HP, GetPlayerMaxVital(Index, HP))
		        Call SetPlayerVital(Index, Vitals.MP, GetPlayerMaxVital(Index, MP))
		        Call SetPlayerVital(Index, Vitals.SP, GetPlayerMaxVital(Index, SP))
		        Call SendVital(Index, HP)
		        Call SendVital(Index, MP)
		        Call SendVital(Index, SP)
		    End If
		End If
	*/
}
