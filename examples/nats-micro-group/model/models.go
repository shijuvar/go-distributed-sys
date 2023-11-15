package model

type Order struct {
	OrderID    string `json:"orderid"`
	CustomerID string `json:"customerid""`
	Status     string `josn:"status,omitempty"`
}
type CreateOrderRequest struct {
	Order
}
type GetOrderRequest struct {
	OrderID string `json:"orderid"`
}
