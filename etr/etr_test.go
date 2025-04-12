package etr_test

import (
	"gec/etr"
	"os"
	"path/filepath"
	"testing"
)

func TestEtr(t *testing.T) {
	buf, err := os.ReadFile(filepath.Join("testdata", "input.bin"))
	if err != nil {
		t.Fatal(err)
	}
	term, err := etr.Decode(buf)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Term: %#v", term)
	t.Fail()
}
