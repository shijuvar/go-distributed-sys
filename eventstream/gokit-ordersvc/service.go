package ordersvc

import (
	"context"
	"errors"

	"github.com/shijuvar/go-distributed-sys/eventstream/order"
)

var (
	ErrNotFound = errors.New("order not found")
)

type Service interface {
	CreateOrder(ctx context.Context, o order.Order) (string, error)
	GetOrderByID(ctx context.Context, id string) (order.Order, error)
}
