package color

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
