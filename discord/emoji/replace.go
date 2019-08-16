package emoji

import (
	"regexp"
)

var (
	colonRegex = regexp.MustCompile(":[A-Za-z0-9-_]+:")
)

// Replace replaces all :emoji: in a text with their full ID
func Replace(input string) string {
	// match all :emoji: in the text
	colonIndexes := colonRegex.FindAllStringIndex(input, -1)

	// reverse index slice, so that the index positions are always correct when injecting text in content
	for i := len(colonIndexes)/2 - 1; i >= 0; i-- {
		opp := len(colonIndexes) - 1 - i
		colonIndexes[i], colonIndexes[opp] = colonIndexes[opp], colonIndexes[i]
	}

	// replace all :emoji: in the text
	var colonContent string
	for _, colonIndex := range colonIndexes {
		colonContent = input[colonIndex[0]:colonIndex[1]]
		input = input[:colonIndex[0]] + Get(colonContent) + input[colonIndex[1]:]
	}

	// return result
	return input
}

// ReplaceWithout replaces all :emoji: in a text with their full ID
func ReplaceWithout(input string) string {
	// match all :emoji: in the text
	colonIndexes := colonRegex.FindAllStringIndex(input, -1)

	// reverse index slice, so that the index positions are always correct when injecting text in content
	for i := len(colonIndexes)/2 - 1; i >= 0; i-- {
		opp := len(colonIndexes) - 1 - i
		colonIndexes[i], colonIndexes[opp] = colonIndexes[opp], colonIndexes[i]
	}

	// replace all :emoji: in the text
	var colonContent string
	for _, colonIndex := range colonIndexes {
		colonContent = input[colonIndex[0]:colonIndex[1]]
		input = input[:colonIndex[0]] + Get(colonContent) + input[colonIndex[1]:]
	}

	// return result
	return input
}
