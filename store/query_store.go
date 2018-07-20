package store

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/cockroach-go/crdb"
	"github.com/pkg/errors"

	"github.com/shijuvar/go-distributed-sys/pb"
)

// QueryStore syncs data model to be used for query operations
// Because it's store for read model, denormalized data can be inserted

type QueryStore struct{}

func (store QueryStore) SyncOrderQueryModel(order pb.OrderCreateCommand) error {

	// Run a transaction to sync the query model.
	// Node: There is an issue with this at this moment
	err := crdb.ExecuteTx(context.Background(), db, nil, func(tx *sql.Tx) error {
		return createOrderQueryModel(tx, order)
	})
	if err != nil {
		return errors.Wrap(err, "Error on syncing query store")
	}
	return nil
}

func createOrderQueryModel(tx *sql.Tx, order pb.OrderCreateCommand) error {

	// Insert order into the "orders" table.
	sql := `
INSERT INTO orders (id, customerid, status, createdon, restaurantid, amount) 
VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := tx.Exec(sql, order.OrderId, order.CustomerId, order.Status, order.CreatedOn, order.RestaurantId, order.Amount)
	if err != nil {
		return errors.Wrap(err, "error on insert into orders")
	}
	// Insert order items into the "orderitems" table.
	// Because it's store for read model, we can insert denormalized data
	for _, v := range order.OrderItems {
		sql = `
INSERT INTO orderitems (orderid, customerid, code, name, unitprice, quantity) 
VALUES ($1,$2,$3,$4,$5,$6)`

		_, err := tx.Exec(sql, order.OrderId, order.CustomerId, v.Code, v.Name, v.UnitPrice, v.Quantity)
		if err != nil {
			return errors.Wrap(err, "error on insert into order items")
		}
	}
	return nil
}

// Approve order
func (store QueryStore) ApproveOrder(orderId string) error {
	sql := `
UPDATE orders 
SET status=$2
WHERE orderid=$1`

	_, err := db.Exec(sql, orderId, "Approve")
	if err != nil {
		return errors.Wrap(err, "error on approving order")
	}
	return nil
}
