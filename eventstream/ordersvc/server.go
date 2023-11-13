package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"github.com/shijuvar/go-distributed-sys/eventstream/eventstore"
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
	conn, err := grpc.Dial(grpcUri, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

type orderHandler struct {
	rpc rpcClient
}

func (h orderHandler) createOrder(w http.ResponseWriter, r *http.Request) {
	var order order.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid Order Data", 500)
		return
	}
	id, _ := uuid.NewUUID()
	aggregateID := id.String()
	order.ID = aggregateID
	order.Status = "Pending"
	order.CreatedOn = time.Now()
	order.Amount = order.GetAmount()
	err = h.rpc.createOrder(order)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to create Order", 500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(order)
	w.Write(j)
}

func initRoutes() *mux.Router {
	router := mux.NewRouter()
	h := orderHandler{
		rpc: grpcClient{},
	}
	router.HandleFunc("/api/orders", h.createOrder).Methods("POST")
	return router
}
func main() {
	// Create the Server
	server := &http.Server{
		Addr:    ":3000",
		Handler: initRoutes(),
	}
	log.Println("HTTP Sever listening...")
	// Running the HTTP Server
	server.ListenAndServe()
}
