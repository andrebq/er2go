package etf_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/andrebq/er2go/etf"
)

func TestEtr(t *testing.T) {
	buf, err := os.ReadFile(filepath.Join("testdata", "input.bin"))
	if err != nil {
		t.Fatal(err)
	}
	term, err := etf.Decode(buf)
	if err != nil {
		t.Fatal(err)
	}

	expected := etf.NewTuple(
		etf.NewTuple(etf.NewAtom("tuple"), 42, 3.14),
		map[any]any{etf.NewBinstr([]byte("key")): []byte("value"), etf.NewAtom("atom_key"): 123},
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
