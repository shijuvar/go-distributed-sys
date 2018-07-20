package store

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	// Connect to the "ordersdb" database
	db, err = sql.Open("postgres", "postgresql://shijuvar@localhost:26257/ordersdb?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	createTables()
}
func createTables() {

	// Create the "events" table.
	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS events (id string PRIMARY KEY, eventtype string, aggregateid string, aggregatetype string, eventdata string, channel string)"); err != nil {
		log.Fatal(err)
	}

	// Create the "orders" table.
	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS orders (id string PRIMARY KEY, customerid string, status string, createdon int, restaurantid string, amount float)"); err != nil {
		log.Fatal(err)
	}

	// Create the "orderitems" table.
	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS orderitems (id serial PRIMARY KEY, orderid string, customerid string, code string, name string, unitprice float, quantity int)"); err != nil {
		log.Fatal(err)
	}

}
