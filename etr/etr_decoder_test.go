package etr_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/andrebq/er2go/etr"
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

	expected := etr.NewTuple(
		etr.NewTuple(etr.NewAtom("tuple"), 42, 3.14),
		map[any]any{etr.NewBinstr([]byte("key")): []byte("value"), etr.NewAtom("atom_key"): 123},
		3.14,
		42,
		[]byte("hello"),
		[]byte{1, 2, 3, 4},
		string([]byte{1, 2, 3, 4, 5}),
	)

	if !reflect.DeepEqual(term, expected) {
		t.Fatalf("decoded term does not match expected value, expected:\n%#v\ngot:\n%#v", expected, term)
	}
}
