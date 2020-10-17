package events

import (
	"go.opentelemetry.io/otel/label"
)

type SpanContext struct {
	Values map[string]string `json:"values"`
}

func (sp *SpanContext) Get(key string) string {
	if sp.Values == nil {
		return ""
	}

	return sp.Values[key]
}

func (sp *SpanContext) Set(key, value string) {
	if sp.Values == nil {
		sp.Values = make(map[string]string)
	}

	sp.Values[key] = value

}

var (
	SpanLabelEventingType      = label.Key("eventing_type")
	SpanLabelEventingIsCommand = label.Key("eventing_is_command")
	SpanLabelEventingCommand   = label.Key("eventing_command")
	SpanLabelDiscordBotUserID  = label.Key("discord_bot_user_id")
	SpanLabelDiscordGuildID    = label.Key("discord_guild_id")
	SpanLabelDiscordChannelID  = label.Key("discord_channel_id")
	SpanLabelDiscordUserID     = label.Key("discord_user_id")
	SpanLabelDiscordMessageID  = label.Key("discord_message_id")
)
