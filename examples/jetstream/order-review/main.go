package main

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"

	"github.com/shijuvar/go-distributed-sys/examples/jetstream/model"
)

const (
	subSubjectName = "ORDERS.created"
	pubSubjectName = "ORDERS.approved"
)

func main() {
	// Connect to NATS
	nc, _ := nats.Connect(nats.DefaultURL)
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	// Create a Pull-based consumer with maximum 128 inflight.
	// PullMaxWaiting defines the max inflight pull requests.
	/*
		In-flight messages are messages that have been delivered by the server
		to a consumer/reader but have not yet been acknowledged (acked) as successfully processed
	*/
	sub, _ := js.PullSubscribe(subSubjectName, "order-review", nats.PullMaxWaiting(128))
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//defer cancel()

	for {
		//select {
		//case <-ctx.Done():
		//	return
		//default:
		//}
		//msgs, _ := sub.Fetch(10, nats.Context(ctx))
		msgs, _ := sub.Fetch(10)

		for _, msg := range msgs {
			msg.Ack()
			var order model.Order
			err := json.Unmarshal(msg.Data, &order)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("order-review service")
			log.Printf("OrderID:%d, CustomerID: %s, Status:%s\n", order.OrderID, order.CustomerID, order.Status)
			reviewOrder(js, order)
		}
	}
}

// reviewOrder reviews the order and publishes ORDERS.approved event
func reviewOrder(js nats.JetStreamContext, order model.Order) {
	// Changing the Order status
	order.Status = "approved"
	orderJSON, _ := json.Marshal(order)
	_, err := js.Publish(pubSubjectName, orderJSON)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Order with OrderID:%d has been %s\n", order.OrderID, order.Status)
}
