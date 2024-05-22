package consumers

import (
	"balance-service/services/rmq/handlers"
	log "github.com/sirupsen/logrus"

	"github.com/streadway/amqp"
)

type Consumer struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue
	Handler    handlers.MessageHandler
}

func NewConsumer(connection *amqp.Connection, exchange, routingKey, queueName string, handler handlers.MessageHandler) *Consumer {
	channel, err := connection.Channel()
	if err != nil {
		return nil
	}
	err = channel.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil
	}
	queue, err := channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil
	}

	err = channel.QueueBind(
		queue.Name,
		routingKey,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil
	}

	return &Consumer{
		Connection: connection,
		Channel:    channel,
		Queue:      queue,
		Handler:    handler,
	}
}

func (c *Consumer) ConsumeMessages() {
	msgs, err := c.Channel.Consume(
		c.Queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumers: %v", err)
	}
	go func() {
		for d := range msgs {
			c.Handler.HandleMessage(d.Body)
		}
	}()
}

func (c *Consumer) Close() {
	c.Channel.Close()
	c.Connection.Close()
}
