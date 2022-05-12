package ordersyncrepository

import (
	"context"
	"database/sql"

	"github.com/cockroachdb/cockroach-go/crdb"

	"github.com/shijuvar/go-distributed-sys/eventstream/order"
)

type repository struct {
	db *sql.DB
}

// New returns a concrete repository backed by CockroachDB
func New(db *sql.DB) (order.Repository, error) {
	// return  repository
	return &repository{
		db: db,
	}, nil
}

// CreateOrder inserts a new ordersyncrepository and its ordersyncrepository items into db
func (repo *repository) CreateOrder(ctx context.Context, order order.Order) error {

	// Run a transaction to sync the query model.
	err := crdb.ExecuteTx(ctx, repo.db, nil, func(tx *sql.Tx) error {
		return createOrder(tx, order)
	})
	if err != nil {
		return err
	}
	return nil
}

func createOrder(tx *sql.Tx, order order.Order) error {

	// Insert ordersyncrepository into the "orders" table.
	sql := `
			INSERT INTO orders (id, customerid, status, createdon, restaurantid)
			VALUES ($1,$2,$3,$4,$5)`
	_, err := tx.Exec(sql, order.ID, order.CustomerID, order.Status, order.CreatedOn, order.RestaurantId)
	if err != nil {
		return err
	}
	// Insert ordersyncrepository items into the "orderitems" table.
	// Because it's store for read model, we can insert denormalized data
	for _, v := range order.OrderItems {
		sql = `
			INSERT INTO orderitems (orderid, customerid, code, name, unitprice, quantity)
			VALUES ($1,$2,$3,$4,$5,$6)`

		_, err := tx.Exec(sql, order.ID, order.CustomerID, v.ProductCode, v.Name, v.UnitPrice, v.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}

// ChangeOrderStatus changes the ordersyncrepository status
func (repo *repository) ChangeOrderStatus(ctx context.Context, command order.ChangeOrderStatusCommand) error {
	sql := `
UPDATE orders
SET status=$2
WHERE id=$1`

	_, err := repo.db.ExecContext(ctx, sql, command.OrderID, command.Status)
	if err != nil {
		return err
	}
	return nil
}

// GetOrderByID query the ordersyncrepository by given id
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
