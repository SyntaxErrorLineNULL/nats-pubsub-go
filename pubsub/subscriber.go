package pubsub

import (
	"context"
	"sync/atomic"

	pubsub "github.com/SyntaxErrorLineNULL/nats-pubsub-go"

	"github.com/nats-io/nats.go"
)

// Subscriber represents a subscription to a NATS server.
// It includes the connection to the server and a flag indicating whether
// the Subscriber has been closed.
type Subscriber struct {
	// conn is the active connection to the NATS server.
	// This connection is responsible for sending and receiving messages, enabling interaction with the server.
	// It serves as the core interface for all subscription-related operations and remains active until the
	// Subscriber is closed or an error occurs that terminates the connection.
	conn *nats.Conn

	// isClose is a thread-safe flag that indicates whether the Subscriber has been closed.
	// This atomic boolean ensures concurrent safety, allowing multiple goroutines to check the
	// Subscriber's state without race conditions. When set to true, it signals that no further
	// operations, such as subscribing or publishing, should be performed on this Subscriber.
	isClose atomic.Bool

	// closeCh is a channel used to signal the closure of the Subscriber.
	// This channel is written to when the Subscriber is closed, notifying any
	// goroutines waiting for this event. It helps coordinate graceful shutdowns
	// and ensures resources tied to the Subscriber are properly cleaned up.
	closeCh chan struct{}

	// errCh is a channel used to communicate errors occurring within the Subscriber.
	// It allows for asynchronous error handling by providing a mechanism to report
	// issues such as failed unsubscriptions or other runtime errors. This enables
	// external processes to respond to errors without blocking the main flow of operations.
	errCh chan error
}

// NewSubscriber creates a new Subscriber instance with the given NATS connection.
// It initializes the Subscriber with the provided connection and sets the
// isClose flag to false, indicating that the Subscriber is open for operations.
func NewSubscriber(conn *nats.Conn) *Subscriber {
	// Return a pointer to a new Subscriber instance initialized with the provided connection.
	// The isClose flag is initialized to its zero value, which is false.
	return &Subscriber{conn: conn, closeCh: make(chan struct{}, 1), errCh: make(chan error)}
}

// Subscriber initializes a subscription to a specific subject and queue (optional) using the provided context.
// It validates the subscriber's state and input parameters, manages the subscription lifecycle,
// and provides a `MessageHandler` for processing messages. This method ensures proper error handling
// and cleanup, allowing seamless message consumption while respecting context cancellation.
func (s *Subscriber) Subscriber(ctx context.Context, subject, queue string) (pubsub.MessageHandlerInterface, error) {
	// Check if the subscriber is closed. If closed, return an ErrCloseConnection error.
	// This prevents a situation where the client has closed the Subscriber but then tries to perform some manipulations afterwards, guaranteeing
	// that no operations will be performed on a closed instance.
	if s.isClose.Load() {
		return nil, pubsub.ErrCloseConnection
	}

	// Check if the provided subject or queue is empty.
	// An empty subject is invalid and cannot be subscribed to.
	// Return an ErrInvalidArgument error to indicate the issue.
	if subject == "" {
		return nil, pubsub.ErrInvalidArgument
	}

	// Create a channel to receive incoming messages asynchronously.
	// This channel will be used to pass messages from the NATS subscription callback
	// to the code that is using the subscription.
	messages := make(chan *nats.Msg)

	// Subscribes to the specified subject and queue using the `subscribe` method.
	// This sets up the actual connection to NATS and binds the subject and queue to the provided channel.
	// If an error occurs during subscription, it is returned to the caller for handling.
	subscription, err := s.subscribe(subject, queue, messages)
	// Checks if an error occurred during the execution of the previous operation.
	// The presence of an error indicates that the subscription process failed due to reasons such as
	// invalid arguments, connection issues, or a misconfigured NATS server.
	if err != nil {
		// Returns `nil` for the subscription object and propagates the error to the caller.
		// This ensures that the caller is informed of the failure and can take appropriate action,
		// such as retrying the subscription or logging the issue for further investigation.
		return nil, err
	}

	// Creates a `MessageHandler` struct that wraps the message channel and the subscription object.
	// This handler provides an abstraction for managing the subscription and processing messages.
	// It is returned to the caller, allowing them to receive messages and control the subscription.
	handler := &MessageHandler{Message: messages, Subscription: subscription}

	// Starts a goroutine to manage the subscription lifecycle using the `subscribeProcess` method.
	// The `subscribeProcess` method listens for cancellation signals and ensures proper cleanup when triggered.
	// This ensures the subscription is properly unsubscribed when the context is canceled or the subscriber is closed.
	go func(handler *MessageHandler) {
		// Passes the context and handler to the subscription lifecycle manager.
		// This decouples the subscription management from the main application logic, improving maintainability.
		s.subscribeProcess(ctx, handler)
	}(handler)

	// Returns the `MessageHandler` to the caller, providing access to the subscription and message channel.
	// The caller can use this handler to process messages or manage the subscription as needed.
	return handler, nil
}

