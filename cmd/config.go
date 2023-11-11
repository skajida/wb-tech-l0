package main

import (
	"github.com/caarlos0/env/v9"
	"go.uber.org/zap"
)

type dbConfig struct {
	User     string `env:"PG_USER" envDefault:"postgres"`
	Password string `env:"PG_PASSWORD"`
}

type appConfig struct {
	Db          dbConfig
	NatsSubject string `env:"NATS_SUBJECT" envDefault:"addOrder"`
	Port        uint16 `env:"SERVICE_PORT,notEmpty"`
}

func initConfig(logger *zap.Logger) (appCfg appConfig) {
	if err := env.Parse(&appCfg); err != nil {
		logger.Fatal("Internal application config initialization error", zap.Error(err))
	}
	return
}
