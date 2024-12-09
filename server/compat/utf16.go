package compat

import (
	"encoding/binary"
	"unicode/utf16"
)

// StringToUtf16 converts a string to a byte array of UTF-16 characters.
func StringToUtf16(s string, maxLen int) []byte {
	const space uint16 = 0x20

	bytes := make([]byte, maxLen*2)

	codes := utf16.Encode([]rune(s))
	codesLen := len(codes)

	for i := 0; i < maxLen; i++ {
		if i < codesLen {
			binary.LittleEndian.PutUint16(bytes[i*2:], codes[i])
		} else {
			binary.LittleEndian.PutUint16(bytes[i*2:], space)
		}
	}

	return bytes
}
