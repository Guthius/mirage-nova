﻿package main

import (
	"fmt"

	"github.com/guthius/mirage-nova/net"
	"github.com/guthius/mirage-nova/server/character"
	"github.com/guthius/mirage-nova/server/color"
	"github.com/guthius/mirage-nova/server/common"
	"github.com/guthius/mirage-nova/server/config"
	"github.com/guthius/mirage-nova/server/data"
	"github.com/guthius/mirage-nova/server/data/equipment"
	"github.com/guthius/mirage-nova/server/data/vitals"
	"github.com/guthius/mirage-nova/server/utils"
)

// Public Sub AttackNpc(ByVal Attacker As Long, ByVal MapNpcNum As Long, ByVal Damage As Long)
//     Dim Name As String
//     Dim Exp As Long
//     Dim n As Long
//     Dim I As Long
//     Dim Str As Long
//     Dim DEF As Long
//     Dim MapNum As Long
//     Dim NpcNum As Long
//     Dim Buffer As clsBuffer

//     ' Check for subscript out of range
//     If IsPlaying(Attacker) = False Or MapNpcNum <= 0 Or MapNpcNum > MAX_MAP_NPCS Or Damage < 0 Then
//         Exit Sub
//     End If

//     MapNum = GetPlayerMap(Attacker)
//     NpcNum = MapNpc(MapNum, MapNpcNum).Num
//     Name = Trim$(Npc(NpcNum).Name)

//     Set Buffer = New clsBuffer

//     Buffer.PreAllocate 6
//     Buffer.WriteInteger SAttack
//     Buffer.WriteLong Attacker
//     ' Send this packet so they can see the person attacking
//     Call SendDataToMapBut(Attacker, MapNum, Buffer.ToArray())

//     Set Buffer = Nothing

//     ' Check for weapon
//     If GetPlayerEquipmentSlot(Attacker, Weapon) > 0 Then
//         n = GetPlayerInvItemNum(Attacker, GetPlayerEquipmentSlot(Attacker, Weapon))
//     End If

//     If Damage >= MapNpc(MapNum, MapNpcNum).Vital(Vitals.HP) Then
//         ' Check for a weapon and say damage
//         If n = 0 Then
//             Call PlayerMsg(Attacker, "You hit a " & Name & " for " & Damage & " hit points, killing it.", BrightRed)
//         Else
//             Call PlayerMsg(Attacker, "You hit a " & Name & " with a " & Trim$(Item(n).Name) & " for " & Damage & " hit points, killing it.", BrightRed)
//         End If

//         ' Calculate exp to give attacker
//         Str = Npc(NpcNum).Stat(Stats.Strength)
//         DEF = Npc(NpcNum).Stat(Stats.Defense)
//         Exp = Str * DEF * 2

//         ' Make sure we dont get less then 0
//         If Exp < 0 Then
//             Exp = 1
//         End If

//         ' Check if in party, if so divide the exp up by 2
//         If TempPlayer(Attacker).InParty = NO Then
//             Call SetPlayerExp(Attacker, GetPlayerExp(Attacker) + Exp)
//             Call PlayerMsg(Attacker, "You have gained " & Exp & " experience points.", BrightBlue)
//         Else
//             Exp = Exp / 2

//             If Exp < 0 Then
//                 Exp = 1
//             End If

//             Call SetPlayerExp(Attacker, GetPlayerExp(Attacker) + Exp)
//             Call PlayerMsg(Attacker, "You have gained " & Exp & " party experience points.", BrightBlue)

//             n = TempPlayer(Attacker).PartyPlayer
//             If n > 0 Then
//                 Call SetPlayerExp(n, GetPlayerExp(n) + Exp)
//                 Call PlayerMsg(n, "You have gained " & Exp & " party experience points.", BrightBlue)
//             End If
//         End If

//         ' Drop the goods if they get it
//         n = Int(Rnd * Npc(NpcNum).DropChance) + 1
//         If n = 1 Then
//             Call SpawnItem(Npc(NpcNum).DropItem, Npc(NpcNum).DropItemValue, MapNum, MapNpc(MapNum, MapNpcNum).X, MapNpc(MapNum, MapNpcNum).Y)
//         End If

//         ' Now set HP to 0 so we know to actually kill them in the server loop (this prevents subscript out of range)
//         MapNpc(MapNum, MapNpcNum).Num = 0
//         MapNpc(MapNum, MapNpcNum).SpawnWait = GetTickCount
//         MapNpc(MapNum, MapNpcNum).Vital(Vitals.HP) = 0

//         Set Buffer = New clsBuffer

//         Buffer.PreAllocate 6 + 4
//         Buffer.WriteInteger SNpcDead
//         Buffer.WriteLong MapNum
//         Buffer.WriteLong MapNpcNum
//         Call SendDataToMap(MapNum, Buffer.ToArray())

//         ' Check for level up
//         Call CheckPlayerLevelUp(Attacker)

//         ' Check for level up party member
//         If TempPlayer(Attacker).InParty = YES Then
//             Call CheckPlayerLevelUp(TempPlayer(Attacker).PartyPlayer)
//         End If

//         ' Check if target is Npc that died and if so set target to 0
//         If TempPlayer(Attacker).TargetType = TARGET_TYPE_NPC Then
//             If TempPlayer(Attacker).Target = MapNpcNum Then
//                 TempPlayer(Attacker).Target = 0
//                 TempPlayer(Attacker).TargetType = TARGET_TYPE_NONE
//             End If
//         End If
//     Else
//         ' Npc not dead, just do the damage
//         MapNpc(MapNum, MapNpcNum).Vital(Vitals.HP) = MapNpc(MapNum, MapNpcNum).Vital(Vitals.HP) - Damage

//         ' Check for a weapon and say damage
//         If n = 0 Then
//             Call PlayerMsg(Attacker, "You hit a " & Name & " for " & Damage & " hit points.", White)
//         Else
//             Call PlayerMsg(Attacker, "You hit a " & Name & " with a " & Trim$(Item(n).Name) & " for " & Damage & " hit points.", White)
//         End If

//         ' Check if we should send a message
//         If MapNpc(MapNum, MapNpcNum).Target = 0 Then
//             If LenB(Trim$(Npc(NpcNum).AttackSay)) > 0 Then
//                 Call PlayerMsg(Attacker, "A " & Trim$(Npc(NpcNum).Name) & " says, '" & Trim$(Npc(NpcNum).AttackSay) & "' to you.", SayColor)
//             End If
//         End If

//         ' Set the Npc target to the player
//         MapNpc(MapNum, MapNpcNum).Target = Attacker

//         ' Now check for guard ai and if so have all onmap guards come after'm
//         If Npc(MapNpc(MapNum, MapNpcNum).Num).Behavior = Npc_BEHAVIOR_GUARD Then
//             For I = 1 To MAX_MAP_NPCS
//                 If MapNpc(MapNum, I).Num = MapNpc(MapNum, MapNpcNum).Num Then
//                     MapNpc(MapNum, I).Target = Attacker
//                 End If
//             Next
//         End If
//     End If

//     ' Reduce durability of weapon
//     Call DamageEquipment(Attacker, Weapon)

//     ' Reset attack timer
//     TempPlayer(Attacker).AttackTimer = GetTickCount
// End Sub

// Public Sub AttackPlayer(ByVal Attacker As Long, ByVal Victim As Long, ByVal Damage As Long)
//     Dim Exp As Long
//     Dim n As Long
//     Dim I As Long
//     Dim Buffer As clsBuffer

//     ' Check for subscript out of range
//     If IsPlaying(Attacker) = False Or IsPlaying(Victim) = False Or Damage < 0 Then
//         Exit Sub
//     End If

//     ' Check for weapon
//     n = 0
//     If GetPlayerEquipmentSlot(Attacker, Weapon) > 0 Then
//         n = GetPlayerInvItemNum(Attacker, GetPlayerEquipmentSlot(Attacker, Weapon))
//     End If

//     ' Send this packet so they can see the person attacking
//     Set Buffer = New clsBuffer
//     Buffer.PreAllocate 6
//     Buffer.WriteInteger SAttack
//     Buffer.WriteLong Attacker
//     Call SendDataToMapBut(Attacker, GetPlayerMap(Attacker), Buffer.ToArray())
//     Set Buffer = Nothing

//     ' reduce dur. on victims equipment
//     Call DamageEquipment(Victim, Armor)
//     Call DamageEquipment(Victim, Helmet)

//     If Damage >= GetPlayerVital(Victim, Vitals.HP) Then
//         ' Check for a weapon and say damage
//         If n = 0 Then
//             Call PlayerMsg(Attacker, "You hit " & GetPlayerName(Victim) & " for " & Damage & " hit points.", White)
//             Call PlayerMsg(Victim, GetPlayerName(Attacker) & " hit you for " & Damage & " hit points.", BrightRed)
//         Else
//             Call PlayerMsg(Attacker, "You hit " & GetPlayerName(Victim) & " with a " & Trim$(Item(n).Name) & " for " & Damage & " hit points.", White)
//             Call PlayerMsg(Victim, GetPlayerName(Attacker) & " hit you with a " & Trim$(Item(n).Name) & " for " & Damage & " hit points.", BrightRed)
//         End If

//         ' Player is dead
//         Call GlobalMsg(GetPlayerName(Victim) & " has been killed by " & GetPlayerName(Attacker), BrightRed)

//         ' Calculate exp to give attacker
//         Exp = (GetPlayerExp(Victim) \ 10)

//         ' Make sure we dont get less then 0
//         If Exp < 0 Then
//             Exp = 0
//         End If

//         If Exp = 0 Then
//             Call PlayerMsg(Victim, "You lost no experience points.", BrightRed)
//             Call PlayerMsg(Attacker, "You received no experience points from that weak insignificant player.", BrightBlue)
//         Else
//             Call SetPlayerExp(Victim, GetPlayerExp(Victim) - Exp)
//             Call PlayerMsg(Victim, "You lost " & Exp & " experience points.", BrightRed)
//             Call SetPlayerExp(Attacker, GetPlayerExp(Attacker) + Exp)
//             Call PlayerMsg(Attacker, "You got " & Exp & " experience points for killing " & GetPlayerName(Victim) & ".", BrightBlue)
//         End If

//         ' Check for a level up
//         Call CheckPlayerLevelUp(Attacker)

//         ' Check if target is player who died and if so set target to 0
//         If TempPlayer(Attacker).TargetType = TARGET_TYPE_PLAYER Then
//             If TempPlayer(Attacker).Target = Victim Then
//                 TempPlayer(Attacker).Target = 0
//                 TempPlayer(Attacker).TargetType = TARGET_TYPE_NONE
//             End If
//         End If

//         If Map(GetPlayerMap(Attacker)).Moral <> MAP_MORAL_ARENA Then
//             If GetPlayerPK(Victim) = NO Then
//                 If GetPlayerPK(Attacker) = NO Then
//                     Call SetPlayerPK(Attacker, YES)
//                     Call SendPlayerData(Attacker)
//                     Call GlobalMsg(GetPlayerName(Attacker) & " has been deemed a Player Killer!!!", BrightRed)
//                 End If
//             Else
//                 Call GlobalMsg(GetPlayerName(Victim) & " has paid the price for being a Player Killer!!!", BrightRed)
//             End If
//         End If

