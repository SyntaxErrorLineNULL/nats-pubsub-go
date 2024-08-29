package nats_pubsub_go

import (
	"time"

	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn    *nats.Conn
	isClose bool
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

	// Flush any buffered messages to the NATS server to ensure they are sent immediately.
	// This helps confirm that all published messages are transmitted without delay.
	if err := p.conn.Flush(); err != nil {
		return err
	}

	// If all messages are successfully published, return nil to indicate success.
	// This means no errors occurred during the publishing process.
	return nil
}

// Request sends a message request using the NATS connection and waits for a response
// within the specified timeout period. It ensures that the message is not nil before
// attempting to send the request. If the request is successful, it returns the
// response message; otherwise, it returns an error.
func (p *Publisher) Request(message *nats.Msg, timeout time.Duration) (*nats.Msg, error) {
	// Check if the provided message is nil. If it is, return an ErrInvalidArgument error.
	// This prevents attempting to send a nil message, which would result in a runtime panic.
	if message == nil {
		return nil, ErrInvalidArgument
	}

	// Send the request message using the NATS connection and wait for a response.
	// The response will be received within the specified timeout period.
	// If the request fails, return the error to the caller.
	msg, err := p.conn.RequestMsg(message, timeout)
	if err != nil {
		return nil, err
	}

	// If the request is successful, return the received response message.
	return msg, nil
}
