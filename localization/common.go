package localization

import "text/template"

func newTemplate(key, value string) (*template.Template, error) {
	return template.New(key).Funcs(getTranslationFuncs()).Parse(value)
}
