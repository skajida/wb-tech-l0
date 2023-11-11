package main

import (
	"context"
	"fmt"
	"net/http"
	"wb-tech-l0/internal/controller"
	getorderinfo "wb-tech-l0/internal/handler/get-order-info"
	"wb-tech-l0/internal/repository"
	"wb-tech-l0/internal/service/cacher"
	"wb-tech-l0/internal/service/updater"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction())
	appCfg := initConfig(logger)

	database := initPostgres(logger, appCfg.Db)
	pgRepository := repository.NewPgDatabase(database)

	updater := updater.NewService(pgRepository)
	subscriber := controller.NewSubscriber(logger, updater)
	connection := mustConnect(logger)
	defer connection.Close()
	subscription := mustSubscribe(logger, appCfg.NatsSubject, connection, subscriber)
	defer subscription.Close()

	cacher := cacher.NewService(repository.NewInMemoryDatabase(), pgRepository)
	handler := getorderinfo.NewHandler(logger, cacher)

	router := mux.NewRouter()
	router.HandleFunc("/get/{order_uid}", handler.GetOrderInfoHandle).Methods(http.MethodGet)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", appCfg.Port),
		Handler: router,
	}

	logger.Info("REST server is listening", zap.Uint16("port", appCfg.Port))
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal("Internal REST server error", zap.Error(err))
	}
	server.Shutdown(context.Background())
}
