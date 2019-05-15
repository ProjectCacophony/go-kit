package events

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/rabbitpubsub"
)

type Publisher struct {
	topic *pubsub.Topic
}

func NewPublisher(
	amqpDSN string,
) (*Publisher, error) {
	p := &Publisher{}

	rabbitConn, err := amqp.Dial(amqpDSN)
	if err != nil {
		return nil, err
	}

	amqpChannel, err := rabbitConn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "cannot open channel")
	}

	err = declareExchange(amqpChannel)
	if err != nil {
		return nil, errors.Wrap(err, "cannot declare exchange")
	}

	p.topic = rabbitpubsub.OpenTopic(
		rabbitConn,
		exchangeName,
		nil,
	)

	return p, nil
}

func (p *Publisher) Publish(
	ctx context.Context,
	event *Event,
) error {
	body, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(
			err,
			"error marshalling event",
		)
	}

	return p.topic.Send(
		ctx,
		&pubsub.Message{
			Body: body,
		},
	)
}
