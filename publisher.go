package nats_pubsub_go

import (
	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(conn *nats.Conn) *Publisher {
	return &Publisher{conn: conn}
}

// Publish sends one or more NATS messages to the connected NATS server.
// It takes a variadic number of *nats.Msg objects and returns an error if
// any message fails to be published. If no messages are provided, it returns an
// ErrInvalidArgument indicating invalid input.
func (p *Publisher) Publish(messages ...*nats.Msg) error {
	// Check if the messages slice is empty or nil. If so, return an error indicating invalid arguments.
	// This handles the case where no messages are provided or an improper call is made.
	if messages == nil || len(messages) == 0 {
		return ErrInvalidArgument
	}

	// Iterate over each message in the messages slice.
	// This loop processes each message individually, attempting to publish it.
	for _, msg := range messages {
		// Attempt to publish the current message using the PublishMsg method.
		// If publishing fails, return the encountered error immediately.
		if err := p.conn.PublishMsg(msg); err != nil {
			return err
		}
	}

	// If all messages are successfully published, return nil to indicate success.
	// This means no errors occurred during the publishing process.
	return nil
}
