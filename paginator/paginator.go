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

func (p *Paginator) getSession(guildID string) (*discordgo.Session, error) {
	botID, err := p.state.BotForGuild(guildID)
	if err != nil {
		return nil, err
	}

	return discord.NewSession(p.tokens, botID)
}

func (p *Paginator) sendComplex(
	guildID, channelID string, send *discordgo.MessageSend,
) ([]*discordgo.Message, error) {
	session, err := p.getSession(guildID)
	if err != nil {
		return nil, err
	}

	return discord.SendComplexWithVars(
		session,
		nil,
		channelID,
		send,
	)
}

func (p *Paginator) editComplex(
	guildID string, edit *discordgo.MessageEdit) (*discordgo.Message, error) {
	session, err := p.getSession(guildID)
	if err != nil {
		return nil, err
	}

	return session.ChannelMessageEditComplex(edit)
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
	session, err := p.getSession(message.GuildID)
	if err != nil {
		return err
	}

	err = session.MessageReactionAdd(message.ChannelID, message.MessageID, LeftArrowEmoji)
	if err != nil {
		return err
	}
	err = session.MessageReactionAdd(message.ChannelID, message.MessageID, RightArrowEmoji)
	if err != nil {
		return err
	}

	if message.TotalNumOfPages > 4 {
		err = session.MessageReactionAdd(message.ChannelID, message.MessageID, NumbersEmoji)
		if err != nil {
			return err
		}
	}

	return session.MessageReactionAdd(message.ChannelID, message.MessageID, CloseEmoji)
}
