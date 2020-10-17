package events

import (
	"gitlab.com/Cacophony/go-kit/localization"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

// Translate translates a given key for the event
func (e *Event) Translate(key string, values ...interface{}) string {
	_, span := global.Tracer("cacophony.dev/kit").Start(e.Context(), "event.Translate",
		trace.WithAttributes(label.String("key", key)),
	)
	defer span.End()

	values = append(values, "prefix", e.prefix)
	return localization.Translate(e.Localizations(), key, values...)
}
