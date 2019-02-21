package events

// Translate translates a given key for the event
func (e *Event) Translate(key string, values ...interface{}) string {
	if key == "" {
		return key
	}

	values = append(values,
		"prefix", e.prefix,
		"translate", e.Translate,
	)

	var translation string
	for _, localisation := range e.Localisations() {
		translation = localisation.Translatef(key, values...)

		// valid translation found
		if translation != key {
			return translation
		}
	}

	return key
}
