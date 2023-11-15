package main

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"

	"github.com/shijuvar/go-distributed-sys/examples/nats-micro-group/model"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	srv, err := micro.AddService(nc, micro.Config{
		Name:        "ordersvc",
		Version:     "0.0.1",
		Description: "Order service",
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Created service: %s (%s)\n", srv.Info().Name, srv.Info().ID)
	orderSvcGroup := srv.AddGroup("ordersvc")
	orderSvcGroup.AddEndpoint("createorder", micro.HandlerFunc(handleCreateOrder))
	orderSvcGroup.AddEndpoint("getorder", micro.HandlerFunc(handleGetOrder))
	runtime.Goexit()
}

func handleCreateOrder(req micro.Request) {
	var orderReq model.Order
	_ = json.Unmarshal([]byte(req.Data()), &orderReq)
	order := model.Order{
		OrderID:    orderReq.OrderID,
		CustomerID: orderReq.CustomerID,
		Status:     "Created",
	}
	req.RespondJSON(order)
}

func handleGetOrder(req micro.Request) {
	var orderReq model.GetOrderRequest
	_ = json.Unmarshal([]byte(req.Data()), &orderReq)
	order := model.Order{
		OrderID: orderReq.OrderID,
	}
	req.RespondJSON(order)
}
