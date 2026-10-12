package main

import (
	"encoding/json"
	"log"
	"runtime"

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
	// Subscribe messages for subject order.created
	nc.Subscribe("order.created", func(m *nats.Msg) {
		//log.Println(m.Data)
		var order Order
		err := json.Unmarshal(m.Data, &order)
		if err == nil {
			log.Printf("Order ID:%s Customer ID:%s Status:%s", order.OrderID, order.CustomerID, order.Status)
		}
	})

	// Subscriber for providing reply to requests
	// A subscriber can optionally give a reply
	nc.Subscribe("discovery", func(m *nats.Msg) {
		log.Printf("Request for:%s", string(m.Data))
		nc.Publish(m.Reply, []byte("http://localhost:8080/orders"))
	})
	runtime.Goexit()
}