//         Call OnDeath(Victim)
//     Else
//         ' Player not dead, just do the damage
//         Call SetPlayerVital(Victim, Vitals.HP, GetPlayerVital(Victim, Vitals.HP) - Damage)
//         Call SendVital(Victim, Vitals.HP)

//         ' Check for a weapon and say damage
//         If n = 0 Then
//             Call PlayerMsg(Attacker, "You hit " & GetPlayerName(Victim) & " for " & Damage & " hit points.", White)
//             Call PlayerMsg(Victim, GetPlayerName(Attacker) & " hit you for " & Damage & " hit points.", BrightRed)
//         Else
//             Call PlayerMsg(Attacker, "You hit " & GetPlayerName(Victim) & " with a " & Trim$(Item(n).Name) & " for " & Damage & " hit points.", White)
//             Call PlayerMsg(Victim, GetPlayerName(Attacker) & " hit you with a " & Trim$(Item(n).Name) & " for " & Damage & " hit points.", BrightRed)
//         End If
//     End If

//     ' Reduce durability of weapon
//     Call DamageEquipment(Attacker, Weapon)

//     ' Reset attack timer
//     TempPlayer(Attacker).AttackTimer = GetTickCount
// End Sub

// Public Function FindOpenMapItemSlot(ByVal MapNum As Long) As Long
//     Dim I As Long

//     ' Check for subscript out of range
//     If MapNum <= 0 Or MapNum > MAX_MAPS Then
//         Exit Function
//     End If

//     For I = 1 To MAX_MAP_ITEMS
//         If MapItem(MapNum, I).Num = 0 Then
//             FindOpenMapItemSlot = I
//             Exit Function
//         End If
//     Next
// End Function

func JoinGame(p *PlayerData) {
	char := p.Character
	if char == nil {
		return
	}

	PlayersOnline++

	UpdateHighIndex()

	if char.Access == character.AccessNone {
		SendGlobalMessage(fmt.Sprintf("%s has joined %s!", char.Name, config.GameName), color.JoinLeftColor)
	} else {
		SendGlobalMessage(fmt.Sprintf("%s has joined %s!", char.Name, config.GameName), color.White)
	}

	// Send an ok to client to start receiving in game data
	SendLoginOk(p)

	// Send some more little goodies, no need to explain these
	CheckEquippedItems(p)

	SendClasses(p)
	SendItems(p)
	SendNpcs(p)
	SendShops(p)
	SendSpells(p)
	SendInventory(p)
	SendEquipment(p)
	SendVital(p, vitals.HP)
	SendVital(p, vitals.MP)
	SendVital(p, vitals.SP)
	SendStats(p)

	// Warp the player to his saved location
	rooms[char.Room].AddPlayer(p)

	// Send welcome messages
	SendWelcome(p)

	// Send the flag so they know they can start doing stuff
	SendInGame(p)
}

// Public Sub LeftGame(ByVal Index As Long)
//     Dim n As Long

//     If TempPlayer(Index).InGame Then
//         TempPlayer(Index).InGame = False

//         ' Check if player was the only player on the map and stop Npc processing if so
//         If GetTotalMapPlayers(GetPlayerMap(Index)) < 1 Then
//             PlayersOnMap(GetPlayerMap(Index)) = NO
//         End If

//         ' Check for boot map
//         If Map(GetPlayerMap(Index)).BootMap > 0 Then
//             Call SetPlayerX(Index, Map(GetPlayerMap(Index)).BootX)
//             Call SetPlayerY(Index, Map(GetPlayerMap(Index)).BootY)
//             Call SetPlayerMap(Index, Map(GetPlayerMap(Index)).BootMap)
//         End If

//         ' Check if the player was in a party, and if so cancel it out so the other player doesn't continue to get half exp
//         If TempPlayer(Index).InParty = YES Then
//             n = TempPlayer(Index).PartyPlayer

//             Call PlayerMsg(n, GetPlayerName(Index) & " has left " & GAME_NAME & ", disbanning party.", Pink)
//             TempPlayer(n).InParty = NO
//             TempPlayer(n).PartyPlayer = 0
//         End If

//         Call SavePlayer(Index)

//         ' Send a global message that he/she left
//         If GetPlayerAccess(Index) <= ADMIN_MONITOR Then
//             Call GlobalMsg(GetPlayerName(Index) & " has left " & GAME_NAME & "!", JoinLeftColor)
//         Else
//             Call GlobalMsg(GetPlayerName(Index) & " has left " & GAME_NAME & "!", White)
//         End If

//         Call TextAdd(GetPlayerName(Index) & " has disconnected from " & GAME_NAME & ".")
//         Call SendLeftGame(Index)

//         TotalPlayersOnline = TotalPlayersOnline - 1
//         Call UpdateHighIndex

//     End If

//     Call ClearPlayer(Index)
// End Sub

// Public Function FindPlayer(ByVal Name As String) As Long
//     Dim I As Long

//     For I = 1 To TotalPlayersOnline
//         ' Make sure we dont try to check a name thats to small
//         If Len(GetPlayerName(PlayersOnline(I))) >= Len(Trim$(Name)) Then
//             If UCase$(Mid$(GetPlayerName(PlayersOnline(I)), 1, Len(Trim$(Name)))) = UCase$(Trim$(Name)) Then
//                 FindPlayer = PlayersOnline(I)
//                 Exit Function
//             End If
//         End If
//     Next

//     FindPlayer = 0
// End Function

// Public Sub SpawnItem(ByVal ItemNum As Long, ByVal ItemVal As Long, ByVal MapNum As Long, ByVal X As Long, ByVal Y As Long)
//     Dim I As Long

//     ' Check for subscript out of range
//     If ItemNum < 1 Or ItemNum > MAX_ITEMS Or MapNum <= 0 Or MapNum > MAX_MAPS Then
//         Exit Sub
//     End If

//     ' Find open map item slot
//     I = FindOpenMapItemSlot(MapNum)

//     Call SpawnItemSlot(I, ItemNum, ItemVal, Item(ItemNum).Data1, MapNum, X, Y)
// End Sub

// Public Sub SpawnItemSlot(ByVal MapItemSlot As Long, ByVal ItemNum As Long, ByVal ItemVal As Long, ByVal ItemDur As Long, ByVal MapNum As Long, ByVal X As Long, ByVal Y As Long)
//     Dim Packet As String
//     Dim I As Long
//     Dim Buffer As clsBuffer

//     ' Check for subscript out of range
//     If MapItemSlot <= 0 Or MapItemSlot > MAX_MAP_ITEMS Or ItemNum < 0 Or ItemNum > MAX_ITEMS Or MapNum <= 0 Or MapNum > MAX_MAPS Then
//         Exit Sub
//     End If

//     I = MapItemSlot

//     If I <> 0 Then
//         If ItemNum >= 0 Then
//             If ItemNum <= MAX_ITEMS Then

//                 MapItem(MapNum, I).Num = ItemNum
//                 MapItem(MapNum, I).value = ItemVal

//                 If ItemNum <> 0 Then
//                     If (Item(ItemNum).Type >= ITEM_TYPE_WEAPON) And (Item(ItemNum).Type <= ITEM_TYPE_SHIELD) Then
//                         MapItem(MapNum, I).Dur = ItemDur
//                     Else
//                         MapItem(MapNum, I).Dur = 0
//                     End If
//                 Else
//                     MapItem(MapNum, I).Dur = 0
//                 End If

//                 MapItem(MapNum, I).X = X
//                 MapItem(MapNum, I).Y = Y

//                 Set Buffer = New clsBuffer

//                 Buffer.PreAllocate 26 + 4
//                 Buffer.WriteInteger SSpawnItem
//                 Buffer.WriteLong MapNum
//                 Buffer.WriteLong I
//                 Buffer.WriteLong ItemNum
//                 Buffer.WriteLong ItemVal
//                 Buffer.WriteLong MapItem(MapNum, I).Dur
//                 Buffer.WriteLong X
//                 Buffer.WriteLong Y
//                 Call SendDataToAll(Buffer.ToArray())

//                 Set Buffer = Nothing
//             End If
//         End If
//     End If

// End Sub

// Public Sub SpawnAllMapsItems()
//     Dim I As Long

//     For I = 1 To MAX_MAPS
//         Call SpawnMapItems(I)
//     Next
// End Sub

// Public Sub SpawnMapItems(ByVal MapNum As Long)
//     Dim X As Long
//     Dim Y As Long

//     ' Check for subscript out of range
//     If MapNum <= 0 Or MapNum > MAX_MAPS Then
//         Exit Sub
//     End If

//     ' Spawn what we have
//     For X = 0 To MAX_MAPX
//         For Y = 0 To MAX_MAPY
//             ' Check if the tile type is an item or a saved tile incase someone drops something
//             If (Map(MapNum).Tile(X, Y).Type = TILE_TYPE_ITEM) Then
//                 ' Check to see if its a currency and if they set the value to 0 set it to 1 automatically
//                 If Item(Map(MapNum).Tile(X, Y).Data1).Type = ITEM_TYPE_CURRENCY And Map(MapNum).Tile(X, Y).Data2 <= 0 Then
//                     Call SpawnItem(Map(MapNum).Tile(X, Y).Data1, 1, MapNum, X, Y)
//                 Else
//                     Call SpawnItem(Map(MapNum).Tile(X, Y).Data1, Map(MapNum).Tile(X, Y).Data2, MapNum, X, Y)
//                 End If
//             End If
//         Next
//     Next
// End Sub

// Public Sub SpawnNpc(ByVal MapNpcNum As Long, ByVal MapNum As Long)
//     Dim Packet As String
//     Dim NpcNum As Long
//     Dim I As Long
//     Dim X As Long
//     Dim Y As Long
//     Dim Spawned As Boolean
//     Dim Buffer As clsBuffer

//     ' Check for subscript out of range
//     If MapNpcNum <= 0 Or MapNpcNum > MAX_MAP_NPCS Or MapNum <= 0 Or MapNum > MAX_MAPS Then
//         Exit Sub
//     End If

//     NpcNum = Map(MapNum).Npc(MapNpcNum)
//     If NpcNum > 0 Then
//         MapNpc(MapNum, MapNpcNum).Num = NpcNum
//         MapNpc(MapNum, MapNpcNum).Target = 0

//         MapNpc(MapNum, MapNpcNum).Vital(Vitals.HP) = GetNpcMaxVital(NpcNum, Vitals.HP)
//         MapNpc(MapNum, MapNpcNum).Vital(Vitals.MP) = GetNpcMaxVital(NpcNum, Vitals.MP)
//         MapNpc(MapNum, MapNpcNum).Vital(Vitals.SP) = GetNpcMaxVital(NpcNum, Vitals.SP)

//         MapNpc(MapNum, MapNpcNum).Dir = Int(Rnd * 4)

