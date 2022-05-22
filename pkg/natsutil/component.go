package natsutil

import (
	"sync"

	nats "github.com/nats-io/nats.go"
)

// NATSComponent is contains reusable logic related to handling
// of the connection to NATS  in the system.
type NATSComponent struct {
	// cmu is the lock from the component.
	cmu sync.Mutex

	// nc is the connection to NATS Streaming.
	nc *nats.Conn

	// name is the name of component.
	name string
}

// NewNATSComponent creates a StreamingComponent
func NewNATSComponent(name string) *NATSComponent {
	return &NATSComponent{
		name: name,
	}
}

// ConnectToServer connects to NATS Server
func (c *NATSComponent) ConnectToServer(url string, options ...nats.Option) error {
	c.cmu.Lock()
	defer c.cmu.Unlock()

	// Connect to NATS with Cluster Id, Client Id and customized options.
	nc, err := nats.Connect(url, options...)
	if err != nil {
		return err
	}
	c.nc = nc
	return err
}

// NATS returns the current NATS connection.
func (c *NATSComponent) NATS() *nats.Conn {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return c.nc
}

// JetStreamContext returns the returns a JetStreamContext
// for messaging and stream management.
func (c *NATSComponent) JetStreamContext(opts ...nats.JSOpt) (nats.JetStreamContext, error) {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	jsContext, err := c.nc.JetStream(opts...)
	return jsContext, err
}

// Name is the label used to identify the NATS connection.
func (c *NATSComponent) Name() string {
	c.cmu.Lock()
	defer c.cmu.Unlock()
	return c.name
}

// Shutdown makes the component go away.
func (c *NATSComponent) Shutdown() error {
	c.NATS().Close()
	return nil
}
