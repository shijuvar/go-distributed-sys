package implementation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/shijuvar/go-distributed-sys/eventstream/eventstore"
	ordersvc "github.com/shijuvar/go-distributed-sys/eventstream/gokit-ordersvc"
	"github.com/shijuvar/go-distributed-sys/eventstream/order"
)

const (
	event     = "ORDERS.created"
	aggregate = "order"
	grpcUri   = "localhost:50051"
)

type rpcClient interface {
	createOrder(order order.Order) error
}
type grpcClient struct {
}

// createOrder calls the CreateEvent RPC
func (gc grpcClient) createOrder(order order.Order) error {
	conn, err := grpc.Dial(grpcUri, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()
	client := eventstore.NewEventStoreClient(conn)
	orderJSON, _ := json.Marshal(order)

	eventid, _ := uuid.NewUUID()
	event := &eventstore.Event{
		EventId:       eventid.String(),
		EventType:     event,
		AggregateId:   order.ID,
		AggregateType: aggregate,
		EventData:     string(orderJSON),
		Stream:        "ORDERS",
	}

	createEventRequest := &eventstore.CreateEventRequest{Event: event}
	resp, err := client.CreateEvent(context.Background(), createEventRequest)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return fmt.Errorf("error from RPC server with: status code:%s message:%s", st.Code().String(), st.Message())
		}
		return fmt.Errorf("error from RPC server: %w", err)
	}
	if resp.IsSuccess {
		return nil
	}
	return errors.New("error from RPC server")
}

type service struct {
	rpc rpcClient
}

func NewServiceWithRPC(rpc rpcClient) ordersvc.Service {
	return &service{
		rpc: rpc,
	}
}
func NewService() ordersvc.Service {
	return &service{
		rpc: grpcClient{},
	}
}

func (s *service) CreateOrder(ctx context.Context, order order.Order) (string, error) {
	id, _ := uuid.NewUUID()
	aggregateID := id.String()
	order.ID = aggregateID
	order.Status = "Pending"
	order.CreatedOn = time.Now()
	order.Amount = order.GetAmount()
	//err := s.rpc.createOrder(order)
	//if err != nil {
	//	return "", err
	//}
	return aggregateID, nil
}

func (s *service) GetOrderByID(ctx context.Context, id string) (order.Order, error) {
	return order.Order{}, nil
}
