package main

import (
	"encoding/json"
	"log"
	"runtime"

	"github.com/nats-io/nats.go"

	"github.com/shijuvar/go-distributed-sys/examples/jetstream/model"
)

func main() {
	// Connect to NATS
	nc, _ := nats.Connect(nats.DefaultURL)
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	/*
		A durable consumer in NATS JetStream is a subscription tracker that saves
		its state—such as its exact position (cursor) in a stream—on the NATS server
		so it survives application crashes, disconnects, and restarts
	*/
	// Create durable consumer monitor
	// ManualAck disables auto ack functionality
	js.Subscribe("ORDERS.*", func(msg *nats.Msg) {
		msg.Ack() // Ack acknowledges a message
		var order model.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("monitor service subscribes from subject:%s\n", msg.Subject)
		log.Printf("OrderID:%d, CustomerID: %s, Status:%s\n", order.OrderID, order.CustomerID, order.Status)
	}, nats.Durable("monitor"), nats.ManualAck())
	/*
		Calling Goexit from the main goroutine terminates that goroutine without func main returning.
		Since func main has not returned, the program continues execution of other goroutines.
		If all other goroutines exit, the program crashes.
		It crashes if called from a thread not created by the Go runtime.
	*/
	runtime.Goexit()

}
