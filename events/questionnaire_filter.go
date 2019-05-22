package events

import (
	"fmt"
)

// QuestionnaireFilter matches the fields with the events
type QuestionnaireFilter struct {
	GuildID   string
	ChannelID string
	UserID    string
	Type      Type
}

func questionnaireRedisKey(eventType Type, guildID, channelID, userID string) string {
	return fmt.Sprintf(
		"cacophony:questionnaire:event-%s:guildid-%s:channelid-%s:userid-%s",
		eventType,
		guildID,
		channelID,
		userID,
	)
}

func questionnaireEventRedisKey(event *Event) string {
	return questionnaireRedisKey(
		event.Type,
		event.GuildID,
		event.ChannelID,
		event.UserID,
	)
}
