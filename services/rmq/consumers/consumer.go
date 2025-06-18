package consumers

import (
	"balance-service/services/rmq/handlers"
	"context"
	log "github.com/sirupsen/logrus"
	"sync"

	"github.com/streadway/amqp"
)

type Consumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
	handler    handlers.MessageHandler
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
		connection: connection,
		channel:    channel,
		queue:      &queue,
		handler:    handler,
	}
}

func (c *Consumer) ConsumeMessages(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	msgs, err := c.channel.Consume(
		c.queue.Name,
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
		<-ctx.Done()
		log.Println("Context cancelled, stopping consumer...")
		_ = c.channel.Cancel("", false)
	}()

	for d := range msgs {
		c.handler.HandleMessage(d.Body)
	}
}

func (c *Consumer) Close() {
	c.channel.Close()
	c.connection.Close()
}
