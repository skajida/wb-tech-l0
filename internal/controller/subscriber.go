package controller

import (
	"context"
	"encoding/json"
	"wb-tech-l0/internal/model"

	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

type NatsSubscriber struct {
	logger  *zap.Logger
	service updater
}

func NewSubscriber(logger *zap.Logger, service updater) *NatsSubscriber {
	return &NatsSubscriber{logger: logger, service: service}
}

func (ns NatsSubscriber) AddOrderHandle(msg *stan.Msg) {
	var order model.Order
	if err := json.Unmarshal(msg.Data, &order); err != nil {
		ns.logger.Error("Unmarshalling error", zap.Error(err))
		return
	}
	switch err := ns.service.AddOrder(context.Background(), order); err {
	case nil:
	case model.ErrOrderConflict:
		ns.logger.Error("Add order conflict", zap.Error(err))
	default:
		ns.logger.Error("Internal error", zap.Error(err))
	}
}
