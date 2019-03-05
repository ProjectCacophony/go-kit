package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/discord"
)

// Respond sends a message to the source channel, translates it if possible
func (e *Event) Respond(message string, values ...interface{}) ([]*discordgo.Message, error) {
	if e.Type != MessageCreateType {
		return nil, errors.New("cannot respond to this event")
	}

	channelID := e.MessageCreate.ChannelID
	if e.DM() {
		channelID = e.MessageCreate.Author.ID
	}

	return e.Send(channelID, message, e.DM(), values...)
}

// RespondComplex sends a message to the source channel, translates it if possible
func (e *Event) RespondComplex(message *discordgo.MessageSend, values ...interface{}) ([]*discordgo.Message, error) {
	if e.Type != MessageCreateType {
		return nil, errors.New("cannot respond to this event")
	}

	channelID := e.MessageCreate.ChannelID
	if e.DM() {
		channelID = e.MessageCreate.Author.ID
	}

	return e.SendComplex(channelID, message, e.DM(), values...)
}

// Send sends a message to the given channel, translates it if possible
// TODO: check language
func (e *Event) Send(channelID, message string, dm bool, values ...interface{}) ([]*discordgo.Message, error) {
	return discord.SendComplexWithVars(
		e.Redis(),
		e.Discord(),
		e.Localisations(),
		channelID,
		&discordgo.MessageSend{
			Content: message,
		},
		dm,
		append(values, "prefix", e.Prefix())...,
	)
}

// SendComplex sends a message to the given channel, translates it if possible
// TODO: check language
func (e *Event) SendComplex(
	channelID string, message *discordgo.MessageSend, dm bool, values ...interface{},
) ([]*discordgo.Message, error) {
	return discord.SendComplexWithVars(
		e.Redis(),
		e.Discord(),
		e.Localisations(),
		channelID,
		message,
		dm,
		append(values, "prefix", e.Prefix())...,
	)
}

// Typing starts typing in the event channel
func (e *Event) Typing() {
	if e.Type != MessageCreateType {
		return
	}

	e.Discord().Client.ChannelTyping(e.MessageCreate.ChannelID) // nolint: errcheck
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
