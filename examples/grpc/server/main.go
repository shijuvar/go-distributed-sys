package main

import (
	"context"
	"crypto/tls"
	"net"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/shijuvar/go-distributed-sys/examples/grpc/pb"
)

const (
	port = ":3000"
)

// server is used to implement customer.CustomerServer.
type server struct {
	pb.UnimplementedCustomerServer
	savedCustomers []*pb.CustomerRequest
}

// CreateCustomer creates a new Customer
func (s *server) CreateCustomer(ctx context.Context, in *pb.CustomerRequest) (*pb.CustomerResponse, error) {
	s.savedCustomers = append(s.savedCustomers, in)
	return &pb.CustomerResponse{Id: in.Id, Success: true}, nil
}

// GetCustomers returns all customers by given filter
func (s *server) GetCustomers(filter *pb.CustomerFilter, stream pb.Customer_GetCustomersServer) error {
	for _, customer := range s.savedCustomers {
		if filter.Keyword != "" {
			if !strings.Contains(customer.Name, filter.Keyword) {
				continue
			}
		}
		if err := stream.Send(customer); err != nil {
			return err
		}
	}
	return nil
}

var log grpclog.LoggerV2

func init() {
	log = grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)
	grpclog.SetLoggerV2(log)
}

func withServerInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(serverInterceptor)
}

// general unary interceptor function to handle auth per RPC call as well as logging
func serverInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	if info.FullMethod == "/customer.Customer/CreateCustomer" {
		log.Info("---------CreateCustomer---------\n")

		if err := authorize(ctx); err != nil {
			return nil, err
		}
	}

	h, err := handler(ctx, req)

	//logging
	log.Infof("request - Method:%s\tDuration:%s\tError:%v\n",
		info.FullMethod,
		time.Since(start),
		err)

	return h, err
}
func authorize(ctx context.Context) error {
	// code from the authorize() function:

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "retrieving metadata failed")
	}

	token, ok := md["authorization"]
	if !ok {
		return status.Errorf(codes.Unauthenticated, "no auth details supplied")
	}

	if token[0] != "my-valid-token" {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}

	log.Info("Authorized to the RPC server")
	return nil
}
func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Creates a new gRPC server
	//gRPC server with TLS
	//tlsCredentials, err := loadTLSCredentials()
	//s := grpc.NewServer(
	//	grpc.Creds(tlsCredentials),
	//	withServerInterceptor(),
	//)

	s := grpc.NewServer(withServerInterceptor())
	// Register v1 server
	pb.RegisterCustomerServer(s, &server{})
	// Serve gRPC server
	s.Serve(lis)
}
