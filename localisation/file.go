package localisation

import (
	"bytes"
	"math/rand"
	"text/template"
	"time"

	toml "github.com/pelletier/go-toml"
)

// FileSource is a localisation source from a toml file
type FileSource struct {
	path     string
	language string

	translations map[string][]*template.Template
	random       *rand.Rand
}

// NewFileSource creates a new FileSource
func NewFileSource(path, language string) (*FileSource, error) {
	l := &FileSource{
		path:     path,
		language: language,

		translations: make(map[string][]*template.Template),
		random:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	err := l.load()
	if err != nil {
		return nil, err
	}

	return l, nil
}

// load loads and parses the file
func (l *FileSource) load() error {
	tree, err := toml.LoadFile(l.path)
	if err != nil {
		return err
	}

	treeMap := tree.ToMap()

	for key, value := range treeMap {

		switch v := value.(type) {
		case string:

			translation, err := template.New(key).Parse(v)
			if err != nil {
				continue
			}

			l.translations[key] = []*template.Template{
				translation,
			}
		case []interface{}:

			for _, vValue := range v {

				vString, ok := vValue.(string)
				if ok {

					translation, err := template.New(key).Parse(vString)
					if err != nil {
						continue
					}

					l.translations[key] = append(
						l.translations[key],
						translation,
					)
				}
			}
		}
	}

	return nil
}

// Language returns the language of the localisations
func (l *FileSource) Language() string {
	return l.language
}

// Translate translates a given key
func (l *FileSource) Translate(key string) string {
	return l.TranslateMap(key, nil)
}

// Translate translates a given key with the given fields
func (l *FileSource) Translatef(key string, fields ...interface{}) string {
	// create map out of fields
	data := make(map[interface{}]interface{})
	for i := range fields {
		if i%2 == 0 && len(fields) > i+1 {
			data[fields[i]] = fields[i+1]
		}
	}

	return l.TranslateMap(key, data)
}

// TranslateMap translates a given key with the given values
func (l *FileSource) TranslateMap(key string, values map[interface{}]interface{}) string {
	translations, ok := l.translations[key]
	if !ok || len(translations) == 0 {
		return key
	}

	translation := translations[l.random.Intn(len(translations))]

	var buffer bytes.Buffer
	err := translation.Execute(&buffer, values)
	if err != nil {
		return key
	}

	return buffer.String()
}