package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"
	"unicode/utf16"

	"github.com/guthius/mirage-nova/net"
	"github.com/guthius/mirage-nova/server/character"
	"github.com/guthius/mirage-nova/server/common"
	"github.com/guthius/mirage-nova/server/config"
	"github.com/guthius/mirage-nova/server/data"
	"github.com/guthius/mirage-nova/server/user"
	"github.com/guthius/mirage-nova/server/utils"
)

type PacketHandler func(player *PlayerData, reader *net.PacketReader)

var PacketHandlers [MaxClientPacketId]PacketHandler

func init() {
	PacketHandlers[ClGetClasses] = HandleGetClasses
	PacketHandlers[ClCreateAccount] = HandleCreateAccount
	PacketHandlers[ClLogin] = HandleLogin
	PacketHandlers[ClCreateCharacter] = HandleCreateCharacter
	PacketHandlers[ClDeleteCharacter] = HandleDeleteCharacter
	PacketHandlers[ClSelectCharacter] = HandleSelectCharacter
	PacketHandlers[ClPlayerMove] = HandlePlayerMove
	PacketHandlers[ClRequestNewLevel] = HandleRequestNewLevel
	PacketHandlers[ClLevelData] = HandleLevelData
	PacketHandlers[ClNeedLevel] = HandleNeedLevel
	PacketHandlers[ClRequestEditLevel] = HandleRequestEditLevel
}

func HandlePacket(player *PlayerData, reader *net.PacketReader) {
	if reader.Remaining() < 2 {
		return
	}

	packetId := reader.ReadInteger()
	if packetId < 0 || packetId >= MaxClientPacketId {
		return
	}

	packetHandler := PacketHandlers[packetId]
	if packetHandler == nil {
		return
	}

	packetHandler(player, reader)
}

// :::::::::::::::::::::::::::::::::::::::::::::::
// :: Requesting classes for making a character ::
// :::::::::::::::::::::::::::::::::::::::::::::::

func HandleGetClasses(player *PlayerData, _ *net.PacketReader) {
	if !player.IsPlaying() {
		SendNewCharClasses(player)
	}
}

// ::::::::::::::::::::::::
// :: New account packet ::
// ::::::::::::::::::::::::

func HandleCreateAccount(player *PlayerData, reader *net.PacketReader) {
	if player.IsLoggedIn() {
		return
	}

	// Get the data
	accountName := reader.ReadString()
	password := reader.ReadString()

	// Make sure the account name length is valid
	if len(accountName) < 3 || len(accountName) > 20 {
		SendAlert(player, "Your account name must be between 3 and 20 characters long.")
		return
	}

	// Make sure the password length is valid
	if len(password) < 3 {
		SendAlert(player, "Your password must be between at least 3 characters long.")
		return
	}

	// Make sure the account name is valid
	if !utils.IsValidName(accountName) {
		SendAlert(player, "Invalid account name, only letters, numbers, spaces, and _ allowed in names.")
		return
	}

	// Make sure the account name is not already taken
	if user.Exists(accountName) {
		SendAlert(player, "Sorry, that account name is already taken!")
		return
	}

	_, ok := user.Create(accountName, password, player.Connection.RemoteAddr())
	if !ok {
		SendAlert(player, "There was an problem creating your account. Please try again later.")
		return
	}

	log.Printf("[%d] Account '%s' has been created\n", player.Id, accountName)

	SendAlert(player, "Your account has been created!")
}

// ::::::::::::::::::
// :: Login packet ::
// ::::::::::::::::::

func HandleLogin(player *PlayerData, reader *net.PacketReader) {
	if player.IsLoggedIn() {
		return
	}

	// Get the account data
	accountName := reader.ReadString()
	password := reader.ReadString()

	// Make sure client version is correct
	if reader.ReadByte() != config.VersionMajor || reader.ReadByte() != config.VersionMinor || reader.ReadByte() != config.VersionRevision {
		SendAlert(player, fmt.Sprintf(
			"Your client is outdated.\n\n"+
				"To continue, please update to the latest version.\n\n"+
				"Download the latest version from %s.", config.GameWebsite))
		return
	}

	// Make sure the account name length is valid
	if len(accountName) < 3 || len(accountName) > 20 {
		SendAlert(player, "Your account name must be between 3 and 20 characters long.")
		return
	}

	// Make sure a password was entered
	if len(password) < 3 {
		SendAlert(player, "Your password must be between at least 3 characters long.")
		return
	}

	// Do not allow players to login while shutting down
	if IsShuttingDown {
		SendAlert(player, "The server is currently undergoing maintenance. Please try again later.")
		return
	}

	// Make sure the account exists and the password is correct
	account := user.Load(accountName)
	if account == nil || !account.IsPasswordCorrect(password) {
		SendAlert(player, "That account name does not exist or the password is incorrect.")
		return
	}

	// Make sure the account is not already logged in
	if IsAccountLoggedIn(accountName) {
		SendAlert(player, "Multiple account logins are not allowed.")
		return
	}

	characters := character.LoadCharactersForAccount(account.Id)
	characterCount := len(characters)

	player.Account = account

	for i := 0; i < config.MaxChars; i++ {
		if i < characterCount {
			player.CharacterList[i] = characters[i]
		} else {
			player.CharacterList[i].Clear()
		}
	}

	SendCharacters(player)
	SendLimits(player)
	SendMapRevisions(player)

	log.Printf("[%d] %s has logged in from %s\n", player.Id, account.Name, player.Connection.RemoteAddr())
}

