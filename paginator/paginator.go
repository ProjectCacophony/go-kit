package paginator

import (
	"fmt"
	"regexp"

	"go.uber.org/zap"

	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/state"

	"github.com/go-redis/redis"

	"github.com/bwmarrin/discordgo"
)

const (
	LeftArrowEmoji  = "⬅"
	RightArrowEmoji = "➡"
	CloseEmoji      = "🇽"
	NumbersEmoji    = "🔢"
)

type Paginator struct {
	logger *zap.Logger
	redis  *redis.Client
	state  *state.State
	tokens map[string]string

	messageRegexp *regexp.Regexp
}

func NewPaginator(
	logger *zap.Logger,
	redis *redis.Client,
	state *state.State,
	tokens map[string]string,
) (*Paginator, error) {
	p := &Paginator{
		logger: logger,
		redis:  redis,
		state:  state,
		tokens: tokens,
	}

	var err error
	p.messageRegexp, err = regexp.Compile("^[0-9]+$") // nolint: gocritic
	return p, err
}

// nolint: gochecknoglobals
var (
	validReactions = map[string]bool{
		LeftArrowEmoji:  true,
		RightArrowEmoji: true,
		CloseEmoji:      true,
		NumbersEmoji:    true,
	}
)

func (p *Paginator) getSession(botID string) (*discord.Session, error) {
	return discord.NewSession(p.tokens, botID)
}

func (p *Paginator) sendComplex(
	botID, channelID string, send *discordgo.MessageSend, dm bool,
) ([]*discordgo.Message, error) {
	session, err := p.getSession(botID)
	if err != nil {
		return nil, err
	}

	return discord.SendComplexWithVars(
		p.redis,
		session,
		nil,
		channelID,
		send,
		dm,
	)
}

func (p *Paginator) editComplex(
	botID string, edit *discordgo.MessageEdit, dm bool) (*discordgo.Message, error) {
	session, err := p.getSession(botID)
	if err != nil {
		return nil, err
	}

	return discord.EditComplexWithVars(
		p.redis,
		session,
		nil,
		edit,
		dm,
	)
}

// getEmbedFooter is a simlple helper function to return the footer for the embed message
func (p *Paginator) getEmbedFooter(message *PagedEmbedMessage) *discordgo.MessageEmbedFooter {
	footer := &discordgo.MessageEmbedFooter{}

	if message.FullEmbed.Footer != nil {
		footer.IconURL = message.FullEmbed.Footer.IconURL
	}

	footerText := fmt.Sprintf(
		"Page %d / %d",
		message.CurrentPage, message.TotalNumOfPages,
	)
	if message.FullEmbed.Footer.Text != "" {
		footerText += " • " + message.FullEmbed.Footer.Text
	}

	footer.Text = footerText
	return footer
}

func (p *Paginator) addReactionsToMessage(message *PagedEmbedMessage) error {
	session, err := p.getSession(message.BotID)
	if err != nil {
		return err
	}

	err = discord.React(
		p.redis, session, message.ChannelID, message.MessageID, message.DM, LeftArrowEmoji)
	if err != nil {
		return err
	}
	err = discord.React(
		p.redis, session, message.ChannelID, message.MessageID, message.DM, RightArrowEmoji)
	if err != nil {
		return err
	}

	if message.TotalNumOfPages > 4 {
		err = discord.React(
			p.redis, session, message.ChannelID, message.MessageID, message.DM, NumbersEmoji)
		if err != nil {
			return err
		}
	}

	return discord.React(
		p.redis, session, message.ChannelID, message.MessageID, message.DM, CloseEmoji)
}
