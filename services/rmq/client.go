package rmq

import (
	"balance-service/services/rmq/consumers"
	"balance-service/services/rmq/handlers"
	"balance-service/services/rmq/senders"
)

type Client struct {
	sender     senders.Sender
	handler    handlers.MessageHandler
	consumer   consumers.Consumer
	deliveryCh chan []byte
}
