﻿package main

// *************
// ** Packets **
// *************

const (
	_ = iota
	SvAlert
	SvCharacters
	SvLoginOk
	SvNewCharClasses
	SvClasses
	SvInGame
	SvPlayerInventory
	SPlayerInvUpdate
	SvPlayerEquipment
	SvPlayerHP
	SvPlayerMP
	SvPlayerSP
	SvPlayerStats
	SvPlayerData
	SvPlayerMove
	SNpcMove
	SPlayerDir
	SNpcDir
	SvPlayerXY
	SAttack
	SNpcAttack
	SvCheckForLevel
	SvLevelData
	SMapItemData
	SMapNpcData
	SvLevelDone
	SSayMsg
	SvGlobalMessage
	SAdminMsg
	SvPlayerMessage
	SvRoomMessage
	SSpawnItem
	SItemEditor
	SvUpdateItem
	SEditItem
	SREditor
	SSpawnNpc
	SNpcDead
	SNpcEditor
	SvUpdateNpc
	SEditNpc
	SvMapKey
	SvEditLevel
	SShopEditor
	SvUpdateShop
	SEditShop
	SSpellEditor
	SvUpdateSpell
	SEditSpell
	STrade
	SSpells
	SvLeft
	SvHighIndex
	SCastSpell
	SvDoor
	SvLimits
	SSync
	SvMapRevisions
)

const (
	_ = iota
	ClGetClasses
	ClCreateAccount
	_
	ClLogin
	ClCreateCharacter
	ClDeleteCharacter
	ClSelectCharacter
	CSayMsg
	CEmoteMsg
	CBroadcastMsg
	CGlobalMsg
	CAdminMsg
	CPlayerMsg
	ClPlayerMove
	CPlayerDir
	CUseItem
	CAttack
	CUseStatPoint
	CPlayerInfoRequest
	CWarpMeTo
	CWarpToMe
	CWarpTo
	CSetSprite
	CGetStats
	ClRequestNewLevel
	ClLevelData
	ClNeedLevel
	CMapGetItem
	CMapDropItem
	CMapRespawn
	CMapReport
	CKickPlayer
	CBanList
	CBanDestroy
	CBanPlayer
	ClRequestEditLevel
	CRequestEditItem
	CEditItem
	CSaveItem
	CRequestEditNpc
	CEditNpc
	CSaveNpc
	CRequestEditShop
	CEditShop
	CSaveShop
	CRequestEditSpell
	CEditSpell
	CSaveSpell
	CDelete
	CSetAccess
	CWhosOnline
	CSetMotd
	CTrade
	CTradeRequest
	CFixItem
	CSearch
	CParty
	CJoinParty
	CLeaveParty
	CSpells
	CCast
	CQuit
	CSync
	CMapReqs
	CSleepinn
	CRemoveFromGuild
	CCreateGuild
	CInviteGuild
	CKickGuild
	CGuildPromote
	CLeaveGuild

	MaxClientPacketId
)
