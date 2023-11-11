package repository

import (
	"context"
	"sync"
	"wb-tech-l0/internal/model"
)

type InMemoryDatabase struct {
	sync.RWMutex
	orders map[model.OrderId]model.Order
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{orders: make(map[model.OrderId]model.Order)}
}

func (imdb *InMemoryDatabase) GetOrderInfo(ctx context.Context, orderId model.OrderId) (model.Order, error) {
	imdb.RLock()
	defer imdb.RUnlock()
	if order, exists := imdb.orders[orderId]; exists {
		return order, nil
	}
	return model.Order{}, model.ErrOrderBadParam
}

func (imdb *InMemoryDatabase) AddOrder(ctx context.Context, order model.Order) error {
	imdb.Lock()
	defer imdb.Unlock()
	if _, exists := imdb.orders[order.Uid]; !exists {
		imdb.orders[order.Uid] = order
		return nil
	}
	return model.ErrOrderConflict
}