//         ' Check if theres a spawn tile for the specific npc
//         For X = 0 To MAX_MAPX
//             For Y = 0 To MAX_MAPY
//                 If Map(MapNum).Tile(X, Y).Type = TILE_TYPE_NPCSPAWN Then
//                     If Map(MapNum).Tile(X, Y).Data1 = MapNpcNum Then
//                         MapNpc(MapNum, MapNpcNum).X = X
//                         MapNpc(MapNum, MapNpcNum).Y = Y
//                         MapNpc(MapNum, MapNpcNum).Dir = Map(MapNum).Tile(X, Y).Data2
//                         'MapNpc(MapNum, MapNpcNum).Moveable = Map(MapNum).Tile(X, y).Data3
//                         Spawned = True
//                         Exit For
//                     End If
//                 End If
//             Next Y
//         Next X

//         ' Well try 100 times to randomly place the sprite
//         If Not Spawned Then
//             For I = 1 To 100
//                 X = Int(Rnd * MAX_MAPX)
//                 Y = Int(Rnd * MAX_MAPY)

//                 ' Check if the tile is walkable
//                 If Map(MapNum).Tile(X, Y).Type = TILE_TYPE_WALKABLE Then
//                     MapNpc(MapNum, MapNpcNum).X = X
//                     MapNpc(MapNum, MapNpcNum).Y = Y
//                     Spawned = True
//                     Exit For
//                 End If
//             Next
//         End If

//         ' Didn't spawn, so now we'll just try to find a free tile
//         If Not Spawned Then
//             For X = 0 To MAX_MAPX
//                 For Y = 0 To MAX_MAPY
//                     If Map(MapNum).Tile(X, Y).Type = TILE_TYPE_WALKABLE Then
//                         MapNpc(MapNum, MapNpcNum).X = X
//                         MapNpc(MapNum, MapNpcNum).Y = Y
//                         Spawned = True
//                     End If
//                 Next
//             Next
//         End If

//         ' If we succeeded in spawning then send it to everyone
//         If Spawned Then
//             Set Buffer = New clsBuffer

//             Buffer.PreAllocate 12 + 4
//             Buffer.WriteInteger SSpawnNpc
//             Buffer.WriteLong MapNum
//             Buffer.WriteLong MapNpcNum
//             Buffer.WriteInteger MapNpc(MapNum, MapNpcNum).Num
//             Buffer.WriteByte MapNpc(MapNum, MapNpcNum).X
//             Buffer.WriteByte MapNpc(MapNum, MapNpcNum).Y
//             Buffer.WriteInteger MapNpc(MapNum, MapNpcNum).Dir

//             Call SendDataToAll(Buffer.ToArray())
//         End If
//     End If
// End Sub

// Public Sub SpawnMapNpcs(ByVal MapNum As Long)
//     Dim I As Long

//     For I = 1 To MAX_MAP_NPCS
//         Call SpawnNpc(I, MapNum)
//     Next
// End Sub

// Public Sub SpawnAllMapNpcs()
//     Dim I As Long

//     For I = 1 To MAX_MAPS
//         Call SpawnMapNpcs(I)
//     Next
// End Sub

// Public Function CanAttackPlayer(ByVal Attacker As Long, ByVal Victim As Long) As Boolean
//     ' Check attack timer
//     If GetTickCount < TempPlayer(Attacker).AttackTimer + 1000 Then Exit Function

//     ' Check for subscript out of range
//     If Not IsPlaying(Victim) Then Exit Function

//     ' Make sure they are on the same map
//     If Not GetPlayerMap(Attacker) = GetPlayerMap(Victim) Then Exit Function

//     ' Make sure we dont attack the player if they are switching maps
//     If TempPlayer(Victim).GettingMap = YES Then Exit Function

//     ' Check if at same coordinates
//     Select Case GetPlayerDir(Attacker)
//     Case DIR_UP
//         If Not ((GetPlayerY(Victim) + 1 = GetPlayerY(Attacker)) And (GetPlayerX(Victim) = GetPlayerX(Attacker))) Then Exit Function
//     Case DIR_DOWN
//         If Not ((GetPlayerY(Victim) - 1 = GetPlayerY(Attacker)) And (GetPlayerX(Victim) = GetPlayerX(Attacker))) Then Exit Function
//     Case DIR_LEFT
//         If Not ((GetPlayerY(Victim) = GetPlayerY(Attacker)) And (GetPlayerX(Victim) + 1 = GetPlayerX(Attacker))) Then Exit Function
//     Case DIR_RIGHT
//         If Not ((GetPlayerY(Victim) = GetPlayerY(Attacker)) And (GetPlayerX(Victim) - 1 = GetPlayerX(Attacker))) Then Exit Function
//     Case Else
//         Exit Function
//     End Select

//     ' Check if map is attackable
//     If (Not Map(GetPlayerMap(Attacker)).Moral = MAP_MORAL_NONE) Or (Not Map(GetPlayerMap(Attacker)).Moral = MAP_MORAL_ARENA) Then
//         If GetPlayerPK(Victim) = NO Then
//             Call PlayerMsg(Attacker, "This is a safe zone!", BrightRed)
//             Exit Function
//         End If
//     End If

//     ' Make sure they have more then 0 hp
//     If GetPlayerVital(Victim, Vitals.HP) <= 0 Then Exit Function

//     ' Check to make sure that they dont have access
//     If GetPlayerAccess(Attacker) > ADMIN_MONITOR Then
//         Call PlayerMsg(Attacker, "You cannot attack any player for thou art an admin!", BrightBlue)
//         Exit Function
//     End If

//     ' Check to make sure the victim isn't an admin
//     If GetPlayerAccess(Victim) > ADMIN_MONITOR Then
//         Call PlayerMsg(Attacker, "You cannot attack " & GetPlayerName(Victim) & "!", BrightRed)
//         Exit Function
//     End If

//     ' Make sure attacker is high enough level
//     If GetPlayerLevel(Attacker) < 10 Then
//         Call PlayerMsg(Attacker, "You are below level 10, you cannot attack another player yet!", BrightRed)
//         Exit Function
//     End If

//     ' Make sure victim is high enough level
//     If GetPlayerLevel(Victim) < 10 Then
//         Call PlayerMsg(Attacker, GetPlayerName(Victim) & " is below level 10, you cannot attack this player yet!", BrightRed)
//         Exit Function
//     End If

//     CanAttackPlayer = True

// End Function

// Public Function CanAttackNpc(ByVal Attacker As Long, ByVal MapNpcNum As Long) As Boolean
//     Dim MapNum As Long
//     Dim NpcNum As Long
//     Dim NpcX As Long
//     Dim NpcY As Long

//     ' Check for subscript out of range
//     If IsPlaying(Attacker) = False Or MapNpcNum <= 0 Or MapNpcNum > MAX_MAP_NPCS Then
//         Exit Function
//     End If

//     ' Check for subscript out of range
//     If MapNpc(GetPlayerMap(Attacker), MapNpcNum).Num <= 0 Then
//         Exit Function
//     End If

//     MapNum = GetPlayerMap(Attacker)
//     NpcNum = MapNpc(MapNum, MapNpcNum).Num

//     ' Make sure the Npc isn't already dead
//     If MapNpc(MapNum, MapNpcNum).Vital(Vitals.HP) <= 0 Then
//         Exit Function
//     End If

//     ' Make sure they are on the same map
//     If IsPlaying(Attacker) Then
//         If NpcNum > 0 And GetTickCount > TempPlayer(Attacker).AttackTimer + 1000 Then
//             ' Check if at same coordinates
//             Select Case GetPlayerDir(Attacker)
//             Case DIR_UP
//                 NpcX = MapNpc(MapNum, MapNpcNum).X
//                 NpcY = MapNpc(MapNum, MapNpcNum).Y + 1
//             Case DIR_DOWN
//                 NpcX = MapNpc(MapNum, MapNpcNum).X
//                 NpcY = MapNpc(MapNum, MapNpcNum).Y - 1
//             Case DIR_LEFT
//                 NpcX = MapNpc(MapNum, MapNpcNum).X + 1
//                 NpcY = MapNpc(MapNum, MapNpcNum).Y
//             Case DIR_RIGHT
//                 NpcX = MapNpc(MapNum, MapNpcNum).X - 1
//                 NpcY = MapNpc(MapNum, MapNpcNum).Y
//             End Select

//             If NpcX = GetPlayerX(Attacker) Then
//                 If NpcY = GetPlayerY(Attacker) Then
//                     If Npc(NpcNum).Behavior <> Npc_BEHAVIOR_FRIENDLY And Npc(NpcNum).Behavior <> Npc_BEHAVIOR_SHOPKEEPER Then
//                         CanAttackNpc = True
//                     Else
//                         Call PlayerMsg(Attacker, "You cannot attack a " & Trim$(Npc(NpcNum).Name) & "!", BrightBlue)
//                     End If
//                 End If
//             End If
//         End If
//     End If
// End Function

// Public Function CanNpcAttackPlayer(ByVal MapNpcNum As Long, ByVal Index As Long) As Boolean
//     Dim MapNum As Long
//     Dim NpcNum As Long

//     ' Check for subscript out of range
//     If MapNpcNum <= 0 Or MapNpcNum > MAX_MAP_NPCS Or Not IsPlaying(Index) Then
//         Exit Function
//     End If

//     ' Check for subscript out of range
//     If MapNpc(GetPlayerMap(Index), MapNpcNum).Num <= 0 Then
//         Exit Function
//     End If

//     MapNum = GetPlayerMap(Index)
//     NpcNum = MapNpc(MapNum, MapNpcNum).Num

//     ' Make sure the Npc isn't already dead
//     If MapNpc(MapNum, MapNpcNum).Vital(Vitals.HP) <= 0 Then
//         Exit Function
//     End If

//     ' Make sure Npcs dont attack more then once a second
//     If GetTickCount < MapNpc(MapNum, MapNpcNum).AttackTimer + 1000 Then
//         Exit Function
//     End If

//     ' Make sure we dont attack the player if they are switching maps
//     If TempPlayer(Index).GettingMap = YES Then
//         Exit Function
//     End If

//     MapNpc(MapNum, MapNpcNum).AttackTimer = GetTickCount

//     ' Make sure they are on the same map
//     If IsPlaying(Index) Then
//         If NpcNum > 0 Then
//             ' Check if at same coordinates
//             If (GetPlayerY(Index) + 1 = MapNpc(MapNum, MapNpcNum).Y) And (GetPlayerX(Index) = MapNpc(MapNum, MapNpcNum).X) Then
//                 CanNpcAttackPlayer = True
//             Else
//                 If (GetPlayerY(Index) - 1 = MapNpc(MapNum, MapNpcNum).Y) And (GetPlayerX(Index) = MapNpc(MapNum, MapNpcNum).X) Then
//                     CanNpcAttackPlayer = True
//                 Else
//                     If (GetPlayerY(Index) = MapNpc(MapNum, MapNpcNum).Y) And (GetPlayerX(Index) + 1 = MapNpc(MapNum, MapNpcNum).X) Then
//                         CanNpcAttackPlayer = True
//                     Else
//                         If (GetPlayerY(Index) = MapNpc(MapNum, MapNpcNum).Y) And (GetPlayerX(Index) - 1 = MapNpc(MapNum, MapNpcNum).X) Then
//                             CanNpcAttackPlayer = True
//                         End If
//                     End If
//                 End If
//             End If

