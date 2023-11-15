package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	msg, err := nc.Request("svc.echo", []byte("Hello"), time.Second*30)
	if err != nil {
		log.Fatalf("error:", err)
	}
	fmt.Println(string(msg.Data))
	runtime.Goexit()
}
