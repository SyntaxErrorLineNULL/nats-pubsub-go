package nats_pubsub_go

import (
	"github.com/nats-io/nats.go"
	"testing"
	"time"

	"github.com/SyntaxErrorLineNULL/nats-pubsub-go/test"
	"github.com/stretchr/testify/assert"
)

func TestPublisher(t *testing.T) {
	t.Parallel()

	// Start a mock NATS server for testing purposes
	err := test.InitNats()
	// Assert that the NATS server started without errors
	assert.NoError(t, err, "Expected no error when starting NATS server")
	// Ensure the NATS server is properly shut down after tests
	defer test.ShutdownNatsServer()

	// Get the NATS connection from the global state
	natsConnection := test.GetNatsConnection()
	// Assert that the NATS connection is initialized correctly
	assert.NotNil(t, natsConnection, "Expected the NATS connection to be initialized, but got nil")
	// Ensure the NATS connection is properly closed after tests
	// defer natsConnection.Close()

	// Create a new Publisher instance using the provided NATS connection.
	// This step initializes the Publisher object with the necessary connection
	// for publishing a messages.
	publisher := NewPublisher(natsConnection)

	// Assert that the Publisher instance is not nil.
	// This checks that the `NewPublisher` function returned a valid object and did not
	// encounter any errors or fail during initialization.
	// The test will fail if `Publisher` is nil, indicating that the initialization
	// did not succeed as expected.
	assert.NotNil(t, publisher, "expected Publisher instance to be initialized and not nil")

	// InvalidMessage tests the behavior of the Publish method when it is provided
	// with a nil message. This ensures that the method properly handles invalid input
	// by returning the expected error. It is essential to confirm that the method
	// does not process nil messages and appropriately signals an error when
	// given such input.
	t.Run("InvalidMessage", func(t *testing.T) {
		// Attempt to publish a nil message using the publisher's Publish method.
		// This simulates a scenario where invalid input is provided (i.e., a nil message).
		err = publisher.Publish(nil...)

		// Assert that an error is returned from the Publish method.
		// This check ensures that the method correctly identifies and handles invalid arguments.
		assert.Error(t, err, "expected an error due to nil message")

		// Assert that the error returned is of type ErrInvalidArgument.
		// This verifies that the specific error indicating invalid arguments is returned.
		// The assertion confirms that the method is behaving as expected when encountering nil input.
		assert.ErrorIs(t, err, ErrInvalidArgument, "expected ErrInvalidArgument due to nil message")
	})

	// SuccessfulPublish tests the ability of the `Publish` method to successfully publish a message to a specified subject.
	// It ensures that the message is correctly sent and can be received by a subscriber.
	// The test verifies that the publish operation completes without errors and that the received message's subject
	// and payload match what was sent. This confirms that the `Publish` method behaves as expected under normal conditions.
	t.Run("SuccessfulPublish", func(t *testing.T) {
		// Define the subject for the message to be published.
		// The subject acts as a channel or topic to which the message will be sent.
		// In this test, the subject is set to "test.subject".
		subject := "test.subject"
		// Define the payload for the message to be published.
		// The payload is the actual data or content of the message that will be sent to the subject.
		// In this test, the payload is a byte slice containing the string "test payload".
		payload := []byte("test payload")

		// Subscribe to the subject using the NATS connection, allowing us to receive the message.
		// The subscription is synchronous, meaning it will block and wait for messages to arrive.
		sub, errSub := natsConnection.SubscribeSync(subject)
		// Assert that the subscription was successful and no error occurred.
		// This ensures that we can receive messages published to this subject.
		assert.NoError(t, errSub, "failed to subscribe to the subject")

		// Publish the message to the specified subject using the Publisher's Publish method.
		// The message contains the subject and the payload defined earlier.
		err = publisher.Publish(&nats.Msg{Subject: subject, Data: payload})
		// Assert that the publish operation was successful and no error occurred.
		// If the `Publish` method encounters an error, this test will fail with the provided message.
		assert.NoError(t, err, "failed to publish message")

		// Retrieve the next message from the subscription with a timeout of 10 milliseconds.
		// This checks if the message was successfully published and received within the given time frame.
		msg, errMsg := sub.NextMsg(10 * time.Millisecond)
		// Assert that there was no error in receiving the message.
		// If no message is received within the timeout or another error occurs, this assertion will fail.
		assert.NoError(t, errMsg, "failed to receive message")

		// Assert that the subject of the received message matches the expected subject.
		// This ensures that the message was published to and received from the correct channel.
		assert.Equal(t, subject, msg.Subject, "expected subject to match")
		// Assert that the payload of the received message matches the expected payload.
		// This confirms that the correct data was transmitted without alteration or loss.
		assert.Equal(t, payload, msg.Data, "expected payload to match")
	})
}
