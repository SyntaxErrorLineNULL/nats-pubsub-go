package pkg

import (
	pubsub "github.com/SyntaxErrorLineNULL/nats-pubsub-go"
	"github.com/nats-io/nats.go"
)

type Subscriber struct {
	// conn holds the connection to the NATS server.
	// This connection is used to publish messages.
	conn *nats.Conn

	// isClose is a flag indicating whether the Subscriber has been closed.
	// Once set to true, the Subscriber should not allow further publishing.
	isClose bool
}

func NewSubscriber(conn *nats.Conn) *Subscriber {
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
	if len(subject) == 0 {
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
	if len(subject) == 0 {
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
