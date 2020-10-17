package events

type SpanContext struct {
	Values map[string]string `json:"values"`
}

func (sp *SpanContext) Get(key string) string {
	if sp == nil || sp.Values == nil {
		return ""
	}

	return sp.Values[key]
}

func (sp *SpanContext) Set(key, value string) {
	if sp == nil || sp.Values == nil {
		sp = &SpanContext{
			Values: make(map[string]string),
		}
	}

	sp.Values[key] = value
}
