package etf_test

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/andrebq/er2go/etf"
)

func TestEncodeDecodeRoundTrip(t *testing.T) {
	inputData, err := os.ReadFile(filepath.Join("testdata", "input.bin"))
	if err != nil {
		t.Fatalf("failed to read input file: %v", err)
	}

	decoder := etf.NewDecoder(bytes.NewReader(inputData))
	decodedValue, err := decoder.Decode()
	if err != nil {
		t.Fatalf("failed to decode input data: %v", err)
	}

	var outputBuffer bytes.Buffer
	encoder := etf.NewEncoder(&outputBuffer)
	if err := encoder.Encode(decodedValue); err != nil {
		t.Fatalf("failed to encode value: %v", err)
	}

	decoder = etf.NewDecoder(bytes.NewBuffer(outputBuffer.Bytes()))
	roundTripValue, err := decoder.Decode()
	if err != nil {
		t.Fatalf("failed to decode input data: %v", err)
	}

	if !reflect.DeepEqual(roundTripValue, decodedValue) {
		t.Fatalf("round trip value does not match original, expected:\n%#v\ngot:\n%#v", decodedValue, roundTripValue)
	}
}