//             '            Select Case MapNpc(MapNum, MapNpcNum).Dir
//             '                Case DIR_UP
//             '                    If (GetPlayerY(Index) + 1 = MapNpc(MapNum, MapNpcNum).y) And (GetPlayerX(Index) = MapNpc(MapNum, MapNpcNum).x) Then
//             '                        CanNpcAttackPlayer = True
//             '                    End If
//             '
//             '                Case DIR_DOWN
//             '                    If (GetPlayerY(Index) - 1 = MapNpc(MapNum, MapNpcNum).y) And (GetPlayerX(Index) = MapNpc(MapNum, MapNpcNum).x) Then
//             '                        CanNpcAttackPlayer = True
//             '                    End If
//             '
//             '                Case DIR_LEFT
//             '                    If (GetPlayerY(Index) = MapNpc(MapNum, MapNpcNum).y) And (GetPlayerX(Index) + 1 = MapNpc(MapNum, MapNpcNum).x) Then
//             '                        CanNpcAttackPlayer = True
//             '                    End If
//             '
//             '                Case DIR_RIGHT
//             '                    If (GetPlayerY(Index) = MapNpc(MapNum, MapNpcNum).y) And (GetPlayerX(Index) - 1 = MapNpc(MapNum, MapNpcNum).x) Then
//             '                        CanNpcAttackPlayer = True
//             '                    End If
//             '            End Select

//         End If
//     End If
// End Function

// Public Sub NpcAttackPlayer(ByVal MapNpcNum As Long, ByVal Victim As Long, ByVal Damage As Long)
//     Dim Name As String
//     Dim Exp As Long
//     Dim MapNum As Long
//     Dim I As Long
//     Dim Buffer As clsBuffer

//     ' Check for subscript out of range
//     If MapNpcNum <= 0 Or MapNpcNum > MAX_MAP_NPCS Or IsPlaying(Victim) = False Or Damage < 0 Then
//         Exit Sub
//     End If

//     ' Check for subscript out of range
//     If MapNpc(GetPlayerMap(Victim), MapNpcNum).Num <= 0 Then
//         Exit Sub
//     End If

//     MapNum = GetPlayerMap(Victim)
//     Name = Trim$(Npc(MapNpc(MapNum, MapNpcNum).Num).Name)

//     ' Send this packet so they can see the person attacking
//     Set Buffer = New clsBuffer
//     Buffer.PreAllocate 6
//     Buffer.WriteInteger SNpcAttack
//     Buffer.WriteLong MapNpcNum
//     Call SendDataToMap(MapNum, Buffer.ToArray())

//     ' reduce dur. on victims equipment
//     Call DamageEquipment(Victim, Armor)
//     Call DamageEquipment(Victim, Helmet)

//     If Damage >= GetPlayerVital(Victim, Vitals.HP) Then
//         ' Say damage
//         Call PlayerMsg(Victim, "A " & Name & " hit you for " & Damage & " hit points.", BrightRed)

//         ' Player is dead
//         Call GlobalMsg(GetPlayerName(Victim) & " has been killed by a " & Name, BrightRed)

//         ' Calculate exp to give attacker
//         Exp = GetPlayerExp(Victim) \ 3

//         ' Make sure we dont get less then 0
//         If Exp < 0 Then Exp = 0

//         If Exp = 0 Then
//             Call PlayerMsg(Victim, "You lost no experience points.", BrightRed)
//         Else
//             Call SetPlayerExp(Victim, GetPlayerExp(Victim) - Exp)
//             Call PlayerMsg(Victim, "You lost " & Exp & " experience points.", BrightRed)
//         End If

//         ' Set Npc target to 0
//         MapNpc(MapNum, MapNpcNum).Target = 0

//         Call OnDeath(Victim)
//     Else
//         ' Player not dead, just do the damage
//         Call SetPlayerVital(Victim, Vitals.HP, GetPlayerVital(Victim, Vitals.HP) - Damage)
//         Call SendVital(Victim, Vitals.HP)

//         ' Say damage
//         Call PlayerMsg(Victim, "A " & Name & " hit you for " & Damage & " hit points.", BrightRed)
//     End If
// End Sub

// Public Function CanNpcMove(ByVal MapNum As Long, ByVal MapNpcNum As Long, ByVal Dir As Byte) As Boolean
//     Dim I As Long
//     Dim n As Long
//     Dim X As Long
//     Dim Y As Long

//     ' Check for subscript out of range
//     If MapNum <= 0 Or MapNum > MAX_MAPS Or MapNpcNum <= 0 Or MapNpcNum > MAX_MAP_NPCS Or Dir < DIR_UP Or Dir > DIR_RIGHT Then
//         Exit Function
//     End If

//     X = MapNpc(MapNum, MapNpcNum).X
//     Y = MapNpc(MapNum, MapNpcNum).Y

//     CanNpcMove = True

//     Select Case Dir
//     Case DIR_UP
//         ' Check to make sure not outside of boundries
//         If Y > 0 Then
//             n = Map(MapNum).Tile(X, Y - 1).Type

//             ' Check to make sure that the tile is walkable
//             If n <> TILE_TYPE_WALKABLE Then
//                 If n <> TILE_TYPE_ITEM Then
//                     CanNpcMove = False
//                     Exit Function
//                 End If
//             End If

//             ' Check to make sure that there is not a player in the way
//             For I = 1 To TotalPlayersOnline
//                 If (GetPlayerMap(PlayersOnline(I)) = MapNum) Then
//                     If (GetPlayerX(PlayersOnline(I)) = MapNpc(MapNum, MapNpcNum).X) Then
//                         If (GetPlayerY(PlayersOnline(I)) = MapNpc(MapNum, MapNpcNum).Y - 1) Then
//                             CanNpcMove = False
//                             Exit Function
//                         End If
//                     End If
//                 End If
//             Next

//             ' Check to make sure that there is not another Npc in the way
//             For I = 1 To MAX_MAP_NPCS
//                 If (I <> MapNpcNum) Then
//                     If (MapNpc(MapNum, I).Num > 0) Then
//                         If (MapNpc(MapNum, I).X = MapNpc(MapNum, MapNpcNum).X) Then
//                             If (MapNpc(MapNum, I).Y = MapNpc(MapNum, MapNpcNum).Y - 1) Then
//                                 CanNpcMove = False
//                                 Exit Function
//                             End If
//                         End If
//                     End If
//                 End If
//             Next
//         Else
//             CanNpcMove = False
//         End If

//     Case DIR_DOWN
//         ' Check to make sure not outside of boundries
//         If Y < MAX_MAPY Then
//             n = Map(MapNum).Tile(X, Y + 1).Type

//             ' Check to make sure that the tile is walkable
//             If n <> TILE_TYPE_WALKABLE Then
//                 If n <> TILE_TYPE_ITEM Then
//                     CanNpcMove = False
//                     Exit Function
//                 End If
//             End If

//             ' Check to make sure that there is not a player in the way
//             For I = 1 To TotalPlayersOnline
//                 If (GetPlayerMap(PlayersOnline(I)) = MapNum) Then
//                     If (GetPlayerX(PlayersOnline(I)) = MapNpc(MapNum, MapNpcNum).X) Then
//                         If (GetPlayerY(PlayersOnline(I)) = MapNpc(MapNum, MapNpcNum).Y + 1) Then
//                             CanNpcMove = False
//                             Exit Function
//                         End If
//                     End If
//                 End If
//             Next

//             ' Check to make sure that there is not another Npc in the way
//             For I = 1 To MAX_MAP_NPCS
//                 If (I <> MapNpcNum) Then
//                     If (MapNpc(MapNum, I).Num > 0) Then
//                         If (MapNpc(MapNum, I).X = MapNpc(MapNum, MapNpcNum).X) Then
//                             If (MapNpc(MapNum, I).Y = MapNpc(MapNum, MapNpcNum).Y + 1) Then
//                                 CanNpcMove = False
//                                 Exit Function
//                             End If
//                         End If
//                     End If
//                 End If
//             Next
//         Else
//             CanNpcMove = False
//         End If

//     Case DIR_LEFT
//         ' Check to make sure not outside of boundries
//         If X > 0 Then
//             n = Map(MapNum).Tile(X - 1, Y).Type

//             ' Check to make sure that the tile is walkable
//             If n <> TILE_TYPE_WALKABLE Then
//                 If n <> TILE_TYPE_ITEM Then
//                     CanNpcMove = False
//                     Exit Function
//                 End If
//             End If

//             ' Check to make sure that there is not a player in the way
//             For I = 1 To TotalPlayersOnline
//                 If (GetPlayerMap(PlayersOnline(I)) = MapNum) Then
//                     If (GetPlayerX(PlayersOnline(I)) = MapNpc(MapNum, MapNpcNum).X - 1) Then
//                         If (GetPlayerY(PlayersOnline(I)) = MapNpc(MapNum, MapNpcNum).Y) Then
//                             CanNpcMove = False
//                             Exit Function
//                         End If
//                     End If
//                 End If
//             Next

//             ' Check to make sure that there is not another Npc in the way
//             For I = 1 To MAX_MAP_NPCS
//                 If (I <> MapNpcNum) Then
//                     If (MapNpc(MapNum, I).Num > 0) Then
//                         If (MapNpc(MapNum, I).X = MapNpc(MapNum, MapNpcNum).X - 1) Then
//                             If (MapNpc(MapNum, I).Y = MapNpc(MapNum, MapNpcNum).Y) Then
//                                 CanNpcMove = False
//                                 Exit Function
//                             End If
//                         End If
//                     End If
//                 End If
//             Next
//         Else
//             CanNpcMove = False
//         End If

//     Case DIR_RIGHT
//         ' Check to make sure not outside of boundries
//         If X < MAX_MAPX Then
//             n = Map(MapNum).Tile(X + 1, Y).Type

//             ' Check to make sure that the tile is walkable
//             If n <> TILE_TYPE_WALKABLE Then
//                 If n <> TILE_TYPE_ITEM Then
//                     CanNpcMove = False
//                     Exit Function
//                 End If
//             End If

//             ' Check to make sure that there is not a player in the way
//             For I = 1 To TotalPlayersOnline
//                 If (GetPlayerMap(PlayersOnline(I)) = MapNum) Then
//                     If (GetPlayerX(PlayersOnline(I)) = MapNpc(MapNum, MapNpcNum).X + 1) Then
//                         If (GetPlayerY(PlayersOnline(I)) = MapNpc(MapNum, MapNpcNum).Y) Then
//                             CanNpcMove = False
//                             Exit Function
//                         End If
//                     End If
//                 End If
//             Next

//             ' Check to make sure that there is not another Npc in the way
//             For I = 1 To MAX_MAP_NPCS
//                 If (I <> MapNpcNum) Then
//                     If (MapNpc(MapNum, I).Num > 0) Then
//                         If (MapNpc(MapNum, I).X = MapNpc(MapNum, MapNpcNum).X + 1) Then
//                             If (MapNpc(MapNum, I).Y = MapNpc(MapNum, MapNpcNum).Y) Then
//                                 CanNpcMove = False
//                                 Exit Function
//                             End If
//                         End If
//                     End If
//                 End If
//             Next
//         Else
//             CanNpcMove = False
//         End If
//     End Select
// End Function

