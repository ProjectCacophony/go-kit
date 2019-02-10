package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// Respond sends a message to the source channel, translates it if possible
func (e *Event) Respond(message string) (*discordgo.Message, error) {
	return e.Respondf(message)
}

// RespondComplex sends a message to the source channel, translates it if possible
func (e *Event) RespondComplex(message *discordgo.MessageSend) (*discordgo.Message, error) {
	if e.Type != MessageCreateType {
		return nil, errors.New("cannot respond to this event")
	}

	return e.SendComplex(e.MessageCreate.ChannelID, message)
}

// Respondf sends a message to the source channel, translates it if possible
func (e *Event) Respondf(message string, values ...interface{}) (*discordgo.Message, error) {
	if e.Type != MessageCreateType {
		return nil, errors.New("cannot respond to this event")
	}

	return e.Sendf(e.MessageCreate.ChannelID, message, values...)
}

// Send sends a message to the given channel, translates it if possible
func (e *Event) Send(channelID, message string) (*discordgo.Message, error) {
	return e.Sendf(channelID, message)
}

// SendComplex sends a message to the given channel, translates it if possible
// TODO: make DMs possible
// TODO: check language
func (e *Event) SendComplex(channelID string, message *discordgo.MessageSend) (*discordgo.Message, error) {
	message.Content = e.Translate(message.Content)

	return e.Discord().ChannelMessageSendComplex(channelID, message)
}

// Sendf sends a message to the given channel, translates it if possible
// TODO: make DMs possible
// TODO: check language
func (e *Event) Sendf(channelID, message string, values ...interface{}) (*discordgo.Message, error) {
	return e.Discord().ChannelMessageSend(channelID, e.Translate(message, values...))
}
