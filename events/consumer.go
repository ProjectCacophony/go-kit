package events

import (
	"context"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"gocloud.dev/pubsub"
	"gocloud.dev/pubsub/rabbitpubsub"
)

// TODO: shutdown logic

// Consumer processes incoming events
type Consumer struct {
	logger                    *zap.Logger
	serviceName               string
	amqpDSN                   string
	amqpExchangeName          string
	amqpRoutingKey            string
	concurrentProcessingLimit int
	handler                   func(*Event) error

	amqpConnection *amqp.Connection
	subscription   *pubsub.Subscription
}

// NewConsumer creates a new processor
func NewConsumer(
	logger *zap.Logger,
	serviceName string,
	amqpDSN string,
	amqpExchangeName string,
	amqpRoutingKey string,
	concurrentProcessingLimit int,
	handler func(*Event) error,
) (*Consumer, error) {
	processor := &Consumer{
		logger:                    logger,
		serviceName:               serviceName,
		amqpDSN:                   amqpDSN,
		amqpExchangeName:          amqpExchangeName,
		amqpRoutingKey:            amqpRoutingKey,
		concurrentProcessingLimit: concurrentProcessingLimit,
		handler:                   handler,
	}

	err := processor.init()
	if err != nil {
		return nil, err
	}

	return processor, nil
}

// init declares the exchange, the queue, and the queue binding
func (c *Consumer) init() error {
	var err error

	c.amqpConnection, err = amqp.Dial(c.amqpDSN)
	if err != nil {
		return errors.Wrap(err, "unable to initialise AMQP session")
	}

	amqpChannel, err := c.amqpConnection.Channel()
	if err != nil {
		return errors.Wrap(err, "cannot open channel")
	}

	err = amqpChannel.ExchangeDeclare(
		c.amqpExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "cannot declare exchange")
	}

	ampqQueue, err := amqpChannel.QueueDeclare(
		c.serviceName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "cannot declare queue")
	}

	err = amqpChannel.QueueBind(
		ampqQueue.Name,
		c.amqpRoutingKey,
		c.amqpExchangeName,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "cannot bind queue")
	}

	c.subscription = rabbitpubsub.OpenSubscription(
		c.amqpConnection,
		c.serviceName,
		nil,
	)

	return nil
}

// Start starts processing events
func (c *Consumer) Start(ctx context.Context) error {
	return c.start(ctx)
}

func (c *Consumer) start(ctx context.Context) error {
	// keep semaphore channel to limit the amount of events being processed concurrently
	semaphore := make(chan interface{}, c.concurrentProcessingLimit)

	for {
		delivery, err := c.subscription.Receive(context.Background())
		if err != nil {
			c.logger.Error(
				"error receiving event",
				zap.Error(err),
			)
			break
		}

		select {
		case semaphore <- nil:
		case <-ctx.Done():
			break
		}

		go func(d *pubsub.Message) {
			defer func() {
				// clear channel when completed
				<-semaphore
			}()

			err := c.handle(d)
			if err != nil {
				c.logger.Error("failed to handle event",
					zap.Error(err),
				)
			}
		}(delivery)
	}

	// try to fill channel with amount of buffer size, this makes sure we will wait for all events to finish processing
	for i := 0; i < cap(semaphore); i++ {
		semaphore <- nil
	}

	c.logger.Info("finished Start()")

	return errors.New("unexpected shutdown")
}
