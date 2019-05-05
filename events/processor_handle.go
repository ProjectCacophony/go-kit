package events

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

func (p *Processor) handle(delivery amqp.Delivery) error {
	var event Event
	err := json.Unmarshal(delivery.Body, &event)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal event")
	}

	err = delivery.Ack(false)
	if err != nil {
		return errors.Wrap(err, "failed to ack event")
	}

	return p.handler(&event)
}
