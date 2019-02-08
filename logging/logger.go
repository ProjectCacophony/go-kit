package logging

import (
	"net/http"

	"go.uber.org/zap"
)

// NewLogger creates a zap.DiscordgoLogger based on the environment
func NewLogger(env Environment, service, discordWebhookURL string, client *http.Client) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	switch env {
	case ProductionEnvironment:

		logger, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	default:

		logger, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	}

	logger = logger.With(zap.String("service", service))

	logger = logger.WithOptions(zap.Hooks(
		NewZapHookDiscord(
			service, discordWebhookURL, client,
		),
	))

	// TODO: add discord hook

	return logger, nil
}
