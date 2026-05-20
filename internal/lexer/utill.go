package lexer

import (
	"strings"
	"unicode"
)

// normalizeAllNewlines handles \r\n and legacy \r, turning both into \n
func normalizeAllNewlines(s string) string {
	replacer := strings.NewReplacer("\r\n", "\n", "\r", "\n")
	return replacer.Replace(s)
}

// FindByteIndex mimics Python's str.find(sub, start) but returns the BYTE index.
// 'start' must be a valid byte index and UTF-8 boundary.
// Returns -1 if the substring is not found.
func FindByteIndex(s, sub string, start int) int {
	// Guard against out-of-bounds start indices
	if start < 0 {
		start = 0
	}
	if start >= len(s) {
		return -1
	}

	// Slice from the start byte and find the substring
	byteIdx := strings.Index(s[start:], sub)
	if byteIdx == -1 {
		return -1
	}

	// The returned index is relative to the slice,
	// so add the 'start' offset to get the absolute byte index.
	return start + byteIdx
}

func IsAllWhitespace(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

func leftWhitespaceByteCount(s string) int {
	return len(s) - len(strings.TrimLeftFunc(s, unicode.IsSpace))
}