// The `subscribeProcess` method manages the subscription lifecycle by listening for cancellation or shutdown signals.
// When triggered, it ensures that the subscription is cleanly unsubscribed to release resources and avoid leaks.
// It also handles errors during the unsubscribe process by sending them to the error channel for further handling.
func (s *Subscriber) subscribeProcess(ctx context.Context, handler *MessageHandler) {
	// The select statement listens for two possible signals: context cancellation or a shutdown signal.
	// It ensures the subscription handler reacts appropriately to these signals by cleaning up resources.
	select {
	// The context cancellation case is triggered when the parent context is explicitly terminated.
	// This could occur due to timeouts, manual cancellation, or other upstream context-related events.
	case <-ctx.Done():
		// Calls the Unsubscribe method on the handler to release resources associated with the subscription.
		// The unsubscribe process ensures that the handler stops receiving messages and clears its state.
		if err := handler.Unsubscribe(); err != nil {
			// Sends the error to the `errCh` channel, which is monitored by other components.
			// This allows errors to be logged, handled, or reported without crashing the application.
			s.errCh <- err
		}
		// Exits the function after handling the cancellation signal to prevent further execution.
		// This ensures that no unintended actions are taken after the subscription is cleaned up.
		return

	// The shutdown case handles the closure of the `closeCh` channel, signaling the need to terminate the subscriber.
	// This occurs when the subscriber is intentionally shut down, often as part of the application's shutdown process.
	case <-s.closeCh:
		// Unsubscribes the handler to stop receiving messages and release any resources tied to the subscription.
		// Ensures that the subscription handler is properly terminated and cleaned up.
		if err := handler.Unsubscribe(); err != nil {
			// Sends any errors encountered during the unsubscribe process to the `errCh` channel.
			// This allows the application to handle errors gracefully and avoid leaving them unaddressed.
			s.errCh <- err
		}
		// Exits the function after processing the shutdown signal to ensure no further operations are performed.
		// Prevents unintended side effects by terminating the function once the handler is unsubscribed.
		return
	}
}

// The `subscribe` method sets up a subscription to a specific subject and queue using the NATS messaging system.
// It connects to a given subject, associates it with a queue, and sends messages to the provided channel.
// This method returns the created subscription or an error if the subscription process fails.
func (s *Subscriber) subscribe(subject, queue string, messageCh chan *nats.Msg) (*nats.Subscription, error) {
	// Creates a queue subscription to the specified subject and queue using the provided channel.
	// The `QueueSubscribeSyncWithChan` method ensures that messages matching the subject are delivered
	// to the channel, with the queue enabling load-balanced processing among multiple subscribers.
	subscription, err := s.conn.QueueSubscribeSyncWithChan(subject, queue, messageCh)
	// Checks if an error occurred during the subscription process.
	// If there is an error, it indicates a failure in establishing the subscription,
	// such as connectivity issues or invalid subject/queue parameters.
	if err != nil {
		// Returns the error to inform the caller of the failure.
		return nil, err
	}

	// Returns the created subscription object to the caller.
	// The subscription can be used for further operations, such as managing the subscription
	// or unsubscribing when no longer needed.
	return subscription, nil
}

// Close terminates the connection associated with the Subscriber and marks it as closed.
// It ensures that any ongoing communication with the NATS server is properly finalized and
// resources are released. This method should be called when the Subscriber is no longer needed
// to prevent resource leaks and ensure graceful shutdown.
func (s *Subscriber) Close() error {
	// Check if the connection is marked as closed by reading the state of the isClose atomic flag.
	// The isClose flag is used to track whether the connection has already been closed in a thread-safe manner.
	// If the flag indicates that the connection is closed, return the predefined ErrConnectionAlreadyClosed error.
	// This ensures that any further operations on a closed connection are prevented, maintaining stability and integrity.
	if s.isClose.Load() {
		return pubsub.ErrConnectionAlreadyClosed
	}

	// Attempt to drain any remaining messages from the NATS connection.
	// The Drain method ensures that all pending messages are processed before closing the connection.
	// If an error occurs during this process, it is returned to signal that the connection
	// could not be properly closed.
	if err := s.conn.Drain(); err != nil {
		return err
	}

	// Mark the Subscriber as closed by setting the isClose flag to true.
	// This flag indicates that the Subscriber is no longer active and should not
	// allow further message subscriptions or publications.
	s.isClose.Store(true)

	// Notify any goroutines or processes waiting on the `closeCh` channel that the Subscriber is being closed.
	// Sending a signal through the channel allows these processes to clean up their resources or terminate gracefully.
	s.closeCh <- struct{}{}

	// Close the `closeCh` channel to release resources and signal that no further notifications will be sent.
	// This step ensures that the channel is not inadvertently written to after the connection is closed.
	close(s.closeCh)

	// Return nil to indicate that the connection was successfully closed.
	// If no errors occurred during the draining process, the Subscriber is now safely closed.
	return nil
}
