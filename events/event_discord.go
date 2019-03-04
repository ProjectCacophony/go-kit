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
