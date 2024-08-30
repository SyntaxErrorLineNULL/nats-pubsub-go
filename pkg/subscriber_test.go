package pkg

import (
	"testing"
	"time"

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

	// SuccessfulSubscription tests the behavior of the AsyncSubscribe method
	// when subscribing with a valid subject. It verifies that the subscription
	// correctly receives and processes messages published to that subject.
	t.Run("SuccessfulSubscription", func(t *testing.T) {
		// Define the subject for the message to be subscribed.
		// The subject acts as a channel or topic to which the message will be sent.
		// In this test, the subject is set to "test.subject".
		subject := "test.subscribe.subject"
		// Define the expected message to be published.
		// This message should match what is later received by the subscription.
		expectedMessage := []byte("test_message")

		// Call AsyncSubscribe with the defined subject to create a new subscription.
		// The method should return a valid subscription and no error if successful.
		subscription, errSubscribe := subscriber.AsyncSubscribe(subject)
		// Assert that no error occurred while creating the subscription.
		// This ensures that the subscription was successfully created.
		assert.NoError(t, errSubscribe, "Expected no error when subscribing to a valid subject")

		// Publish a message to the subject to test if the subscription receives it.
		// The message should match the expected message defined above.
		publishErr := natsConnection.Publish(subject, expectedMessage)
		// Assert that no error occurred while publishing the message.
		// This confirms that the message was sent successfully.
		assert.NoError(t, publishErr, "Expected no error when publishing to the subject")

		// Use a select statement to wait for the message to be received or time out.
		select {
		// Case when a message is received from the subscription's message channel.
		// This block ensures that the received message matches the expected message.
		case receivedMessage := <-subscription.GetMessage():
			// Assert that the received message data matches the expected message.
			// This verifies that the subscription correctly received and processed the message.
			assert.Equal(t, expectedMessage, receivedMessage.Data, "received message does not match expected")

		// Case when waiting times out after 2 seconds.
		// This block handles the scenario where no message is received within the timeout period.
		case <-time.After(2 * time.Second):
			// Fail the test with a timeout error if no message is received.
			// This indicates that the message was not received as expected.
			t.Fatal("Timed out waiting for message")
		}
	})
}
