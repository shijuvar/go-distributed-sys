package main

import (
	"log"

	"github.com/nats-io/nats.go"

	"github.com/shijuvar/go-distributed-sys/eventstream/sqldb"
	"github.com/shijuvar/go-distributed-sys/pkg/natsutil"
)

const (
	orderStream         = "ORDERS"
	orderStreamSubjects = "ORDERS.*"
)

func main() {
	createJStream(orderStream, orderStreamSubjects)
	//createEventStoreDB()
	//createOrdersDB()
}

// createStream creates a stream by using JetStreamContext
func createJStream(streamName string, streamSubjects string) error {
	natsComponent := natsutil.NewNATSComponent("bootstrapper")
	err := natsComponent.ConnectToServer(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	var js nats.JetStreamContext
	js, err = natsComponent.JetStreamContext()
	if err != nil {
		log.Fatalln(err)
	}
	// Check if the ORDERS stream already exists; if not, create it.
	stream, err := js.StreamInfo(streamName)
	if err != nil {
		log.Println(err)
	}
	if stream == nil {
		log.Printf("creating stream %q and subjects %q", streamName, streamSubjects)
		_, err = js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamSubjects},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func createEventStoreDB() {
	eventStoreDB, err := sqldb.NewEventStoreDB()
	if err != nil {
		log.Fatalln("Error on eventstore db:", err)
	}
	defer eventStoreDB.Close()

	err = eventStoreDB.CreateEventStoreDBSchema()
	if err != nil {
		log.Fatalln("Error:", err)
	}
	log.Println("eventstore DB has ben created")
}

func createOrdersDB() {
	ordersDB, err := sqldb.NewOrdersDB()
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer ordersDB.Close()
	err = ordersDB.CreateOrdersDBSchema()
	if err != nil {
		log.Fatalln("Error:", err)
	}
	log.Println("ordersdb DB has ben created")
}
