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

	return e.Send(e.MessageCreate.ChannelID, message, values...)
}

// RespondComplex sends a message to the source channel, translates it if possible
func (e *Event) RespondComplex(message *discordgo.MessageSend, values ...interface{}) ([]*discordgo.Message, error) {
	if e.Type != MessageCreateType {
		return nil, errors.New("cannot respond to this event")
	}

	return e.SendComplex(e.MessageCreate.ChannelID, message, values...)
}

// Send sends a message to the given channel, translates it if possible
// TODO: check language
func (e *Event) Send(channelID, message string, values ...interface{}) ([]*discordgo.Message, error) {
	return discord.SendComplexWithVars(
		e.Discord(),
		e.Localisations(),
		channelID,
		&discordgo.MessageSend{
			Content: message,
		},
		append(values, "prefix", e.Prefix())...,
	)
}

// SendComplex sends a message to the given channel, translates it if possible
// TODO: check language
func (e *Event) SendComplex(
	channelID string, message *discordgo.MessageSend, values ...interface{},
) ([]*discordgo.Message, error) {
	return discord.SendComplexWithVars(
		e.Discord(),
		e.Localisations(),
		channelID,
		message,
		append(values, "prefix", e.Prefix())...,
	)
}
