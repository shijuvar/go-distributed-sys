package main

import (
	"encoding/json"
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

	defer nc.Close()

	// Send the request for getting a Reply
	msg, err := nc.Request("discovery", []byte("ordersvc"), 10*time.Millisecond)
	fmt.Println(string(msg.Data))

	order := Order{
		OrderID:    "1234-5678-90",
		CustomerID: "shijuvar",
		Status:     "Placed",
	}
	jsonOrder, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	// Publish message
	nc.Publish("order.created", jsonOrder)
	log.Println("Message published")

	// Close connection
	nc.Close()
}
