//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package updater

import (
	"context"
	"wb-tech-l0/internal/model"
)

type repository interface {
	AddOrder(context.Context, model.Order) error
}
