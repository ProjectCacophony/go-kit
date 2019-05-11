package interfaces

// Localization represents a localization service
type Localization interface {
	Language() string
	Translate(string) string
	Translatef(string, ...interface{}) string
	TranslateMap(string, map[interface{}]interface{}) string
}
