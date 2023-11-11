package cacher

import (
	"context"
	"wb-tech-l0/internal/model"
)

type OrderInfoService struct {
	cache cacheRepository
	repo  repository
}

func NewService(cache cacheRepository, repo repository) *OrderInfoService {
	return &OrderInfoService{cache: cache, repo: repo}
}

func (ois OrderInfoService) GetOrderInfo(ctx context.Context, orderId model.OrderId) (model.Order, error) {
	if order, err := ois.cache.GetOrderInfo(ctx, orderId); err == nil {
		return order, nil
	}
	switch order, err := ois.repo.GetOrderInfo(ctx, orderId); err {
	case nil:
		ois.cache.AddOrder(ctx, order)
		return order, nil
	case model.ErrOrderBadParam:
		fallthrough
	default:
		return model.Order{}, err
	}
}
