package emoji

import (
	"strings"
)

// Get returns a random emoji for the type with < > around it, can be with or without colons around it
// Example: robyulblush => <:robyulblush:327206930437373952>
func Get(name string) string {
	emoji, ok := List[strings.ToLower(strings.Trim(name, ":"))]
	if ok {
		// return requested emoji
		return emoji.Get()
	}

	// return input
	return name
}

// GetWithout returns a random emoji for the type without < > around it, can be with or without colons around it
// Example: robyulblush => :robyulblush:327206930437373952
func GetWithout(name string) string {
	emoji, ok := List[strings.ToLower(strings.Trim(name, ":"))]
	if ok {
		// return requested emoji
		return emoji.GetWithout()
	}

	name = strings.ReplaceAll(name, "<", "")
	name = strings.ReplaceAll(name, ">", "")

	// return input
	return name
}
