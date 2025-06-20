package services

import (
	"balance-service/config"
	"balance-service/services/pg"
	"balance-service/services/rmq/consumers"
	"balance-service/services/rmq/handlers"
	"balance-service/services/rmq/senders"
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"sync"
)

type Server struct {
	pgClient *pg.Client
	sender   *senders.Sender
	conn     *amqp.Connection

	consumers map[string]*consumers.Consumer
}

func NewServer2(pg *pg.Client, s *senders.Sender, c *amqp.Connection) *Server {
	return &Server{
		pgClient:  pg,
		sender:    s,
		conn:      c,
		consumers: make(map[string]*consumers.Consumer),
	}
}

func (s *Server) Run(ctx context.Context, wg *sync.WaitGroup) {
	log.Info("Starting server...")
	//handlers
	emitBalanceHandler := handlers.NewEmitBalanceByUserIdHandler(s.pgClient, s.sender)
	getBalanceByUserIdHandler := handlers.NewGetBalanceByUserIdHandler(s.pgClient, s.sender)
	log.Info("Handlers was prepared")

	//consumers
	s.consumers["GetBalanceByUserIdConsumer"] = consumers.NewConsumer(s.conn, config.RabbitBalanceExchange, config.GetBalanceByUserIdRequestRoutingKey, config.GetBalanceByUserIdRequestQueueName, getBalanceByUserIdHandler)
	s.consumers["EmitBalanceByUserIdConsumer"] = consumers.NewConsumer(s.conn, config.RabbitBalanceExchange, config.EmitBalanceByUserIdRequestRoutingKey, config.EmitUserBalanceRequestQueueName, emitBalanceHandler)
	log.Info("Consumers was prepared")

	s.runAllConsumers(ctx, wg)
	log.Info("Server is running...")

	select {
	case <-ctx.Done():
		log.Info("Stopped server")
		return
	}

}

func (s *Server) runAllConsumers(ctx context.Context, wg *sync.WaitGroup) {
	for _, client := range s.consumers {
		wg.Add(1)
		go client.ConsumeMessages(ctx, wg)
	}
}
