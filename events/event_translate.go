package events

import (
	"gitlab.com/Cacophony/go-kit/localization"
	"go.opentelemetry.io/otel/api/global"
)

// Translate translates a given key for the event
func (e *Event) Translate(key string, values ...interface{}) string {
	_, span := global.Tracer("cacophony.dev/processor").Start(e.Context(), "event.Translate")
	defer span.End()

	values = append(values, "prefix", e.prefix)
	return localization.Translate(e.Localizations(), key, values...)
}
