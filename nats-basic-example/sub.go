package main


import (
	"log"
	"runtime"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	nc.Subscribe("order.created", func(m *nats.Msg) {
		log.Printf("[Orer] %s", string(m.Data))
	})
	runtime.Goexit()
}
