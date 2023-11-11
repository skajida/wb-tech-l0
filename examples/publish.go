package main

import (
	"io"
	"os"

	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

func main() {
	const natsSubject = "test-cluster"

	logger := zap.Must(zap.NewProduction())
	connection, err := stan.Connect(natsSubject, "pub", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		logger.Fatal("Publisher connection error")
	}
	defer connection.Close()

	msg, err := io.ReadAll(os.Stdin)
	if err != nil {
		logger.Fatal("Publisher message reading error")
	}

	if err = connection.Publish("addOrder", msg); err != nil {
		logger.Fatal("Publisher publishing error")
	}
}
