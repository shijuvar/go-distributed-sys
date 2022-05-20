package main

import (
	"fmt"
	//"encoding/json"
	"log"
	//"runtime"
	"time"

	nats "github.com/nats-io/nats.go"
)

type Order struct {
	OrderID    string
	CustomerID string
	Status     string
}

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	// Encoded connection
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer ec.Close()

	// Send the request for getting a Reply
	msg, err := nc.Request("discovery", []byte("orderservice"), 10*time.Millisecond)
	fmt.Println(string(msg.Data))

	// Publish without encoder
	//nc.Publish("order.created", []byte("OrderID: 1234-5678-80"))
	//log.Println("Message published")

	order := Order{
		OrderID:    "1234-5678-90",
		CustomerID: "shijuvar",
		Status:     "Placed",
	}
	// Publish using JSON encoder
	ec.Publish("order.created", order)
	log.Println("Message published")

	// Close connection
	nc.Close()
}
