package nats_pubsub_go

import (
	"context"
	"github.com/nats-io/nats.go"
	"sync"
)

type Container []byte

type Header map[string][]string

type Message struct {
	RequestID string
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

	parentCtx context.Context
}

func NewMessage() *Message {
	return &Message{}
}
