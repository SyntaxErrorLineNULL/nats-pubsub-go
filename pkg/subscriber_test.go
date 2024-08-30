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
	assert.NotNil(t, subscriber, "Expected Subscriber instance to be created successfully")

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

	// SyncSubscribeInvalidArguments tests the behavior of the SyncSubscribe method
	// when an invalid argument (an empty subject) is provided. It verifies that
	// the method returns an appropriate error indicating the issue with the argument.
	t.Run("SyncSubscribeInvalidArguments", func(t *testing.T) {
		// Attempt to call SyncSubscribe with an empty subject string.
		// This simulates an invalid subscription request to test how the method handles such cases.
		_, err = subscriber.SyncSubscribe("")

		// Assert that an error is returned when calling SyncSubscribe with an empty subject.
		// This checks if the method correctly identifies and handles the invalid input.
		assert.Error(t, err, "Expected an error when subscribing with an empty subject")

		// Assert that the specific error returned is ErrInvalidArgument.
		// This verifies that the method returns the correct type of error for the given invalid argument.
		assert.ErrorIs(t, err, pubsub.ErrInvalidArgument, "Expected error to be ErrInvalidArgument when subscribing with an empty subject")
	})

	// SuccessSyncSubscribe tests the behavior of the SyncSubscribe method
	// when creating a synchronous subscription and receiving a published message.
	// It verifies that the subscription correctly receives and processes messages
	// published to the specified subject.
	t.Run("SuccessSyncSubscribe", func(t *testing.T) {
		// Define the subject for the message to be published.
		// The subject acts as a channel or topic to which the message will be sent.
		subject := "test.sync.subject"
		// Define the payload for the message to be published.
		// The payload is the actual data or content of the message that will be sent to the subject.
		payload := []byte("test payload")

		// Create a synchronous subscription to the defined subject using the subscriber instance.
		// The SyncSubscribe method sets up a subscription that waits for messages on the given subject.
		// It returns a MessageHandler and any error encountered during subscription setup.
		sub, subErr := subscriber.SyncSubscribe(subject)
		// Assert that no error occurred while creating the synchronous subscription.
		// This ensures that the subscription was established successfully.
		assert.NoError(t, subErr, "Expected no error when subscribing to the subject")

		// Publish a message to the subject to test if the subscription receives it.
		// The message should match the expected message defined above.
		publishErr := natsConnection.Publish(subject, payload)
		// Assert that no error occurred while publishing the message.
		// This confirms that the message was sent successfully.
		assert.NoError(t, publishErr, "Expected no error when publishing to the subject")

		// Retrieve the next message from the subscription with a timeout of 10 milliseconds.
		// This checks if the message was successfully published and received within the given time frame.
		msg, errMsg := sub.ReceiveMessage(10 * time.Millisecond)
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

	// AsyncQueueSubscribeInvalidArguments tests the behavior of the AsyncQueueSubscribe method
	// when provided with invalid arguments. It verifies that the method correctly handles
	// cases where either the subject or queue is empty, returning the appropriate error.
	t.Run("AsyncQueueSubscribeInvalidArguments", func(t *testing.T) {
		// Attempt to subscribe to a queue with an empty subject and a valid queue name "q".
		// This simulates an invalid input scenario where the subject is missing, which is
		// expected to result in an error.
		_, err = subscriber.AsyncQueueSubscribe("", "q")
		// Assert that an error was returned when attempting to subscribe with an empty subject.
		// The test expects an error to be returned in this scenario.
		assert.Error(t, err, "Expected error when subscribing with an empty subject")
		// Assert that the specific error returned matches pubsub.ErrInvalidArgument.
		// This ensures that the method correctly identifies the empty subject as an invalid argument.
		assert.ErrorIs(t, err, pubsub.ErrInvalidArgument, "Expected error to be ErrInvalidArgument when subscribing with an empty subject")

		// Attempt to subscribe to a queue with a valid subject "s" but an empty queue name.
		// This simulates an invalid input scenario where the queue name is missing, which is
		// also expected to result in an error.
		_, err = subscriber.AsyncQueueSubscribe("s", "")
		// Assert that an error was returned when attempting to subscribe with an empty queue name.
		// The test expects an error to be returned in this case.
		assert.Error(t, err, "Expected error when subscribing with an empty queue name")
		// Assert that the specific error returned matches pubsub.ErrInvalidArgument.
		// This ensures that the method correctly identifies the empty queue name as an invalid argument.
		assert.ErrorIs(t, err, pubsub.ErrInvalidArgument, "Expected error to be ErrInvalidArgument when subscribing with an empty queue name")
	})

	// SuccessfulQueueSubscription tests the behavior of the AsyncQueueSubscribe method
	// when valid arguments are provided. It verifies that the method can successfully
	// subscribe to a queue, receive a message, and handle it correctly.
	t.Run("SuccessfulQueueSubscription", func(t *testing.T) {
		// Define the subject for the message to be subscribed.
		// The subject acts as a channel or topic to which the message will be sent.
		// In this test, the subject is set to "test.subject".
		subject := "test_subscribe_subject"
		// Define the expected message to be published.
		// This message should match what is later received by the subscription.
		expectedMessage := []byte("test_message")
		// Define the queue name for the subscription.
		// This specifies the queue group to which messages will be delivered.
		queue := "test_queue"

		// Create an asynchronous queue subscription with the defined subject and queue.
		// This sets up the subscription to listen for messages on the specified subject and queue.
		subscription, err := subscriber.AsyncQueueSubscribe(subject, queue)
		// Assert that no error occurred during the subscription setup.
		// This ensures that the subscription was created successfully.
		assert.NoError(t, err, "Expected no error when creating queue subscription")

		// Publish a message to the defined subject.
		// This message will be sent to the specified subject and should be received by the subscription.
		err = natsConnection.Publish(subject, expectedMessage)
		// Assert that no error occurred while publishing the message.
		// This confirms that the message was sent successfully to the NATS server.
		assert.NoError(t, err, "Failed to publish message")

		// Wait for the message to be received by the subscription.
		// The select statement will either receive the message from the subscription or timeout after 2 seconds.
		select {
		case receivedMessage := <-subscription.GetMessage():
			// Assert that the received message data matches the expected message.
			// This confirms that the message received by the subscription is correct.
			assert.Equal(t, expectedMessage, receivedMessage.Data, "Received message does not match expected")
		case <-time.After(2 * time.Second): // Timeout
			// Fail the test if no message was received within the timeout period.
			// This ensures that the test will fail if the message is not received in a timely manner.
			t.Fatal("Timed out waiting for message")
		}
	})

	// Unsubscribe tests the behavior of the Unsubscribe method
	// to ensure that a subscription can be properly unsubscribed
	// and that the data channel is closed accordingly.
	t.Run("Unsubscribe", func(t *testing.T) {
		// Define the subject for the message to be subscribed.
		// The subject acts as a channel or topic to which the message will be sent.
		subject := "test_unsubscribe_subject"
		// Define the queue name for the subscription.
		// This specifies the queue group to which messages will be delivered.
		queue := "test_queue"

		// Create an asynchronous queue subscription with the defined subject and queue.
		// This sets up the subscription to listen for messages on the specified subject and queue.
		subscription, errSub := subscriber.AsyncQueueSubscribe(subject, queue)
		// Assert that no error occurred during the subscription setup.
		// This ensures that the subscription was created successfully.
		assert.NoError(t, errSub, "Expected no error when creating queue subscription")

		// Assert that the data channel is not nil after subscription.
		// This confirms that the subscription has an active channel for receiving messages.
		assert.NotNil(t, subscription.GetMessage(), "data channel is nil after subscription")

		// Call the Unsubscribe method on the subscription.
		// This should remove the subscription from the NATS server and close the data channel.
		errUnsubscribe := subscription.Unsubscribe()
		// Assert that no error occurred during the unsubscription process.
		// This ensures that the Unsubscribe method completed successfully.
		assert.NoError(t, errUnsubscribe, "failed to unsubscribe")

		// Attempt to receive from the data channel after unsubscribing.
		// This checks if the data channel has been closed as expected.
		_, ok := <-subscription.GetMessage()
		// Assert that the data channel is closed after unsubscribing.
		// The channel should be closed, so this assertion ensures that no more messages can be received.
		assert.False(t, ok, "data channel is not closed after unsubscribe")
	})

	// GetData tests the behavior of retrieving messages from a queue subscription.
	// It verifies that a message published to a subject is correctly received
	// through the queue subscription.
	t.Run("GetData", func(t *testing.T) {
		// Define the subject for the message to be subscribed.
		// The subject acts as a channel or topic to which the message will be sent.
		// In this test, the subject is set to "test_subject1".
		subject := "test_subject1"

		// Define the queue name for the subscription.
		// This specifies the queue group to which messages will be delivered.
		// In this test, the queue is named "test_queue1".
		queue := "test_queue1"

		// Create an asynchronous queue subscription with the defined subject and queue.
		// This sets up the subscription to listen for messages on the specified subject and queue.
		subscription, errSub := subscriber.AsyncQueueSubscribe(subject, queue)
		// Assert that no error occurred during the subscription setup.
		// This ensures that the subscription was created successfully and is ready to receive messages.
		assert.NoError(t, errSub, "Expected no error when creating queue subscription")

		// Ensure that the subscription is unsubscribed after the test completes.
		// This cleans up resources and avoids potential interference with other tests.
		defer subscription.Unsubscribe()

		// Define the expected message to be published.
		// This is the message that should be received through the subscription.
		expectedMessage := []byte("test_message")

		// Publish a message to the subject.
		// This sends the expected message to the subject so that the subscription can receive it.
		errPublish := natsConnection.Publish(subject, expectedMessage)
		// Assert that no error occurred while publishing the message.
		// This ensures that the message was sent successfully.
		assert.NoError(t, errPublish, "Failed to publish message")

		// Receive a message from the subscription's data channel.
		// This retrieves the next message received by the subscription.
		receivedMessage := <-subscription.GetMessage()

		// Assert that the received message data matches the expected message.
		// This verifies that the message published was correctly received by the subscription.
		assert.Equal(t, expectedMessage, receivedMessage.Data, "received message does not match expected")
	})
}
