package localization

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	toml "github.com/pelletier/go-toml"
)

// FileSource is a localization source from a toml file
type FileSource struct {
	path     string
	language string

	template *template.Template
}

// NewFileSource creates a new FileSource
func NewFileSource(path, language string) (*FileSource, error) {
	l := &FileSource{
		path:     path,
		language: language,

		template: (&template.Template{}).Funcs(getTranslationFuncs()),
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

			l.template, err = addTemplate(l.template, key, v)
			if err != nil {
				panic(err) // TODO: handle error better
			}

		case []interface{}:

			var values []string

			for _, vValue := range v {

				vString, ok := vValue.(string)
				if ok {

					values = append(values, vString)
				}
			}

			l.template, err = addTemplate(l.template, key, values...)
			if err != nil {
				panic(err) // TODO: handle error better
			}
		}
	}

	return nil
}

// Language returns the language of the localizations
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
	var buffer bytes.Buffer

	err := l.template.ExecuteTemplate(&buffer, key, values)
	if err != nil {
		if !strings.Contains(err.Error(), "associated with template") {
			fmt.Println(err.Error())
		}

		return key
	}

	return buffer.String()
}
