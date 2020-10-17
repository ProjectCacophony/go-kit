package events

type SpanContext map[string]string

func (sp SpanContext) Get(key string) string {
	return sp[key]
}

func (sp SpanContext) Set(key, value string) {
	sp[key] = value
}
