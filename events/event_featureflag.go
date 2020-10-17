package events

import (
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/label"
)

func (e *Event) IsEnabled(key string, fallback bool) bool {
	_, span := global.Tracer("cacophony.dev/processor").Start(e.Context(), "event.IsEnabled",
		trace.WithAttributes(label.String("featureflag_key", key), label.Bool("featureflag_fallback", fallback)),
	)
	defer span.End()

	if e.FeatureFlagger() == nil {
		return fallback
	}

	return e.FeatureFlagger().IsEnabledFor(key, fallback, e.UserID)
}
