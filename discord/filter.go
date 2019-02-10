package discord

import (
	"strings"
)

const (
	zeroWidthSpace = "\u200B"
)

// EscapeDiscordContent escapes @ everyone, and @ here
func EscapeDiscordContent(content string) (output string) {
	return strings.Replace(
		strings.Replace(
			content, "@everyone", "@"+zeroWidthSpace+"everyone", -1),
		"@here", "@"+zeroWidthSpace+"here", -1,
	)
}
