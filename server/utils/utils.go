package utils

import "unicode"

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
