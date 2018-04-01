// Package bytex contains byte (string) processing functions.
package bytex

import (
	"unicode"
	"unicode/utf8"
)

// ByteContainsFold is like bytes.Contains but uses Unicode case-folding.
func ByteContainsFold(s, substr []byte) bool {
	return ByteIndexFold(s, substr) >= 0
}

// EqualFoldRune compares a and b runes whether they fold equally.
// The code comes from strings.EqualFold, but shortened to only one rune.
func EqualFoldRune(sr, tr rune) bool {
	if sr == tr {
		return true
	}

	// Make sr < tr to simplify what follows.
	if tr < sr {
		sr, tr = tr, sr
	}

	// Fast check for ASCII.
	if tr < utf8.RuneSelf && 'A' <= sr && sr <= 'Z' {
		// ASCII, and sr is upper case.  tr must be lower case.
		return tr == sr+'a'-'A'
	}

	// General case.  SimpleFold(x) returns the next equivalent rune > x
	// or wraps around to smaller values.
	r := unicode.SimpleFold(sr)
	for r != sr && r < tr {
		r = unicode.SimpleFold(r)
	}

	return r == tr
}

// ByteIndexFold is like bytes.Contains but uses Unicode case-folding.
func ByteIndexFold(s, substr []byte) int {
	if len(substr) == 0 {
		return 0
	}
	if len(s) == 0 {
		return -1
	}

	firstRune := rune(substr[0])
	if firstRune >= utf8.RuneSelf {
		firstRune, _ = utf8.DecodeRune(substr)
	}

	pos := 0
	for {
		rune, size := utf8.DecodeRune(s)
		if EqualFoldRune(rune, firstRune) && ByteHasPrefixFold(s, substr) {
			return pos
		}
		pos += size
		s = s[size:]
		if len(s) == 0 {
			break
		}
	}

	return -1
}

// ByteHasPrefixFold is like strings.HasPrefix but uses Unicode case-folding.
func ByteHasPrefixFold(s, prefix []byte) bool {
	if len(prefix) == 0 {
		return true
	}

	for {
		pr, prSize := utf8.DecodeRune(prefix)
		prefix = prefix[prSize:]
		if len(s) == 0 {
			return false
		}

		// Step with s, too.
		sr, size := utf8.DecodeRune(s)
		if sr == utf8.RuneError {
			return false
		}

		s = s[size:]
		if !EqualFoldRune(sr, pr) {
			return false
		}

		if len(prefix) == 0 {
			break
		}
	}

	return true
}
