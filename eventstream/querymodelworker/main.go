package main

import (
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/shijuvar/go-distributed-sys/eventstream/cockroachdb/ordersyncrepository"
	ordermodel "github.com/shijuvar/go-distributed-sys/eventstream/order"
	"github.com/shijuvar/go-distributed-sys/eventstream/sqldb"
	"github.com/shijuvar/go-distributed-sys/pkg/natsutil"
	"log"
	"runtime"
)

const (
	clientID         = "query-model-worker"
	subscribeSubject = "ORDERS.created"
	queueGroup       = "query-model-worker"
)

func main() {
	natsComponent := natsutil.NewNATSComponent(clientID)
	err := natsComponent.ConnectToServer(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	jetStreamContext, err := natsComponent.JetStreamContext()
	if err != nil {
		log.Fatal(err)
	}
	// Create durable consumer monitor
	jetStreamContext.QueueSubscribe(subscribeSubject, queueGroup, func(msg *nats.Msg) {
		msg.Ack()
		var order ordermodel.Order
		// Unmarshal JSON that represents the Order data
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Print(err)
			return
		}
		orderDB, _ := sqldb.NewOrdersDB()
		repository, _ := ordersyncrepository.New(orderDB.DB)
		// Sync query model with event data
		if err := repository.CreateOrder(context.Background(), order); err != nil {
			log.Printf("Error while replicating the query model %+v", err)
		}

	}, nats.Durable(clientID), nats.ManualAck())
	runtime.Goexit()
}
