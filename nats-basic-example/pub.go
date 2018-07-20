package main


import (
	"encoding/json"
	"log"
	"runtime"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	payload := struct {
		OrderID string
		Status      string
	}{
		OrderID: "1234-5678-90",
		Status:      "Placed",
	}
	msg, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	nc.Publish("order.created", msg)
	log.Println("Message published")
	runtime.Goexit()
}
