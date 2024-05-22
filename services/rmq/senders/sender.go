package senders

import (
	"balance-service/util"
	"github.com/streadway/amqp"
)

type Sender struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewSender(connection *amqp.Connection) Sender {
	channel, err := connection.Channel()
	util.IsError(err, "Errror Chanel")
	return Sender{
		Connection: connection,
		Channel:    channel,
	}
}

func (s *Sender) SendMessage(exchange, routingKey string, message []byte) {
	err := s.Channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
	util.IsError(err, "err send message")
}

func (s *Sender) Close() {
	s.Channel.Close()
	s.Connection.Close()
}
