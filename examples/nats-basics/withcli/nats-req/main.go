package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

// NOTE: Can test with demo servers.
// nats-req -s demo.nats.io <subject> <msg>
// nats-req -s demo.nats.io:4443 <subject> <msg> (TLS version)

func usage() {
	log.Printf("Usage: nats-req [-s server] [-creds file] <subject> <msg>\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var userCreds = flag.String("creds", "", "User Credentials File")
	var showHelp = flag.Bool("h", false, "Show help message")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	args := flag.Args()
	if len(args) < 2 {
		showUsageAndExit(1)
	}

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Requestor")}

	// Use UserCredentials
	if *userCreds != "" {
		opts = append(opts, nats.UserCredentials(*userCreds))
	}

	// Connect to NATS
	nc, err := nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	subj, payload := args[0], []byte(args[1])

	msg, err := nc.Request(subj, payload, 2*time.Second)
	if err != nil {
		if nc.LastError() != nil {
			log.Fatalf("%v for request", nc.LastError())
		}
		log.Fatalf("%v for request", err)
	}

	log.Printf("Published [%s] : '%s'", subj, payload)
	log.Printf("Received  [%v] : '%s'", msg.Subject, string(msg.Data))
}
