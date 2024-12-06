package main

// Game Constants
const (
	GameName    = "Mirage Nova"
	GameWebsite = "https://www.miragenova.com"
	GameAddr    = ":7777"
)

// Current Version
const (
	VersionMajor    = 7
	VersionMinor    = 0
	VersionRevision = 0
)

const (
	MaxPlayers = 100
)

const (
	MinAccountNameLength   = 3
	MaxAccountNameLength   = 20
	MinPasswordLength      = 3
	MinCharacterNameLength = 3
)

type Color byte

const (
	Black Color = iota
	Blue
	Green
	Cyan
	Red
	Magenta
	Brown
	Grey
	DarkGrey
	BrightBlue
	BrightGreen
	BrightCyan
	BrightRed
	Pink
	Yellow
	White
)

const (
	SayColor       = DarkGrey
	GlobalColor    = BrightBlue
	BroadcastColor = Pink
	TellColor      = Green
	EmoteColor     = Cyan
	AdminColor     = Cyan
	HelpColor      = Pink
	WhoColor       = Pink
	JoinLeftColor  = Black
	NpcColor       = Brown
	AlertColor     = Red
	NewMapColor    = Pink
)
