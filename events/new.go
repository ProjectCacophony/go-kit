package events

import (
	"time"

	"github.com/google/uuid"
)

func New(eventType Type) (*Event, error) {
	eventID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Event{
		ID:          eventID.String(),
		Type:        eventType,
		ReceivedAt:  time.Now(),
		SpanContext: make(map[string]string),
	}, nil
}
