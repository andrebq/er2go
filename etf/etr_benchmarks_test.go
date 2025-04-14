package etf_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/andrebq/er2go/etf"
)

func prepareBenchmarkData(b *testing.B) any {
	inputData, err := os.ReadFile(filepath.Join("testdata", "input.bin"))
	if err != nil {
		b.Fatalf("failed to read input file: %v", err)
	}

	decoder := etf.NewDecoder(bytes.NewReader(inputData))
	decodedValue, err := decoder.Decode()
	if err != nil {
		b.Fatalf("failed to decode input data: %v", err)
	}

	return decodedValue
}

func BenchmarkEncoding(b *testing.B) {
	b.StopTimer()
	decodedValue := prepareBenchmarkData(b)
	b.ResetTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		err := etf.NewEncoder(io.Discard).Encode(decodedValue)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodingParallel(b *testing.B) {
	b.StopTimer()
	decodedValue := prepareBenchmarkData(b)
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err := etf.NewEncoder(io.Discard).Encode(decodedValue)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkDecoding(b *testing.B) {
	b.StopTimer()
	buf, err := os.ReadFile(filepath.Join("testdata", "input.bin"))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.StartTimer()
	br := bytes.NewReader(buf)
	for i := 0; i < b.N; i++ {
		// avoiding new allocations since this is a serial test
		// consistently saves few dozen nanoseconds on my machine.
		br.Seek(0, io.SeekStart)
		_, err := etf.NewDecoder(br).Decode()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodingParallel(b *testing.B) {
	b.StopTimer()
	buf, err := os.ReadFile(filepath.Join("testdata", "input.bin"))
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.StartTimer()
	b.RunParallel(func(p *testing.PB) {
		br := bytes.NewReader(buf)
		for p.Next() {
			br.Seek(0, io.SeekStart)
			_, err := etf.NewDecoder(br).Decode()
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
