package nats_pubsub_go

import (
	"context"
	"github.com/nats-io/nats.go"
	"sync"
)

// Container represents a byte slice used to store data within the Message structure.
// It acts as a container for the payload associated with the message.
type Container []byte

// Header represents a map of string keys to slices of string values.
// It is used to store metadata or headers associated with a message.
type Header map[string][]string

// Message represents a communication unit in a NATS-based messaging system.
// It encapsulates the necessary components for processing messages,
// including the payload, metadata, and underlying NATS-specific details.
type Message struct {
	// RequestID is a unique identifier for the message.
	// This ID is used to track and correlate requests and responses in the messaging system.
	RequestID string

	// Container holds the payload of the message as a byte slice.
	// It represents the data being transmitted or processed in the communication.
	Container Container

	// Message is a NATS message data.
	// This channel allows consumers to process incoming messages
	// that are published to the NATS subject this handler is subscribed to.
	message *nats.Msg

	// Subscription is the underlying NATS subscription.
	// This represents the subscription to a NATS subject or subjects,
	// allowing the handler to receive messages from NATS.
	subscription *nats.Subscription

	// once is used to ensure certain operations are performed only once.
	// It uses sync.Once to guarantee that specific actions, such as closing
	// the channel and unsubscribing, are executed only a single time.
	once sync.Once

	// parentCtx is the context associated with the message's parent operation.
	// It provides a way to propagate cancellation, timeouts, or deadlines across operations.
	parentCtx context.Context
}

func NewMessage(parentCtx context.Context) *Message {
	return &Message{parentCtx: parentCtx}
}

// GetContainer retrieves the container payload from the underlying NATS message data.
// This method returns the raw data associated with the message, allowing consumers
// to access the payload for further processing or handling.
func (msg *Message) GetContainer() Container {
	// Access and return the data field from the underlying NATS message.
	// This represents the payload of the message that was received from NATS.
	return msg.message.Data
}
