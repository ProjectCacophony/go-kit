package localization

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	polr "github.com/Seklfreak/polr-go"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"

	humanize "github.com/dustin/go-humanize"
)

var (
	polrClient             *polr.Polr
	shortenedLinkCache     = make(map[string]string)
	shortenedLinkCacheLock sync.Mutex
)

func init() {

	if os.Getenv("POLR_BASE_URL") != "" &&
		os.Getenv("POLR_API_KEY") != "" {

		polrClient, _ = polr.New(
			os.Getenv("POLR_BASE_URL"),
			os.Getenv("POLR_API_KEY"),
			&http.Client{
				Timeout: 30 * time.Second,
			},
		)
	}
}

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
			default:
				return "", errors.New("unable to convert into string")
			}
		},

		"Shorten": func(value string) string {
			if polrClient == nil {
				return value
			}

			shortenedLinkCacheLock.Lock()
			defer shortenedLinkCacheLock.Unlock()

			if shortenedLinkCache[value] != "" {
				return shortenedLinkCache[value]
			}

			shortened, err := polrClient.Shorten(value, "", false)
			if err != nil {
				return value
			}

			shortenedLinkCache[value] = shortened

			return shortened
		},

		"Contains": strings.Contains,

		"Join": strings.Join,

		"MessageLink": func(message *discordgo.Message) string {
			if message == nil {
				return ""
			}

			guildID := message.GuildID
			if message.GuildID == "" {
				guildID = "@me"
			}

			return fmt.Sprintf(
				"https://discordapp.com/channels/%s/%s/%s",
				guildID,
				message.ChannelID,
				message.ID,
			)
		},
	}
)
