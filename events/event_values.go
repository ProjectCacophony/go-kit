package events

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

// Context returns the context for the event
func (e *Event) Context() context.Context {
	if e.ctx == nil {
		e.ctx = context.Background()

		return e.ctx
	}

	return e.ctx
}

// WithContext sets the context for the event
func (e *Event) WithContext(ctx context.Context) {
	e.ctx = ctx
}

func (e *Event) WithTokens(tokens map[string]string) {
	e.tokens = tokens
}

// Discord gets the Discord API client for the events bot
func (e *Event) Discord() *discordgo.Session {
	if e.discordSession != nil {
		return e.discordSession
	}

	if e.BotUserID == "" {
		panic("could not create discordgo session, no bot user ID set")
	}

	session, err := discord.NewSession(e.tokens, e.BotUserID)
	if err != nil {
		panic(errors.Wrap(err, "could not create discordgo session"))
	}

	e.discordSession = session

	return e.discordSession
}

// WithLogger stores a logger in the event
func (e *Event) WithLogger(logger *zap.Logger) {
	e.logger = logger
}

// Logger retrieves the logger from the event
func (e *Event) Logger() *zap.Logger {
	return e.logger
}

// WithLocalisations stores the localisations in the event
func (e *Event) WithLocalisations(localisations []interfaces.Localisation) {
	e.localisations = localisations
}

// Localisations retrieves the localisations from the event
func (e *Event) Localisations() []interfaces.Localisation {
	return e.localisations
}

// WithState stores the state in the event
func (e *Event) WithState(state *state.State) {
	e.state = state
}

// State retrieves the state from the event
func (e *Event) State() *state.State {
	return e.state
}
