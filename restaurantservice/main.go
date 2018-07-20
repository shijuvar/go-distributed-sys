package main

import (
	"encoding/json"
	"log"
	"runtime"
	"time"

	stan "github.com/nats-io/go-nats-streaming"

	"github.com/shijuvar/go-distributed-sys/pb"
	"github.com/shijuvar/go-distributed-sys/natsutil"
)

const (
	clusterID = "test-cluster"
	clientID  = "restaurant-service"
	channel   = "order-payment-debited"
	durableID = "restaurant-service-durable"
)

func main() {
	// Register new component within the system.
	comp := natsutil.NewStreamingComponent(clientID)

	// Connect to NATS Streaming server
	err := comp.ConnectToNATSStreaming(
		clusterID,
		stan.NatsURL(stan.DefaultNatsURL),
	)
	if err != nil {
		log.Fatal(err)
	}
	// Get the NATS Streaming Connection
	sc := comp.NATS()
	// Subscribe with manual ack mode, and set AckWait to 60 seconds
	aw, _ := time.ParseDuration("60s")
	// Subscribe the channel
	sc.Subscribe(channel, func(msg *stan.Msg) {
		msg.Ack() // Manual ACK
		paymentDebited := pb.OrderPaymentDebitedCommand{}
		// Unmarshal JSON that represents the Order data
		err := json.Unmarshal(msg.Data, &paymentDebited)
		if err != nil {
			log.Print(err)
			return
		}
		// Handle the message
		log.Printf("Order approved for Order ID: %s for Customer: %s\n", paymentDebited.OrderId, paymentDebited.CustomerId)
		// ToDo: Publish event to Event Store

	}, stan.DurableName(durableID),
		stan.MaxInflight(25),
		stan.SetManualAckMode(),
		stan.AckWait(aw),
	)
	runtime.Goexit()
}
