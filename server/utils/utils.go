package utils

import (
	"time"
	"unicode"

	"github.com/guthius/mirage-nova/server/character"
)

// isValidName checks if the specified name is valid.
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
func GetAdjacentTile(x int, y int, dir character.Direction) (int, int) {
	switch dir {
	case character.Up:
		return x, y - 1
	case character.Down:
		return x, y + 1
	case character.Left:
		return x - 1, y
	case character.Right:
		return x + 1, y
	}
	return x, y
}
