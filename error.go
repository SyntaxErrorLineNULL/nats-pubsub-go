package nats_pubsub_go

import "errors"

var (
	ErrInvalidArgument = errors.New("invalid argument")

	ErrCloseConnection = errors.New("connection is close")
)
