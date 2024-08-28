package nats_pubsub_go

import (
	"testing"

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
}
