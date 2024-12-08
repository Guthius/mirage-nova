package main

import (
	"fmt"
	"log"

	"github.com/guthius/mirage-nova/net"
	"github.com/guthius/mirage-nova/server/character"
	"github.com/guthius/mirage-nova/server/config"
	"github.com/guthius/mirage-nova/server/data"
	"github.com/guthius/mirage-nova/server/user"
	"github.com/guthius/mirage-nova/server/utils"
)

type PacketHandler func(player *PlayerData, packet *net.PacketReader)

var PacketHandlers [MaxClientPacketId]PacketHandler

func init() {
	PacketHandlers[ClpGetClasses] = HandleGetClasses
	PacketHandlers[ClpCreateAccount] = HandleCreateAccount
	PacketHandlers[ClpLogin] = HandleLogin
	PacketHandlers[ClpCreateCharacter] = HandleCreateCharacter
	PacketHandlers[ClpDeleteCharacter] = HandleDeleteCharacter
	PacketHandlers[ClpSelectCharacter] = HandleSelectCharacter
	PacketHandlers[ClRequestNewMap] = HandleRequestNewMap
	PacketHandlers[ClNeedMap] = HandleNeedMap
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

func HandleCreateAccount(player *PlayerData, packet *net.PacketReader) {
	if player.IsLoggedIn() {
		return
	}

	// Get the data
	accountName := packet.ReadString()
	password := packet.ReadString()

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

	_, ok := user.Create(accountName, password)
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

func HandleLogin(player *PlayerData, packet *net.PacketReader) {
	if player.IsLoggedIn() {
		return
	}

	// Get the data
	accountName := packet.ReadString()
	password := packet.ReadString()

	// Make sure client version is correct
	if packet.ReadByte() != config.VersionMajor || packet.ReadByte() != config.VersionMinor || packet.ReadByte() != config.VersionRevision {
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

func HandleCreateCharacter(player *PlayerData, packet *net.PacketReader) {
	if !player.IsLoggedIn() {
		return
	}

	characterName := packet.ReadString()
	gender := character.CharacterGender(packet.ReadLong())
	classId := packet.ReadLong() - 1
	slot := packet.ReadLong() - 1

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

	if character.CharacterExists(characterName) {
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

func HandleDeleteCharacter(player *PlayerData, packet *net.PacketReader) {
	if !player.IsLoggedIn() {
		return
	}

	slot := packet.ReadLong() - 1
	if slot < 0 || slot >= len(player.CharacterList) {
		return
	}

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

func HandleSelectCharacter(player *PlayerData, packet *net.PacketReader) {
	if !player.IsLoggedIn() || player.Character != nil {
		return
	}

	slot := packet.ReadLong() - 1
	if slot < 0 || slot >= len(player.CharacterList) {
		ReportHack(player, "character slot out of range")
	}

	if player.CharacterList[slot].Id == 0 {
		SendAlert(player, "character does not exist")
	}

	player.Character = &player.CharacterList[slot]

	JoinGame(player)

	log.Printf("[%d] %s(%s) started playing\n",
		player.Id, player.Account.Name,
		player.Character.Name)
}

// ::::::::::::::::::::::::::::::::::
// :: Player request for a new map ::
// ::::::::::::::::::::::::::::::::::

func HandleRequestNewMap(player *PlayerData, packet *net.PacketReader) {
	dir := character.Direction(packet.ReadLong())
	if dir < character.Down || dir >= character.Right {
		return
	}

	// TODO: player.Move(dir, 1)
}

// ::::::::::::::::::::::::::::
// :: Need map yes/no packet ::
// ::::::::::::::::::::::::::::

func HandleNeedMap(player *PlayerData, reader *net.PacketReader) {
	//  Check if map data is needed to be sent
	needMap := reader.ReadByte()
	if needMap != 0 {
		// Call SendMap(Index, GetPlayerMap(Index))
	}

	// For I = 1 To MAX_MAPS
	//     Call SendMapItemsTo(Index, I)
	//     Call SendMapNpcsTo(Index, I)
	// Next I
	// Call SendJoinMap(Index)

	player.GettingMap = false

	writer := net.NewWriter()
	writer.WriteInteger(SvMapDone)

	player.Send(writer.Bytes())

	//  Call SendDoorData(Index)
}
