package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"

	"github.com/shijuvar/go-distributed-sys/pb"
)

const (
	event     = "order-created"
	aggregate = "order"
	grpcUri   = "localhost:50051"
)

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

func initRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/orders", createOrder).Methods("POST")
	return router
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order pb.OrderCreateCommand
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, "Invalid Order Data", 500)
		return
	}
	aggregateID := uuid.NewV4().String()
	order.OrderId = aggregateID
	order.Status = "Pending"
	order.CreatedOn = time.Now().Unix()
	err = createOrderRPC(order)
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

// createOrderRPC calls the CreateEvent RPC
func createOrderRPC(order pb.OrderCreateCommand) error {

	conn, err := grpc.Dial(grpcUri, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Unable to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewEventStoreClient(conn)
	orderJSON, _ := json.Marshal(order)

	event := &pb.Event{
		EventId:       uuid.NewV4().String(),
		EventType:     event,
		AggregateId:   order.OrderId,
		AggregateType: aggregate,
		EventData:     string(orderJSON),
		Channel:       event,
	}

	resp, err := client.CreateEvent(context.Background(), event)
	if err != nil {
		return errors.Wrap(err, "Error from RPC server")
	}
	if resp.IsSuccess {
		return nil
	} else {
		return errors.Wrap(err, "Error from RPC server")
	}
}
