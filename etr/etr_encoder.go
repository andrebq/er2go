package etr

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
)

type (
	etrWriter struct {
		*bufio.Writer
	}

	Encoder interface {
		Encode(value any) error
	}
)

func NewEncoder(output io.Writer) Encoder {
	return &etrWriter{bufio.NewWriter(output)}
}

func (e *etrWriter) Encode(value any) error {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("failed to encode value: %v", err))
		}
	}()
	e.mustWriteByte(etrVersion)
	e.mustEncodeEtr(value)
	return e.Flush()
}

func (e *etrWriter) mustWriteNil() {
	e.mustWriteByte(106)
}

func (e *etrWriter) mustWriteByte(b byte) {
	if err := e.WriteByte(b); err != nil {
		panic(err)
	}
}

func (e *etrWriter) mustWriteN(data []byte) {
	if _, err := e.Write(data); err != nil {
		panic(err)
	}
}

func (e *etrWriter) mustWriteUint16(value uint16) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, value)
	e.mustWriteN(buf)
}

func (e *etrWriter) mustWriteUint32(value uint32) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, value)
	e.mustWriteN(buf)
}

func (e *etrWriter) mustWriteUint64(value uint64) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, value)
	e.mustWriteN(buf)
}

func (e *etrWriter) mustEncodeEtr(value any) {
	switch v := value.(type) {
	case int:
		if v >= 0 && v <= 255 {
			e.mustWriteByte(97) // SMALL_INTEGER_EXT
			e.mustWriteByte(byte(v))
		} else if v >= -2147483648 && v <= 2147483647 {
			e.mustWriteByte(98) // INTEGER_EXT
			e.mustWriteUint32(uint32(int32(v)))
		} else {
			e.mustEncodeBigInt(big.NewInt(int64(v)))
		}
	case Atom:
		if !v.validutf8() {
			panic(fmt.Errorf("invalid UTF-8 string: %q", v.v))
		}
		if v.small() {
			e.mustWriteByte(119)
			e.mustWriteByte(byte(len(v.v)))
			e.mustWriteN([]byte(v.v))
		} else {
			e.mustWriteByte(118)
			e.mustWriteUint16(uint16(len(v.v)))
			e.mustWriteN([]byte(v.v))
		}
	case float64:
		bits := math.Float64bits(v)
		e.mustWriteByte(70)
		e.mustWriteUint64(bits)
	case float32:
		bits := math.Float64bits(float64(v))
		e.mustWriteByte(70)
		e.mustWriteUint64(bits)
	case Binstr:
		e.mustWriteByte(109)
		e.mustWriteUint32(uint32(len(v)))
		e.mustWriteN([]byte(v))
	case string:
		e.mustWriteByte(107)
		e.mustWriteUint16(uint16(len(v)))
		e.mustWriteN([]byte(v))
	case []byte:
		e.mustWriteByte(109)
		e.mustWriteUint32(uint32(len(v)))
		e.mustWriteN(v)
	case Tuple:
		if v.small() {
			e.mustWriteByte(104)
			e.mustWriteByte(byte(len(v.v)))
			for _, elem := range v.v {
				e.mustEncodeEtr(elem)
			}
		} else {
			e.mustWriteByte(105)
			e.mustWriteUint32(uint32(len(v.v)))
			for _, elem := range v.v {
				e.mustEncodeEtr(elem)
			}
		}
	case []any:
		e.mustWriteByte(108)
		e.mustWriteUint32(uint32(len(v)))
		for _, elem := range v {
			e.mustEncodeEtr(elem)
		}
		e.mustWriteNil()
	case map[any]any:
		e.mustWriteByte(116)
		e.mustWriteUint32(uint32(len(v)))
		for key, val := range v {
			e.mustEncodeEtr(key)
			e.mustEncodeEtr(val)
		}
	case *big.Int:
		e.mustEncodeBigInt(v)
	default:
		panic(fmt.Errorf("unsupported type: %T", v))
	}
}

func (e *etrWriter) mustEncodeBigInt(value *big.Int) {
	bytes := value.Bytes()
	if len(bytes) <= 255 {
		e.mustWriteByte(110)
		e.mustWriteByte(byte(len(bytes)))
	} else {
		e.mustWriteByte(111)
		e.mustWriteUint32(uint32(len(bytes)))
	}
	sign := byte(0)
	if value.Sign() < 0 {
		sign = 1
	}
	e.mustWriteByte(sign)
	e.mustWriteN(reverseBytes(bytes))
}

func reverseBytes(data []byte) []byte {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
	return data
}
