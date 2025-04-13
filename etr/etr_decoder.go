package etr

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
)

type (
	etrReader struct {
		*bufio.Reader
	}

	Decoder interface {
		Decode() (any, error)
	}

	Atom struct {
		v string
	}

	Tuple struct {
		v []any
	}
)

func NewAtom(v string) Atom { return Atom{v} }
func NewTuple(items ...any) Tuple {
	return Tuple{v: items}
}

func init() {
	tagDecode = map[byte]func(*etrReader) any{
		82:  reject("atom ext is not supported"),
		102: reject("port ext is not supported"),
		89:  reject("new port ext is not supported"),
		120: reject("v4 port ext is not supported"),
		103: reject("pid ext is not supported"),
		88:  reject("new pid ext is not supported"),
		99:  reject("float ext is not supported"),
		101: reject("reference ext is not supported"),
		114: reject("new reference ext is not supported"),
		90:  reject("newer reference ext is not supported"),
		117: reject("fun ext is not supported"),
		112: reject("new fun ext is not supported"),
		113: reject("export fun ext is not supported"),
		77:  reject("bit binary ext is not supported"),
		100: reject("atom (deprecated) ext is not supported"),
		115: reject("small atom ext is not supported"),
		121: reject("local ext is not supported"),

		97:  decodeSmallInt,
		98:  decodeLargeInt,
		104: decodeSmallTuple,
		105: decodeLargeTuple,
		116: decodeMap,
		106: decodeNil,
		107: decodeString,
		108: decodeList,
		109: decodeBinary,
		110: decodeSmallBig,
		111: decodeLargeBig,
		70:  decodeFloat,
		118: decodeAtom,
		119: decodeSmallAtom,
	}
}

var (
	tagDecode = map[byte]func(*etrReader) any{}
)

func reject(reason string) func(*etrReader) any {
	return func(r *etrReader) any {
		panic(fmt.Errorf("unsupported type: %s", reason))
	}
}

const (
	etrVersion = 131
)

func NewDecoder(input io.Reader) Decoder {
	return &etrReader{bufio.NewReader(input)}
}

func (e *etrReader) Decode() (value any, err error) {
	defer func() {
		perr := recover()
		if perr != nil {
			err = fmt.Errorf("invalid stream: %v", perr)
		}
	}()
	switch b := e.mustByte(); b {
	case etrVersion:
		value = mustDecodeEtr(e)
	default:
		return nil, fmt.Errorf("unsupported version: %d", b)
	}
	return
}

func Decode(buf []byte) (value any, err error) {
	return NewDecoder(bytes.NewReader(buf)).Decode()
}

func (r *etrReader) mustByte() byte {
	b, err := r.ReadByte()
	if err != nil {
		panic(err)
	}
	return b
}

func (r *etrReader) mustReadN(n int) []byte {
	buf := make([]byte, n)
	if _, err := r.Read(buf); err != nil {
		panic(err)
	}
	return buf
}

func (r *etrReader) mustReadUint16() uint16 {
	buf := r.mustReadN(2)
	return binary.BigEndian.Uint16(buf)
}

func (r *etrReader) mustReadUint32() uint32 {
	buf := r.mustReadN(4)
	return binary.BigEndian.Uint32(buf)
}

func (r *etrReader) mustReadUint64() uint64 {
	buf := r.mustReadN(8)
	return binary.BigEndian.Uint64(buf)
}

func (r *etrReader) eof() bool {
	return r.Buffered() == 0
}

func mustDecodeEtr(r *etrReader) any {
	if r.eof() {
		return nil
	}
	tag := r.mustByte()
	if decode, ok := tagDecode[tag]; ok {
		return decode(r)
	}
	panic(fmt.Errorf("unknown tag: %d", tag))
}

func decodeSmallInt(r *etrReader) any {
	return int(r.mustByte())
}

func decodeLargeInt(r *etrReader) any {
	val := int32(r.mustReadUint32())
	return int(val)
}

func decodeSmallTuple(r *etrReader) any {
	arity := int(r.mustByte())
	return decodeTuple(r, arity)
}

func decodeLargeTuple(r *etrReader) any {
	arity := int(r.mustReadUint32())
	return decodeTuple(r, arity)
}

func decodeTuple(r *etrReader, arity int) Tuple {
	elements := make([]any, arity)
	for i := 0; i < arity; i++ {
		elements[i] = mustDecodeEtr(r)
	}
	return Tuple{v: elements}
}

func decodeMap(r *etrReader) any {
	arity := int(r.mustReadUint32())
	m := make(map[any]any, arity)
	for i := 0; i < arity; i++ {
		key := mustDecodeEtr(r)
		value := mustDecodeEtr(r)
		switch key := key.(type) {
		case string:
			m[key] = value
		case []uint8:
			m[NewBinstr(key)] = value
		case Atom:
			m[key] = value
		default:
			panic("unsupported map key type: " + fmt.Sprintf("%T", key))
		}
	}
	return m
}

func decodeNil(r *etrReader) any {
	return []any{}
}

func decodeString(r *etrReader) any {
	length := int(r.mustReadUint16())
	data := r.mustReadN(length)

	return string(data)
}

func decodeList(r *etrReader) any {
	length := int(r.mustReadUint32())
	elements := make([]any, length)

	for i := 0; i < length; i++ {
		elements[i] = mustDecodeEtr(r)
	}

	tail := mustDecodeEtr(r)

	if _, ok := tail.([]any); !ok || len(tail.([]any)) > 0 {
		panic("imporper list are not supported")
	}

	return elements
}

func decodeBinary(r *etrReader) any {
	length := int(r.mustReadUint32())
	return r.mustReadN(length)
}

func decodeSmallBig(r *etrReader) any {
	n := int(r.mustByte())
	return decodeBigInt(r, n)
}

func decodeLargeBig(r *etrReader) any {
	n := int(r.mustReadUint32())
	return decodeBigInt(r, n)
}

func decodeBigInt(r *etrReader, n int) any {
	sign := r.mustByte()

	data := r.mustReadN(n)

	bigint := new(big.Int)
	for i := n - 1; i >= 0; i-- {
		bigint.Lsh(bigint, 8)
		bigint.Or(bigint, big.NewInt(int64(data[i])))
	}

	if sign > 0 {
		bigint.Neg(bigint)
	}

	return bigint
}

func decodeFloat(r *etrReader) any {
	bits := r.mustReadUint64()
	return math.Float64frombits(bits)
}

func decodeAtom(r *etrReader) any {
	length := int(r.mustReadUint16())
	data := r.mustReadN(length)
	return Atom{string(data)}
}

func decodeSmallAtom(r *etrReader) any {
	length := int(r.mustByte())
	data := r.mustReadN(length)
	return Atom{string(data)}
}
