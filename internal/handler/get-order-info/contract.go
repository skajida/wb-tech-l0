//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package getorderinfo

import (
	"context"
	"wb-tech-l0/internal/model"
)

type cacher interface {
	GetOrderInfo(context.Context, model.OrderId) (model.Order, error)
}
