package consumers

import (
	"balance-service/services/rmq/handlers"
	"context"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"

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
		amqp.ExchangeTopic,
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

	stopChant := make(chan struct{})

	go func() {
		<-ctx.Done()

		log.Infof("Context cancelled, stopping consumer [%s] (waiting time is 10 seconds)...", c.queue.Name)
		_ = c.channel.Cancel("", false)

		<-time.After(10 * time.Second)
		stopChant <- struct{}{}
	}()

	for {
		select {
		case val, ok := <-msgs:
			if !ok {
				log.Println("RabbitMQ channel closed.")
				return
			}
			c.handler.HandleMessage(val.Body)

		case <-stopChant:
			log.Println("Stopped consuming messages.")
			return
		}
	}
}

func (c *Consumer) Close() {
	c.channel.Close()
	c.connection.Close()
}
