package main

import (
	"encoding/json"
	"log"
	"runtime"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"

	"github.com/shijuvar/go-distributed-sys/pb"
	"github.com/shijuvar/go-distributed-sys/natsutil"
	"github.com/shijuvar/go-distributed-sys/store"
	"context"
	"github.com/pkg/errors"
)

const (
	clusterID = "test-cluster"
	clientID  = "restaurant-service"
	channel   = "order-payment-debited"
	durableID = "restaurant-service-durable"

	event     = "order-approved"
	aggregate = "order"

	grpcUri   = "localhost:50051"
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
		store := store.QueryStore{}
		if err :=store.ChangeOrderStatus(paymentDebited.OrderId, "Approved"); err!=nil {
			log.Println(err)
			return
		}
		log.Printf("Order approved for Order ID: %s for Customer: %s\n", paymentDebited.OrderId, paymentDebited.CustomerId)
		// Publish event to Event Store
		if err:= createOrderApprovedCommand(paymentDebited.OrderId); err!=nil {
			log.Println("error occured while executing the OrderApproved command")
		}

	}, stan.DurableName(durableID),
		stan.MaxInflight(25),
		stan.SetManualAckMode(),
		stan.AckWait(aw),
	)
	runtime.Goexit()
}

// createOrderApprovedCommand calls the event store RPC to create an event
// OrderApproved command is created on Event Store
func createOrderApprovedCommand(orderId string) error {

	conn, err := grpc.Dial(grpcUri, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewEventStoreClient(conn)

	event := &pb.Event{
		EventId:       uuid.NewV4().String(),
		EventType:     event,
		AggregateId:   orderId,
		AggregateType: aggregate,
		EventData:     "",
		Channel:       event,
	}

	resp, err := client.CreateEvent(context.Background(), event)
	if err != nil {
		return errors.Wrap(err, "error from RPC server")
	}
	if resp.IsSuccess {
		return nil
	} else {
		return errors.Wrap(err, "error from RPC server")
	}
}

