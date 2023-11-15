package main

import (
	"fmt"
	"log"
	"runtime"

	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/micro"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)

	// request handler
	echoHandler := func(req micro.Request) {
		req.Respond(req.Data())
	}

	srv, err := micro.AddService(nc, micro.Config{
		Name:    "EchoService",
		Version: "1.0.0",
		// DoneHandler can be set to customize behavior on stopping a service.
		DoneHandler: func(srv micro.Service) {
			info := srv.Info()
			fmt.Printf("stopped service %q with ID %q\n", info.Name, info.ID)
		},
		// ErrorHandler can be used to customize behavior on service execution error.
		ErrorHandler: func(srv micro.Service, err *micro.NATSError) {
			info := srv.Info()
			fmt.Printf("Service %q returned an error on subject %q: %s", info.Name, err.Subject, err.Description)
		},
		// base handler
		Endpoint: &micro.EndpointConfig{
			Subject: "svc.echo",
			Handler: micro.HandlerFunc(echoHandler),
		},
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Created service: %s (%s)\n", srv.Info().Name, srv.Info().ID)
	//defer srv.Stop()

	runtime.Goexit()
}
