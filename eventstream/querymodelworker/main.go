package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"

	"github.com/shijuvar/go-distributed-sys/eventstream/cockroachdb/ordersyncrepository"
	ordermodel "github.com/shijuvar/go-distributed-sys/eventstream/order"
	"github.com/shijuvar/go-distributed-sys/eventstream/sqldb"
	"github.com/shijuvar/go-distributed-sys/pkg/natsutil"
)

const (
	clientID         = "query-model-worker"
	subscribeSubject = "ORDERS.created"
	queueGroup       = "query-model-worker"
	batch            = 1 // just for the same of demo. Use bigger numbers
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
	pullSubscribeOnOrder(jetStreamContext)
}

func pushSubscribeOnOrder(js nats.JetStreamContext) {
	// Create durable push consumer
	js.QueueSubscribe(subscribeSubject, queueGroup, func(msg *nats.Msg) {
		msg.Ack()
		var order ordermodel.Order
		// Unmarshal JSON that represents the Order data
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Print(err)
			return
		}
		log.Printf("Message subscribed on subject:%s, from:%s,  data:%v", subscribeSubject, clientID, order)
		orderDB, _ := sqldb.NewOrdersDB()
		repository, _ := ordersyncrepository.New(orderDB.DB)
		// Sync query model with event data
		if err := repository.CreateOrder(context.Background(), order); err != nil {
			log.Printf("Error while replicating the query model %+v", err)
		}

	}, nats.Durable(clientID), nats.ManualAck())
}

func pullSubscribeOnOrder(js nats.JetStreamContext) {
	// Create Pull based consumer with maximum 128 inflight.
	// PullMaxWaiting defines the max inflight pull requests.
	sub, _ := js.PullSubscribe(subscribeSubject, clientID, nats.PullMaxWaiting(128))
	for {

		msgs, _ := sub.Fetch(batch)
		for _, msg := range msgs {
			msg.Ack()
			var order ordermodel.Order
			// Unmarshal JSON that represents the Order data
			err := json.Unmarshal(msg.Data, &order)
			if err != nil {
				log.Print(err)
				return
			}
			log.Printf("Message subscribed on subject:%s, from:%s,  data:%v", subscribeSubject, clientID, order)
			orderDB, _ := sqldb.NewOrdersDB()
			repository, _ := ordersyncrepository.New(orderDB.DB)
			// Sync query model with event data
			if err := repository.CreateOrder(context.Background(), order); err != nil {
				log.Printf("Error while replicating the query model %+v", err)
			}
		}
	}
}
