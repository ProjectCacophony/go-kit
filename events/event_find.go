package events

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// FindUser finds any kind of target user in the command
func (e *Event) FindUser() (*discordgo.User, error) {
	if e.Type != MessageCreateType {
		return nil, errors.New("event has to be MessageCreate")
	}

	// try any mentions in the command
	for _, mention := range e.MessageCreate.Mentions {
		user, err := e.State().User(mention.ID)
		if err == nil {
			return user, nil
		}
	}

	for _, field := range e.Fields() {
		user, err := e.State().UserFromMention(field)
		if err == nil {
			return user, nil
		}
	}

	return e.State().User(e.UserID)
}

// FindChannel finds a target text channel in the command
// channels have to be on the current guild
func (e *Event) FindChannel() (*discordgo.Channel, error) {
	if e.Type != MessageCreateType {
		return nil, errors.New("event has to be MessageCreate")
	}

	for _, field := range e.Fields() {
		channel, err := e.State().ChannelFromMention(e.GuildID, field)
		if err == nil {
			return channel, nil
		}
	}

	return e.State().Channel(e.ChannelID)
}

// FindAnyChannel finds any kind of target channel in the command
// channels have to be on the current guild
func (e *Event) FindAnyChannel() (*discordgo.Channel, error) {
	if e.Type != MessageCreateType {
		return nil, errors.New("event has to be MessageCreate")
	}

	for _, field := range e.Fields() {
		channel, err := e.State().ChannelFromMentionTypes(e.GuildID, field)
		if err == nil {
			return channel, nil
		}
	}

	return e.State().Channel(e.ChannelID)
}
