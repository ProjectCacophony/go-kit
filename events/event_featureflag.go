package events

import (
	"go.opentelemetry.io/otel/api/global"
)

func (e *Event) IsEnabled(key string, fallback bool) bool {
	_, span := global.Tracer("cacophony.dev/processor").Start(e.Context(), "event.IsEnabled")
	defer span.End()

	if e.FeatureFlagger() == nil {
		return fallback
	}

	return e.FeatureFlagger().IsEnabledFor(key, fallback, e.UserID)
}
