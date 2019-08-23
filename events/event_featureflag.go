package events

func (e *Event) IsEnabled(key string, fallback bool) bool {
	if e.FeatureFlagger() == nil {
		return fallback
	}

	return e.FeatureFlagger().IsEnabledFor(key, fallback, e.UserID)
}
