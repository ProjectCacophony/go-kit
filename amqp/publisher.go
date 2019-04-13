package amqp

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

const (
	retryLimit = 24
	retryWait  = 5 * time.Second
)

// Publisher publishes Events to AMQP
type Publisher struct {
	logger           *zap.Logger
	amqpDSN          string
	amqpExchangeName string
	eventTTL         time.Duration

	amqpConnection *amqp.Connection
	amqpChannel    *amqp.Channel
	amqpExpiration string

	sendLock sync.Mutex
}

// NewPublisher creates new Publisher
func NewPublisher(
	logger *zap.Logger,
	amqpDSN string,
	amqpExchangeName string,
	eventTTL time.Duration,
) (*Publisher, error) {
	publisher := &Publisher{
		logger:           logger,
		amqpDSN:          amqpDSN,
		amqpExchangeName: amqpExchangeName,
		eventTTL:         eventTTL,
	}

	if eventTTL > 0 {
		publisher.amqpExpiration = strconv.Itoa(int(eventTTL.Seconds() * 1000))
	}

	err := publisher.init()
	if err != nil {
		return nil, err
	}

	return publisher, nil
}

// init initialises the channel and exchange
func (p *Publisher) init() error {
	var err error

	p.amqpConnection, err = amqp.Dial(p.amqpDSN)
	if err != nil {
		return err
	}

	p.amqpChannel, err = p.amqpConnection.Channel()
	if err != nil {
		return errors.Wrap(err, "cannot open channel")
	}

	err = p.amqpChannel.ExchangeDeclare(
		p.amqpExchangeName,
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

	return nil
}

// reconnect attempts to reconnect to the AMQP Broker
func (p *Publisher) reconnect() {
	p.amqpChannel.Close()
	p.amqpConnection.Close()

	for i := 0; i < retryLimit; i++ {
		time.Sleep(retryWait)

		p.logger.Info("attempting to reconnect")
		err := p.init()
		if err != nil {
			continue
		}

		p.logger.Info("reconnected successfully")

		return
	}

	p.logger.Fatal("could not reconnect")
}

// Publish publishes a specific event
func (p *Publisher) Publish(routingKey string, body json.RawMessage) error {
	// sendLock to: avoid multiple reconnecting at once, and sending to a connection that is known to be dead
	// TODO: instead use channel in the future to queue messages?
	p.sendLock.Lock()
	defer p.sendLock.Unlock()

	if p.amqpChannel == nil {
		return errors.New("no channel setup")
	}

	err := p.publish(routingKey, body)
	if err != nil {
		return err
	}

	return nil
}

// publisher actually sends the event to the AMQP Broker
// it will attempt to reconnect if the sending fails with the error 504
func (p *Publisher) publish(routingKey string, body json.RawMessage) error {
	err := p.amqpChannel.Publish(
		p.amqpExchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "",
			Body:            body,
			DeliveryMode:    amqp.Transient,
			Priority:        0,
			Expiration:      p.amqpExpiration,
		},
	)
	if err != nil {
		// should we attempt to reconnect?
		amqpErr, ok := err.(*amqp.Error)
		if ok && amqpErr.Code == 504 {
			p.logger.Warn("looks like we lost the connection to the AMQP Broker, will attempt to reconnect",
				zap.Error(err),
			)

			p.reconnect()

			// resend our message
			return p.publish(routingKey, body)
		}
	}

	return err
}
