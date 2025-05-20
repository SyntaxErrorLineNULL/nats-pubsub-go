package nats_pubsub_go

import (
	"context"
	"errors"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	// ErrInvalidArgument is an error returned when an invalid argument is provided.
	// This is used to indicate that a function or method has been called with
	// arguments that do not meet the required criteria or format.
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrCloseConnection is an error returned when an operation is attempted
	// on a closed connection. It signifies that the connection has been
	// terminated and cannot be used for further operations.
	ErrCloseConnection = errors.New("connection is close")

	// ErrConnectionAlreadyClosed indicates that a connection closure was attempted on an already closed connection.
	// This error helps differentiate between the connection being in a valid state versus being redundantly closed.
	// By defining this error, the code provides a specific signal to handle such redundant closure attempts gracefully.
	ErrConnectionAlreadyClosed = errors.New("connection is already closed")
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

// Subscriber defines the interface for managing message subscriptions.
// It provides methods to subscribe to specific subjects and queues, ensuring proper lifecycle management,
// error handling, and asynchronous message processing. This interface abstracts the implementation details,
// enabling flexibility and easier testing of components that rely on subscriptions.
type Subscriber interface {
	// Subscriber creates a subscription to a specified subject and queue using a context for lifecycle control.
	// The method validates input parameters, initializes the subscription, and returns a MessageHandler for
	// message processing. It ensures proper cleanup and resource management when the context is canceled or errors occur.
	//
	// Parameters:
	//
	// - ctx: The context used to manage the subscription lifecycle, including cancellation signals.
	//
	// - subject: The subject to subscribe to, which acts as the topic for message delivery.
	//
	// - queue: Optional. The queue group for the subscription, enabling message load balancing among subscribers.
	//
	// Returns:
	//
	// - MessageHandlerInterface: An abstraction for handling the subscription and processing messages.
	//
	// - error: An error object if the subscription fails, such as due to invalid arguments or connection issues.
	Subscriber(ctx context.Context, subject, queue string) (MessageHandlerInterface, error)
}

// MessageHandlerInterface defines methods for handling and managing messages from a subscription.
// It provides functionality to receive messages, manage subscription lifecycle, and interact
// with the message channel.
type MessageHandlerInterface interface {
	// Unsubscribe stops receiving messages from the subscription and closes the connection
	// to the NATS server for this subscription. This method ensures that resources associated
	// with the subscription are properly released. It returns an error if there was an issue
	// during the unsubscription process.
	Unsubscribe() error

	// ReceiveMessage waits for a message to arrive at the subscription channel within the
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
