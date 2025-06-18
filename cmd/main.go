package main

import (
	"balance-service/config"
	"balance-service/services"
	"balance-service/services/pg"
	"balance-service/services/rmq/senders"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Error load configuration: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Infoln("Shutting down gracefully...")
		cancel()
	}()

	conn, err := amqp.Dial(fmt.Sprintf(config.RmqUrlConnectionPattern, cfg.RabbitConfig.Username, cfg.RabbitConfig.Password, cfg.RabbitConfig.Host, cfg.RabbitConfig.Port))
	if err != nil {
		log.Fatalf("Error connect RabbitMq: %v", err)
	}

	sender, err := senders.NewSender(conn)
	if err != nil {
		log.Fatalf("Error create rmq sender: %v", err)
	}

	pgClient, err := pg.NewClient(&cfg.PostreSqlConfig)
	if err != nil {
		log.Fatalf("Error create pg client: %v", err)
	}

	svr := services.NewServer2(pgClient, sender, conn)
	wg := &sync.WaitGroup{}
	svr.Run(ctx, wg)

	wg.Wait()
	log.Infoln("Shutdown complete.")

}
