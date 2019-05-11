package events

import (
	"gitlab.com/Cacophony/go-kit/localization"
)

// Translate translates a given key for the event
func (e *Event) Translate(key string, values ...interface{}) string {
	values = append(values, "prefix", e.prefix)
	return localization.Translate(e.Localizations(), key, values...)
}
