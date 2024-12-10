package pkg

import "github.com/nats-io/nats.go"

type Encoder interface {
	Encode(*nats.Msg) (*Message, error)
}
