package localization

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	polr "github.com/Seklfreak/polr-go"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"mvdan.cc/xurls/v2"

	humanize "github.com/dustin/go-humanize"
)

const nonBreakingSpace = "\u00A0"

var (
	polrClient             *polr.Polr
	shortenedLinkCache     = make(map[string]string)
	shortenedLinkCacheLock sync.Mutex
	xurlsStrict            = xurls.Strict()
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

func getTranslationFuncs() map[string]interface{} {
	// start with sprig
	methods := sprig.FuncMap()

	// remove dangerous methods
	for _, unsafe := range []string{
		"env",
		"expandenv",
		"typeOf",
		"kindIs",
		"typeIsLike",
	} {
		delete(methods, unsafe)
	}

	// add custom methods
	for key, method := range translationFuncs {
		methods[key] = method
	}

	return methods
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

		"TimeIsZero": func(t time.Time) bool {
			return t.IsZero()
		},

		"TimeFormat": func(t time.Time, zone *time.Location) string {
			if t.IsZero() {
				return "`Never`"
			}

			return fmt.Sprintf(
				"`%s` (%s)",
				t.In(zone).Format(time.RFC1123),
				humanize.Time(t),
			)
		},

		"TimeFormatShort": func(t time.Time, zone *time.Location) string {
			if t.IsZero() {
				return "`Never`"
			}

			return fmt.Sprintf(
				"`%s`",
				t.In(zone).Format(time.RFC822),
			)
		},

		"RandIntn": func(n int) int {
			return rand.Intn(n)
		},

		"HideEmbeds": func(text string) string {
			indexes := xurlsStrict.FindAllStringIndex(text, -1)

			var index []int
			for i := len(indexes) - 1; i >= 0; i-- {
				index = indexes[i]
				text = text[:index[0]] + "<" + text[index[0]:index[1]] + ">" + text[index[1]:]
			}

			return text
		},

		"ReplaceWithNonBreakingSpace": func(input string) string {
			return strings.Replace(input, " ", nonBreakingSpace, -1)
		},

		"NonBreakingSpace": func() string {
			return nonBreakingSpace
		},
	}
)