// ::::::::::::::::::::::::::
// :: Add character packet ::
// ::::::::::::::::::::::::::

func HandleCreateCharacter(player *PlayerData, reader *net.PacketReader) {
	if !player.IsLoggedIn() {
		return
	}

	characterName := reader.ReadString()
	gender := character.Gender(reader.ReadLong())
	classId := reader.ReadLong() - 1
	slot := reader.ReadLong() - 1

	if slot < 0 || slot >= len(player.CharacterList) {
		ReportHack(player, "character slot out of range")
		return
	}

	if classId < 0 || classId >= data.GetClassCount() {
		ReportHack(player, "class id out of range")
		return
	}

	if gender != character.GenderMale && gender != character.GenderFemale {
		ReportHack(player, "invalid gender")
		return
	}

	if len(characterName) < 3 {
		SendAlert(player, "Character name must be at least 3 characters in length.")
		return
	}

	char := &player.CharacterList[slot]
	if char.Id != 0 {
		SendAlert(player, "Character already exists!")
		return
	}

	if !utils.IsValidName(characterName) {
		SendAlert(player, "Invalid character name, only letters, numbers, spaces, and _ allowed in names.")
		return
	}

	if character.Exists(characterName) {
		SendAlert(player, "Sorry, but that name is in use!")
		return
	}

	_, ok := character.CreateCharacter(player.Account.Id, characterName, gender, classId)
	if !ok {
		SendAlert(player, "There was an problem creating the character. Please try again later.")
		return
	}

	log.Printf("[%d] Character '%s' has been created by '%s' from %s\n",
		player.Id,
		characterName,
		player.Account.Name,
		player.Connection.RemoteAddr())

	SendAlert(player, "Character has been created!")
}

func HandleDeleteCharacter(player *PlayerData, reader *net.PacketReader) {
	if !player.IsLoggedIn() {
		return
	}

	// Get the index of the character slot to delete
	slot := reader.ReadLong() - 1
	if slot < 0 || slot >= len(player.CharacterList) {
		return
	}

	// Check whether the character exists
	character := &player.CharacterList[slot]
	if character.Id == 0 {
		SendAlert(player, "There is no character in this slot.")
		return
	}

	character.Delete()

	log.Printf("[%d] Character '%s' has been deleted by '%s' from %s\n",
		player.Id,
		character.Name,
		player.Account.Name,
		player.Connection.RemoteAddr())

	SendAlert(player, "Character has been deleted!")
}

// ::::::::::::::::::::::::::::
// :: Using character packet ::
// ::::::::::::::::::::::::::::

func HandleSelectCharacter(player *PlayerData, reader *net.PacketReader) {
	if !player.IsLoggedIn() || player.Character != nil {
		return
	}

	// Get the index of the selected character slot
	slot := reader.ReadLong() - 1
	if slot < 0 || slot >= len(player.CharacterList) {
		ReportHack(player, "character slot out of range")
	}

	// Check whether the character exists
	if player.CharacterList[slot].Id == 0 {
		SendAlert(player, "character does not exist")
	}

	player.Character = &player.CharacterList[slot]

	JoinGame(player)

	log.Printf("[%d] %s(%s) started playing\n",
		player.Id, player.Account.Name,
		player.Character.Name)
}

// :::::::::::::::::::::::::::::
// :: Moving character packet ::
// :::::::::::::::::::::::::::::

const (
	MoveWalk = 1
	MoveRun  = 2
)

func HandlePlayerMove(player *PlayerData, reader *net.PacketReader) {
	if player.GettingLevel {
		return
	}

	dir := common.Direction(reader.ReadLong())

	movement := reader.ReadLong()
	if movement != MoveWalk && movement != MoveRun {
		ReportHack(player, "invalid movement")
		return
	}

	// Prevent player from moving if they have cast a spell
	if player.CastSpell {
		if utils.GetTickCount() > player.AttackTimer+1000 {
			player.CastSpell = false
		} else {
			SendPlayerXY(player)
			return
		}
	}

	MovePlayer(player, dir, movement)
}

