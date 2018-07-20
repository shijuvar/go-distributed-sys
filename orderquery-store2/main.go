package main

import (
	"encoding/json"
	"log"
	"runtime"

	stan "github.com/nats-io/go-nats-streaming"

	"github.com/shijuvar/go-distributed-sys/pb"
	"github.com/shijuvar/go-distributed-sys/store"
	"github.com/shijuvar/go-distributed-sys/natsutil"
)

const (
	clusterID  = "test-cluster"
	clientID   = "order-query-store2"
	channel    = "order-created"
	durableID  = "store-durable"
	queueGroup = "order-query-store-group"
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

	sc.QueueSubscribe(channel, queueGroup, func(msg *stan.Msg) {
		order := pb.OrderCreateCommand{}
		err := json.Unmarshal(msg.Data, &order)
		if err == nil {
			// Handle the message
			log.Printf("Subscribed message from clientID - %s: %+v\n", clientID, order)
			queryStore := store.QueryStore{}
			// Perform data replication for query model into CockroachDB
			err := queryStore.SyncOrderQueryModel(order)
			if err != nil {
				log.Printf("Error while replicating the query model %+v", err)
			}
		}
	}, stan.DurableName(durableID),
	)
	runtime.Goexit()
}
