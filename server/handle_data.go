package main

import (
	"encoding/binary"
	"fmt"
	"mirage/internal/packet"
)

type PacketHandler func(player *Player, packet *packet.Reader)

var PacketHandlers = func() [MaxClientPacketId]PacketHandler {
	var handlers [MaxClientPacketId]PacketHandler

	handlers[CGetClasses] = HandleGetClasses
	handlers[CNewAccount] = HandleNewAccount
	handlers[CDelAccount] = HandleDelAccount
	handlers[cLogin] = HandleLogin

	return handlers
}()

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

func HandleNewAccount(player *Player, packet *packet.Reader) {
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
	if len(password) < MinPasswordLength || len(password) > MaxPasswordLength {
		player.SendAlert(fmt.Sprintf("Your password must be between %d and %d characters long.",
			MinPasswordLength, MaxPasswordLength))
		return
	}

	// Make sure the account name is valid
	if !IsValidAccountName(accountName) {
		player.SendAlert("Invalid account name, only letters, numbers, spaces, and _ allowed in names.")
		return
	}

	// Make sure the account name is not already taken
	if AccountExists(accountName) {
		player.SendAlert("Sorry, that account name is already taken!")
		return
	}

	//Call AddAccount(Index, Name, Password)

	PlayerLog.Printf("Account %s has been created\n", accountName)

	player.SendAlert("Your account has been created!")
}

// :::::::::::::::::::::::::::
// :: Delete account packet ::
// :::::::::::::::::::::::::::

func HandleDelAccount(_ *Player, _ *packet.Reader) {
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
	if len(password) == 0 {
		player.SendAlert("Please enter your password.")
		return
	}

	// Do not allow players to login while shutting down
	if IsShuttingDown {
		player.SendAlert("The server is currently undergoing maintenance. Please try again later.")
		return
	}

	// Make sure the account exists and the password is correct
	account := LoadAccount(accountName)
	if account == nil || !account.IsPasswordCorrect(password) {
		player.SendAlert("That account name does not exist or the password is incorrect.")
		return
	}

	// Make sure the account is not already logged in
	if IsAccountLoggedIn(accountName) {
		player.SendAlert("Multiple account logins are not allowed.")
		return
	}

	player.Account = account

	// ' Load the player
	// Call LoadPlayer(Index, Name)
	// Call SendChars(Index)
	// Call SendMaxes(Index)
	// Call SendMapRevs(Index)

	PlayerLog.Printf("[%s] has logged in from %s\n", accountName, player.Connection.RemoteAddr())
}
