package events

import (
	"github.com/streadway/amqp"
)

const exchangeName = "cacophony"

func declareExchange(channel *amqp.Channel) error {
	return channel.ExchangeDeclare(
		exchangeName,
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
}
