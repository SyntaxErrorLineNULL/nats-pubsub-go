package pkg

import (
	"testing"

	pubsub "github.com/SyntaxErrorLineNULL/nats-pubsub-go"
	"github.com/SyntaxErrorLineNULL/nats-pubsub-go/test"
	"github.com/stretchr/testify/assert"
)

func TestSubscriber(t *testing.T) {
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

	// Create a new Subscriber instance using the provided NATS connection.
	// The NewSubscriber function initializes a Subscriber with the given connection.
	// This sets up the subscriber to handle incoming messages from the NATS server.
	// Here, `natsConnection` is expected to be an established connection to the NATS server.
	subscriber := NewSubscriber(natsConnection)
	// Assert that the Subscriber instance is not nil.
	// This check ensures that the NewSubscriber function was successful in creating a valid Subscriber instance.
	// If `subscriber` is nil, it indicates a failure in the creation process, which would mean the setup is not correct.
	// The message provided in the assert statement will be displayed if the `subscriber` is found to be nil.
	assert.NotNil(t, subscriber)

	// SubscribeInvalidArguments tests the behavior of the AsyncSubscribe method
	// when an invalid argument (an empty subject) is provided. It verifies that
	// the method returns an appropriate error indicating the issue with the argument.
	t.Run("SubscribeInvalidArguments", func(t *testing.T) {
		// Attempt to call AsyncSubscribe with an empty subject string.
		// This simulates an invalid subscription request to test how the method handles such cases.
		_, err = subscriber.AsyncSubscribe("")

		// Assert that an error is returned when calling AsyncSubscribe with an empty subject.
		// This checks if the method correctly identifies and handles the invalid input.
		assert.Error(t, err, "Expected an error when subscribing with an empty subject")

		// Assert that the specific error returned is ErrInvalidArgument.
		// This verifies that the method returns the correct type of error for the given invalid argument.
		assert.ErrorIs(t, err, pubsub.ErrInvalidArgument, "Expected error to be ErrInvalidArgument when subscribing with an empty subject")
	})
}
