package nats_pubsub_go

import "errors"

var (
	// ErrInvalidArgument is an error returned when an invalid argument is provided.
	// This is used to indicate that a function or method has been called with
	// arguments that do not meet the required criteria or format.
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrCloseConnection is an error returned when an operation is attempted
	// on a closed connection. It signifies that the connection has been
	// terminated and cannot be used for further operations.
	ErrCloseConnection = errors.New("connection is close")

	// ErrConnectionAlreadyClosed indicates that a connection closure was attempted on an already closed connection.
	// This error helps differentiate between the connection being in a valid state versus being redundantly closed.
	// By defining this error, the code provides a specific signal to handle such redundant closure attempts gracefully.
	ErrConnectionAlreadyClosed = errors.New("connection is already closed")
)