// Public Sub NpcMove(ByVal MapNum As Long, ByVal MapNpcNum As Long, ByVal Dir As Long, ByVal Movement As Long)
//     Dim Buffer As clsBuffer

//     ' Check for subscript out of range
//     If MapNum <= 0 Or MapNum > MAX_MAPS Or MapNpcNum <= 0 Or MapNpcNum > MAX_MAP_NPCS Or Dir < DIR_UP Or Dir > DIR_RIGHT Or Movement < 1 Or Movement > 2 Then
//         Exit Sub
//     End If

//     MapNpc(MapNum, MapNpcNum).Dir = Dir

//     Select Case Dir
//     Case DIR_UP
//         MapNpc(MapNum, MapNpcNum).Y = MapNpc(MapNum, MapNpcNum).Y - 1

//     Case DIR_DOWN
//         MapNpc(MapNum, MapNpcNum).Y = MapNpc(MapNum, MapNpcNum).Y + 1

//     Case DIR_LEFT
//         MapNpc(MapNum, MapNpcNum).X = MapNpc(MapNum, MapNpcNum).X - 1

//     Case DIR_RIGHT
//         MapNpc(MapNum, MapNpcNum).X = MapNpc(MapNum, MapNpcNum).X + 1
//     End Select

//     Set Buffer = New clsBuffer

//     Buffer.PreAllocate 12 + 4
//     Buffer.WriteInteger SNpcMove
//     Buffer.WriteLong MapNum
//     Buffer.WriteInteger MapNpcNum
//     Buffer.WriteByte MapNpc(MapNum, MapNpcNum).X
//     Buffer.WriteByte MapNpc(MapNum, MapNpcNum).Y
//     Buffer.WriteInteger MapNpc(MapNum, MapNpcNum).Dir
//     Buffer.WriteLong Movement

//     Call SendDataToAll(Buffer.ToArray())
// End Sub

// Public Sub NpcDir(ByVal MapNum As Long, ByVal MapNpcNum As Long, ByVal Dir As Long)
//     Dim Buffer As clsBuffer

//     ' Check for subscript out of range
//     If MapNum <= 0 Or MapNum > MAX_MAPS Or MapNpcNum <= 0 Or MapNpcNum > MAX_MAP_NPCS Or Dir < DIR_UP Or Dir > DIR_RIGHT Then
//         Exit Sub
//     End If

//     Set Buffer = New clsBuffer

//     Buffer.PreAllocate 12 + 4
//     Buffer.WriteInteger SNpcDir
//     Buffer.WriteLong MapNum
//     Buffer.WriteInteger MapNpcNum
//     Buffer.WriteLong Dir

//     Call SendDataToAll(Buffer.ToArray())
// End Sub

// Public Function GetTotalMapPlayers(ByVal MapNum As Long) As Long
//     Dim I As Long
//     Dim n As Long

//     n = 0

//     For I = 1 To High_Index
//         If IsPlaying(I) Then
//             If GetPlayerMap(I) = MapNum Then
//                 n = n + 1
//             End If
//         End If
//     Next

//     GetTotalMapPlayers = n
// End Function

// Public Function GetNpcMaxVital(ByVal NpcNum As Long, ByVal Vital As Vitals) As Long
//     Dim X As Long
//     Dim Y As Long

//     ' Prevent subscript out of range
//     If NpcNum <= 0 Or NpcNum > MAX_NPCS Then
//         GetNpcMaxVital = 0
//         Exit Function
//     End If

//     Select Case Vital
//     Case HP
//         X = Npc(NpcNum).Stat(Stats.Strength)
//         Y = Npc(NpcNum).Stat(Stats.Defense)
//         GetNpcMaxVital = X * Y
//     Case MP
//         GetNpcMaxVital = Npc(NpcNum).Stat(Stats.Magic) * 2
//     Case SP
//         GetNpcMaxVital = Npc(NpcNum).Stat(Stats.Speed) * 2
//     End Select
// End Function

// Public Function GetNpcVitalRegen(ByVal NpcNum As Long, ByVal Vital As Vitals) As Long
//     Dim I As Long

//     'Prevent subscript out of range
//     If NpcNum <= 0 Or NpcNum > MAX_NPCS Then
//         GetNpcVitalRegen = 0
//         Exit Function
//     End If

//     Select Case Vital
//     Case HP
//         I = Npc(NpcNum).Stat(Stats.Defense) \ 3
//         If I < 1 Then I = 1
//         GetNpcVitalRegen = I
//         'Case MP

//         'Case SP

//     End Select
// End Function

// Public Sub ClearTempTile()
//     Dim I As Long
//     Dim Y As Long
//     Dim X As Long

//     For I = 1 To MAX_MAPS
//         TempTile(I).DoorTimer = 0

//         For X = 0 To MAX_MAPX
//             For Y = 0 To MAX_MAPY
//                 TempTile(I).DoorOpen(X, Y) = NO
//             Next
//         Next

//     Next
// End Sub

// Public Function GetPlayerDamage(ByVal Index As Long) As Long
//     Dim WeaponSlot As Long

//     GetPlayerDamage = 0

//     ' Check for subscript out of range
//     If IsPlaying(Index) = False Or Index <= 0 Or Index > High_Index Then
//         Exit Function
//     End If

//     GetPlayerDamage = (GetPlayerStat(Index, Stats.Strength) \ 2)

//     If GetPlayerDamage <= 0 Then
//         GetPlayerDamage = 1
//     End If

//     If GetPlayerEquipmentSlot(Index, Weapon) > 0 Then
//         WeaponSlot = GetPlayerEquipmentSlot(Index, Weapon)

//         GetPlayerDamage = GetPlayerDamage + Item(GetPlayerInvItemNum(Index, WeaponSlot)).Data2
//     End If
// End Function

// Public Function GetPlayerProtection(ByVal Index As Long) As Long
//     Dim ArmorSlot As Long
//     Dim HelmSlot As Long

//     GetPlayerProtection = 0

//     ' Check for subscript out of range
//     If IsPlaying(Index) = False Or Index <= 0 Or Index > High_Index Then
//         Exit Function
//     End If

//     ArmorSlot = GetPlayerEquipmentSlot(Index, Armor)
//     HelmSlot = GetPlayerEquipmentSlot(Index, Helmet)

//     GetPlayerProtection = (GetPlayerStat(Index, Stats.Defense) \ 5)

//     If ArmorSlot > 0 Then
//         GetPlayerProtection = GetPlayerProtection + Item(GetPlayerInvItemNum(Index, ArmorSlot)).Data2
//     End If

//     If HelmSlot > 0 Then
//         GetPlayerProtection = GetPlayerProtection + Item(GetPlayerInvItemNum(Index, HelmSlot)).Data2
//     End If
// End Function

// Public Function CanPlayerCriticalHit(ByVal Index As Long) As Boolean
//     Dim I As Long
//     Dim n As Long

//     If GetPlayerEquipmentSlot(Index, Weapon) > 0 Then
//         n = Int(Rnd * 2)
//         If n = 1 Then
//             I = (GetPlayerStat(Index, Stats.Strength) \ 2) + (GetPlayerLevel(Index) \ 2)

//             n = Int(Rnd * 100) + 1
//             If n <= I Then
//                 CanPlayerCriticalHit = True
//             End If
//         End If
//     End If
// End Function

// Public Function CanPlayerBlockHit(ByVal Index As Long) As Boolean
//     Dim I As Long
//     Dim n As Long
//     Dim ShieldSlot As Long

//     ShieldSlot = GetPlayerEquipmentSlot(Index, Shield)

//     If ShieldSlot > 0 Then
//         n = Int(Rnd * 2)
//         If n = 1 Then
//             I = (GetPlayerStat(Index, Stats.Defense) \ 2) + (GetPlayerLevel(Index) \ 2)

//             n = Int(Rnd * 100) + 1
//             If n <= I Then
//                 CanPlayerBlockHit = True
//             End If
//         End If
//     End If
// End Function

// Public Sub CastSpell(ByVal Index As Long, ByVal SpellSlot As Long)
//     Dim SpellNum As Long
//     Dim MPReq As Long
//     Dim I As Long
//     Dim n As Long
//     Dim Damage As Long
//     Dim Casted As Boolean
//     Dim CanCast As Boolean
//     Dim TargetType As Byte
//     Dim TargetName As String
//     Dim Buffer As clsBuffer

//     ' Prevent subscript out of range
//     If SpellSlot <= 0 Or SpellSlot > MAX_PLAYER_SPELLS Then
//         Exit Sub
//     End If

//     SpellNum = GetPlayerSpell(Index, SpellSlot)

//     ' Make sure player has the spell
//     If Not HasSpell(Index, SpellNum) Then
//         Call PlayerMsg(Index, "You do not have this spell!", BrightRed)
//         Exit Sub
//     End If

//     ' (does not check for level requirement)
//     ' Make sure they are the right level
//     'If ?? > GetPlayerLevel(Index) Then
//     '    Call PlayerMsg(Index, "You must be level " & ??? & " to cast this spell.", BrightRed)
//     '    Exit Sub
//     'End If

//     MPReq = Spell(SpellNum).MPReq

//     ' Check if they have enough MP
//     If GetPlayerVital(Index, Vitals.MP) < MPReq Then
//         Call PlayerMsg(Index, "Not enough mana points!", BrightRed)
//         Exit Sub
//     End If

//     ' Check if timer is ok
//     If GetTickCount < TempPlayer(Index).AttackTimer + 1000 Then
//         Exit Sub
//     End If

//     ' *** Self Cast Spells ***
//     ' Check if the spell is a give item and do that instead of a stat modification
//     If Spell(SpellNum).Type = SPELL_TYPE_GIVEITEM Then
//         n = FindOpenInvSlot(Index, Spell(SpellNum).Data1)

//         If n > 0 Then
//             Call GiveItem(Index, Spell(SpellNum).Data1, Spell(SpellNum).Data2)
//             Call MapMsg(GetPlayerMap(Index), GetPlayerName(Index) & " casts " & Trim$(Spell(SpellNum).Name) & ".", BrightBlue)

//             ' Take away the mana points
//             Call SetPlayerVital(Index, Vitals.MP, GetPlayerVital(Index, Vitals.MP) - MPReq)
//             Call SendVital(Index, Vitals.MP)
//             Casted = True
//         Else
//             Call PlayerMsg(Index, "Your inventory is full!", BrightRed)
//         End If

//         Exit Sub
//     End If

//     n = TempPlayer(Index).Target
//     TargetType = TempPlayer(Index).TargetType

//     Select Case TargetType
//     Case TARGET_TYPE_PLAYER

//         If IsPlaying(n) Then

