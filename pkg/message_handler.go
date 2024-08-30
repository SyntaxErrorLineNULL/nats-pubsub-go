package pkg

import (
	"sync"

	"github.com/nats-io/nats.go"
)

// MessageHandler manages the subscription to NATS messages and
// provides a channel to receive those messages. It also ensures
// that certain operations like unsubscribing are only performed once.
type MessageHandler struct {
	// Message is a channel for receiving NATS message data.
	// This channel allows consumers to process incoming messages
	// that are published to the NATS subject this handler is subscribed to.
	Message chan *nats.Msg

	// Subscription is the underlying NATS subscription.
	// This represents the subscription to a NATS subject or subjects,
	// allowing the handler to receive messages from NATS.
	Subscription *nats.Subscription

	// once is used to ensure certain operations are performed only once.
	// It uses sync.Once to guarantee that specific actions, such as closing
	// the channel and unsubscribing, are executed only a single time.
	once sync.Once
}

// Unsubscribe terminates the subscription and closes the data channel.
// It ensures that the channel is closed only once and that the subscription
// is properly unsubscribed from. This method helps clean up resources
// and prevent memory leaks or dangling subscriptions.
func (msg *MessageHandler) Unsubscribe() (err error) {
	// Ensure the Data channel is closed only once by using the sync.Once mechanism.
	// The sync.Once type ensures that the provided function is executed only once,
	// regardless of how many times it's called.
	msg.once.Do(func() {
		// Check if the Message channel is not nil before attempting to close it.
		// This prevents potential panics or runtime errors when trying to close a nil channel.
		if msg.Message != nil {
			// Close the Message channel to stop receiving any further messages.
			close(msg.Message)
		}

		// Unsubscribe from the current subscription to stop receiving messages.
		// The Unsubscribe method call removes the subscription and cleans up resources.
		err = msg.Subscription.Unsubscribe()
		// Return the error if Unsubscribe() failed.
		return
	})

	// Return any error encountered during the Unsubscribe process.
	return err
}
