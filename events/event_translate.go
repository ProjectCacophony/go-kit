package events

import (
	"gitlab.com/Cacophony/go-kit/localisation"
)

// Translate translates a given key for the event
func (e *Event) Translate(key string, values ...interface{}) string {
	values = append(values, "prefix", e.prefix)
	return localisation.Translate(e.Localisations(), key, values)
}
