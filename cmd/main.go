package main

import (
	"balance-service/config"
	"balance-service/services"
	"balance-service/services/pg"
	"balance-service/util"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.GetConfig()
	util.IsError(err, config.ErrLoadConfig)

	initLogs(&cfg.LoggerConfig)

	initServer(cfg)

}
func initServer(cfg *config.Config) {
	server := services.NewServer(&cfg.PostreSqlConfig, &cfg.RabbitConfig)
	server.InitSenders()
	server.InitHandlers()
	server.InitConsumers()
	server.Run()
	log.Printf("qwertyuiop[]asdfghjkl;'zxcvbnm,./")
}

func initLogs(cfg *config.LoggerConfig) {
	logLvl, err := log.ParseLevel(cfg.Level)
	util.IsError(err, config.ErrParseLog)
	log.SetLevel(logLvl)
}

func initPgClient(cfg *config.PostreSqlConfig) pg.PgClient {
	dbClient := pg.NewClient(cfg)
	return *dbClient
}
