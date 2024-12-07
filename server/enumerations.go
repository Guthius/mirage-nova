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
	SPlayerData
	SPlayerMove
	SNpcMove
	SPlayerDir
	SNpcDir
	SPlayerXY
	SAttack
	SNpcAttack
	SvCheckForMap
	SMapData
	SMapItemData
	SMapNpcData
	SMapDone
	SSayMsg
	SvGlobalMessage
	SAdminMsg
	SvPlayerMessage
	SMapMsg
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
	SMapKey
	SEditMap
	SShopEditor
	SvUpdateShop
	SEditShop
	SSpellEditor
	SvUpdateSpell
	SEditSpell
	STrade
	SSpells
	SLeft
	SHighIndex
	SCastSpell
	SDoor
	SvLimits
	SSync
	SvMapRevisions
)

const (
	_ = iota
	ClpGetClasses
	ClpCreateAccount
	ClpDeleteAccount
	ClpLogin
	ClpCreateCharacter
	ClpDeleteCharacter
	ClpSelectCharacter
	CSayMsg
	CEmoteMsg
	CBroadcastMsg
	CGlobalMsg
	CAdminMsg
	CPlayerMsg
	CPlayerMove
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
	CRequestNewMap
	CMapData
	CNeedMap
	CMapGetItem
	CMapDropItem
	CMapRespawn
	CMapReport
	CKickPlayer
	CBanList
	CBanDestroy
	CBanPlayer
	CRequestEditMap
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
