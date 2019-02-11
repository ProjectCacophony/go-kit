package emoji

import (
	"regexp"
	"strings"
)

var (
	// nolint: gochecknoglobals
	colonRegex = regexp.MustCompile(":[A-Za-z0-9-_]+:")
)

// Replace replaces all :emoji: in a text with their full ID
func Replace(input string) string {
	// match all :emoji: in the text
	colonContents := colonRegex.FindAllString(input, -1)

	// replace all :emoji: in the text
	for _, colonContent := range colonContents {
		input = strings.Replace(input, colonContent, Get(colonContent), -1)
	}

	// return result
	return input
}
