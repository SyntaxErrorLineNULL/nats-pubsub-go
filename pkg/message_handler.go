package pkg

import (
	"sync"
	"time"

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
	})

	// Return any error encountered during the Unsubscribe process.
	return err
}

// ReceiveMessage waits for the next message on the subscription with the specified timeout duration.
// It returns the received NATS message or an error if the operation times out.
// Note: This method can only be used with SyncSubscribe.
func (msg *MessageHandler) ReceiveMessage(timeout time.Duration) (*nats.Msg, error) {
	// Wait for the next message on the subscription with the given timeout duration.
	// The NextMsg method blocks until a message is received or the timeout is reached.
	// If a message is received, it is returned; otherwise, an error is returned.
	nextMessage, err := msg.Subscription.NextMsg(timeout)
	if err != nil {
		return nil, err
	}
	// Return the received message along with a nil error if NextMsg succeeds.
	// This means the message was successfully retrieved within the specified timeout.
	return nextMessage, nil
}

// GetMessage retrieves the channel for received messages from the subscription.
// It returns a channel through which messages can be received.
// Note: If the SyncSubscribe method is used, the channel will default to nil.
func (msg *MessageHandler) GetMessage() chan *nats.Msg {
	// Return the channel for received messages
	return msg.Message
}
