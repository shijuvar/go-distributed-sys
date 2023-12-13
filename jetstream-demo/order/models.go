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
	Quantity    int     `json:"quantity,omitempty"`
}

/* Domain Services */

// GetAmount returns total amount of the order
func (order Order) GetAmount() float64 {
	var amount float64
	for _, v := range order.OrderItems {
		amount += v.UnitPrice * float64(v.Quantity)
	}
	return amount
}
