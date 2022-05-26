package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	ordersvc "github.com/shijuvar/go-distributed-sys/eventstream/gokit-ordersvc"
	"github.com/shijuvar/go-distributed-sys/eventstream/order"
)

// Endpoints holds all Go kit endpoints for the Order service.
type Endpoints struct {
	CreateOrder  endpoint.Endpoint
	GetOrderByID endpoint.Endpoint
}

// MakeServerEndpoints initializes all Go kit endpoints for the Order service.
func MakeServerEndpoints(s ordersvc.Service) Endpoints {
	return Endpoints{
		CreateOrder:  makeCreateOrderEndpoint(s),
		GetOrderByID: makeGetOrderByIDEndpoint(s),
	}
}

func makeCreateOrderEndpoint(s ordersvc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateOrderRequest) // type assertion
		id, err := s.CreateOrder(ctx, req.Order)
		return CreateOrderResponse{ID: id, Err: err}, nil
	}
}

func makeGetOrderByIDEndpoint(s ordersvc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetOrderByIDRequest)
		order, err := s.GetOrderByID(ctx, req.ID)
		return GetOrderByIDResponse{Order: order, Err: err}, nil
	}
}

// CreateOrderRequest holds the request parameters for the CreateOrder method.
type CreateOrderRequest struct {
	Order order.Order
}

// CreateOrderResponse holds the response values for the CreateOrder method.
type CreateOrderResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

func (r CreateOrderResponse) Error() error { return r.Err }

// GetOrderByIDRequest holds the request parameters for the GetOrderByID method.
type GetOrderByIDRequest struct {
	ID string
}

// GetOrderByIDResponse holds the response values for the GetOrderByID method.
type GetOrderByIDResponse struct {
	Order order.Order `json:"order"`
	Err   error       `json:"error,omitempty"`
}

func (r GetOrderByIDResponse) Error() error { return r.Err }
