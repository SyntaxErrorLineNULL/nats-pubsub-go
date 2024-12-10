package pkg

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/segmentio/encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestEncodingDecode(t *testing.T) {
	// Test the Decode function of the Encoding type for different message scenarios.
	// This test case is designed to validate that the Decode function correctly handles various
	// cases, including valid messages, nil messages, and messages with invalid headers.
	// The tests check for correct error handling, message transformation, and matching expected output.
	cases := []struct {
		name        string
		input       *Message
		expectError bool
		expectedMsg *nats.Msg
	}{
		{name: "Valid Message", input: &Message{Subject: "test.subject", Header: Header{"key": {"value1", "value2"}}, RequestID: "12345", Container: []byte("test payload")}, expectError: false, expectedMsg: &nats.Msg{Subject: "test.subject", Header: nats.Header{"key": []string{"value1", "value2"}}, Data: json.RawMessage(`{"subject":"test.subject","header":{"key":["value1","value2"]},"request_id":"12345","container":"dGVzdCBwYXlsb2Fk"}`)}},
		{name: "Nil Message", input: nil, expectError: true, expectedMsg: nil},
		{name: "Invalid Header", input: &Message{Subject: "test.subject", Header: Header{"key": nil}, RequestID: "12345", Container: []byte("test payload")}, expectError: true, expectedMsg: nil},
		{name: "Empty Header", input: &Message{Subject: "test.subject", RequestID: "12345", Container: []byte("test payload")}, expectError: false, expectedMsg: &nats.Msg{Subject: "test.subject", Data: json.RawMessage(`{"subject":"test.subject","request_id":"12345","container":"dGVzdCBwYXlsb2Fk"}`)}},
	}

	// Declare a variable of type Encoding. This variable will hold the instance of the encoding/decoding functionality
	// that will be tested in the test cases.
	var encoder Encoding

	// Iterate over the slice of test cases, 'cases', to execute the test logic for each test case individually.
	// This loop allows us to run the same set of assertions on different inputs and expected outputs, ensuring
	// that the function works across a variety of scenarios.
	for _, tt := range cases {
		// Start a subtest for the current test case, using 'tt.name' as the name of the subtest.
		// This ensures that each test case is isolated, and we can easily identify which test case failed
		// if any errors occur. The subtest allows for independent results and better debugging.
		t.Run(tt.name, func(t *testing.T) {
			// Call the 'Decode' method of the encoder, passing the 'input' from the current test case (tt.input).
			// This is the core functionality being tested, where the input data is processed by the decoder.
			// The result of the decoding operation is stored in the 'result' variable, and any errors that occur
			// are captured in the 'err' variable.
			result, err := encoder.Decode(tt.input)

			// Check if the current test case expects an error. The 'expectError' field in the test case struct
			// indicates whether an error is expected during the 'Decode' function call. If 'expectError' is true,
			// the assertions will verify that the function behaves as expected by producing an error.
			if tt.expectError {
				// Assert that an error has occurred during the decoding process. This assertion ensures that the error
				// returned from 'Decode' is not nil when the test case expects an error. If no error occurs, the test will fail.
				assert.Error(t, err, "Expected an error but got none for test case: "+tt.name)
				// Assert that the 'result' is nil when an error is expected. This checks that the decoder does not
				// return a valid result when an error occurs. If a valid result is returned when an error was expected,
				// the test will fail.
				assert.Nil(t, result, "Expected nil result but got a valid result for test case: "+tt.name)
			} else {
				// Assert that no error occurred during the decoding process. This check ensures that the function behaves as expected
				// when an error is not anticipated. If an error is returned when it shouldn't be, the test will fail.
				assert.NoError(t, err, fmt.Sprintf("Expected no error but got one for test case: %s", tt.name))

				// Assert that the 'result' is not nil when no error is expected. This checks that a valid decoded result
				// is returned by the 'Decode' function. If the result is nil, the test will fail.
				assert.NotNil(t, result, fmt.Sprintf("Expected a valid result but got nil for test case: %s", tt.name))

				// Assert that the 'Subject' field of the result matches the expected value from the test case.
				// This checks that the 'Decode' function correctly transfers the 'Subject' from the input to the output.
				assert.Equal(t, tt.expectedMsg.Subject, result.Subject, fmt.Sprintf("Subject mismatch for test case: %s", tt.name))

				// Assert that the 'Header' field of the result matches the expected value from the test case.
				// This ensures that the 'Header' is properly decoded and matches the expected format.
				assert.Equal(t, tt.expectedMsg.Header, result.Header, fmt.Sprintf("Header mismatch for test case: %s", tt.name))

				// Create a buffer to hold the expected JSON-encoded data. This buffer will store the result of encoding the input
				// message into a JSON format for comparison against the 'Data' field of the result.
				var expectedData bytes.Buffer

				// Encode the input message into the 'expectedData' buffer. This step ensures that the expected JSON data
				// is correctly generated before comparing it with the actual 'Data' field in the result.
				err = json.NewEncoder(&expectedData).Encode(tt.input)

				// Assert that no error occurred during the encoding process. If an error is returned, it will indicate a
				// problem with the encoding step, causing the test to fail.
				assert.NoError(t, err, fmt.Sprintf("Expected no encoding error but got one for test case: %s", tt.name))

				// Assert that the 'Data' field of the result matches the expected JSON data generated above.
				// This verifies that the 'Decode' function correctly processes the input and generates the expected data.
				assert.Equal(t, expectedData.Bytes(), result.Data, fmt.Sprintf("Data mismatch for test case: %s", tt.name))
			}
		})
	}
}
