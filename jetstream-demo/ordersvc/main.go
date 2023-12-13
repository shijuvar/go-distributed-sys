package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
	"log/slog"
	"runtime"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"

	"github.com/shijuvar/go-distributed-sys/jetstream-demo/order"
)

const (
	orderStream = "ORDERS"
	subject     = "ORDERS.created"
)

var (
	logger *slog.Logger
	ctx    context.Context
	nc     *nats.Conn
)

func main() {
	ctx = context.Background()
	logger = slog.Default()
	nc, _ = nats.Connect(nats.DefaultURL)
	srv, err := micro.AddService(nc, micro.Config{
		Name:        "ordersvc",
		Version:     "0.0.1",
		Description: "Order service",
	})

	if err != nil {
		logger.Error("Failed to create service", slog.String("error", err.Error()))
		return
	}
	logger.Info("Created service: %s (%s)\n", slog.String("name", srv.Info().Name), slog.String("ID", srv.Info().ID))
	orderSvcGroup := srv.AddGroup("ordersvc")
	orderSvcGroup.AddEndpoint("createorder", micro.HandlerFunc(handleCreateOrder))
	orderSvcGroup.AddEndpoint("getorder", micro.HandlerFunc(handleGetOrder))
	runtime.Goexit()
}

func handleCreateOrder(req micro.Request) {
	var order order.Order
	_ = json.Unmarshal([]byte(req.Data()), &order)
	id, _ := uuid.NewUUID()
	order.ID = id.String()
	order.Status = "Created"
	order.Amount = order.GetAmount()
	publishOrderCreatedEvent(order)
	req.RespondJSON(order)
}
func publishOrderCreatedEvent(command order.Order) {
	js, _ := jetstream.New(nc)
	orderJSON, _ := json.Marshal(command)
	// Publish message on subject
	// The subject has to belong to a stream
	ack, err := js.PublishMsg(ctx, &nats.Msg{
		Data:    orderJSON,
		Subject: subject, // "ORDERS.created"
	})
	if err != nil {
		logger.Error("Failed to publish message", slog.String("subject", subject))
		return
	}
	logger.Info("Published message",
		slog.Uint64("sequence", ack.Sequence),
		slog.String("stream", ack.Stream),
	)
}

func handleGetOrder(req micro.Request) {

}
