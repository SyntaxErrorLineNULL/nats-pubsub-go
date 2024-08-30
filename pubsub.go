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

// Subscriber defines an interface for managing various types of subscriptions to NATS subjects.
// It supports both synchronous and asynchronous subscriptions, with or without queue groups,
// allowing for flexible message handling strategies.
type Subscriber interface {
	// AsyncSubscribe creates an asynchronous subscription to a specified subject.
	// This method sets up a subscription that does not block the caller; instead, it allows
	// the caller to continue processing while messages are received in the background.
	// It returns a MessageHandler for managing the subscription and receiving messages, and
	// an error if there was an issue creating the subscription.
	AsyncSubscribe(subject string) (MessageHandler, error)

	// SyncSubscribe creates a synchronous subscription to a specified subject.
	// This method sets up a blocking subscription that will wait until a message is received.
	// The caller will be blocked until a message arrives or an error occurs.
	// It returns a MessageHandler for managing the subscription and receiving messages, and
	// an error if there was an issue creating the subscription.
	SyncSubscribe(subject string) (MessageHandler, error)

	// AsyncQueueSubscribe creates an asynchronous subscription to a specified subject and queue group.
	// Messages published to the subject will be distributed among the members of the queue group
	// in a round-robin fashion, providing load balancing for message processing.
	// This method sets up the subscription to listen for messages in a non-blocking manner.
	// It returns a MessageHandler for managing the subscription and receiving messages, and
	// an error if there was an issue creating the subscription.
	AsyncQueueSubscribe(subject, queue string) (MessageHandler, error)

	// SyncQueueSubscribe creates a synchronous subscription to a specified subject and queue group.
	// This method sets up a blocking subscription where messages are distributed among the queue
	// group members, with the caller being blocked until a message is received or an error occurs.
	// It returns a MessageHandler for managing the subscription and receiving messages, and
	// an error if there was an issue creating the subscription.
	SyncQueueSubscribe(subject, queue string) (MessageHandler, error)
}

// MessageHandler defines methods for handling and managing messages from a subscription.
// It provides functionality to receive messages, manage subscription lifecycle, and interact
// with the message channel.
type MessageHandler interface {
	// Unsubscribe stops receiving messages from the subscription and closes the connection
	// to the NATS server for this subscription. This method ensures that resources associated
	// with the subscription are properly released. It returns an error if there was an issue
	// during the unsubscription process.
	Unsubscribe() error

	// ReceiveMessage waits for a message to arrive on the subscription channel within the
	// specified timeout period. If a message arrives within the timeout, it is returned
	// along with any error that occurred. If the timeout elapses without receiving a message,
	// this method will return an error indicating a timeout.
	ReceiveMessage(timeout time.Duration) (*nats.Msg, error)

	// GetMessage returns a channel through which messages from the subscription are received.
	// The channel will provide messages asynchronously as they arrive, allowing the caller
	// to process messages as they come in. This method does not block and provides a way
	// to continuously receive messages in a non-blocking manner.
	GetMessage() chan *nats.Msg
}
