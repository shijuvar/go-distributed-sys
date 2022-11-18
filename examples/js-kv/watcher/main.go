package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/nats-io/nats.go"
)

func main() {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = nats.DefaultURL
	}

	nc, _ := nats.Connect(url)
	defer nc.Drain()

	js, _ := nc.JetStream()
	var kv nats.KeyValue
	if stream, _ := js.StreamInfo("KV_discovery"); stream == nil {
		// A key-value (KV) bucket is created by specifying a bucket name.
		kv, _ = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: "discovery",
		})
	} else {
		kv, _ = js.KeyValue("discovery")
	}
	// KeyWatcher for the wildcard "services.*"
	w, _ := kv.Watch("services.*")
	defer w.Stop()
	for kve := range w.Updates() {
		if kve != nil {
			fmt.Printf("%s @ %d -> %q (op: %s)\n", kve.Key(), kve.Revision(), string(kve.Value()), kve.Operation())
		}

	}
	runtime.Goexit()
}
