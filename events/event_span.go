package events

type SpanContext struct {
	Values map[string]string `json:"values"`
}

func (sp *SpanContext) Get(key string) string {
	if sp.Values == nil {
		return ""
	}

	return sp.Values[key]
}

func (sp *SpanContext) Set(key, value string) {
	if sp.Values == nil {
		sp.Values = make(map[string]string)
	}

	sp.Values[key] = value
}
