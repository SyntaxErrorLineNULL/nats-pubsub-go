package nats_pubsub_go

import (
	"time"

	"github.com/nats-io/nats.go"
)

// Publisher is a structure that encapsulates a NATS connection and provides
// methods to publish messages to a NATS server. It also tracks whether the
// publisher has been closed to prevent further operations after closure.
type Publisher struct {
	// conn holds the connection to the NATS server.
	// This connection is used to publish messages.
	conn *nats.Conn

	// isClose is a flag indicating whether the Publisher has been closed.
	// Once set to true, the Publisher should not allow further publishing.
	isClose bool
}

// NewPublisher creates and returns a new instance of Publisher.
// It takes a *nats.Conn as an argument, which represents the active connection
// to the NATS server that the Publisher will use for sending messages.
func NewPublisher(conn *nats.Conn) *Publisher {
	// Initialize a new Publisher with the provided NATS connection.
	// The isClose flag is set to false by default, indicating the Publisher is active.
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

// Close terminates the Publisher instance, marking it as closed and
// closing the underlying NATS connection. This method ensures that
// no further publishing can occur and that resources are properly released.
func (p *Publisher) Close() {
	// Set the isClose flag to true, indicating that the Publisher is closed.
	// This flag can be used to prevent further publishing operations.
	p.isClose = true

	// Close the underlying NATS connection. This releases any resources associated
	// with the connection and ensures that the Publisher is properly shut down.
	p.conn.Close()
}
