package pkg

import (
	"github.com/nats-io/nats.go"
	"sync"
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
