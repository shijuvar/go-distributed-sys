package sqldb

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type EventStoreDB struct {
	DB *sql.DB
}

func NewEventStoreDB() (*EventStoreDB, error) {
	// Connect to the "eventstoredb" database
	dbEventStore, err := sql.Open("postgres", "postgresql://shijuvar@localhost:26257/eventstoredb?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &EventStoreDB{
		DB: dbEventStore,
	}, nil
}

func (eventstore *EventStoreDB) CreateEventStoreDBSchema() error {
	// Create the "events" table.
	if _, err := eventstore.DB.Exec(
		"CREATE TABLE IF NOT EXISTS events (id string PRIMARY KEY, eventtype STRING, aggregateid STRING, aggregatetype STRING, eventdata JSONB, stream STRING)"); err != nil {
		return err
	}
	return nil
}

func (eventstore *EventStoreDB) Close() {
	eventstore.DB.Close()
}

type OrdersDB struct {
	DB *sql.DB
}

func NewOrdersDB() (*OrdersDB, error) {
	// Connect to the "ordersdb" database
	dbOrders, err := sql.Open("postgres", "postgresql://shijuvar@localhost:26257/ordersdb?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &OrdersDB{
		DB: dbOrders,
	}, nil
}
func (orders *OrdersDB) CreateOrdersDBSchema() error {
	// Create the "orders" table.
	if _, err := orders.DB.Exec(
		"CREATE TABLE IF NOT EXISTS orders (id STRING PRIMARY KEY, customerid STRING, status STRING, createdon DATE, restaurantid STRING, amount FLOAT)"); err != nil {
		return err
	}

	// Create the "orderitems" table.
	if _, err := orders.DB.Exec(
		"CREATE TABLE IF NOT EXISTS orderitems (id serial PRIMARY KEY, orderid STRING, customerid STRING, code STRING, name STRING, unitprice FLOAT, quantity INT )"); err != nil {
		return err
	}
	return nil
}
func (orders *OrdersDB) Close() {
	orders.DB.Close()
}
