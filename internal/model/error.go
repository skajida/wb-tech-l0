package model

import "errors"

var (
	ErrOrderConflict = errors.New("order already exists")
	ErrOrderBadParam = errors.New("order doesn't exist")
)