//             If GetPlayerVital(n, Vitals.HP) > 0 Then
//                 If GetPlayerMap(Index) = GetPlayerMap(n) Then
//                     'If GetPlayerLevel(Index) >= 10 Then
//                     'If GetPlayerLevel(n) >= 10 Then
//                     If (Map(GetPlayerMap(Index)).Moral = MAP_MORAL_NONE) Or (Map(GetPlayerMap(Index)).Moral = MAP_MORAL_ARENA) Then
//                         If GetPlayerAccess(Index) <= 0 Then
//                             If GetPlayerAccess(n) <= 0 Then
//                                 If n <> Index Then
//                                     CanCast = True
//                                 End If
//                             End If
//                         End If
//                     End If
//                     'End If
//                     'End If
//                 End If
//             End If

//             TargetName = GetPlayerName(n)

//             If Spell(SpellNum).Type = SPELL_TYPE_SUBHP Or _
//             Spell(SpellNum).Type = SPELL_TYPE_SUBMP Or _
//             Spell(SpellNum).Type = SPELL_TYPE_SUBSP Then

//             If CanCast Then
//                 Select Case Spell(SpellNum).Type
//                 Case SPELL_TYPE_SUBHP
//                     Damage = (GetPlayerStat(Index, Stats.Magic) \ 4) + Spell(SpellNum).Data1 - GetPlayerProtection(n)
//                     If Damage > 0 Then
//                         Call AttackPlayer(Index, n, Damage)
//                     Else
//                         Call PlayerMsg(Index, "The spell was to weak to hurt " & GetPlayerName(n) & "!", BrightRed)
//                     End If

//                 Case SPELL_TYPE_SUBMP
//                     Call SetPlayerVital(n, Vitals.MP, GetPlayerVital(n, Vitals.MP) - Spell(SpellNum).Data1)
//                     Call SendVital(n, Vitals.MP)

//                 Case SPELL_TYPE_SUBSP
//                     Call SetPlayerVital(n, Vitals.SP, GetPlayerVital(n, Vitals.SP) - Spell(SpellNum).Data1)
//                     Call SendVital(n, Vitals.SP)
//                 End Select

//                 Casted = True

//             End If

//         ElseIf Spell(SpellNum).Type = SPELL_TYPE_ADDHP Or _
//         Spell(SpellNum).Type = SPELL_TYPE_ADDMP Or _
//         Spell(SpellNum).Type = SPELL_TYPE_ADDSP Then

//         If GetPlayerMap(Index) = GetPlayerMap(n) Then
//             CanCast = True
//         End If

//         If CanCast Then
//             Select Case Spell(SpellNum).Type
//             Case SPELL_TYPE_ADDHP
//                 Call SetPlayerVital(n, Vitals.HP, GetPlayerVital(n, Vitals.HP) + Spell(SpellNum).Data1)
//                 Call SendVital(n, Vitals.HP)

//             Case SPELL_TYPE_ADDMP
//                 Call SetPlayerVital(n, Vitals.MP, GetPlayerVital(n, Vitals.MP) + Spell(SpellNum).Data1)
//                 Call SendVital(n, Vitals.MP)

//             Case SPELL_TYPE_ADDSP
//                 Call SetPlayerVital(n, Vitals.SP, GetPlayerVital(n, Vitals.SP) + Spell(SpellNum).Data1)
//                 Call SendVital(n, Vitals.SP)
//             End Select

//             Casted = True
//         End If

//     End If
// End If

// Case TARGET_TYPE_NPC

//     If Npc(MapNpc(GetPlayerMap(Index), n).Num).Behavior <> Npc_BEHAVIOR_FRIENDLY Then
//         If Npc(MapNpc(GetPlayerMap(Index), n).Num).Behavior <> Npc_BEHAVIOR_SHOPKEEPER Then
//             CanCast = True
//         End If
//     End If

//     TargetName = Npc(MapNpc(GetPlayerMap(Index), n).Num).Name

//     If CanCast Then
//         Select Case Spell(SpellNum).Type
//         Case SPELL_TYPE_ADDHP
//             MapNpc(GetPlayerMap(Index), n).Vital(Vitals.HP) = MapNpc(GetPlayerMap(Index), n).Vital(Vitals.HP) + Spell(SpellNum).Data1

//         Case SPELL_TYPE_SUBHP

//             Damage = (GetPlayerStat(Index, Stats.Magic) \ 4) + Spell(SpellNum).Data1 - (Npc(MapNpc(GetPlayerMap(Index), n).Num).Stat(Stats.Defense) \ 2)
//             If Damage > 0 Then
//                 Call AttackNpc(Index, n, Damage)
//             Else
//                 Call PlayerMsg(Index, "The spell was to weak to hurt " & Trim$(Npc(MapNpc(GetPlayerMap(Index), n).Num).Name) & "!", BrightRed)
//             End If

//         Case SPELL_TYPE_ADDMP
//             MapNpc(GetPlayerMap(Index), n).Vital(Vitals.MP) = MapNpc(GetPlayerMap(Index), n).Vital(Vitals.MP) + Spell(SpellNum).Data1

//         Case SPELL_TYPE_SUBMP
//             MapNpc(GetPlayerMap(Index), n).Vital(Vitals.MP) = MapNpc(GetPlayerMap(Index), n).Vital(Vitals.MP) - Spell(SpellNum).Data1

//         Case SPELL_TYPE_ADDSP
//             MapNpc(GetPlayerMap(Index), n).Vital(Vitals.SP) = MapNpc(GetPlayerMap(Index), n).Vital(Vitals.SP) + Spell(SpellNum).Data1

//         Case SPELL_TYPE_SUBSP
//             MapNpc(GetPlayerMap(Index), n).Vital(Vitals.SP) = MapNpc(GetPlayerMap(Index), n).Vital(Vitals.SP) - Spell(SpellNum).Data1
//         End Select

//         Casted = True
//     End If

// End Select

// If Casted Then
//     Call MapMsg(GetPlayerMap(Index), GetPlayerName(Index) & " casts " & Trim$(Spell(SpellNum).Name) & " on " & Trim$(TargetName) & ".", BrightBlue)

//     Set Buffer = New clsBuffer
//     Buffer.PreAllocate 11
//     Buffer.WriteInteger SCastSpell
//     Buffer.WriteByte TargetType
//     Buffer.WriteLong n
//     Buffer.WriteLong SpellNum
//     Call SendDataToMap(GetPlayerMap(Index), Buffer.ToArray())

//     ' Take away the mana points
//     Call SetPlayerVital(Index, Vitals.MP, GetPlayerVital(Index, Vitals.MP) - MPReq)
//     Call SendVital(Index, Vitals.MP)

//     TempPlayer(Index).AttackTimer = GetTickCount
//     TempPlayer(Index).CastedSpell = YES
// Else
//     Call PlayerMsg(Index, "Could not cast spell!", BrightRed)
// End If

// End Sub

// Public Sub PlayerChangeMap(ByVal Index As Long, ByVal MapNum As Long, ByVal X As Long, ByVal Y As Long)
//     Dim ShopNum As Long
//     Dim OldMap As Long
//     Dim I As Long
//     Dim Buffer As clsBuffer

//     ' Check for subscript out of range
//     If IsPlaying(Index) = False Or MapNum <= 0 Or MapNum > MAX_MAPS Then
//         Exit Sub
//     End If

//     TempPlayer(Index).Target = 0
//     TempPlayer(Index).TargetType = TARGET_TYPE_NONE

//     ' Check if there was a shop on the map the player is leaving, and if so say goodbye
//     ShopNum = Map(GetPlayerMap(Index)).Shop
//     If ShopNum > 0 Then
//         If LenB(Trim$(Shop(ShopNum).LeaveSay)) > 0 Then
//             Call PlayerMsg(Index, Trim$(Shop(ShopNum).Name) & " says, '" & Trim$(Shop(ShopNum).LeaveSay) & "'", SayColor)
//         End If
//     End If

//     ' Save old map to send erase player data to
//     OldMap = GetPlayerMap(Index)

//     If OldMap <> MapNum Then
//         'Call SendLeaveMap(Index, OldMap)
//     End If

//     Call SetPlayerMap(Index, MapNum)
//     Call SetPlayerX(Index, X)
//     Call SetPlayerY(Index, Y)

//     Call SendDataToAllBut(Index, PlayerData(Index))

//     For I = 1 To High_Index
//         If IsPlaying(I) And I <> Index Then
//             Call SendDataTo(Index, PlayerData(I))
//         End If
//     Next

//     ' Check if there is a shop on the map and say hello if so
//     ShopNum = Map(GetPlayerMap(Index)).Shop
//     If ShopNum > 0 Then
//         If LenB(Trim$(Shop(ShopNum).JoinSay)) > 0 Then
//             Call PlayerMsg(Index, Trim$(Shop(ShopNum).Name) & " says, '" & Trim$(Shop(ShopNum).JoinSay) & "'", SayColor)
//         End If
//     End If

//     ' Now we check if there were any players left on the map the player just left, and if not stop processing Npcs
//     If GetTotalMapPlayers(OldMap) = 0 Then
//         PlayersOnMap(OldMap) = NO

//         ' Regenerate all Npcs' health
//         For I = 1 To MAX_MAP_NPCS
//             If MapNpc(OldMap, I).Num > 0 Then
//                 MapNpc(OldMap, I).Vital(Vitals.HP) = GetNpcMaxVital(MapNpc(OldMap, I).Num, Vitals.HP)
//             End If
//         Next

//     End If

//     ' Sets it so we know to process Npcs on the map
//     PlayersOnMap(MapNum) = YES
// End Sub

// MovePlayerToRoom moves the player to the specified room and position.
func MovePlayerToRoom(player *PlayerData, roomId int, dx int, dy int) {
	if roomId < 0 || roomId >= config.MaxMaps {
		return
	}

	room := &rooms[roomId]
	if room == player.Room {
		return
	}

	room.AddPlayerAt(player, dx, dy)
}

// MovePlayer moves the player in the specified direction.
func MovePlayer(player *PlayerData, dir common.Direction, movement int) {
	if player.Room == nil || player.Character == nil {
		return
	}

	player.Character.Dir = dir

	dx, dy := utils.GetAdjacentTile(player.Character.X, player.Character.Y, dir)

	// If the player is trying to move out of bounds move them to the adjacent room
	if !player.Room.Level.Contains(dx, dy) {
		switch dir {
		case common.DirUp:
			MovePlayerToRoom(player, player.Room.Level.Up, dx, player.Room.Level.Height-1)
		case common.DirDown:
			MovePlayerToRoom(player, player.Room.Level.Down, dx, 0)
		case common.DirLeft:
			MovePlayerToRoom(player, player.Room.Level.Left, player.Room.Level.Width-1, dy)
		case common.DirRight:
			MovePlayerToRoom(player, player.Room.Level.Right, 0, dy)
		}
		return
	}

	// Move the player to the new position
	writer := net.NewWriter()

	writer.WriteInteger(SvPlayerMove)
	writer.WriteLong(player.Id + 1)
	writer.WriteLong(player.Character.X)
	writer.WriteLong(player.Character.Y)
	writer.WriteLong(int(player.Character.Dir))
	writer.WriteLong(movement)

	player.Room.SendExclude(writer.Bytes(), player)

	TriggerTileEffect(player)
}

