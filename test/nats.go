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
func InitNats() (err error) {
	natsInitOnce.Do(func() {
		// Start a mock NATS server
		natsServer = StartNats()

		// Connect to the mock NATS server on localhost using DefaultNatsPort
		sURL := fmt.Sprintf("nats://127.0.0.1:%d", DefaultNatsPort)
		natsConnection, err = nats.Connect(sURL)
	})

	return err
}

// GetNatsConnection returns the established NATS connection, or nil if not initialized.
func GetNatsConnection() *nats.Conn {
	return natsConnection
}

// ShutdownNatsServer gracefully shuts down the mock NATS server (if running)
// using a sync.Once to ensure it happens only once.
func ShutdownNatsServer() {
	natsShutdownOnce.Do(func() {
		natsServer.Shutdown()
	})
}

// StartNats initializes a NATS service instance for testing purposes.
// It creates and returns a NATS server instance with default test options,
// including the port specified by DefaultNatsPort constant.
func StartNats() *server.Server {
	// Set default test options for NATS server
	opts := natsserver.DefaultTestOptions
	opts.Port = DefaultNatsPort

	// Run the NATS server with the provided options and return the server instance
	return RunServerWithOptions(&opts)
}

// RunServerWithOptions runs a NATS server with the provided options.
// It starts a NATS server using the given configuration options and returns the server instance.
func RunServerWithOptions(opts *server.Options) *server.Server {
	// Run the NATS server with the provided options and return the server instance
	return natsserver.RunServer(opts)
}
