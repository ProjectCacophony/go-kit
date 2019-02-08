package events

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

// Type defines the type for a Event
type Type string

// defines various Event Types
const (
	MessageCreateEventType Type = "message_create"
	MessageUpdateEventType Type = "message_update"
	MessageDeleteEventType Type = "message_delete"
)

// Event represents an Event
type Event struct {
	Type       Type
	ReceivedAt time.Time
	BotUserID  string

	// discordgo event data
	MessageCreate *discordgo.MessageCreate
	MessageUpdate *discordgo.MessageUpdate
	MessageDelete *discordgo.MessageDelete
}

// GenerateRoutingKey generates an Routing Key for AMQP based on a Event Type
func GenerateRoutingKey(eventType Type) string {
	return fmt.Sprintf("cacophony.discord.%s", eventType)
}

// GenerateEventFromDiscordgoEvent generates an Event from a Discordgo Event
func GenerateEventFromDiscordgoEvent(botUserID string, eventItem interface{}) (*Event, error) {
	event := &Event{
		ReceivedAt: time.Now(),
		BotUserID:  botUserID,
	}

	switch t := eventItem.(type) {
	case *discordgo.MessageCreate:
		event.Type = MessageCreateEventType
		event.MessageCreate = t

		return event, nil
	case *discordgo.MessageUpdate:
		event.Type = MessageUpdateEventType
		event.MessageUpdate = t

		return event, nil
	case *discordgo.MessageDelete:
		event.Type = MessageDeleteEventType
		event.MessageDelete = t

		return event, nil
	}

	return nil, errors.New("event type is not supported")
}