// CheckEquippedItems checks wether the type of the items equipped by the specified player match the slots in which they are equipped.
// If the item type does not match the slot, the item is removed from that slot.
func CheckEquippedItems(p *PlayerData) {
	character := p.Character
	if character == nil {
		return
	}

	CheckSlot := func(itemId int, slot equipment.Slot) int {
		if itemId < 0 || itemId >= config.MaxItems {
			return -1
		}

		item := data.GetItem(itemId)
		if item == nil {
			return -1
		}

		switch slot {
		case equipment.Weapon:
			if item.Type != data.ItemWeapon {
				return -1
			}

		case equipment.Armor:
			if item.Type != data.ItemArmor {
				return -1
			}

		case equipment.Helmet:
			if item.Type != data.ItemHelmet {
				return -1
			}

		case equipment.Shield:
			if item.Type != data.ItemShield {
				return -1
			}
		}

		return itemId
	}

	character.Equipment.Weapon = CheckSlot(character.Equipment.Weapon, equipment.Weapon)
	character.Equipment.Armor = CheckSlot(character.Equipment.Armor, equipment.Armor)
	character.Equipment.Helmet = CheckSlot(character.Equipment.Helmet, equipment.Helmet)
	character.Equipment.Shield = CheckSlot(character.Equipment.Shield, equipment.Shield)
}

// Public Function FindOpenInvSlot(ByVal Index As Long, ByVal ItemNum As Long) As Long
//     Dim I As Long

//     ' Check for subscript out of range
//     If IsPlaying(Index) = False Or ItemNum <= 0 Or ItemNum > MAX_ITEMS Then
//         Exit Function
//     End If

//     If Item(ItemNum).Type = ITEM_TYPE_CURRENCY Then
//         ' If currency then check to see if they already have an instance of the item and add it to that
//         For I = 1 To MAX_INV
//             If GetPlayerInvItemNum(Index, I) = ItemNum Then
//                 FindOpenInvSlot = I
//                 Exit Function
//             End If
//         Next
//     End If

//     For I = 1 To MAX_INV
//         ' Try to find an open free slot
//         If GetPlayerInvItemNum(Index, I) = 0 Then
//             FindOpenInvSlot = I
//             Exit Function
//         End If
//     Next
// End Function

// Public Function HasItem(ByVal Index As Long, ByVal ItemNum As Long) As Long
//     Dim I As Long

//     ' Check for subscript out of range
//     If IsPlaying(Index) = False Or ItemNum <= 0 Or ItemNum > MAX_ITEMS Then
//         Exit Function
//     End If

//     For I = 1 To MAX_INV
//         ' Check to see if the player has the item
//         If GetPlayerInvItemNum(Index, I) = ItemNum Then
//             If Item(ItemNum).Type = ITEM_TYPE_CURRENCY Then
//                 HasItem = GetPlayerInvItemValue(Index, I)
//             Else
//                 HasItem = 1
//             End If
//             Exit Function
//         End If
//     Next
// End Function

// Public Sub TakeItem(ByVal Index As Long, ByVal ItemNum As Long, ByVal ItemVal As Long)
//     Dim I As Long
//     Dim n As Long
//     Dim TakeItem As Boolean

//     ' Check for subscript out of range
//     If IsPlaying(Index) = False Or ItemNum <= 0 Or ItemNum > MAX_ITEMS Then
//         Exit Sub
//     End If

//     For I = 1 To MAX_INV
//         ' Check to see if the player has the item
//         If GetPlayerInvItemNum(Index, I) = ItemNum Then
//             If Item(ItemNum).Type = ITEM_TYPE_CURRENCY Then
//                 ' Is what we are trying to take away more then what they have?  If so just set it to zero
//                 If ItemVal >= GetPlayerInvItemValue(Index, I) Then
//                     TakeItem = True
//                 Else
//                     Call SetPlayerInvItemValue(Index, I, GetPlayerInvItemValue(Index, I) - ItemVal)
//                     Call SendInventoryUpdate(Index, I)
//                 End If
//             Else
//                 ' Check to see if its any sort of ArmorSlot/WeaponSlot
//                 Select Case Item(GetPlayerInvItemNum(Index, I)).Type
//                 Case ITEM_TYPE_WEAPON
//                     If GetPlayerEquipmentSlot(Index, Weapon) > 0 Then
//                         If I = GetPlayerEquipmentSlot(Index, Weapon) Then
//                             Call SetPlayerEquipmentSlot(Index, 0, Weapon)
//                             Call SendWornEquipment(Index)
//                             TakeItem = True
//                         Else
//                             ' Check if the item we are taking isn't already equipped
//                             If ItemNum <> GetPlayerInvItemNum(Index, GetPlayerEquipmentSlot(Index, Weapon)) Then
//                                 TakeItem = True
//                             End If
//                         End If
//                     Else
//                         TakeItem = True
//                     End If

//                 Case ITEM_TYPE_ARMOR
//                     If GetPlayerEquipmentSlot(Index, Armor) > 0 Then
//                         If I = GetPlayerEquipmentSlot(Index, Armor) Then
//                             Call SetPlayerEquipmentSlot(Index, 0, Armor)
//                             Call SendWornEquipment(Index)
//                             TakeItem = True
//                         Else
//                             ' Check if the item we are taking isn't already equipped
//                             If ItemNum <> GetPlayerInvItemNum(Index, GetPlayerEquipmentSlot(Index, Armor)) Then
//                                 TakeItem = True
//                             End If
//                         End If
//                     Else
//                         TakeItem = True
//                     End If

//                 Case ITEM_TYPE_HELMET
//                     If GetPlayerEquipmentSlot(Index, Helmet) > 0 Then
//                         If I = GetPlayerEquipmentSlot(Index, Helmet) Then
//                             Call SetPlayerEquipmentSlot(Index, 0, Helmet)
//                             Call SendWornEquipment(Index)
//                             TakeItem = True
//                         Else
//                             ' Check if the item we are taking isn't already equipped
//                             If ItemNum <> GetPlayerInvItemNum(Index, GetPlayerEquipmentSlot(Index, Helmet)) Then
//                                 TakeItem = True
//                             End If
//                         End If
//                     Else
//                         TakeItem = True
//                     End If

//                 Case ITEM_TYPE_SHIELD
//                     If GetPlayerEquipmentSlot(Index, Shield) > 0 Then
//                         If I = GetPlayerEquipmentSlot(Index, Shield) Then
//                             Call SetPlayerEquipmentSlot(Index, 0, Shield)
//                             Call SendWornEquipment(Index)
//                             TakeItem = True
//                         Else
//                             ' Check if the item we are taking isn't already equipped
//                             If ItemNum <> GetPlayerInvItemNum(Index, GetPlayerEquipmentSlot(Index, Shield)) Then
//                                 TakeItem = True
//                             End If
//                         End If
//                     Else
//                         TakeItem = True
//                     End If
//                 End Select

//                 n = Item(GetPlayerInvItemNum(Index, I)).Type
//                 ' Check if its not an equipable weapon, and if it isn't then take it away
//                 If (n <> ITEM_TYPE_WEAPON) And (n <> ITEM_TYPE_ARMOR) And (n <> ITEM_TYPE_HELMET) And (n <> ITEM_TYPE_SHIELD) Then
//                     TakeItem = True
//                 End If
//             End If

//             If TakeItem Then
//                 Call SetPlayerInvItemNum(Index, I, 0)
//                 Call SetPlayerInvItemValue(Index, I, 0)
//                 Call SetPlayerInvItemDur(Index, I, 0)

//                 ' Send the inventory update
//                 Call SendInventoryUpdate(Index, I)
//                 Exit Sub
//             End If
//         End If
//     Next
// End Sub

// Public Sub GiveItem(ByVal Index As Long, ByVal ItemNum As Long, ByVal ItemVal As Long)
//     Dim I As Long

//     ' Check for subscript out of range
//     If IsPlaying(Index) = False Or ItemNum <= 0 Or ItemNum > MAX_ITEMS Then
//         Exit Sub
//     End If

//     I = FindOpenInvSlot(Index, ItemNum)

//     ' Check to see if inventory is full
//     If I <> 0 Then
//         Call SetPlayerInvItemNum(Index, I, ItemNum)
//         Call SetPlayerInvItemValue(Index, I, GetPlayerInvItemValue(Index, I) + ItemVal)

//         If (Item(ItemNum).Type = ITEM_TYPE_ARMOR) Or (Item(ItemNum).Type = ITEM_TYPE_WEAPON) Or (Item(ItemNum).Type = ITEM_TYPE_HELMET) Or (Item(ItemNum).Type = ITEM_TYPE_SHIELD) Then
//             Call SetPlayerInvItemDur(Index, I, Item(ItemNum).Data1)
//         End If

//         Call SendInventoryUpdate(Index, I)
//     Else
//         Call PlayerMsg(Index, "Your inventory is full.", BrightRed)
//     End If
// End Sub

// Public Function HasSpell(ByVal Index As Long, ByVal SpellNum As Long) As Boolean
//     Dim I As Long

//     For I = 1 To MAX_PLAYER_SPELLS
//         If GetPlayerSpell(Index, I) = SpellNum Then
//             HasSpell = True
//             Exit Function
//         End If
//     Next
// End Function

// Public Function FindOpenSpellSlot(ByVal Index As Long) As Long
//     Dim I As Long

//     For I = 1 To MAX_PLAYER_SPELLS
//         If GetPlayerSpell(Index, I) = 0 Then
//             FindOpenSpellSlot = I
//             Exit Function
//         End If
//     Next
// End Function

// Public Sub PlayerMapGetItem(ByVal Index As Long)
//     Dim I As Long
//     Dim n As Long
//     Dim MapNum As Long
//     Dim Msg As String

//     If Not IsPlaying(Index) Then Exit Sub

//     MapNum = GetPlayerMap(Index)

//     For I = 1 To MAX_MAP_ITEMS
//         ' See if theres even an item here
//         If (MapItem(MapNum, I).Num > 0) Then
//             If (MapItem(MapNum, I).Num <= MAX_ITEMS) Then

//                 ' Check if item is at the same location as the player
//                 If (MapItem(MapNum, I).X = GetPlayerX(Index)) Then

//                     If (MapItem(MapNum, I).Y = GetPlayerY(Index)) Then

//                         ' Find open slot
//                         n = FindOpenInvSlot(Index, MapItem(MapNum, I).Num)

//                         ' Open slot available?
//                         If n <> 0 Then
//                             ' Set item in players inventor
//                             Call SetPlayerInvItemNum(Index, n, MapItem(MapNum, I).Num)
//                             If Item(GetPlayerInvItemNum(Index, n)).Type = ITEM_TYPE_CURRENCY Then
//                                 Call SetPlayerInvItemValue(Index, n, GetPlayerInvItemValue(Index, n) + MapItem(MapNum, I).value)
//                                 Msg = "You picked up " & MapItem(MapNum, I).value & " " & Trim$(Item(GetPlayerInvItemNum(Index, n)).Name) & "."
//                             Else
//                                 Call SetPlayerInvItemValue(Index, n, 0)
//                                 Msg = "You picked up a " & Trim$(Item(GetPlayerInvItemNum(Index, n)).Name) & "."
//                             End If
//                             Call SetPlayerInvItemDur(Index, n, MapItem(MapNum, I).Dur)

