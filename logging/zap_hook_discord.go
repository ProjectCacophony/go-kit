package logging

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap/zapcore"
)

// NewZapHookDiscord sends Zap log messages to a Discord Webhook
// TODO: ratelimit
func NewZapHookDiscord(serviceName, webhookURL string, client *http.Client) func(zapcore.Entry) error {
	if webhookURL == "" || client == nil {
		return nil
	}

	return func(entry zapcore.Entry) error {
		if entry.Level == zapcore.DebugLevel ||
			entry.Level == zapcore.InfoLevel {
			return nil
		}

		body, err := json.Marshal(discordgo.WebhookParams{
			Username: strings.Title(serviceName),
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Logging message: " + strings.ToUpper(entry.Level.String()),
					Description: entry.Message,
					Timestamp:   entry.Time.Format(time.RFC3339),
					Color:       0, // TODO: color per log level
					Footer: &discordgo.MessageEmbedFooter{
						Text: "powered by Cacophony",
					},
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "caller",
							Value: entry.Caller.String(),
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		_, err = client.Post(
			webhookURL, "application/json", bytes.NewBuffer(body),
		)
		return err
	}

}
