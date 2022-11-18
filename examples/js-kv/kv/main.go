package main

import (
	"fmt"
	"os"

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
	// Creates a Key/Value store with a bucket named discovery if it doesn't exist.
	// Otherwise, it creates a Key/Value store by providing the existing bucket name.
	if stream, _ := js.StreamInfo("KV_discovery"); stream == nil {
		// A key-value (KV) bucket is created by specifying a bucket name.
		kv, _ = js.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: "discovery",
		})
	} else {
		kv, _ = js.KeyValue("discovery")
	}
	// `KeyValue` interface provides the
	// standard `Put` and `Get` methods, with a revision number of the entry.
	kv.Put("services.orders", []byte("https://localhost:8080/orders"))
	entry, _ := kv.Get("services.orders")
	fmt.Printf("%s @ %d -> %q\n", entry.Key(), entry.Revision(), string(entry.Value()))

	kv.Put("services.orders", []byte("https://localhost:8080/v1/orders"))
	entry, _ = kv.Get("services.orders")
	fmt.Printf("%s @ %d -> %q\n", entry.Key(), entry.Revision(), string(entry.Value()))

	_, err := kv.Update("services.orders", []byte("https://localhost:8080/v1/orders"), 1)
	fmt.Printf("expected error because of wrong revision: %s\n", err)

	kv.Update("services.orders", []byte("https://localhost:8080/v2/orders"), 2)
	entry, _ = kv.Get("services.orders")
	fmt.Printf("%s @ %d -> %q\n", entry.Key(), entry.Revision(), string(entry.Value()))

	name := <-js.StreamNames()
	fmt.Printf("KV stream name: %s\n", name)

	kv.Put("services.customers", []byte("https://localhost:8080/v2/customers"))
	entry, _ = kv.Get("services.customers")
	fmt.Printf("%s @ %d -> %q\n", entry.Key(), entry.Revision(), string(entry.Value()))

	kv.Delete("services.customers")
	entry, err = kv.Get("services.customers")
	if err == nil {
		fmt.Printf("%s @ %d -> %q\n", entry.Key(), entry.Revision(), string(entry.Value()))
	}
	//if err := js.DeleteKeyValue("discovery"); err != nil {
	//	fmt.Println(err)
	//}
}
