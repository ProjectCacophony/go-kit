package localisation

import (
	"strings"
	"text/template"
)

// nolint: gochecknoglobals
var (
	// the additional functions to use in the template engine
	translationFuncs = template.FuncMap{

		// ToUpper returns a uppercase version of a string
		// example: {{ToUpper foobar}} => FOOBAR
		"ToUpper": strings.ToUpper,

		// ToLower returns a lowercase version of a string
		// example: {{ToLower FOOBAR}} => foobar
		"ToLower": strings.ToLower,

		// Title returns a titleised version of a string
		// example: {{Title foo bar}} => Foo Bar
		"Title": strings.Title,
	}
)
