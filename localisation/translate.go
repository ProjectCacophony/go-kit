package localisation

import (
	"gitlab.com/Cacophony/go-kit/interfaces"
)

// Translate translates a given key for the event
func Translate(localisations []interfaces.Localisation, key string, values ...interface{}) string {
	if key == "" {
		return key
	}

	var translation string
	for _, localisation := range localisations {
		translation = localisation.Translatef(key, values...)

		// valid translation found
		if translation != key {
			return translation
		}
	}

	return key
}
