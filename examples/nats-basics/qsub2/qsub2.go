package main

import (
	"log"
	"runtime"

	nats "github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	nc.QueueSubscribe("order.created", "worker-group", func(m *nats.Msg) {
		log.Printf("[Orer] %s", string(m.Data))
	})
	runtime.Goexit()
}
