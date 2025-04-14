package etf

import (
	"math"
	"unicode/utf8"
)

func (a Atom) String() string { return a.v }

func (a Atom) validutf8() bool {
	return len(a.v) < math.MaxUint16 && utf8.ValidString(a.v)
}

func (a Atom) small() bool {
	return len(a.v) < 256
}
