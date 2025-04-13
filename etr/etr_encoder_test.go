package etr

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	// Read the input binary file
	inputData, err := os.ReadFile(filepath.Join("testdata", "input.bin"))
	if err != nil {
		t.Fatalf("failed to read input file: %v", err)
	}

	// Decode the input data
	decoder := NewDecoder(bytes.NewReader(inputData))
	decodedValue, err := decoder.Decode()
	if err != nil {
		t.Fatalf("failed to decode input data: %v", err)
	}

	// Re-encode the decoded value
	var outputBuffer bytes.Buffer
	encoder := NewEncoder(&outputBuffer)
	if err := encoder.Encode(decodedValue); err != nil {
		t.Fatalf("failed to encode value: %v", err)
	}

	decoder = NewDecoder(bytes.NewBuffer(outputBuffer.Bytes()))
	roundTripValue, err := decoder.Decode()
	if err != nil {
		t.Fatalf("failed to decode input data: %v", err)
	}

	if !reflect.DeepEqual(roundTripValue, decodedValue) {
		t.Fatalf("round trip value does not match original, expected:\n%#v\ngot:\n%#v", decodedValue, roundTripValue)
	}
}
