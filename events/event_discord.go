package events

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/discord/emoji"
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

// Typing starts typing in the event channel
func (e *Event) Typing() {
	if e.Type != MessageCreateType {
		return
	}

	e.Discord().ChannelTyping(e.MessageCreate.ChannelID) // nolint: errcheck
}

func (e *Event) React(emojiID string, emojiIDs ...string) error {
	if e.Type != MessageCreateType {
		return nil
	}

	if len(emojiIDs) > 0 {
		emojiID = append(emojiIDs, emojiID)[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(emojiIDs)+1)]
	}

	return e.Discord().MessageReactionAdd( // nolint: errcheck
		e.MessageCreate.ChannelID, e.MessageCreate.ID, emoji.GetWithout(emojiID),
	)
}
