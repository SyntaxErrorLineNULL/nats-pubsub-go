package pkg

import (
	pubsub "github.com/SyntaxErrorLineNULL/nats-pubsub-go"
	"github.com/nats-io/nats.go"
)

// Subscriber represents a subscription to a NATS server.
// It includes the connection to the server and a flag indicating whether
// the Subscriber has been closed.
type Subscriber struct {
	// conn holds the connection to the NATS server.
	// This connection is used to publish messages.
	conn *nats.Conn

	// isClose is a flag indicating whether the Subscriber has been closed.
	// Once set to true, the Subscriber should not allow further publishing.
	isClose bool
}

// NewSubscriber creates a new Subscriber instance with the given NATS connection.
// It initializes the Subscriber with the provided connection and sets the
// isClose flag to false, indicating that the Subscriber is open for operations.
func NewSubscriber(conn *nats.Conn) *Subscriber {
	// Return a pointer to a new Subscriber instance initialized with the provided connection.
	// The isClose flag is initialized to its zero value, which is false.
	return &Subscriber{conn: conn}
}

// AsyncSubscribe subscribes to a subject asynchronously and returns a Subscription object
// that allows receiving messages on the subscribed subject. It handles errors and provides
// a channel for receiving messages.
// Docs: https://docs.nats.io/using-nats/developer/receiving/async
func (s *Subscriber) AsyncSubscribe(subject string) (pubsub.MessageHandler, error) {
	// Check if the provided subject is empty.
	// An empty subject is invalid and cannot be subscribed to.
	// Return an ErrInvalidArgument error to indicate the issue.
	if subject == "" {
		return nil, pubsub.ErrInvalidArgument
	}

	// Create a channel to receive incoming messages asynchronously.
	// This channel will be used to pass messages from the NATS subscription callback
	// to the code that is using the subscription.
	messages := make(chan *nats.Msg)

	// Subscribe to the specified subject with a callback function that sends
	// received messages to the messages channel.
	sub, err := s.conn.Subscribe(subject, func(msg *nats.Msg) { messages <- msg })

	// Check for errors that occurred during subscription.
	// If an error is returned, it indicates that the subscription could not be created.
	// Return nil for the Subscription and the error to signal failure.
	if err != nil {
		return nil, err
	}

	// Return a new MessageHandler instance with the created messages channel
	// and the NATS subscription. This allows the caller to receive and handle
	// messages asynchronously.
	return &MessageHandler{Message: messages, Subscription: sub}, nil
}

// SyncSubscribe creates a synchronous subscription to the provided subject,
// returning a MessageHandler that can be used to receive messages. It validates
// the subject and ensures proper error handling for subscription issues.
func (s *Subscriber) SyncSubscribe(subject string) (pubsub.MessageHandler, error) {
	// Check if the provided subject is empty.
	// An empty subject is invalid and cannot be subscribed to.
	// Return an ErrInvalidArgument error to indicate the issue.
	if subject == "" {
		return nil, pubsub.ErrInvalidArgument
	}

	// Attempting to create a synchronous subscription to a provided object using a NATS connection.
	// The SubscribeSync method from the NATS library sets up a blocking subscription.
	// which will wait for messages coming to the specified object.
	sub, err := s.conn.SubscribeSync(subject)

	// Check if an error occurred while attempting to subscribe.
	// If so, return nil for the Subscription and the encountered error.
	// This ensures that the caller is aware of the failure and can handle it accordingly.
	if err != nil {
		return nil, err
	}

	// If the subscription is successfully created, wrap it in a MessageHandler struct.
	// The MessageHandler struct embeds the actual subscription object and provides
	// additional methods or properties for managing the subscription lifecycle or processing messages.
	return &MessageHandler{Subscription: sub}, nil
}

// AsyncQueueSubscribe creates a new subscription to the specified NATS subject and queue.
// It returns a Subscription object or an error if the subject or queue is empty.
// The Subscribe method registers a callback function to handle incoming messages asynchronously.
func (s *Subscriber) AsyncQueueSubscribe(subject, queue string) (pubsub.MessageHandler, error) {
	// Check if the provided subject or queue is empty.
	// An empty subject or queue is invalid and cannot be subscribed to.
	// Return an ErrInvalidArgument error to indicate the issue.
	if subject == "" || queue == "" {
		return nil, pubsub.ErrInvalidArgument
	}

	// Create a channel to receive incoming messages asynchronously.
	// This channel will be used to pass messages from the NATS subscription callback
	// to the code that is using the subscription.
	messages := make(chan *nats.Msg)

	// Attempt to create a queue subscription to the provided subject and queue group.
	// The QueueSubscribe method from the NATS library establishes a subscription where
	// multiple subscribers can share the load of message processing.
	sub, err := s.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		// When a message is received on the subscription, it is sent to the messages channel.
		messages <- msg
	})

	// Check if an error occurred during the subscription setup.
	// If an error is returned, it indicates that the subscription could not be created.
	// Return nil for the Subscription and the encountered error to signal failure.
	if err != nil {
		return nil, err
	}

	// If the subscription is created successfully, return a MessageHandler with
	// the messages channel and the subscription object.
	// The MessageHandler allows receiving messages from the subscription asynchronously.
	return &MessageHandler{Message: messages, Subscription: sub}, nil
}

// SyncQueueSubscribe subscribes to a NATS subject with a queue group in a synchronous manner.
// This means it will block until a message is received from the subject within the specified queue group.
// It returns a MessageHandler to manage the subscription and receive messages.
func (s *Subscriber) SyncQueueSubscribe(subject, queue string) (pubsub.MessageHandler, error) {
	// Check if the provided subject or queue is empty.
	// An empty subject or queue is invalid and cannot be subscribed to.
	// Return an ErrInvalidArgument error to indicate the issue.
	if subject == "" || queue == "" {
		return nil, pubsub.ErrInvalidArgument
	}

	// Attempt to create a synchronous queue subscription to the provided subject and queue using the NATS connection.
	// The QueueSubscribeSync method from the NATS library is used to establish a blocking subscription
	// that will wait for messages to arrive at the given subject within the specified queue group.
	sub, err := s.conn.QueueSubscribeSync(subject, queue)
	// Check for errors that occurred during the subscription attempt.
	// If an error is returned, it indicates that the subscription could not be created.
	// Return nil for the subscription and the error to signal failure.
	if err != nil {
		return nil, err
	}

	// If the subscription is successfully created, wrap it in a MessageHandler struct and return it.
	// The MessageHandler struct includes the actual subscription object and can be used to receive messages.
	return &MessageHandler{Subscription: sub}, nil
}

// Close terminates the connection associated with the Subscriber and marks it as closed.
// It ensures that any ongoing communication with the NATS server is properly finalized and
// resources are released. This method should be called when the Subscriber is no longer needed
// to prevent resource leaks and ensure graceful shutdown.
func (s *Subscriber) Close() error {
	// Mark the Subscriber as closed by setting the isClose flag to true.
	// This flag indicates that the Subscriber is no longer active and should not
	// allow further message subscriptions or publications.
	s.isClose = true

	// Attempt to drain any remaining messages from the NATS connection.
	// The Drain method ensures that all pending messages are processed before closing the connection.
	// If an error occurs during this process, it is returned to signal that the connection
	// could not be properly closed.
	if err := s.conn.Drain(); err != nil {
		return err
	}

	// Return nil to indicate that the connection was successfully closed.
	// If no errors occurred during the draining process, the Subscriber is now safely closed.
	return nil
}
