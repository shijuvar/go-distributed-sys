package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/shijuvar/go-distributed-sys/eventstream/cockroachdb/ordersyncrepository"
	"github.com/shijuvar/go-distributed-sys/eventstream/eventstore"
	"github.com/shijuvar/go-distributed-sys/eventstream/sqldb"
	"google.golang.org/grpc"
	"log"
	"runtime"

	"github.com/nats-io/nats.go"
	ordermodel "github.com/shijuvar/go-distributed-sys/eventstream/order"
	"github.com/shijuvar/go-distributed-sys/pkg/natsutil"
)

const (
	clientID         = "order-review-worker"
	subscribeSubject = "ORDERS.created"
	queueGroup       = "order-review-worker"
	event            = "ORDERS.approved"
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
	// Create durable consumer monitor
	jetStreamContext.QueueSubscribe(subscribeSubject, queueGroup, func(msg *nats.Msg) {
		msg.Ack()
		var paymentDebited ordermodel.PaymentDebitedCommand
		// Unmarshal JSON that represents the PaymentDebited data
		err := json.Unmarshal(msg.Data, &paymentDebited)
		if err != nil {
			log.Print(err)
			return
		}
		command := ordermodel.ChangeOrderStatusCommand{
			OrderID: paymentDebited.OrderID,
			Status:  "Approved",
		}
		orderDB, _ := sqldb.NewOrdersDB()
		repository, _ := ordersyncrepository.New(orderDB.DB)
		if err := repository.ChangeOrderStatus(context.Background(), command); err != nil {
			log.Println(err)
			return
		}
		// Publish event to Event Store
		if err := createOrderApprovedCommand(command); err != nil {
			log.Println("error occured while executing the OrderApproved command")
			return
		}

	}, nats.Durable(clientID), nats.ManualAck())
	runtime.Goexit()
}

// createOrderApprovedCommand calls the event store RPC to create an event
// OrderApproved command is created on Event Store
func createOrderApprovedCommand(command ordermodel.ChangeOrderStatusCommand) error {

	conn, err := grpc.Dial(grpcUri, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()
	client := eventstore.NewEventStoreClient(conn)
	eventid, _ := uuid.NewUUID()
	event := &eventstore.Event{
		EventId:       eventid.String(),
		EventType:     event,
		AggregateId:   command.OrderID,
		AggregateType: aggregate,
		EventData:     "",
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