//                             ' Erase item from the map
//                             MapItem(MapNum, I).Num = 0
//                             MapItem(MapNum, I).value = 0
//                             MapItem(MapNum, I).Dur = 0
//                             MapItem(MapNum, I).X = 0
//                             MapItem(MapNum, I).Y = 0

//                             Call SendInventoryUpdate(Index, n)
//                             Call SpawnItemSlot(I, 0, 0, 0, GetPlayerMap(Index), 0, 0)
//                             Call PlayerMsg(Index, Msg, Yellow)
//                             Exit For
//                         Else
//                             Call PlayerMsg(Index, "Your inventory is full.", BrightRed)
//                             Exit For
//                         End If

//                     End If

//                 End If

//             End If

//         End If
//     Next
// End Sub

// Public Sub PlayerMapDropItem(ByVal Index As Long, ByVal InvNum As Long, ByVal Amount As Long)
//     Dim I As Long

//     ' Check for subscript out of range
//     If Not IsPlaying(Index) Or InvNum <= 0 Or InvNum > MAX_INV Then
//         Exit Sub
//     End If

//     If (GetPlayerInvItemNum(Index, InvNum) > 0) Then
//         If (GetPlayerInvItemNum(Index, InvNum) <= MAX_ITEMS) Then

//             I = FindOpenMapItemSlot(GetPlayerMap(Index))

//             If I <> 0 Then
//                 MapItem(GetPlayerMap(Index), I).Dur = 0

//                 ' Check to see if its any sort of ArmorSlot/WeaponSlot
//                 Select Case Item(GetPlayerInvItemNum(Index, InvNum)).Type
//                 Case ITEM_TYPE_ARMOR
//                     If InvNum = GetPlayerEquipmentSlot(Index, Armor) Then
//                         Call SetPlayerEquipmentSlot(Index, 0, Armor)
//                         Call SendWornEquipment(Index)
//                     End If
//                     MapItem(GetPlayerMap(Index), I).Dur = GetPlayerInvItemDur(Index, InvNum)

//                 Case ITEM_TYPE_WEAPON
//                     If InvNum = GetPlayerEquipmentSlot(Index, Weapon) Then
//                         Call SetPlayerEquipmentSlot(Index, 0, Weapon)
//                         Call SendWornEquipment(Index)
//                     End If
//                     MapItem(GetPlayerMap(Index), I).Dur = GetPlayerInvItemDur(Index, InvNum)

//                 Case ITEM_TYPE_HELMET
//                     If InvNum = GetPlayerEquipmentSlot(Index, Helmet) Then
//                         Call SetPlayerEquipmentSlot(Index, 0, Helmet)
//                         Call SendWornEquipment(Index)
//                     End If
//                     MapItem(GetPlayerMap(Index), I).Dur = GetPlayerInvItemDur(Index, InvNum)

//                 Case ITEM_TYPE_SHIELD
//                     If InvNum = GetPlayerEquipmentSlot(Index, Shield) Then
//                         Call SetPlayerEquipmentSlot(Index, 0, Shield)
//                         Call SendWornEquipment(Index)
//                     End If
//                     MapItem(GetPlayerMap(Index), I).Dur = GetPlayerInvItemDur(Index, InvNum)
//                 End Select

//                 MapItem(GetPlayerMap(Index), I).Num = GetPlayerInvItemNum(Index, InvNum)
//                 MapItem(GetPlayerMap(Index), I).X = GetPlayerX(Index)
//                 MapItem(GetPlayerMap(Index), I).Y = GetPlayerY(Index)

//                 If Item(GetPlayerInvItemNum(Index, InvNum)).Type = ITEM_TYPE_CURRENCY Then
//                     ' Check if its more then they have and if so drop it all
//                     If Amount >= GetPlayerInvItemValue(Index, InvNum) Then
//                         MapItem(GetPlayerMap(Index), I).value = GetPlayerInvItemValue(Index, InvNum)
//                         Call MapMsg(GetPlayerMap(Index), GetPlayerName(Index) & " drops " & GetPlayerInvItemValue(Index, InvNum) & " " & Trim$(Item(GetPlayerInvItemNum(Index, InvNum)).Name) & ".", Yellow)
//                         Call SetPlayerInvItemNum(Index, InvNum, 0)
//                         Call SetPlayerInvItemValue(Index, InvNum, 0)
//                         Call SetPlayerInvItemDur(Index, InvNum, 0)
//                     Else
//                         MapItem(GetPlayerMap(Index), I).value = Amount
//                         Call MapMsg(GetPlayerMap(Index), GetPlayerName(Index) & " drops " & Amount & " " & Trim$(Item(GetPlayerInvItemNum(Index, InvNum)).Name) & ".", Yellow)
//                         Call SetPlayerInvItemValue(Index, InvNum, GetPlayerInvItemValue(Index, InvNum) - Amount)
//                     End If
//                 Else
//                     ' Its not a currency object so this is easy
//                     MapItem(GetPlayerMap(Index), I).value = 0
//                     If Item(GetPlayerInvItemNum(Index, InvNum)).Type >= ITEM_TYPE_WEAPON And Item(GetPlayerInvItemNum(Index, InvNum)).Type <= ITEM_TYPE_SHIELD Then
//                         Call MapMsg(GetPlayerMap(Index), GetPlayerName(Index) & " drops a " & Trim$(Item(GetPlayerInvItemNum(Index, InvNum)).Name) & " " & GetPlayerInvItemDur(Index, InvNum) & "/" & Item(GetPlayerInvItemNum(Index, InvNum)).Data1 & ".", Yellow)
//                     Else
//                         Call MapMsg(GetPlayerMap(Index), GetPlayerName(Index) & " drops a " & Trim$(Item(GetPlayerInvItemNum(Index, InvNum)).Name) & ".", Yellow)
//                     End If

//                     Call SetPlayerInvItemNum(Index, InvNum, 0)
//                     Call SetPlayerInvItemValue(Index, InvNum, 0)
//                     Call SetPlayerInvItemDur(Index, InvNum, 0)
//                 End If

//                 ' Send inventory update
//                 Call SendInventoryUpdate(Index, InvNum)
//                 ' Spawn the item before we set the num or we'll get a different free map item slot
//                 Call SpawnItemSlot(I, MapItem(GetPlayerMap(Index), I).Num, Amount, MapItem(GetPlayerMap(Index), I).Dur, GetPlayerMap(Index), GetPlayerX(Index), GetPlayerY(Index))
//             Else
//                 Call PlayerMsg(Index, "To many items already on the ground.", BrightRed)
//             End If
//         End If
//     End If
// End Sub

// Public Sub CheckPlayerLevelUp(ByVal Index As Long)
//     Dim I As Long
//     Dim expRollOver As Long

//     ' Check if attacker got a level up
//     If GetPlayerExp(Index) >= GetPlayerNextLevel(Index) Then
//         expRollOver = CLng(GetPlayerExp(Index) - GetPlayerNextLevel(Index))
//         Call SetPlayerLevel(Index, GetPlayerLevel(Index) + 1)

//         ' Get the amount of skill points to add
//         I = (GetPlayerStat(Index, Stats.Speed) \ 10)
//         If I < 1 Then I = 1
//         If I > 3 Then I = 3

//         Call SetPlayerPOINTS(Index, GetPlayerPOINTS(Index) + I)
//         Call SetPlayerExp(Index, expRollOver)
//         Call GlobalMsg(GetPlayerName(Index) & " has gained a level!", Brown)
//         Call PlayerMsg(Index, "You have gained a level!  You now have " & GetPlayerPOINTS(Index) & " stat points to distribute.", BrightBlue)
//     End If

// End Sub

// Public Function GetPlayerVitalRegen(ByVal Index As Long, ByVal Vital As Vitals) As Long
//     Dim I As Long

//     ' Prevent subscript out of range
//     If Not IsPlaying(Index) Or Index <= 0 Or Index > High_Index Then
//         GetPlayerVitalRegen = 0
//         Exit Function
//     End If

//     Select Case Vital
//     Case HP
//         I = (GetPlayerStat(Index, Stats.Defense) \ 2)
//     Case MP
//         I = (GetPlayerStat(Index, Stats.Magic) \ 2)
//     Case SP
//         I = (GetPlayerStat(Index, Stats.Speed) \ 2)
//     End Select

//     If I < 2 Then I = 2

//     GetPlayerVitalRegen = I
// End Function

// ' ToDo
// Public Sub OnDeath(ByVal Index As Long)
//     Dim I As Long

//     ' Set HP to nothing
//     Call SetPlayerVital(Index, Vitals.HP, 0)

//     ' Drop all worn items
//     If Map(GetPlayerMap(Index)).Moral <> MAP_MORAL_ARENA Then
//         For I = 1 To Equipment.Equipment_Count - 1
//             If GetPlayerEquipmentSlot(Index, I) > 0 Then
//                 PlayerMapDropItem Index, GetPlayerEquipmentSlot(Index, I), 0
//             End If
//         Next
//     End If

//     ' Warp player away
//     Call PlayerWarp(Index, START_MAP, START_X, START_Y)

//     ' Restore vitals
//     Call SetPlayerVital(Index, Vitals.HP, GetPlayerMaxVital(Index, Vitals.HP))
//     Call SetPlayerVital(Index, Vitals.MP, GetPlayerMaxVital(Index, Vitals.MP))
//     Call SetPlayerVital(Index, Vitals.SP, GetPlayerMaxVital(Index, Vitals.SP))
//     Call SendVital(Index, Vitals.HP)
//     Call SendVital(Index, Vitals.MP)
//     Call SendVital(Index, Vitals.SP)

//     ' If the player the attacker killed was a pk then take it away
//     If GetPlayerPK(Index) = YES Then
//         Call SetPlayerPK(Index, NO)
//         Call SendPlayerData(Index)
//     End If

// End Sub

// Public Sub DamageEquipment(ByVal Index As Long, ByVal EquipmentSlot As Equipment)
//     Dim Slot As Long

//     Slot = GetPlayerEquipmentSlot(Index, EquipmentSlot)

//     If Slot > 0 Then
//         Call SetPlayerInvItemDur(Index, Slot, GetPlayerInvItemDur(Index, Slot) - 1)

//         If GetPlayerInvItemDur(Index, Slot) <= 0 Then
//             Call PlayerMsg(Index, "Your " & Trim$(Item(GetPlayerInvItemNum(Index, Slot)).Name) & " has broken.", Yellow)
//             Call TakeItem(Index, GetPlayerInvItemNum(Index, Slot), 0)
//         Else
//             If GetPlayerInvItemDur(Index, Slot) <= 5 Then
//                 Call PlayerMsg(Index, "Your " & Trim$(Item(GetPlayerInvItemNum(Index, Slot)).Name) & " is about to break!", Yellow)
//             End If
//         End If
//     End If
// End Sub

func UpdateHighIndex() {
	index := 0

	for i := 0; i < config.MaxPlayers; i++ {
		if players[i].IsLoggedIn() {
			index = i + 1
		}
	}

	writer := net.NewWriter()
	writer.WriteInteger(SvHighIndex)
	writer.WriteLong(index)

	SendDataToAll(writer.Bytes())
}
