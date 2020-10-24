package discord

import (
	"strings"
)

const (
	zeroWidthSpace = "\u200B"
)

// EscapeDiscordStrict escapes all discord formatting options
func EscapeDiscordStrict(text string) string {
	text = strings.Replace(text, "`", "\\`", -1)
	text = strings.Replace(text, "*", "\\*", -1)
	text = strings.Replace(text, "_", "\\_", -1)
	text = strings.Replace(text, "~", "\\~", -1)
	text = strings.Replace(text, "#", "\\#", -1)
	text = strings.Replace(text, "@", "\\@", -1)

	return EscapeDiscordLax(text)
}

// EscapeDiscordLax escapes @ everyone, and @ here
func EscapeDiscordLax(text string) (output string) {
	return strings.Replace(
		strings.Replace(
			text, "@everyone", "@"+zeroWidthSpace+"everyone", -1),
		"@here", "@"+zeroWidthSpace+"here", -1,
	)
}
