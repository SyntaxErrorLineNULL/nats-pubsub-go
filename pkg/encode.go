package pkg

import (
	"bytes"
	"github.com/nats-io/nats.go"
	"github.com/segmentio/encoding/json"
)

// Decoder defines an interface for decoding a custom Message into a NATS message.
// It abstracts the process of converting application-level message representations
// into the NATS message format to facilitate communication within a NATS-based system.
type Decoder interface {
	// Decode is responsible for transforming a custom Message into a NATS message.
	// The resulting NATS message contains serialized data and metadata, formatted
	// for compatibility with the NATS messaging system.
	// Returns the NATS message if successful or an error if the transformation fails.
	Decode(msg *Message) (*nats.Msg, error)
}

// Encoding implements the Decoder interface, providing functionality to convert
// a custom Message structure into a format compatible with the NATS messaging system.
// It uses JSON serialization to encode the Message and populate the fields of a NATS message.
type Encoding struct{}

// Decode converts a custom Message into a NATS message.
// This method serializes the Message's content, including its payload and metadata,
// into JSON format and assigns the serialized data to the NATS message fields.
func (Encoding) Decode(msg *Message) (*nats.Msg, error) {
	// Initializes a buffer to temporarily hold the JSON-encoded data.
	var buffer bytes.Buffer

	// Creates a JSON encoder that writes to the buffer.
	// The encoder is responsible for serializing the Message into JSON format.
	encoder := json.NewEncoder(&buffer)

	// Attempts to encode the custom Message into JSON and store it in the buffer.
	// If an error occurs during encoding, the function immediately returns the error.
	if err := encoder.Encode(msg); err != nil {
		// Returns an error indicating that the encoding process failed.
		return nil, err
	}

	// Constructs a NATS message using the encoded data and message metadata.
	// The Subject, Data, and Header fields of the NATS message are populated
	// using the corresponding fields from the custom Message.
	return &nats.Msg{Subject: msg.Subject, Data: buffer.Bytes(), Header: nats.Header(msg.Header)}, nil
}
