package pkg

import "github.com/nats-io/nats.go"

type Subscriber struct {
	// conn holds the connection to the NATS server.
	// This connection is used to publish messages.
	conn *nats.Conn

	// isClose is a flag indicating whether the Subscriber has been closed.
	// Once set to true, the Subscriber should not allow further publishing.
	isClose bool
}

func NewSubscriber(conn *nats.Conn) *Subscriber {
	return &Subscriber{conn: conn}
}
