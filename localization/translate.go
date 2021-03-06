package localization

import (
	"gitlab.com/Cacophony/go-kit/interfaces"
)

// Translate translates a given key for the event
func Translate(localizations []interfaces.Localization, key string, values ...interface{}) string {
	if key == "" {
		return key
	}

	var translation string
	for _, localization := range localizations {
		translation = localization.Translatef(key, values...)

		// valid translation found
		if translation != key {
			return translation
		}
	}

	return key
}
