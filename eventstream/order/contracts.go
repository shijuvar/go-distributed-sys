package order

import "context"

type Repository interface {
	CreateOrder(context.Context, Order) error
	ChangeOrderStatus(context.Context, ChangeOrderStatusCommand) error
}

type QueryRepository interface {
	GetOrderByID(context.Context, string) (Order, error)
	GetOrdersByCustomer(context.Context, string) ([]Order, error)
	GetOrderItems(context.Context, string) ([]OrderItem, error)
}
