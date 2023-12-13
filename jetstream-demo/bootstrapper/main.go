package main

import (
	"context"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	orderStream         = "ORDERS"
	orderStreamSubjects = "ORDERS.*"
)

func main() {
	logger := slog.Default()
	// Creating a `context.Context` object
	ctx := context.Background()
	nc, _ := nats.Connect(nats.DefaultURL)

	// Create a JetStream instance
	js, _ := jetstream.New(nc)
	// Create a stream
	if _, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     orderStream,                   // "ORDERS"
		Subjects: []string{orderStreamSubjects}, // []string{"ORDERS.*"}
	}); err != nil {
		logger.Error("Failed to create Stream", slog.String("error", err.Error()))
	}
	logger.Info("Created a Stream", slog.String("stream", orderStream))

}
