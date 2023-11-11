//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package controller

import (
	"context"
	"wb-tech-l0/internal/model"
)

type updater interface {
	AddOrder(context.Context, model.Order) error
}
