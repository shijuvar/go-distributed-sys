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
	return orderRow, nil
}

// GetOrderItems query the order items by given order id
func (repo *repository) GetOrderItems(ctx context.Context, id string) ([]order.OrderItem, error) {
	rows, err := repo.db.QueryContext(ctx,
		"SELECT code, name, unitprice, quantity FROM orderitems WHERE orderid = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// An OrderItem slice to hold data from returned rows.
	var oitems []order.OrderItem

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var item order.OrderItem
		if err := rows.Scan(&item.ProductCode, &item.Name, &item.UnitPrice,
			&item.Quantity); err != nil {
			return oitems, err
		}
		oitems = append(oitems, item)
	}
	if err = rows.Err(); err != nil {
		return oitems, err
	}
	return oitems, nil
}

// GetOrdersByCustomer query the orders by given customer id
func (repo *repository) GetOrdersByCustomer(ctx context.Context, cid string) ([]order.Order, error) {
	rows, err := repo.db.QueryContext(ctx,
		"SELECT id, status, createdon, restaurantid, amount FROM orderitems WHERE orderid = $1", cid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// An Order slice to hold data from returned rows.
	var orders []order.Order

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var o order.Order
		if err := rows.Scan(&o.ID, &o.Status, &o.CreatedOn,
			&o.RestaurantId, o.Amount); err != nil {
			return orders, err
		}
		orders = append(orders, o)
	}
	if err = rows.Err(); err != nil {
		return orders, err
	}
	return orders, nil
}

// Close implements DB.Close
func (repo *repository) Close() error {
	return repo.db.Close()
}
