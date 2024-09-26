# nats-pubsub-go
nats-pubsub-go is a Go library that provides a simple and efficient interface for interacting with NATS, a lightweight, high-performance messaging system. This library simplifies the process of publishing and subscribing to messages on NATS subjects, making it easy to build scalable and distributed applications.

### Installation
```bash
go get -u github.com/SyntaxErrorLineNULL/nats-pubsub-go
```

# Publisher Package

This package provides a Publisher struct for sending messages and requests to a NATS server. It ensures safe and efficient communication, handling connection state and message validation.

### Features

* **Reliable Publishing**: Publishes messages to a NATS server with error handling and flushing for immediate delivery.
* **Request-Response**: Sends message requests with a timeout and retrieves responses.
* **Connection Management**: Tracks the connection state to prevent operations on closed connections.
* **Message Validation**: Ensures messages are valid and not nil before publishing.
* **Error Handling**: Returns specific errors for different failure scenarios, aiding in debugging.

## Usage

### Create a Publisher:
```go
// Establish a connection to the NATS server (replace with your connection details)
conn, err := nats.Connect(...)
if err != nil {
    // Handle connection error
}

// Create a Publisher instance
pub := publisher.NewPublisher(conn)
```

### Publish Messages:
```go
// Message to be published
msg := &nats.Msg{Subject: "my_subject", Data: []byte("Hello, world!")}

// Publish a single message
err = pub.Publish(msg)
if err != nil {
    // Handle publishing error (e.g., connection closed, invalid message)
}

// Publish multiple messages in a batch
messages := []*nats.Msg{msg, anotherMessage}
err = pub.Publish(messages...)
if err != nil {
    // Handle publishing error
}
```

### Send Request and Receive Response:
```go
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
    // Publish the expected response message to the reply subject specified in the incoming message.
    // The `msg.Reply` field contains the subject where the response should be sent.
    // `expectedPublishMessage` is the data payload that will be sent as the response.
    err = natsConnection.Publish(msg.Reply, expectedPublishMessage)
    if err != nil { 
        log.Println("failed publish")
        return
    }

    // Signal that the response has been published by sending a value to the result channel (resCh).
    resCh <- struct{}{}
})

// Send a request message to the NATS server using the `publisher.Request` method.
// The request message is constructed with the specified `subject` and `expectedRequestMessage` data.
// A timeout of 1 second is provided for the request, meaning the request will wait up to 1 second for a response.
response, err := publisher.Request(&nats.Msg{Subject: subject, Data: expectedRequestMessage}, 1*time.Second)
if err != nil { log.Println("failed request") }

// Use a select statement to handle multiple cases: receiving a response or timing out.
select {
// Case when a response message is received on the `resCh` channel.
// The response channel `resCh` is signaled when the expected publish message is received.
case <-resCh:
    ...
    // Exit the loop and the test since the expected response has been successfully received and validated.
    return

// Case when the timeout duration (5 seconds) is reached without receiving a response.
// The `<-time.After(5 * time.Second)` statement waits for 5 seconds and then proceeds.
case <-time.After(5 * time.Second):
    // Fail the test with a timeout error if no response was received within the 5-second window.
    // The `t.Fatal` function logs the provided message and stops the test execution.
    log.Println("Timed out waiting for message")
}
```

### Close the Publisher:
```go
// Close the connection and mark the Publisher as closed
pub.Close()

// Do not attempt to publish after closing the Publisher
```

## Error Handling

The package returns specific errors for different failure scenarios:

```go 
nats_pubsub_go.ErrCloseConnection // The Publisher is closed and cannot be used for further operations.
nats_pubsub_go.ErrInvalidArgument // An invalid argument was provided, such as empty messages or a nil request message.
```

# Subscribe Package

This package provides robust tools for subscribing to messages on a NATS server. It offers functions for both asynchronous and synchronous subscriptions, ensuring flexibility and control over message processing.

### Features

* **Multiple Subscription Modes**: Supports asynchronous and synchronous subscriptions depending on your needs.
* **Error Handling**: Handles errors during subscription creation and message retrieval, providing informative error messages.
* **Safe Unsubscription**: Guarantees that the subscription is unsubscribed from and the message channel is closed only once.
* **Clean Channel Access**: Provides methods to retrieve the message channel for asynchronous subscriptions and receive messages for synchronous subscriptions.

### Create a Subscriber:

```go
// Establish a connection to the NATS server (replace with your connection details)
conn, err := nats.Connect(...)
if err != nil {
    // Handle connection error
}

// Create a new Subscriber instance with the NATS connection
subscriber := pkg.NewSubscriber(conn)
```
