package localization

import (
	"errors"
	"strconv"
	"text/template"
)

func addTemplate(template *template.Template, key string, values ...string) (*template.Template, error) {
	if len(values) <= 0 {
		return template, errors.New("no template given")
	}

	if len(values) == 1 {
		return template.New(key).Parse(values[0])
	}

	content := "{{$rand := RandIntn " + strconv.Itoa(len(values)) + "}}"
	for i, value := range values {
		content += "{{if eq $rand " + strconv.Itoa(i) + "}}" + value + "{{end}}"
	}

	return template.New(key).Parse(content)
}
