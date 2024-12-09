package utils

import (
	"time"
	"unicode"

	"github.com/guthius/mirage-nova/server/common"
)

// IsValidName checks if the specified name is valid.
// A name is considered valid if it only contains letters, digits, spaces and underscores.
func IsValidName(name string) bool {
	for _, ch := range name {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != ' ' && ch != '_' {
			return false
		}
	}
	return true
}

// GetTickCount returns the current time in milliseconds.
func GetTickCount() int64 {
	return time.Now().UnixMilli()
}

// GetAdjacentTile returns the coordinates of the tile adjacent to the specified tile in the specified direction.
func GetAdjacentTile(x int, y int, dir common.Direction) (int, int) {
	switch dir {
	case common.DirUp:
		return x, y - 1
	case common.DirDown:
		return x, y + 1
	case common.DirLeft:
		return x - 1, y
	case common.DirRight:
		return x + 1, y
	}
	return x, y
}
