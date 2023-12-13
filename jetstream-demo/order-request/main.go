package main

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/shijuvar/go-distributed-sys/jetstream-demo/order"
)

func main() {
	logger := slog.Default()
	nc, _ := nats.Connect(nats.DefaultURL)
	order := order.Order{
		CustomerID:   "shijuvar",
		RestaurantId: "cgh",
		OrderItems: []order.OrderItem{
			order.OrderItem{
				ProductCode: "vm",
				Name:        "Veg Meals",
				UnitPrice:   200,
				Quantity:    2,
			},
		},
	}

	orderReq, _ := json.Marshal(order)
	msg, err := nc.Request("ordersvc.createorder", orderReq, time.Second*30)
	if err != nil {
		logger.Error("Error on request:", slog.String("error", err.Error()))
		return
	}
	response := decode(msg)
	logger.Info("Received message from service", slog.Any("response", response))
}

func decode(msg *nats.Msg) order.Order {
	var res order.Order
	json.Unmarshal(msg.Data, &res)
	return res
}