// ::::::::::::::::::::::::::::::::::
// :: Player request for a new map ::
// ::::::::::::::::::::::::::::::::::

func HandleRequestNewLevel(player *PlayerData, reader *net.PacketReader) {
	dir := common.Direction(reader.ReadLong())

	MovePlayer(player, dir, MoveWalk)
}

// :::::::::::::::::::::
// :: Map data packet ::
// :::::::::::::::::::::

func utf16ToString(src []byte) string {
	codes := make([]uint16, len(src)/2)
	for i := 0; i < len(codes); i++ {
		codes[i] = binary.LittleEndian.Uint16(src[i*2:])
	}

	runes := utf16.Decode(codes)
	str := string(runes)

	return strings.TrimSpace(str)
}

func HandleLevelData(player *PlayerData, reader *net.PacketReader) {
	if player.Room == nil || player.Character == nil {
		return
	}

	// Make sure the player has mapper access
	if player.Character.Access < character.AccessMapper {
		return
	}

	levelId := player.Room.Id
	level := player.Room.Level
	newRevision := level.Revision + 1

	level.Name = utf16ToString(reader.Read(config.NameLength * 2))
	level.Revision = reader.ReadLong()
	level.Type = data.LevelType(reader.ReadInteger())
	level.TileSet = reader.ReadInteger()
	level.Up = reader.ReadInteger() - 1
	level.Down = reader.ReadInteger() - 1
	level.Left = reader.ReadInteger() - 1
	level.Right = reader.ReadInteger() - 1
	level.Music = reader.ReadInteger()
	level.BootMap = reader.ReadInteger() - 1
	level.BootX = int(reader.ReadByte())
	level.BootY = int(reader.ReadByte())
	level.Shop = reader.ReadInteger() - 1

	for i := 0; i < len(level.Tiles); i++ {
		for j := 0; j < len(level.Tiles[i].Num); j++ {
			level.Tiles[i].Num[j] = reader.ReadInteger()
		}

		level.Tiles[i].Type = data.TileType(reader.ReadInteger())
		level.Tiles[i].Data1 = reader.ReadInteger()
		level.Tiles[i].Data2 = reader.ReadInteger()
		level.Tiles[i].Data3 = reader.ReadInteger()
	}

	for i := 0; i < config.MaxMapNpcs; i++ {
		level.Npcs[i] = int(reader.ReadByte()) - 1
	}

	level.Revision = newRevision

	for i := 0; i < config.MaxMapNpcs; i++ {
		// TODO: Call ClearMapNpc(I, MapNum)
	}

	/*
	   Call SendMapNpcsToMap(MapNum)
	   Call SpawnMapNpcs(MapNum)

	   ' Clear out it all
	   For I = 1 To MAX_MAP_ITEMS
	       Call SpawnItemSlot(I, 0, 0, 0, GetPlayerMap(Index), MapItem(GetPlayerMap(Index), I).X, MapItem(GetPlayerMap(Index), I).Y)
	       Call ClearMapItem(I, GetPlayerMap(Index))
	   Next

	   ' Respawn
	   Call SpawnMapItems(GetPlayerMap(Index))
	*/

	data.SaveLevel(levelId - 1)

	// Rebuild the level cache
	player.Room.LevelCache = buildLevelCache(levelId, level)

	// Refresh level data for all players in the room
	for _, p := range player.Room.Players {
		if p.IsPlaying() {
			SendLevelData(p)
		}
	}
}

// ::::::::::::::::::::::::::::
// :: Need map yes/no packet ::
// ::::::::::::::::::::::::::::

func HandleNeedLevel(player *PlayerData, reader *net.PacketReader) {
	//  Check if map data is needed to be sent
	needMap := reader.ReadByte()
	if needMap != 0 {
		SendLevelData(player)
	}

	// For I = 1 To MAX_MAPS
	//     Call SendMapItemsTo(Index, I)
	//     Call SendMapNpcsTo(Index, I)
	// Next I

	player.GettingLevel = false

	// Tell the player all map data has been sent
	writer := net.NewWriter()
	writer.WriteInteger(SvLevelDone)

	player.Send(writer.Bytes())

	//  Call SendDoorData(Index)
}

// :::::::::::::::::::::::::::::
// :: Request edit map packet ::
// :::::::::::::::::::::::::::::

func HandleRequestEditLevel(player *PlayerData, reader *net.PacketReader) {
	if player.Character == nil {
		return
	}

	if player.Character.Access < character.AccessMapper {
		return
	}

	writer := net.NewWriter()
	writer.WriteInteger(SvEditLevel)

	player.Send(writer.Bytes())
}
