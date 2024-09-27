package messaging

import amqp "github.com/rabbitmq/amqp091-go"

// DeclareExchangeAndQueue creates an exchange and binds a queue to it
func DeclareExchangeAndQueue(channel *amqp.Channel, exchangeName, exchangeType, queueName, routingKey string) error {
	// Declare exchange
	err := channel.ExchangeDeclare(
		exchangeName, // Name of the exchange
		exchangeType, // Type of exchange
		true,         // Durable
		false,        // Auto-deleted
		false,        // Internal
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		return err
	}

	// Declare queue
	_, err = channel.QueueDeclare(
		queueName, // Name of the queue
		true,      // Durable
		false,     // Auto-deleted
		false,     // Exclusive
		false,     // No-wait
		nil,       // Arguments
	)
	if err != nil {
		return err
	}

	// Bind the queue to the exchange
	err = channel.QueueBind(
		queueName,    // Name of the queue
		routingKey,   // Routing key
		exchangeName, // Name of the exchange
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		return err
	}

	return nil
}
