package main

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

func initPostgres(logger *zap.Logger, dbCfg dbConfig) *sql.DB {
	database, err := sql.Open(
		"postgres",
		fmt.Sprintf(
			"user=%s password=%s sslmode=disable",
			dbCfg.User,
			dbCfg.Password,
		),
	)
	if err != nil {
		logger.Fatal("Internal PostgreSQL initialization error", zap.Error(err))
	}
	return database
}
