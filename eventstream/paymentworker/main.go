package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"runtime"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"

	"github.com/shijuvar/go-distributed-sys/eventstream/eventstore"
	ordermodel "github.com/shijuvar/go-distributed-sys/eventstream/order"
	"github.com/shijuvar/go-distributed-sys/pkg/natsutil"
)

const (
	clientID         = "payment-worker"
	subscribeSubject = "ORDERS.created"
	queueGroup       = "payment-worker"
	event            = "ORDERS.paymentdebited"
	aggregate        = "order"
	stream           = "ORDERS"
	grpcUri          = "localhost:50051"
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
	// Create durable consumer
	jetStreamContext.QueueSubscribe(subscribeSubject, queueGroup, func(msg *nats.Msg) {
		msg.Ack()
		var order ordermodel.Order
		// Unmarshal JSON that represents the Order data
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Print(err)
			return
		}
		log.Printf("Message subscribed on subject:%s, from:%s, data:%v", subscribeSubject, clientID, order)

		// Create OrderPaymentDebitedCommand from Order
		command := ordermodel.PaymentDebitedCommand{
			OrderID:    order.ID,
			CustomerID: order.CustomerID,
			Amount:     order.Amount,
		}
		if err := createPaymentDebitedCommand(command); err != nil {
			log.Println("error occured while executing the PaymentDebited command")
		}

	}, nats.Durable(clientID), nats.ManualAck())
	runtime.Goexit()
}

func createPaymentDebitedCommand(command ordermodel.PaymentDebitedCommand) error {

	conn, err := grpc.Dial(grpcUri, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()
	client := eventstore.NewEventStoreClient(conn)
	paymentJSON, _ := json.Marshal(command)
	eventid, _ := uuid.NewUUID()
	event := &eventstore.Event{
		EventId:       eventid.String(),
		EventType:     event,
		AggregateId:   command.OrderID,
		AggregateType: aggregate,
		EventData:     string(paymentJSON),
		Stream:        stream,
	}
	createEventRequest := &eventstore.CreateEventRequest{Event: event}

	resp, err := client.CreateEvent(context.Background(), createEventRequest)
	if err != nil {
		return fmt.Errorf("error from RPC server: %w", err)
	}
	if resp.IsSuccess {
		return nil
	}
	return errors.New("error from RPC server")
}
