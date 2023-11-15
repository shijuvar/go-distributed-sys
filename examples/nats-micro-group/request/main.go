package main

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/shijuvar/go-distributed-sys/examples/nats-micro-group/model"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	order := model.CreateOrderRequest{
		Order: model.Order{
			OrderID:    "10001",
			CustomerID: "shijuvar",
		},
	}

	orderReq, _ := json.Marshal(order)
	fmt.Println("ordersvc.createorder")
	msg, err := nc.Request("ordersvc.createorder", orderReq, time.Second*30)
	if err != nil {
		log.Fatalf("error:", err)
	}
	response := decode(msg)
	fmt.Printf("%v\n", response)

	getOrderReq, _ := json.Marshal(model.GetOrderRequest{
		OrderID: "10001",
	})
	fmt.Println("ordersvc.getorder")
	msg, err = nc.Request("ordersvc.getorder", getOrderReq, time.Second*30)
	if err != nil {
		log.Fatalf("error:", err)
	}
	response = decode(msg)
	fmt.Printf("%v\n", response)
	runtime.Goexit()
}

func decode(msg *nats.Msg) model.Order {
	var res model.Order
	json.Unmarshal(msg.Data, &res)
	return res
}
