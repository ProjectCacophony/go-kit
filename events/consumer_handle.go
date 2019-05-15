package events

import (
	"encoding/json"

	"github.com/pkg/errors"
	"gocloud.dev/pubsub"
)

func (c *Consumer) handle(delivery *pubsub.Message) error {
	var event Event
	err := json.Unmarshal(delivery.Body, &event)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal event")
	}

	delivery.Ack()

	return c.handler(&event)
}
