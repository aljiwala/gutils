package utilsx

import "strings"

var (
	// QuoteEscaper should escape quotes from applied string.
	QuoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
)

const (
	// minURLRuneCount represents threshold URL rune count (min threshold value).
	minURLRuneCount = 3

	// maxURLRuneCount represents threshold URL rune count (max threshold value).
	maxURLRuneCount = 2083
)
