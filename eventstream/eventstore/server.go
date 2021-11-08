package eventstore

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/shijuvar/go-distributed-sys/eventstream/pb"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedEventStoreServer
}

// Create a new event to the event repository
func (s *server) CreateEvent(context.Context, *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	return nil, nil
}

// Get all events for the given aggregate and event
func (s *server) GetEvents(context.Context, *pb.GetEventsRequest) (*pb.GetEventsResponse, error) {
	return nil, nil
}

//    Get stream of events for the given event
func (s *server) GetEventsStream(*pb.GetEventsRequest, pb.EventStore_GetEventsStreamServer) error {
	return nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterEventStoreServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
