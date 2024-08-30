package nats_pubsub_go

import (
	"time"

	"github.com/nats-io/nats.go"
)

// Publisher defines the interface for a publisher that can send messages,
// handle requests, and be closed. This interface abstracts the operations
// that a concrete publisher implementation must provide.
type Publisher interface {
	// Publish sends one or more messages to the publisher.
	// Each message should be of type *nats.Msg. If the publishing fails,
	// an error is returned indicating the failure.
	Publish(messages ...*nats.Msg) error

	// Request sends a message and waits for a response within the specified timeout.
	// The message should be of type *nats.Msg. If the request is successful,
	// it returns the response message of type *nats.Msg along with a nil error.
	// If the request fails or times out, an error is returned.
	Request(message *nats.Msg, timeout time.Duration) (*nats.Msg, error)

	// Close terminates the Publisher instance. It marks the Publisher as closed
	// and closes any underlying connections. After calling Close, the Publisher
	// should not be used for further operations, and an error is returned if
	// any operation is attempted after closure.
	Close()
}
