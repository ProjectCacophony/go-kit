package events

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"gocloud.dev/gcerrors"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/rabbitpubsub"
)

type Publisher struct {
	topic *pubsub.Topic
	db    *gorm.DB
}

func NewPublisher(
	amqpDSN string,
	db *gorm.DB,
) (*Publisher, error) {
	p := &Publisher{
		db: db,
	}

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

func (p *Publisher) Publish( // nolint: golint
	ctx context.Context,
	event *Event,
) (err error, recoverable bool) {
	recoverable = true

	body, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(
			err,
			"error marshalling event",
		), recoverable
	}

	return p.PublishRaw(ctx, body)
}

func (p *Publisher) PublishRaw( // nolint: golint
	ctx context.Context,
	body []byte,
) (err error, recoverable bool) {
	recoverable = true

	err = p.topic.Send(
		ctx,
		&pubsub.Message{
			Body: body,
		},
	)
	if err != nil {
		if gcerrors.Code(err) == gcerrors.FailedPrecondition {
			recoverable = false
		}
	}

	return err, recoverable
}

// scheduledEventModel is maintained by the Worker
type scheduledEventModel struct {
	gorm.Model
	Body        []byte
	PublishAt   time.Time
	PublishedAt *time.Time
}

func (*scheduledEventModel) TableName() string {
	return "events_scheduled"
}

func (p *Publisher) PublishAt(
	ctx context.Context,
	event *Event,
	publishAt time.Time,
) error {
	body, err := json.Marshal(event)
	if err != nil {
		return errors.Wrap(
			err,
			"error marshalling event",
		)
	}

	return p.db.Create(&scheduledEventModel{
		Body:        body,
		PublishAt:   publishAt,
		PublishedAt: nil,
	}).Error
}
