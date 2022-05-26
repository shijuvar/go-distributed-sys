package middleware

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	ordersvc "github.com/shijuvar/go-distributed-sys/eventstream/gokit-ordersvc"
	"github.com/shijuvar/go-distributed-sys/eventstream/order"
)

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next ordersvc.Service) ordersvc.Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   ordersvc.Service
	logger log.Logger
}

func (mw loggingMiddleware) CreateOrder(ctx context.Context, order order.Order) (id string, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "CreateOrder", "CustomerID", order.CustomerID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.CreateOrder(ctx, order)
}

func (mw loggingMiddleware) GetOrderByID(ctx context.Context, id string) (order order.Order, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetOrderByID", "OrderID", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetOrderByID(ctx, id)
}
