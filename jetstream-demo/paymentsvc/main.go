package main

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	ordermodel "github.com/shijuvar/go-distributed-sys/jetstream-demo/order"
	"log/slog"
)

const (
	orderStream    = "ORDERS"
	consumeSubject = "ORDERS.created"
	publishSubject = "ORDERS.paymentdebited"

	batch = 1 // small batch just for the same of demo
)

var (
	logger *slog.Logger
	ctx    context.Context
	nc     *nats.Conn
	js     jetstream.JetStream
)

func main() {
	logger = slog.Default()
	// Creating a `context.Context` object
	ctx = context.Background()
	nc, _ = nats.Connect(nats.DefaultURL)
	// Create a JetStream instance
	js, _ = jetstream.New(nc)
	// create a consumer on stream by filtering subjects
	cons, _ := js.CreateOrUpdateConsumer(ctx, orderStream, jetstream.ConsumerConfig{
		Durable:        "paymentsvc",
		AckPolicy:      jetstream.AckExplicitPolicy,
		FilterSubjects: []string{consumeSubject},
	})
	// Fetching a batch of messages
	msgs, _ := cons.Fetch(batch)
	for msg := range msgs.Messages() {
		msg.Ack()
		var orderCommand ordermodel.Order
		// Unmarshal JSON that represents the Order data
		err := json.Unmarshal(msg.Data(), &orderCommand)
		if err != nil {
			logger.Error("Error on decoding JSON")
			return
		}
		logger.Info("Received payment command", slog.Any("command", orderCommand))
		// TODO: Payment transaction processor actions
		command := ordermodel.PaymentDebitedCommand{
			OrderID:    orderCommand.ID,
			CustomerID: orderCommand.CustomerID,
			Amount:     orderCommand.Amount,
		}
		// publish PaymentDebited event
		publishPaymentDebitedEvent(command)
	}
}

func publishPaymentDebitedEvent(command ordermodel.PaymentDebitedCommand) {
	commandJSON, _ := json.Marshal(command)
	// Publish message on subject
	ack, err := js.PublishMsg(ctx, &nats.Msg{
		Data:    commandJSON,
		Subject: publishSubject,
	})
	if err != nil {
		logger.Error("Failed to publish message", slog.String("subject", publishSubject))
		return
	}
	logger.Info("Published message",
		slog.Uint64("sequence", ack.Sequence),
		slog.String("subject", publishSubject),
		slog.String("stream", ack.Stream),
	)
}
