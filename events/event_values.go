package events

import (
	"context"
	"net/http"

	"github.com/jinzhu/gorm"

	"github.com/go-redis/redis"

	"gitlab.com/Cacophony/go-kit/paginator"

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
func (e *Event) Discord() *discord.Session {
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

// WithStorage stores the storage bucket in the event
func (e *Event) WithStorage(storage *Storage) {
	e.storage = storage
}

func (e *Event) Storage() *Storage {
	return e.storage
}

// WithLocalizations stores the localizations in the event
func (e *Event) WithLocalizations(localizations []interfaces.Localization) {
	e.localizations = localizations
}

// Localizations retrieves the localizations from the event
func (e *Event) Localizations() []interfaces.Localization {
	return e.localizations
}

// WithState stores the state in the event
func (e *Event) WithState(state *state.State) {
	e.state = state
}

// State retrieves the state from the event
func (e *Event) State() *state.State {
	return e.state
}

// WithPaginator stores the Paginator in the event
func (e *Event) WithPaginator(paginator *paginator.Paginator) {
	e.paginator = paginator
}

// Paginator retrieves the Paginator from the event
func (e *Event) Paginator() *paginator.Paginator {
	return e.paginator
}

// WithRedis stores the Redis Client in the event
func (e *Event) WithRedis(redisClient *redis.Client) {
	e.redisClient = redisClient
}

// Redis retrieves the Redis Client from the event
func (e *Event) Redis() *redis.Client {
	return e.redisClient
}

// WithDB stores the DB Client in the event
func (e *Event) WithDB(db *gorm.DB) {
	e.db = db
}

// DB retrieves the DB Client from the event
func (e *Event) DB() *gorm.DB {
	return e.db
}

func (e *Event) WithHTTPClient(client *http.Client) {
	e.httpClient = client
}

func (e *Event) HTTPClient() *http.Client {
	return e.httpClient
}

func (e *Event) WithQuestionnaire(questionnaire *Questionnaire) {
	e.questionnaire = questionnaire
}

func (e *Event) Questionnaire() *Questionnaire {
	return e.questionnaire
}
