package events

import (
	"github.com/bwmarrin/discordgo"
)

// Send sends a message to the given channel, translates it if possible
func (e *Event) Send(channelID, message string) (*discordgo.Message, error) {
	return e.Sendf(channelID, message)
}

// Sendf sends a message to the given channel, translates it if possible
// TODO: make DMs possible
// TODO: check language
func (e *Event) Sendf(channelID, message string, values ...interface{}) (*discordgo.Message, error) {
	var translation string
	for _, localisation := range e.Localisations() {
		translation = localisation.Translatef(message, values...)

		// valid translation found
		if translation != message {
			break
		}
	}

	return e.Discord().ChannelMessageSend(channelID, translation)
}
