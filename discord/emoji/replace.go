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
	indexes := colonRegex.FindAllStringIndex(input, -1)

	// replace all :emoji: in the text
	var index []int
	for i := len(indexes) - 1; i >= 0; i-- {
		index = indexes[i]
		input = input[:index[0]] + Get(input[index[0]:index[1]]) + input[index[1]:]
	}

	// return result
	return input
}

// ReplaceWithout replaces all :emoji: in a text with their full ID
func ReplaceWithout(input string) string {
	// match all :emoji: in the text
	indexes := colonRegex.FindAllStringIndex(input, -1)

	// replace all :emoji: in the text
	var index []int
	for i := len(indexes) - 1; i >= 0; i-- {
		index = indexes[i]
		input = input[:index[0]] + GetWithout(input[index[0]:index[1]]) + input[index[1]:]
	}

	// return result
	return input
}
