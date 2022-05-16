package orderqueryrepository

import (
	"context"
	"database/sql"

	"github.com/shijuvar/go-distributed-sys/eventstream/order"
)

type repository struct {
	db *sql.DB
}

// New returns a concrete repository backed by CockroachDB
func New(db *sql.DB) (order.QueryRepository, error) {
	// return  repository
	return &repository{
		db: db,
	}, nil
}

// GetOrderByID query the Orders by given id
func (repo *repository) GetOrderByID(ctx context.Context, id string) (order.Order, error) {
	var orderRow = order.Order{}
	if err := repo.db.QueryRowContext(ctx,
		"SELECT id, customerid, status, createdon, restaurantid FROM orders WHERE id = $1",
		id).
		Scan(
			&orderRow.ID, &orderRow.CustomerID, &orderRow.Status, &orderRow.CreatedOn, &orderRow.RestaurantId,
		); err != nil {
		return orderRow, err
	}
	// ToDo: Query ordersyncrepository items from orderitems table
	return orderRow, nil
}

// Close implements DB.Close
func (repo *repository) Close() error {
	return repo.db.Close()
}
