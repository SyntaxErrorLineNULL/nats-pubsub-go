package pubsub

import (
	"context"
	"testing"
	"time"

	pubsub "github.com/SyntaxErrorLineNULL/nats-pubsub-go"
	"github.com/SyntaxErrorLineNULL/nats-pubsub-go/test"

	"github.com/stretchr/testify/assert"
)

func TestSubscriber(t *testing.T) {
	t.Parallel()

	// Initialize the NATS server with a port of 18222.
	// Passing 0 as the port indicates that the server should choose a default port.
	// The InitNats function is expected to return an error if there was a failure in starting the NATS server.
	natsServer := test.NewNatsServer(19222)
	// Assert that the newly created NatsServer instance is not nil.
	// This check verifies that the NatsServer was successfully initialized.
	// If natsServer is nil, the test will fail, indicating an issue with the server creation.
	assert.NotNil(t, natsServer, "Expected NatsServer to be initialized, but got nil")

	// Initialize the NATS server with a port of 18222.
	// Passing 0 as the port indicates that the server should choose a default port.
	// The InitNats function is expected to return an error if there was a failure in starting the NATS server.
	err := natsServer.InitNats()
	// Assert that there is no error from the NATS server initialization.
	// Assert.NoError function checks that `err` (the actual value) is nil.
	// If `err` is not nil, it means there was an error during the initialization of the NATS server,
	// and the provided message "Expected no error when starting NATS server" will indicate the nature of the failure.
	assert.NoError(t, err, "Expected no error when starting NATS server")

	// Schedule the ShutdownNatsServer function to be called when the surrounding function (test) returns.
	// This ensures that the NATS server will be properly shut down after the test completes,
	// preventing resource leaks and ensuring a clean state for subsequent tests.
	defer natsServer.ShutdownNatsServer()

	// Retrieve the NATS connection instance from the test package.
	// The GetNatsConnection function is expected to return a pointer to an initialized NATS connection object.
	// This connection will be used for interacting with the NATS messaging system during the tests.
	natsConnection := natsServer.GetNatsConnection()
	// Assert that the retrieved NATS connection is not nil.
	// Assert.NotNil function checks that `natsConnection` (the actual value) is not nil.
	// If it is nil, the test will fail, and the provided message "Expected the NATS connection to be initialized, but got nil"
	// will indicate that the connection was not properly initialized before this point in the test.
	assert.NotNil(t, natsConnection, "Expected the NATS connection to be initialized, but got nil")

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

	// SuccessfulSubscription tests the behavior of the AsyncSubscribe method
	// when subscribing with a valid subject. It verifies that the subscription
	// correctly receives and processes messages published to that subject.
	t.Run("SuccessfulSubscription", func(t *testing.T) {
		// Define the subject for the Message to be subscribed.
		// The subject acts as a channel or topic to which the Message will be sent.
		subject := "test_subscribe_subject"
		// Define the queue name for the subscription.
		// This specifies the queue group to which messages will be delivered.
		queue := "test_queue"
		// Define the expected message to be published.
		// This message should match what is later received by the subscription.
		expectedMessage := []byte("test_message")

		// Create a subscription to the specified subject and queue using the provided context.
		// This initializes a subscription for receiving messages published to the given subject,
		// enabling message processing within the defined queue group.
		subscription, errSub := subscriber.Subscriber(context.Background(), subject, queue)
		// Assert that no error occurred during the subscription creation.
		// This ensures that the subscription was successfully initialized and is ready for use.
		assert.NoError(t, errSub, "Subscription creation should not return an error")

		// Publish a Message to the subject to test if the subscription receives it.
		// The Message should match the expected Message defined above.
		publishErr := natsConnection.Publish(subject, expectedMessage)
		// Assert that no error occurred while publishing the Message.
		// This confirms that the Message was sent successfully.
		assert.NoError(t, publishErr, "Expected no error when publishing to the subject")

		// Use a select statement to wait for the Message to be received or time out.
		select {
		// Case when a Message is received from the subscription's Message channel.
		// This block ensures that the received Message matches the expected Message.
		case receivedMessage := <-subscription.GetMessage():
			// Assert that the received message data matches the expected message.
			// This verifies that the subscription correctly received and processed the message.
			assert.Equal(t, expectedMessage, receivedMessage.Data, "received message does not match expected")

			err = subscription.Unsubscribe()
			assert.NoError(t, err)

		// Case when waiting times out after 2 seconds.
		// This block handles the scenario where no Message is received within the timeout period.
		case <-time.After(2 * time.Second):
			// Fail the test with a timeout error if no Message is received.
			// This indicates that the Message was not received as expected.
			t.Fatal("Timed out waiting for Message")
		}
	})

	// SubscribeWithContextIsDone tests the behavior of a subscriber's message processing functionality
	// within the context of a managed lifecycle. It verifies that the subscriber can correctly receive
	// messages published to a specific subject while respecting context cancellation. The test ensures
	// that the subscription behaves as expected, including proper cleanup of resources when the context
	// is cancelled, and validates that messages are no longer processed once the subscription is terminated.
	t.Run("SubscribeWithContextIsDone", func(t *testing.T) {
		// Create a context with cancellation to manage the lifecycle of the subscription.
		// This context will allow controlled termination of the subscription process.
		ctx, cancel := context.WithCancel(context.Background())

		// Define the subject for the Message to be subscribed.
		// The subject acts as a channel or topic to which the Message will be sent.
		subject := "test_subscribe_subject"
		// Define the queue name for the subscription.
		// This specifies the queue group to which messages will be delivered.
		queue := "test_queue"
		// Define the expected message to be published.
		// This message should match what is later received by the subscription.
		expectedMessage := []byte("test_message")

		// Create a subscription to the specified subject and queue using the provided context.
		// This initializes a subscription for receiving messages published to the given subject,
		// enabling message processing within the defined queue group.
		subscription, errSub := subscriber.Subscriber(ctx, subject, queue)
		// Assert that no error occurred during the subscription creation.
		// This ensures that the subscription was successfully initialized and is ready for use.
		assert.NoError(t, errSub, "Subscription creation should not return an error")
		// Assert that the subscription instance is not nil.
		// This verifies that a valid subscription object was returned and is usable for receiving messages.
		assert.NotNil(t, subscription, "Subscription should be successfully created and not nil")

		// Publish a Message to the subject to test if the subscription receives it.
		// The Message should match the expected Message defined above.
		publishErr := natsConnection.Publish(subject, expectedMessage)
		// Assert that no error occurred while publishing the Message.
		// This confirms that the Message was sent successfully.
		assert.NoError(t, publishErr, "Expected no error when publishing to the subject")

		// Use a select statement to wait for the Message to be received or time out.
		select {
		// Case when a Message is received from the subscription's Message channel.
		// This block ensures that the received Message matches the expected Message.
		case receivedMessage := <-subscription.GetMessage():
			// Assert that the received message data matches the expected message.
			// This verifies that the subscription correctly received and processed the message.
			assert.Equal(t, expectedMessage, receivedMessage.Data, "received message does not match expected")

			// Cancel the context to signal that the subscription should stop processing messages.
			// This triggers the cleanup process for the subscription and ensures it ceases operations.
			cancel()

			// Introduce a brief delay to allow the cancellation signal to propagate and the subscription to clean up.
			// This ensures that the subscription has enough time to process the cancellation before further checks.
			<-time.After(10 * time.Millisecond)

			// Use a select statement to check the status of the subscription's message channel.
			// This block verifies that the channel has been properly closed after cancellation.
			select {
			case _, ok := <-subscription.GetMessage():
				// Assert that the channel is closed (ok is false).
				// This confirms that the subscription has been terminated and is no longer receiving messages.
				assert.False(t, ok, "Subscription's message channel should be closed after cancellation")
			default:
				// Default case ensures that the select statement does not block if the channel is already inactive.
			}

		// Case when waiting times out after 2 seconds.
		// This block handles the scenario where no Message is received within the timeout period.
		case <-time.After(2 * time.Second):
			// Fail the test with a timeout error if no Message is received.
			// This indicates that the Message was not received as expected.
			t.Fatal("Timed out waiting for Message")
		}

	})

	// SubscribeInvalidArguments tests the behavior of the AsyncSubscribe method
	// when an invalid argument (an empty subject) is provided. It verifies that
	// the method returns an appropriate error indicating the issue with the argument.
	t.Run("SubscribeInvalidArguments", func(t *testing.T) {
		// Attempt to call AsyncSubscribe with an empty subject string.
		// This simulates an invalid subscription request to test how the method handles such cases.
		_, err = subscriber.Subscriber(context.Background(), "", "")

		// Assert that an error is returned when calling AsyncSubscribe with an empty subject.
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
	t.Run("SuccessReceiveMessage", func(t *testing.T) {
		// Define the subject for the message to be published.
		// The subject acts as a channel or topic to which the message will be sent.
		subject := "test.sync.subject"
		// Define the payload for the message to be published.
		// The payload is the actual data or content of the message that will be sent to the subject.
		payload := []byte("test payload")

		// Create a subscription to the specified subject and queue using the provided context.
		// This initializes a subscription for receiving messages published to the given subject,
		// enabling message processing within the defined(empty) queue group.
		sub, errSub := subscriber.Subscriber(context.Background(), subject, "")
		// Assert that no error occurred during the subscription creation.
		// This ensures that the subscription was successfully initialized and is ready for use.
		assert.NoError(t, errSub, "Subscription creation should not return an error")

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
		subscription, errSub := subscriber.Subscriber(context.Background(), subject, queue)
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

	// SubscribeWithCloseSubscriber tests the behavior of a subscriber's message processing and lifecycle
	// when the subscriber is explicitly closed. It ensures that the subscriber properly terminates,
	// closes its message channel, and prevents further message processing after being closed. This test
	// validates the correctness of the subscriber's resource cleanup and its response to the close operation.
	t.Run("SubscribeWithCloseSubscriber", func(t *testing.T) {
		// Define the subject for the Message to be subscribed.
		// The subject acts as a channel or topic to which the Message will be sent.
		subject := "test_subscribe_subject"
		// Define the queue name for the subscription.
		// This specifies the queue group to which messages will be delivered.
		queue := "test_queue"

		// Create a subscription to the specified subject and queue using the provided context.
		// This initializes a subscription for receiving messages published to the given subject,
		// enabling message processing within the defined queue group.
		subscription, errSub := subscriber.Subscriber(context.Background(), subject, queue)
		// Assert that no error occurred during the subscription creation.
		// This ensures that the subscription was successfully initialized and is ready for use.
		assert.NoError(t, errSub, "Subscription creation should not return an error")
		// Assert that the subscription instance is not nil.
		// This verifies that a valid subscription object was returned and is usable for receiving messages.
		assert.NotNil(t, subscription, "Subscription should be successfully created and not nil")

		// Close the subscriber to test its behavior when explicitly terminated.
		// This step ensures the subscriber shuts down cleanly, releasing all resources.
		err = subscriber.Close()
		// Assert that no error occurred while closing the subscriber.
		// This confirms the subscriber's resources were properly released without issues.
		assert.NoError(t, err, "Closing the subscriber should not return an error")

		// Attempt to create a subscription after the subscriber has been closed.
		// This tests whether the Subscriber method correctly prevents such operations.
		_, err := subscriber.Subscriber(context.Background(), subject, queue)
		// Verify that an error is returned, indicating the subscription creation failed as expected.
		// This confirms the method enforces lifecycle constraints effectively.
		// it is impossible to create a subscription after the Close method is called
		assert.Error(t, err, "An error should be returned when attempting to subscribe after the subscriber is closed")
		// Check that the returned error matches the specific ErrCloseConnection error.
		// This ensures the error type is consistent with the expected behavior.
		assert.ErrorIs(t, err, pubsub.ErrCloseConnection, "The error should be ErrCloseConnection when subscribing after closure")

		// Assert that the subscriber's state is correctly marked as closed.
		// This checks that the internal flag isClose is set to true, indicating the subscriber is no longer active.
		assert.True(t, subscriber.isClose.Load(), "Subscriber should be marked as closed after calling Close")

		// Use a select statement to check the status of the subscription's message channel.
		// This block verifies that the channel has been properly closed after cancellation.
		select {
		case _, ok := <-subscription.GetMessage():
			// Assert that the channel is closed (ok is false).
			// This confirms that the subscription has been terminated and is no longer receiving messages.
			assert.False(t, ok, "Subscription's message channel should be closed after cancellation")
		default:
			// Default case ensures that the select statement does not block if the channel is already inactive.
		}
	})
}

func TestUnsubscribe(t *testing.T) {
	t.Parallel()

	// Initialize the NATS server with a port of 18222.
	// Passing 0 as the port indicates that the server should choose a default port.
	// The InitNats function is expected to return an error if there was a failure in starting the NATS server.
	natsServer := test.NewNatsServer(18221)
	// Assert that the newly created NatsServer instance is not nil.
	// This check verifies that the NatsServer was successfully initialized.
	// If natsServer is nil, the test will fail, indicating an issue with the server creation.
	assert.NotNil(t, natsServer, "Expected NatsServer to be initialized, but got nil")

	// Initialize the NATS server with a port of 18222.
	// Passing 0 as the port indicates that the server should choose a default port.
	// The InitNats function is expected to return an error if there was a failure in starting the NATS server.
	err := natsServer.InitNats()
	// Assert that there is no error from the NATS server initialization.
	// Assert.NoError function checks that `err` (the actual value) is nil.
	// If `err` is not nil, it means there was an error during the initialization of the NATS server,
	// and the provided message "Expected no error when starting NATS server" will indicate the nature of the failure.
	assert.NoError(t, err, "Expected no error when starting NATS server")

	// Schedule the ShutdownNatsServer function to be called when the surrounding function (test) returns.
	// This ensures that the NATS server will be properly shut down after the test completes,
	// preventing resource leaks and ensuring a clean state for subsequent tests.
	defer natsServer.ShutdownNatsServer()

	// Retrieve the NATS connection instance from the test package.
	// The GetNatsConnection function is expected to return a pointer to an initialized NATS connection object.
	// This connection will be used for interacting with the NATS messaging system during the tests.
	natsConnection := natsServer.GetNatsConnection()
	// Assert that the retrieved NATS connection is not nil.
	// Assert.NotNil function checks that `natsConnection` (the actual value) is not nil.
	// If it is nil, the test will fail, and the provided message "Expected the NATS connection to be initialized, but got nil"
	// will indicate that the connection was not properly initialized before this point in the test.
	assert.NotNil(t, natsConnection, "Expected the NATS connection to be initialized, but got nil")

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
		subscription, errSub := subscriber.Subscriber(context.Background(), subject, queue)
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
}
