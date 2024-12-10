package pkg

import (
	"context"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
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
	// Subject defines the NATS subject associated with this message.
	// It serves as the primary routing key for message delivery within the NATS system.
	Subject string

	// Header contains metadata associated with the message in key-value format.
	// It provides additional context or configuration for the message,
	// allowing consumers to interpret or process it effectively.
	Header Header

	// RequestID is a unique identifier for the message.
	// This ID is used to track and correlate requests and responses in the messaging system.
	RequestID string `json:"request_id,omitempty"`

	// Container holds the payload of the message as a byte slice.
	// It represents the data being transmitted or processed in the communication.
	Container Container `json:"container,omitempty"`

	// RequestTime records the time when the request was created or received.
	// This timestamp is valuable for tracking message lifecycle and processing logic of some tasks.
	RequestTime time.Time `json:"request_time,omitempty"`

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

// ReceiveMessage waits for the next message on the subscription with the specified timeout duration.
// It returns the received NATS message or an error if the operation times out.
// Note: This method can only be used with SyncSubscribe.
func (msg *Message) ReceiveMessage(timeout time.Duration) (*nats.Msg, error) {
	// Wait for the next message on the subscription with the given timeout duration.
	// The NextMsg method blocks until a message is received or the timeout is reached.
	// If a message is received, it is returned; otherwise, an error is returned.
	nextMessage, err := msg.subscription.NextMsg(timeout)
	if err != nil {
		return nil, err
	}
	// Return the received message along with a nil error if NextMsg succeeds.
	// This means the message was successfully retrieved within the specified timeout.
	return nextMessage, nil
}

// Unsubscribe terminates the subscription and closes the data channel.
// It ensures that the channel is closed only once and that the subscription
// is properly unsubscribed from. This method helps clean up resources
// and prevent memory leaks or dangling subscriptions.
func (msg *Message) Unsubscribe() (err error) {
	// Ensure the Data channel is closed only once by using the sync.Once mechanism.
	// The sync.Once type ensures that the provided function is executed only once,
	// regardless of how many times it's called.
	msg.once.Do(func() {

		// Unsubscribe from the current subscription to stop receiving messages.
		// The Unsubscribe method call removes the subscription and cleans up resources.
		err = msg.subscription.Unsubscribe()
	})

	// Return any error encountered during the Unsubscribe process.
	return err
}

// GetContainer retrieves the container payload from the underlying NATS message data.
// This method returns the raw data associated with the message, allowing consumers
// to access the payload for further processing or handling.
func (msg *Message) GetContainer() Container {
	// Access and return the data field from the underlying NATS message.
	// This represents the payload of the message that was received from NATS.
	return msg.message.Data
}

// GetHeader retrieves the header information from the underlying NATS message.
// This method converts the header from the NATS message into the custom Header type,
// allowing consumers to access metadata associated with the message in a structured manner.
func (msg *Message) GetHeader() Header {
	// Access and convert the header field from the underlying NATS message.
	// The header contains key-value pairs representing metadata about the message.
	return Header(msg.message.Header)
}

// Ack acknowledges the receipt of the message, notifying the NATS system that it has been processed.
// It supports an optional timeout for acknowledgment and cancels if the parent context is done.
// If a timeout is not specified, the message is acknowledged immediately.
// If the parent context is canceled before the timeout, the acknowledgment is aborted, and the context error is returned.
func (msg *Message) Ack(timeout time.Duration) error {
	// Check if no timeout is specified.
	// If timeout is zero, acknowledge the message immediately without delay.
	if timeout == 0 {
		// Acknowledge the message immediately when no timeout is set.
		return msg.message.Ack()
	}

	// Create a ticker that will emit an event after the specified timeout duration.
	// This provides a mechanism to handle the delay in acknowledging the message, based on the given timeout.
	ticker := time.NewTicker(timeout)
	// Ensure that the ticker is properly stopped after use to release any associated resources.
	// This is important to avoid potential resource leaks or unnecessary background work.
	defer ticker.Stop()

	// Use a select statement to wait for either the timeout or a cancellation signal from the parent context.
	// This allows the function to handle both the timeout event and the context cancellation in a non-blocking manner.
	select {
	// Case for when the ticker triggers, signaling the timeout has elapsed.
	// Acknowledge the message at this point.
	case <-ticker.C:
		// Acknowledge the message after the timeout has elapsed.
		// This ensures that the message is acknowledged only after the specified waiting period.
		return msg.message.Ack()

	// Case for when the parent context is canceled before the timeout.
	// Return the context's error to indicate that the acknowledgment was not completed.
	case <-msg.parentCtx.Done():
		// Return the error associated with the context cancellation.
		// This ensures that the operation is properly terminated if the context is canceled.
		return msg.parentCtx.Err()
	}
}

// Nak sends a negative acknowledgment for the message, signaling that it was not processed successfully.
// This can inform the NATS system to requeue the message for further processing or take alternative action.
// It provides support for an optional delay before sending the negative acknowledgment.
// If a timeout is specified, it uses NakWithDelay to apply the delay; otherwise, it sends an immediate Nak.
func (msg *Message) Nak(timeout time.Duration) error {
	// Check if no timeout is specified for the negative acknowledgment.
	// If timeout is zero, immediately send the negative acknowledgment without any delay.
	if timeout == 0 {
		// Directly send a negative acknowledgment for the message, indicating immediate rejection.
		return msg.message.Nak()
	}

	// Send a negative acknowledgment with the specified delay.
	// This uses the provided timeout to delay the rejection, which can be useful in specific scenarios.
	return msg.message.NakWithDelay(timeout)
}

// Respond sends a response back to the sender of the message.
// This is typically used in a request-response pattern where the sender expects a reply to the message it sent.
// The data parameter contains the payload to be sent as the response.
func (msg *Message) Respond(data []byte) error {
	// Use the NATS Respond method to send the provided data as a response to the message.
	// This operation communicates the reply to the message's sender, adhering to the NATS messaging protocol.
	return msg.message.Respond(data)
}
