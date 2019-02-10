package localisation

import "text/template"

func newTemplate(key, value string) (*template.Template, error) {
	return template.New(key).Funcs(translationFuncs).Parse(value)
}
