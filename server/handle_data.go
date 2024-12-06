package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"mirage/internal/database"
	"mirage/internal/packet"
)

type PacketHandler func(player *Player, packet *packet.Reader)

var PacketHandlers [MaxClientPacketId]PacketHandler

func init() {
	PacketHandlers[ClpGetClasses] = HandleGetClasses
	PacketHandlers[ClpCreateAccount] = HandleCreateAccount
	PacketHandlers[ClpDeleteAccount] = HandleDeleteAccount
	PacketHandlers[ClpLogin] = HandleLogin
	PacketHandlers[ClpCreateCharacter] = HandleCreateCharacter
	PacketHandlers[ClpDeleteCharacter] = HandleDeleteCharacter
}

func HandleData(player *Player, bytes []byte) {
	player.Buffer = append(player.Buffer, bytes...)
	if len(player.Buffer) < 2 {
		return
	}

	buf := player.Buffer
	off := 0

	// Handle all packets in the buffer
	for len(buf) >= 2 {
		size := int(binary.LittleEndian.Uint16(buf))
		if len(buf) < size+2 {
			return
		}
		off += 2
		buf = buf[2:]

		reader := packet.NewReader(buf[:size])
		HandlePacket(player, reader)

		off += size
		buf = buf[size:]
	}

	// Move the bytes that are remaining to the front of the buffer
	bytesLeft := len(player.Buffer) - off
	if bytesLeft > 0 {
		copy(player.Buffer, player.Buffer[off:])
	}

	player.Buffer = player.Buffer[:bytesLeft]
}

