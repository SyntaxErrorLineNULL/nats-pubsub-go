package nats_pubsub_go

import "github.com/nats-io/nats.go"

type Publisher struct {
	conn *nats.Conn
}

func NewPublisher(conn *nats.Conn) *Publisher {
	return &Publisher{conn: conn}
}
