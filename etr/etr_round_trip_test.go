package etr_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unicode"

	"github.com/andrebq/er2go/etr"
)

func TestRoundTrip(t *testing.T) {
	files, err := os.ReadDir(filepath.Join("testdata", "round-trip"))
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		t.Run(file.Name(), func(t *testing.T) {
			filePath := filepath.Join("testdata", "round-trip", file.Name())
			originalData, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("failed to read file %q: %v", filePath, err)
				return
			}

			// Decode the original data
			decoder := etr.NewDecoder(bytes.NewReader(originalData))
			decodedValue, err := decoder.Decode()
			if err != nil {
				t.Errorf("failed to decode file %q: %v", filePath, err)
				return
			}

			// Re-encode the decoded value
			var buffer bytes.Buffer
			encoder := etr.NewEncoder(&buffer)
			if err := encoder.Encode(decodedValue); err != nil {
				t.Errorf("failed to encode value from file %q: %v", filePath, err)
				return
			}

			// Compare the re-encoded data with the original data
			reEncodedData := buffer.Bytes()
			if !bytes.Equal(originalData, reEncodedData) {
				t.Logf("original data: \n%v\n---\n%v", toHexTable(originalData), toHexTable(reEncodedData))
				t.Errorf("round-trip mismatch for file %q", filePath)
			}
		})

	}
}

func toHexTable(data []byte) string {
	const columns = 16
	var result string
	for i, b := range data {
		if i > 0 && i%columns == 0 {
			result += "\n"
		}
		r := rune(b)
		switch {
		case unicode.IsDigit(r), unicode.IsLetter(r),
			unicode.IsPunct(r):
			result += fmt.Sprintf("%03d%c ", b, b)
		default:
			result += fmt.Sprintf("%04d ", b)
		}
	}
	return result
}
