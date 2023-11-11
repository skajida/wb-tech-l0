package main

import (
	"wb-tech-l0/internal/controller"

	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

func mustConnect(logger *zap.Logger) stan.Conn {
	conn, err := stan.Connect("test-cluster", "sub", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		logger.Fatal("Subscriber connection error", zap.Error(err))
	}
	return conn
}

func mustSubscribe(logger *zap.Logger, subject string, connection stan.Conn, subscriber *controller.NatsSubscriber) stan.Subscription {
	sub, err := connection.Subscribe(subject, subscriber.AddOrderHandle)
	if err != nil {
		logger.Fatal("Subscriber subscription error", zap.Error(err))
	}
	return sub
}
