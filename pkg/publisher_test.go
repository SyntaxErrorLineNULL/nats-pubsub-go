package pkg

import (
	"testing"
	"time"

	"github.com/SyntaxErrorLineNULL/nats-pubsub-go"
	"github.com/SyntaxErrorLineNULL/nats-pubsub-go/test"
	"github.com/nats-io/nats.go"
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
		assert.ErrorIs(t, err, nats_pubsub_go.ErrInvalidArgument, "expected ErrInvalidArgument due to nil message")
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

	// RequestEmptyMessage tests the behavior of the Request method when it is called
	// with a nil message. It verifies that the method returns the appropriate error
	// when the input is invalid, specifically ensuring that an empty message results
	// in an ErrInvalidArgument error being returned.
	t.Run("RequestEmptyMessage", func(t *testing.T) {
		// Call the `Request` method on the publisher instance with a nil message
		// and a timeout of 0. This tests how the method handles a nil message,
		// which is considered an invalid argument.
		_, err = publisher.Request(nil, 0)

		// Assert that an error is returned, as the method should not accept a nil message.
		// This check ensures that the function correctly identifies and rejects invalid input.
		assert.Error(t, err, "Expected an error when calling Request with a nil message")

		// Assert that the error returned by the `Request` method is of type
		// `ErrInvalidArgument`. This verifies that the method correctly identifies
		// and returns an error for the invalid input scenario (nil message).
		assert.ErrorIs(t, err, nats_pubsub_go.ErrInvalidArgument, "expected ErrInvalidArgument for nil message")
	})

	// SuccessfulRequest tests the behavior of the Request method in the Publisher.
	// It verifies that a request message can be sent and that a response message is received
	// correctly within the expected time frame. The test ensures that the system under test (SUT)
	// is capable of handling request-response interactions as expected.
	t.Run("SuccessfulRequest", func(t *testing.T) {
		// Define the subject for the message to be published.
		// The subject acts as a channel or topic to which the message will be sent.
		// In this test, the subject is set to "test.subject".
		subject := "test.subject"

		// Create a buffered channel with a capacity of 1, which will be used to signal when the response is received.
		// The channel `resCh` is of type `chan struct{}`, which is commonly used for signaling without carrying any data.
		resCh := make(chan struct{}, 1)

		// Ensure that the channel `resCh` is closed when the test function exits.
		// `defer close(resCh)` schedules the closing of the channel to happen after the function completes,
		// which helps in cleaning up resources and avoiding potential memory leaks.
		defer close(resCh)

		// Define the expected request message to be sent to the PubSub system.
		// This byte slice represents the data that should be sent with the request message.
		expectedRequestMessage := []byte("test_request")

		// Define the expected publish message that should be received in response to the request.
		// This byte slice represents the data that the publisher should send as a reply.
		expectedPublishMessage := []byte("test_publish")

		// Subscribe to the subject using the NATS connection.
		// This simulates a service that listens for the request and sends back a response.
		_, _ = natsConnection.Subscribe(subject, func(msg *nats.Msg) {
			// Verify that the received message data matches the expected request message.
			// `assert.Equal` checks if `msg.Data` (the data in the received message) is equal to `expectedRequestMessage`.
			// If the values are not equal, the test will fail, and the provided message will be displayed.
			assert.Equal(t, expectedRequestMessage, msg.Data, "received message does not match expected")

			// Publish the expected response message to the reply subject specified in the incoming message.
			// The `msg.Reply` field contains the subject where the response should be sent.
			// `expectedPublishMessage` is the data payload that will be sent as the response.
			err = natsConnection.Publish(msg.Reply, expectedPublishMessage)

			// Check if there was an error while attempting to publish the response message.
			// The `assert.NoError` function ensures that the publish operation was successful.
			// If an error occurred, the test will fail, and the provided message will be displayed.
			assert.NoError(t, err, "failed to publish response message")

			// Signal that the response has been published by sending a value to the result channel (resCh).
			resCh <- struct{}{}
		})

		// Send a request message to the NATS server using the `publisher.Request` method.
		// The request message is constructed with the specified `subject` and `expectedRequestMessage` data.
		// A timeout of 1 second is provided for the request, meaning the request will wait up to 1 second for a response.
		response, err := publisher.Request(&nats.Msg{Subject: subject, Data: expectedRequestMessage}, 1*time.Second)
		// Check if the `publisher.Request` method returned an error.
		// The `assert.NoError` function verifies that `err` is nil. If `err` is not nil (indicating an error occurred),
		// the test will fail, and the provided message ("failed to send request message") will be displayed.
		assert.NoError(t, err, "failed to send request message")

		// Use a select statement to handle multiple cases: receiving a response or timing out.
		select {
		// Case when a response message is received on the `resCh` channel.
		// The response channel `resCh` is signaled when the expected publish message is received.
		case <-resCh:
			// Assert that the data in the received response message matches the expected publish message.
			// The `assert.Equal` function compares `expectedPublishMessage` with `response.Data` to ensure they are equal.
			// If they do not match, the test will fail with the provided message.
			assert.Equal(t, expectedPublishMessage, response.Data, "received message does not match expected")

			// Exit the loop and the test since the expected response has been successfully received and validated.
			return

		// Case when the timeout duration (5 seconds) is reached without receiving a response.
		// The `<-time.After(5 * time.Second)` statement waits for 5 seconds and then proceeds.
		case <-time.After(5 * time.Second):
			// Fail the test with a timeout error if no response was received within the 5-second window.
			// The `t.Fatal` function logs the provided message and stops the test execution.
			t.Fatal("Timed out waiting for message")
		}
	})

	// Close tests the behavior of the Publisher's Close method. It verifies that the
	// Close method correctly marks the Publisher as closed and prevents further
	// publishing operations by ensuring that an appropriate error is returned when
	// attempting to publish after the Publisher has been closed.
	t.Run("Close", func(t *testing.T) {
		// Create a new Publisher instance using a mock NATS connection.
		// This sets up a Publisher that can be tested for its closing behavior.
		closePublish := NewPublisher(natsConnection)
		// Assert that the Publisher instance is not nil to ensure that it was created successfully.
		assert.NotNil(t, closePublish, "Expected Publisher instance to be created successfully")

		// Call the Close method on the Publisher instance to mark it as closed
		// and to close the underlying NATS connection. This simulates the action of
		// terminating the Publisher's activity.
		closePublish.Close()

		// Assert that the isClose flag is set to true after calling Close.
		// This confirms that the Publisher has been marked as closed and no further
		// publishing operations should be allowed.
		assert.True(t, closePublish.isClose, "Expected Publisher to be marked as closed")

		// Attempt to publish a message using the now-closed Publisher.
		// This should fail because the Publisher's connection has been closed.
		err = closePublish.Publish(&nats.Msg{Subject: "close-method"})

		// Assert that an error is returned when attempting to publish with a closed Publisher.
		// This verifies that the Close method correctly prevents further publishing operations
		// and returns an appropriate error when the Publisher is closed.
		assert.Error(t, err, "Expected error when publishing with a closed Publisher")

		// Assert that the error returned is of type ErrCloseConnection.
		// This checks that the specific error returned when trying to publish with a closed Publisher
		// is the expected ErrCloseConnection error. This ensures that the Publisher correctly handles
		// the situation where an operation is attempted after it has been closed.
		assert.ErrorIs(t, err, nats_pubsub_go.ErrCloseConnection, "Expected error to be ErrCloseConnection when publishing with a closed Publisher")
	})
}
