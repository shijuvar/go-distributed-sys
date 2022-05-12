package order

import (
	"time"
)

// Order aggregate
type Order struct {
	ID           string      `json:"order_id,omitempty"`
	CustomerID   string      `json:"customer_id,omitempty"`
	Status       string      `json:"status,omitempty"`
	CreatedOn    time.Time   `json:"created_on,omitempty"`
	RestaurantId string      `json:"restaurant_id,omitempty"`
	Amount       float64     `json:"amount,omitempty"`
	OrderItems   []OrderItem `json:"order_items,omitempty"`
}

// OrderItem value type
type OrderItem struct {
	ProductCode string  `json:"code,omitempty"`
	Name        string  `json:"name,omitempty"`
	UnitPrice   float64 `json:"unit_price,omitempty"`
	Quantity    int32   `json:"quantity,omitempty"`
}
