package messaging

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageEmitter interface {
	SendMessage(ctx context.Context, message []byte, exchange string, topic string, routingKey string) error
}

type RabbitmqEmitter struct {
	channel *amqp.Channel
}

func New(channel *amqp.Channel) *RabbitmqEmitter {
	return &RabbitmqEmitter{
		channel: channel,
	}
}

func (r *RabbitmqEmitter) SendMessage(ctx context.Context, message []byte, exchange string, topic string, rountingKey string) error {
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        message,
	}
	return r.channel.PublishWithContext(ctx, exchange, rountingKey, false, false, msg)
}
