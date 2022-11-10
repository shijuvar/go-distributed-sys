package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"os"
	"runtime"
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
	if stream, _ := js.StreamInfo("KV_sdkv"); stream == nil {
		fmt.Println("nill")
		// A key-value (KV) bucket is created by specifying a bucket name.
		kv, _ = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: "sdkv",
		})
	} else {
		kv, _ = js.KeyValue("sdkv")
	}
	w, _ := kv.Watch("services.*")
	defer w.Stop()
	for kve := range w.Updates() {
		if kve != nil {
			fmt.Println("before watcher")
			fmt.Printf("%s @ %d -> %q (op: %s)\n", kve.Key(), kve.Revision(), string(kve.Value()), kve.Operation())
			fmt.Println("watcher")
		}

	}
	runtime.Goexit()
}
