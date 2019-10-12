package events

import (
	builtinErrors "errors"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/discord"
)

// Respond sends a message to the source channel, translates it if possible
func (e *Event) Respond(message string, values ...interface{}) ([]*discordgo.Message, error) {
	if e.Type != MessageCreateType {
		return nil, errors.New("cannot respond to this event")
	}

	if e.DM() {
		return e.SendDM(e.MessageCreate.Author.ID, message, values...)
	}

	return e.Send(e.MessageCreate.ChannelID, message, values...)
}

func (e *Event) RespondDM(message string, values ...interface{}) ([]*discordgo.Message, error) {
	if e.Type != MessageCreateType {
		return nil, errors.New("cannot respond to this event")
	}

	return e.SendDM(e.UserID, message, values...)
}

// RespondComplex sends a message to the source channel, translates it if possible
func (e *Event) RespondComplex(message *discordgo.MessageSend, values ...interface{}) ([]*discordgo.Message, error) {
	if e.Type != MessageCreateType {
		return nil, errors.New("cannot respond to this event")
	}

	if e.DM() {
		return e.SendComplexDM(e.MessageCreate.Author.ID, message, values...)
	}

	return e.SendComplex(e.MessageCreate.ChannelID, message, values...)
}

// Send sends a message to the given channel, translates it if possible
// TODO: check language
func (e *Event) Send(channelID, message string, values ...interface{}) ([]*discordgo.Message, error) {
	return discord.SendComplexWithVars(
		e.Discord(),
		e.Localizations(),
		channelID,
		&discordgo.MessageSend{
			Content: message,
		},
		append(values, e.defaultParams()...)...,
	)
}

func (e *Event) SendDM(userID, message string, values ...interface{}) ([]*discordgo.Message, error) {
	return e.SendComplexDM(userID, &discordgo.MessageSend{Content: message}, values...)
}

// SendComplex sends a message to the given channel, translates it if possible
// TODO: check language
func (e *Event) SendComplex(channelID string, message *discordgo.MessageSend, values ...interface{}) ([]*discordgo.Message, error) {
	return discord.SendComplexWithVars(
		e.Discord(),
		e.Localizations(),
		channelID,
		message,
		append(values, e.defaultParams()...)...,
	)
}

func (e *Event) SendComplexDM(userID string, message *discordgo.MessageSend, values ...interface{}) ([]*discordgo.Message, error) {
	channelID, err := discord.DMChannel(e.Redis(), e.Discord(), userID)
	if err != nil {
		return nil, err
	}

	messages, err := e.SendComplex(channelID, message, values...)
	if err != nil {
		var discordError *discordgo.RESTError
		if builtinErrors.As(err, &discordError) &&
			discordError.Message != nil &&
			discordError.Message.Code == discordgo.ErrCodeCannotSendMessagesToThisUser {
			discord.BlockDMChannel(e.Redis(), e.Discord(), userID)
		}
	}

	return messages, err
}

// Typing starts typing in the event channel
func (e *Event) Typing() {
	if e.Type != MessageCreateType {
		return
	}

	e.Discord().Client.ChannelTyping(e.MessageCreate.ChannelID)
}

func (e *Event) React(emojiID string, emojiIDs ...string) error {
	if e.Type != MessageCreateType {
		return nil
	}

	channelID := e.MessageCreate.ChannelID
	if e.DM() {
		channelID = e.MessageCreate.Author.ID
	}

	return discord.React(
		e.Redis(),
		e.Discord(),
		channelID,
		e.MessageCreate.Message.ID,
		e.DM(),
		emojiID, emojiIDs...,
	)
}

func (e *Event) defaultParams() []interface{} {
	return []interface{}{
		"prefix", e.Prefix(),
		"timezone", e.Timezone(),
	}
}

func (e *Event) FindMessageLink(link string) (*discordgo.Message, error) {
	message, err := discord.LookupMessageLink(e.state, e.Discord(), link)
	if err != nil {
		return nil, err
	}

	if message.GuildID != e.GuildID {
		return nil, discord.ErrMessageOnWrongServer
	}

	return message, nil
}
