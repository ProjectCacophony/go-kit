package events

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/state"
	"go.opentelemetry.io/otel/api/global"
)

// FindUser finds any kind of target user in the command
func (e *Event) FindUser(opts ...optionFunc) (*discordgo.User, error) {
	_, span := global.Tracer("cacophony.dev/kit").Start(e.Context(), "event.FindUser")
	defer span.End()

	options := getOptions(opts)

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

	if options.disableFallbackToSelf {
		return nil, state.ErrUserNotFound
	}

	return e.State().User(e.UserID)
}

// FindMember finds any kind of member in the command
func (e *Event) FindMember(opts ...optionFunc) (*discordgo.User, error) {
	_, span := global.Tracer("cacophony.dev/kit").Start(e.Context(), "event.FindMember")
	defer span.End()

	options := getOptions(opts)

	if e.Type != MessageCreateType {
		return nil, errors.New("event has to be MessageCreate")
	}

	// try any mentions in the command
	for _, mention := range e.MessageCreate.Mentions {
		user, err := e.State().Member(e.GuildID, mention.ID)
		if err == nil {
			return user.User, nil
		}
	}

	for _, field := range e.Fields() {
		user, err := e.State().UserFromMention(field)
		if err != nil {
			continue
		}

		isMember, err := e.State().IsMember(e.GuildID, user.ID)
		if err != nil || !isMember {
			continue
		}

		return user, nil
	}

	if options.disableFallbackToSelf {
		return nil, state.ErrUserNotFound
	}

	return e.State().User(e.UserID)
}

// FindChannel finds a target text channel in the command
// channels have to be on the current guild
func (e *Event) FindChannel(opts ...optionFunc) (*discordgo.Channel, error) {
	_, span := global.Tracer("cacophony.dev/kit").Start(e.Context(), "event.FindChannel")
	defer span.End()

	options := getOptions(opts)

	if e.Type != MessageCreateType {
		return nil, errors.New("event has to be MessageCreate")
	}

	for _, field := range e.Fields() {
		channel, err := e.State().ChannelFromMention(e.GuildID, field)
		if err == nil {
			return channel, nil
		}
	}

	if options.disableFallbackToSelf {
		return nil, state.ErrChannelNotFound
	}

	return e.State().Channel(e.ChannelID)
}

// FindAnyChannel finds any kind of target channel in the command
// channels have to be on the current guild
func (e *Event) FindAnyChannel(opts ...optionFunc) (*discordgo.Channel, error) {
	_, span := global.Tracer("cacophony.dev/kit").Start(e.Context(), "event.FindAnyChannel")
	defer span.End()

	options := getOptions(opts)

	if e.Type != MessageCreateType {
		return nil, errors.New("event has to be MessageCreate")
	}

	for _, field := range e.Fields() {
		channel, err := e.State().ChannelFromMentionTypes(e.GuildID, field)
		if err == nil {
			return channel, nil
		}
	}

	if options.disableFallbackToSelf {
		return nil, state.ErrChannelNotFound
	}

	return e.State().Channel(e.ChannelID)
}

// FindRole finds a target role in the command
// the role has to be on the current guild
func (e *Event) FindRole(opts ...optionFunc) (*discordgo.Role, error) {
	_, span := global.Tracer("cacophony.dev/kit").Start(e.Context(), "event.FindRole")
	defer span.End()

	if e.Type != MessageCreateType {
		return nil, errors.New("event has to be MessageCreate")
	}

	for _, field := range e.Fields() {
		channel, err := e.State().RoleFromMention(e.GuildID, field)
		if err == nil {
			return channel, nil
		}
	}

	return nil, state.ErrRoleNotFound
}

type options struct {
	disableFallbackToSelf bool
}

type optionFunc func(*options)

// nolint: golint
func WithoutFallbackToSelf() optionFunc {
	return optionFunc(func(o *options) {
		o.disableFallbackToSelf = true
	})
}

func getOptions(opts []optionFunc) options {
	value := options{}

	for _, o := range opts {
		o(&value)
	}

	return value
}
