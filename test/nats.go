package test

import (
	"fmt"
	"sync"

	"github.com/nats-io/gnatsd/server"
	natsserver "github.com/nats-io/nats-server/test"
	"github.com/nats-io/nats.go"
)

// DefaultNatsPort specifies the port number used for testing NATS server.
const DefaultNatsPort = 8369

var (
	// natsServer represents the mock NATS server instance.
	natsServer *server.Server
	// natsConnection represents the NATS connection instance.
	natsConnection *nats.Conn
	// natsInitOnce ensures that InitNats is executed only once.
	natsInitOnce sync.Once
	// natsShutdownOnce ensures that ShutdownNatsServer is executed only once.
	natsShutdownOnce sync.Once
)

// InitNats initializes the NATS connection by starting the mock server (if not already running)
// and then connecting to it. It ensures this initialization happens only once using a sync.Once.
func InitNats(port int) (err error) {
	// Ensure that the initialization of NATS is only performed once.
	// `natsInitOnce.Do` ensures that even if this function is called multiple times,
	// the NATS server will only be initialized a single time.
	natsInitOnce.Do(func() {
		// Check if the provided port is 0.
		// If no port is specified (port is 0), set it to the default NATS port.
		// This guarantees that the server will start on a valid port.
		if port == 0 {
			port = DefaultNatsPort
		}

		// Start a mock NATS server on the specified or default port.
		// This mock server is used for testing or isolated environments
		// where a full NATS infrastructure is not necessary.
		natsServer = StartNats(port)

		// Create the connection URL for the NATS server using the provided port.
		// This URL is constructed to connect to a locally running NATS server at `127.0.0.1`.
		sURL := fmt.Sprintf("nats://127.0.0.1:%d", port)

		// Attempt to establish a connection to the NATS server using the generated URL.
		// `nats.Connect` connects to the server and returns a connection object,
		// which can be used to publish and subscribe to messages.
		natsConnection, err = nats.Connect(sURL)
	})

	// Return any error that occurred during the connection attempt.
	// If the connection was successful, `err` will be nil, signaling that
	// the NATS server was initialized correctly.
	return err
}

// GetNatsConnection returns the globally initialized NATS connection instance.
// This function allows other components to access the NATS connection without
// needing to manage the connection setup or initialization. It ensures that
// all parts of the system use a single, shared NATS connection.
func GetNatsConnection() *nats.Conn {
	// Return the current NATS connection.
	// This function simply provides access to the globally initialized `natsConnection` instance.
	// It allows other parts of the application to retrieve the established NATS connection
	// without directly interacting with the initialization logic.
	return natsConnection
}

// ShutdownNatsServer gracefully shuts down the running NATS server.
// It ensures the server is only shut down once using a `sync.Once` mechanism
// to prevent multiple shutdown attempts, which could lead to errors or resource
// management issues. This function is typically used when the application or
// system is shutting down and no further NATS operations are required.
func ShutdownNatsServer() {
	// Use `sync.Once` to ensure that the shutdown procedure is only performed once.
	// This avoids issues that might arise from attempting to shut down the server multiple times.
	natsShutdownOnce.Do(func() {
		// Trigger the NATS server's shutdown process.
		// This method stops the server and releases any resources associated with it.
		natsServer.Shutdown()
	})
}

// StartNats initializes a NATS service instance for testing purposes.
// It creates and returns a NATS server instance with default test options,
// including the port specified by DefaultNatsPort constant.
func StartNats(port int) *server.Server {
	// Set default options for a test NATS server.
	// These options are predefined in the NATS library and are suitable for testing purposes.
	opts := natsserver.DefaultTestOptions

	// Assign the specified or default port to the NATS server options.
	// This configures the server to listen on the given port.
	opts.Port = port

	// Start the NATS server using the specified options and return the server instance.
	// The `RunServerWithOptions` function launches the NATS server with the provided configuration.
	// The returned server instance can be used to interact with and control the server.
	return RunServerWithOptions(&opts)
}

// RunServerWithOptions runs a NATS server with the provided options.
// It starts a NATS server using the given configuration options and returns the server instance.
func RunServerWithOptions(opts *server.Options) *server.Server {
	// Run the NATS server with the provided options and return the server instance
	return natsserver.RunServer(opts)
}
