package interfaces

// Localisation represents a localisation service
type Localisation interface {
	Language() string
	Translate(string) string
	Translatef(string, ...interface{}) string
	TranslateMap(string, map[interface{}]interface{}) string
}
