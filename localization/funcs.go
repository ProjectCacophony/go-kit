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
	humanize "github.com/dustin/go-humanize"
	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"mvdan.cc/xurls/v2"
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

// MapURL allows predefining a specific short URL for a long URL
// can be used to define vanity URLs for example
func MapURL(from, to string) {
	shortenedLinkCacheLock.Lock()
	defer shortenedLinkCacheLock.Unlock()

	shortenedLinkCache[from] = to
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

// the additional functions to use in the template engine
var translationFuncs = template.FuncMap{
	// ToUpper returns a uppercase version of a string
	// example: {{ToUpper foobar}} => FOOBAR
	"ToUpper": strings.ToUpper,

	// ToLower returns a lowercase version of a string
	// example: {{ToLower FOOBAR}} => foobar
	"ToLower": strings.ToLower,

	// Title returns a titleised version of a string
	// example: {{Title foo bar}} => Foo Bar
	"Title": cases.Title(language.English).String,

	// Replace replaces
	"Replace": strings.Replace,

	// Escape returns a text to be escaped to be used in Discord Embeds
	// example: {{Escape "`"}} => \`
	"Escape": func(text string) string {
		// TODO: same as discord.EscapeDiscordStrict, cannot import here due to import cycle, fix this
		text = strings.Replace(text, "`", "\\`", -1)
		text = strings.Replace(text, "*", "\\*", -1)
		text = strings.Replace(text, "_", "\\_", -1)
		text = strings.Replace(text, "~", "\\~", -1)
		text = strings.Replace(text, "#", "\\#", -1)
		text = strings.Replace(text, "@", "\\@", -1)

		return text
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

	"HumanizeBytes": func(size int64) string {
		return humanize.Bytes(uint64(size))
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
		shortenedLinkCacheLock.Lock()
		defer shortenedLinkCacheLock.Unlock()

		if shortenedLinkCache[value] != "" {
			return shortenedLinkCache[value]
		}

		if polrClient == nil {
			return value
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

	"ParseTimestamp": func(i int64) time.Time {
		return time.Unix(i, 0)
	},

	"QuoteText": func(input string) string {
		return "> " + strings.Replace(input, "\n", "\n> ", -1)
	},

	"Pad": func(text string, length int) string {
		return fmt.Sprintf("%-"+strconv.Itoa(length)+"v", text)
	},
}
