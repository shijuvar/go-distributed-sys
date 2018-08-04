package main

import (
	"log"
	"runtime"

	"github.com/nats-io/go-nats"
)

func main() {
	nc, err := nats.Connect("nats://127.0.0.1:4222",
		nats.UserInfo("shijuvar", "gopher"),
	)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	nc.QueueSubscribe("order.created", "worker-group", func(m *nats.Msg) {
		log.Printf("[Orer] %s", string(m.Data))
	})
	runtime.Goexit()
}
