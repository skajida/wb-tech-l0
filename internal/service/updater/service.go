package updater

import (
	"context"
	"wb-tech-l0/internal/model"
)

type NatsService struct {
	repo repository
}

func NewService(repo repository) *NatsService {
	return &NatsService{repo: repo}
}

func (ns NatsService) AddOrder(ctx context.Context, order model.Order) error {
	return ns.repo.AddOrder(ctx, order)
}
