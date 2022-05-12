package order

type PaymentDebitedCommand struct {
	OrderID    string
	CustomerID string
	Amount     float64
}

type ChangeOrderStatusCommand struct {
	OrderID string
	Status  string
}
