package localisation

import (
	"strconv"
	"strings"
	"text/template"
	"time"

	"gitlab.com/Cacophony/go-kit/feed"

	"github.com/pkg/errors"

	humanize "github.com/dustin/go-humanize"
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

		// Replace replaces
		"Replace": strings.Replace,

		// Escape returns a text to be escaped to be used in Discord Embeds
		// example: {{Escape "`"}} => \`
		"Escape": func(text string) string {
			return strings.Replace(text, "`", "\\`", -1)
		},

		// EscapeLink returns a link to be escaped to be used in Discord Embeds
		// example: {{EscapeLink "https://example.org/A+(B)"}} => https://example.org/A+%28B%29
		"EscapeLink": func(text string) string {
			text = strings.Replace(text, ")", "%29", -1)
			text = strings.Replace(text, "(", "%28", -1)
			return text
		},

		// HumanizeNumber adds commas after every three orders of magnitude
		"HumanizeNumber": func(number int) string {
			return humanize.Comma(int64(number))
		},

		// HumanizeNumber64 adds commas after every three orders of magnitude
		"HumanizeNumber64": humanize.Comma,

		// HumanizeTime formats a time into a relative string, eg 3 days ago
		"HumanizeTime": func(then time.Time) string {
			if then.IsZero() {
				return "Never"
			}

			return humanize.Time(then)
		},

		"ToString": func(value interface{}) (string, error) {
			switch v := value.(type) {
			case string:
				return v, nil
			case int:
				return strconv.Itoa(v), nil
			case feed.Status:
				return string(v), nil
			default:
				return "", errors.New("unable to convert into string")
			}
		},
	}
)
