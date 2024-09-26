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

// NatsServer represents a mock NATS server used for testing or isolated environments.
// This struct encapsulates all necessary components to manage the server's lifecycle and connection.
type NatsServer struct {
	// port is the port number on which the NATS server will listen for connections.
	// This ensures that the server operates on a specified or default port.
	port int
	// natsServer represents the mock NATS server instance.
	// This field holds the server object created by the NATS library,
	// allowing for management of the server's operations.
	natsServer *server.Server
	// natsConnection represents the NATS connection instance.
	// This field stores the connection to the NATS server,
	// enabling the publication and subscription of messages.
	natsConnection *nats.Conn
	// natsInitOnce ensures that InitNats is executed only once.
	// This sync.Once field is used to prevent multiple concurrent initialization of the NATS server,
	// ensuring that the server setup occurs in a thread-safe manner.
	natsInitOnce sync.Once
	// natsShutdownOnce ensures that ShutdownNatsServer is executed only once.
	// Similar to natsInitOnce, this field prevents multiple concurrent shutdowns,
	// maintaining safe and controlled termination of the NATS server.
	natsShutdownOnce sync.Once
}

// NewNatsServer creates and returns a new instance of NatsServer.
// It accepts a port number as a parameter and ensures that the server
// is initialized with a valid port for operation.
func NewNatsServer(port int) *NatsServer {
	// Check if the provided port is 0.
	// If no port is specified (indicated by port being 0), set it to the default NATS port.
	// This ensures that the server will start on a valid port, preventing issues related to invalid or unassigned ports.
	if port == 0 {
		port = DefaultNatsPort
	}

	// Create and return a new instance of NatsServer,
	// initializing it with the specified or default port.
	return &NatsServer{port: port}
}

// InitNats initializes the NATS server and establishes a connection to it.
// This method ensures that the initialization process is only carried out once,
// preventing multiple initializations that could lead to errors or unexpected behavior.
func (n *NatsServer) InitNats() (err error) {
	// Ensure that the initialization of NATS is only performed once.
	// `natsInitOnce.Do` ensures that even if this function is called multiple times,
	// the NATS server will only be initialized a single time.
	n.natsInitOnce.Do(func() {
		// Start a mock NATS server on the specified or default port.
		// This mock server is used for testing or isolated environments
		// where a full NATS infrastructure is not necessary.
		n.natsServer = n.StartNats()

		// Create the connection URL for the NATS server using the provided port.
		// This URL is constructed to connect to a locally running NATS server at `127.0.0.1`.
		sURL := fmt.Sprintf("nats://127.0.0.1:%d", n.port)

		// Attempt to establish a connection to the NATS server using the generated URL.
		// `nats.Connect` connects to the server and returns a connection object,
		// which can be used to publish and subscribe to messages.
		natsConnection, err := nats.Connect(sURL)
		if err != nil {
			return
		}

		// Assign the established NATS connection to the instance variable,
		// allowing other methods in the NatsServer to utilize this connection.
		n.natsConnection = natsConnection
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
func (n *NatsServer) GetNatsConnection() *nats.Conn {
	// Return the current NATS connection.
	// This function simply provides access to the globally initialized `natsConnection` instance.
	// It allows other parts of the application to retrieve the established NATS connection
	// without directly interacting with the initialization logic.
	return n.natsConnection
}

// ShutdownNatsServer gracefully shuts down the running NATS server.
// It ensures the server is only shut down once using a `sync.Once` mechanism
// to prevent multiple shutdown attempts, which could lead to errors or resource
// management issues. This function is typically used when the application or
// system is shutting down and no further NATS operations are required.
func (n *NatsServer) ShutdownNatsServer() {
	// Use `sync.Once` to ensure that the shutdown procedure is only performed once.
	// This avoids issues that might arise from attempting to shut down the server multiple times.
	n.natsShutdownOnce.Do(func() {
		// Trigger the NATS server's shutdown process.
		// This method stops the server and releases any resources associated with it.
		n.natsServer.Shutdown()
	})
}

// StartNats initializes and starts a NATS server using predefined test options.
// This method is part of the NatsServer struct, which manages the NATS server instance.
// It sets up the server with default configurations suitable for testing environments.
func (n *NatsServer) StartNats() *server.Server {
	// Set default options for a test NATS server.
	// These options are predefined in the NATS library and are suitable for testing purposes.
	opts := natsserver.DefaultTestOptions

	// Assign the desired port to the options.
	// This port will be used by the NATS server for incoming connections.
	opts.Port = n.port

	// Start the NATS server using the specified options and return the server instance.
	// The `RunServerWithOptions` function launches the NATS server with the provided configuration.
	// The returned server instance can be used to interact with and control the server.
	return n.runServerWithOptions(&opts)
}

// runServerWithOptions starts a NATS server using the provided options.
// This method belongs to the NatsServer struct, which is responsible for managing the NATS server instance.
// The opts parameter contains configuration options for the NATS server, such as the port, authorization settings, etc.
func (n *NatsServer) runServerWithOptions(opts *server.Options) *server.Server {
	// Start the NATS server using the specified options.
	// The natsserver.RunServer function initializes a new NATS server instance with the given options.
	// The returned value is a pointer to the server instance that is now running.
	n.natsServer = natsserver.RunServer(opts)

	// Return the reference to the running NATS server instance.
	// This allows other parts of the application to interact with or control the server as needed.
	return n.natsServer
}
