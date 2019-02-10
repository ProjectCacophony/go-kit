package localisation

import (
	"strings"
	"text/template"
)

var (
	// the additional functions to use in the template engine
	// nolint: gochecknoglobals
	translationFuncs = template.FuncMap{

		// ToUpper returns a uppercase version of a string
		// example: {{ToUpper foobar}} => FOOBAR
		"ToUpper": strings.ToUpper,

		// ToLower returns a lowercase version of a string
		// example: {{ToLower FOOBAR}} => foobar
		"ToLower": strings.ToLower,

		// Title returns a titleised version of a string
		// example: {{ToLower foo bar}} => Foo Bar
		"Title": strings.Title,
	}
)