func HandlePacket(player *Player, reader *packet.Reader) {
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

func HandleGetClasses(player *Player, _ *packet.Reader) {
	if !player.IsPlaying() {
		player.SendNewCharClasses()
	}
}

// ::::::::::::::::::::::::
// :: New account packet ::
// ::::::::::::::::::::::::

func HandleCreateAccount(player *Player, packet *packet.Reader) {
	if player.IsLoggedIn() {
		return
	}

	// Get the data
	accountName := packet.ReadString()
	password := packet.ReadString()

	// Make sure the account name length is valid
	if len(accountName) < MinAccountNameLength || len(accountName) > MaxAccountNameLength {
		player.SendAlert(fmt.Sprintf("Your account name must be between %d and %d characters long.",
			MinAccountNameLength, MaxAccountNameLength))
		return
	}

	// Make sure the password length is valid
	if len(password) < MinPasswordLength {
		player.SendAlert(fmt.Sprintf("Your password must be between at least %d characters long.", MinPasswordLength))
		return
	}

	// Make sure the account name is valid
	if !database.IsValidName(accountName) {
		player.SendAlert("Invalid account name, only letters, numbers, spaces, and _ allowed in names.")
		return
	}

	// Make sure the account name is not already taken
	if database.AccountExists(accountName) {
		player.SendAlert("Sorry, that account name is already taken!")
		return
	}

	_, ok := database.CreateAccount(accountName, password)
	if !ok {
		player.SendAlert("There was an problem creating your account. Please try again later.")
		return
	}

	log.Printf("[%d] Account '%s' has been created\n", player.Id, accountName)

	player.SendAlert("Your account has been created!")
}

// :::::::::::::::::::::::::::
// :: Delete account packet ::
// :::::::::::::::::::::::::::

func HandleDeleteAccount(_ *Player, _ *packet.Reader) {
	// Not Supported
}

// ::::::::::::::::::
// :: Login packet ::
// ::::::::::::::::::

func HandleLogin(player *Player, packet *packet.Reader) {
	if player.IsLoggedIn() {
		return
	}

	// Get the data
	accountName := packet.ReadString()
	password := packet.ReadString()

	// Make sure client version is correct
	if packet.ReadByte() != VersionMajor || packet.ReadByte() != VersionMinor || packet.ReadByte() != VersionRevision {
		player.SendAlert(fmt.Sprintf(
			"Your client is outdated.\n\n"+
				"To continue, please update to the latest version.\n\n"+
				"Download the latest version from %s.", GameWebsite))
		return
	}

	// Make sure the account name length is valid
	if len(accountName) < MinAccountNameLength || len(accountName) > MaxAccountNameLength {
		player.SendAlert(fmt.Sprintf("Your account name must be between %d and %d characters long.",
			MinAccountNameLength, MaxAccountNameLength))
		return
	}

	// Make sure a password was entered
	if len(password) < MinPasswordLength {
		player.SendAlert(fmt.Sprintf("Your password must be between at least %d characters long.", MinPasswordLength))
		return
	}

	// Do not allow players to login while shutting down
	if IsShuttingDown {
		player.SendAlert("The server is currently undergoing maintenance. Please try again later.")
		return
	}

	// Make sure the account exists and the password is correct
	account := database.LoadAccount(accountName)
	if account == nil || !account.IsPasswordCorrect(password) {
		player.SendAlert("That account name does not exist or the password is incorrect.")
		return
	}

	// Make sure the account is not already logged in
	if IsAccountLoggedIn(accountName) {
		player.SendAlert("Multiple account logins are not allowed.")
		return
	}

	characters := database.LoadCharactersForAccount(account.Id)
	characterCount := len(characters)

	player.Account = account

	for i := 0; i < MaxChars; i++ {
		if i < characterCount {
			player.Characters[i] = characters[i]
		} else {
			player.Characters[i].Clear()
		}
	}

	player.SendCharacters()
	// Call SendMaxes(Index)
	// Call SendMapRevs(Index)

	log.Printf("[%d] %s has logged in from %s\n", player.Id, account.Name, player.Connection.RemoteAddr())
}

func HandleCreateCharacter(player *Player, packet *packet.Reader) {
	if !player.IsLoggedIn() {
		return
	}

	characterName := packet.ReadString()
	gender := database.CharacterGender(packet.ReadLong())
	classId := packet.ReadLong() - 1
	slot := packet.ReadLong() - 1

	if slot < 0 || slot >= len(player.Characters) {
		player.ReportHack("character slot out of range")
		return
	}

	if classId < 0 || classId >= len(database.Classes) {
		player.ReportHack("class id out of range")
		return
	}

	if gender != database.GenderMale && gender != database.GenderFemale {
		player.ReportHack("invalid gender")
		return
	}

	if len(characterName) < MinCharacterNameLength {
		player.SendAlert(fmt.Sprintf("Character name must be at least %d characters in length.", MinCharacterNameLength))
		return
	}

	character := &player.Characters[slot]
	if character.Id != 0 {
		player.SendAlert("Character already exists!")
		return
	}

	if !database.IsValidName(characterName) {
		player.SendAlert("Invalid character name, only letters, numbers, spaces, and _ allowed in names.")
		return
	}

	if database.CharacterExists(characterName) {
		player.SendAlert("Sorry, but that name is in use!")
		return
	}

	_, ok := database.CreateCharacter(player.Account.Id, characterName, gender, classId)
	if !ok {
		player.SendAlert("There was an problem creating the character. Please try again later.")
		return
	}

	log.Printf("[%d] Character '%s' has been created by '%s' from %s\n",
		player.Id,
		characterName,
		player.Account.Name,
		player.Connection.RemoteAddr())

	player.SendAlert("Character has been created!")
}

func HandleDeleteCharacter(player *Player, packet *packet.Reader) {
	if !player.IsLoggedIn() {
		return
	}

	slot := packet.ReadLong() - 1
	if slot < 0 || slot >= len(player.Characters) {
		return
	}

	character := &player.Characters[slot]
	if character.Id == 0 {
		player.SendAlert("There is no character in this slot.")
		return
	}

	character.Delete()

	log.Printf("[%d] Character '%s' has been deleted by '%s' from %s\n",
		player.Id,
		character.Name,
		player.Account.Name,
		player.Connection.RemoteAddr())

	player.SendAlert("Character has been deleted!")
}

// ' ::::::::::::::::::::::::::::
// ' :: Using character packet ::
// ' ::::::::::::::::::::::::::::
// Private Sub HandleUseChar(ByVal Index As Long, ByRef Data() As Byte, ByVal StartAddr As Long, ByVal ExtraVar As Long)
//     Dim CharNum As Long
//     Dim F As Long
//     Dim Buffer As clsBuffer

//     If Not IsPlaying(Index) Then
//         Set Buffer = New clsBuffer

//         Buffer.WriteBytes Data()

//         CharNum = Buffer.ReadLong

//         ' Prevent hacking
//         If CharNum < 1 Or CharNum > MAX_CHARS Then
//             Call HackingAttempt(Index, "Invalid CharNum")
//             Exit Sub
//         End If

//         ' Check to make sure the character exists and if so, set it as its current char
//         If CharExist(Index, CharNum) Then
//             TempPlayer(Index).CharNum = CharNum
//             Call JoinGame(Index)

//             CharNum = TempPlayer(Index).CharNum
//             Call AddLog(GetPlayerLogin(Index) & "/" & GetPlayerName(Index) & " has began playing " & GAME_NAME & ".", PLAYER_LOG)
//             Call TextAdd(GetPlayerLogin(Index) & "/" & GetPlayerName(Index) & " has began playing " & GAME_NAME & ".")
//             Call UpdateCaption
//         Else
//             Call AlertMsg(Index, "Character does not exist!")
//         End If
//     End If
// End Sub
