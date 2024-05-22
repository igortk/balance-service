package services

import (
	"balance-service/config"
	"balance-service/dto/proto"
	"balance-service/services/pg"
	"balance-service/services/rmq/consumers"
	"balance-service/services/rmq/handlers"
	"balance-service/services/rmq/senders"
	"balance-service/util"
	"fmt"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Server struct {
	pgCl      *pg.PgClient
	rmqConn   *amqp.Connection
	handlers  map[handlers.MessageHandler]handlers.MessageHandler
	sender    senders.Sender
	consumers map[protoreflect.ProtoMessage]*consumers.Consumer
}

func NewServer(cfgPos *config.PostreSqlConfig, cfgRmq *config.RabbitConfig) *Server {
	conn, err := amqp.Dial(fmt.Sprintf(config.RmqUrlConnectionPattern, cfgRmq.Username, cfgRmq.Password, cfgRmq.Host, cfgRmq.Port))
	util.IsError(err, "err")

	return &Server{
		pgCl:      pg.NewClient(cfgPos),
		rmqConn:   conn,
		handlers:  map[handlers.MessageHandler]handlers.MessageHandler{},
		sender:    senders.NewSender(conn),
		consumers: map[protoreflect.ProtoMessage]*consumers.Consumer{},
	}
}
func (s *Server) InitPgClient(cfg *config.PostreSqlConfig) {
	s.pgCl = pg.NewClient(cfg)
}

func (s *Server) InitRmqConnection(cfg *config.RabbitConfig) {
	conn, err := amqp.Dial(fmt.Sprintf(config.RmqUrlConnectionPattern, cfg.Username, cfg.Password, cfg.Host, cfg.Port))
	util.IsError(err, "err")
	s.rmqConn = conn
}

func (s *Server) InitHandlers() {
	s.handlers[handlers.UpdateOrderEventHandler{}] = handlers.NewUpdateOrderEventHandler(s.pgCl)
	s.handlers[handlers.GetBalanceByUserIdHandler{}] = handlers.NewGetBalanceByUserIdHandler(s.pgCl, s.sender)
	s.handlers[handlers.EmmitBalanceByUserIdHandler{}] = handlers.NewEmmitBalanceByUserIdHandler(s.pgCl, s.sender)
}

func (s *Server) InitSenders() {
	//s.senders[&proto.GetBalanceByUserIdResponse{}] = senders.NewSender(s.rmqConn)
}

func (s *Server) InitConsumers() {
	s.consumers[&proto.OrderUpdateEvent{}] = consumers.NewConsumer(
		s.rmqConn,
		config.RabbitEventsExchange,
		config.UpdatedOrderEventRoutingKey,
		config.UpdatedOrderEventQueueName,
		s.handlers[handlers.UpdateOrderEventHandler{}])

	s.consumers[&proto.GetBalanceByUserIdRequest{}] = consumers.NewConsumer(
		s.rmqConn,
		config.RabbitBalanceExchange,
		config.GetBalanceByUserIdRequestRoutingKey,
		config.GetBalanceByUserIdRequestQueueName,
		s.handlers[handlers.GetBalanceByUserIdHandler{}])

	s.consumers[&proto.GetBalanceByUserIdRequest{}] = consumers.NewConsumer(
		s.rmqConn,
		config.RabbitBalanceExchange,
		config.EmitBalanceByUserIdRequestRoutingKey,
		config.EmitUserBalanceRequestQueueName,
		s.handlers[handlers.EmmitBalanceByUserIdHandler{}])
}

func (s *Server) Run() {
	forever := make(chan bool)
	s.runAllConsumers()
	<-forever
}

func (s *Server) runAllConsumers() {
	for _, client := range s.consumers {
		go client.ConsumeMessages()
	}
}
