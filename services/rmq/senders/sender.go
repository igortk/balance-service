package senders

import (
	"fmt"
	"github.com/streadway/amqp"
)

type Sender struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewSender(connection *amqp.Connection) (*Sender, error) {
	channel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed open cannel for sender: %v", err)
	}
	return &Sender{
		connection: connection,
		channel:    channel,
	}, nil
}

func (s *Sender) SendMessage(exchange, routingKey string, message []byte) error {
	err := s.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)

	if err != nil {
		return fmt.Errorf("failed send message: %v", err)
	}

	return nil
}

func (s *Sender) Close() {
	s.channel.Close()
	s.connection.Close()
}
