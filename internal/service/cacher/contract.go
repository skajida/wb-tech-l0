//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package cacher

import (
	"context"
	"wb-tech-l0/internal/model"
)

type repository interface {
	GetOrderInfo(context.Context, model.OrderId) (model.Order, error)
}

type cacheRepository interface {
	repository
	AddOrder(context.Context, model.Order) error
}
