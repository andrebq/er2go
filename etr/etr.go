package etr

import (
	"bufio"
	"bytes"
	"fmt"
)

type (
	etrReader struct {
		*bufio.Reader
	}
)

var (
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
)

func reject(reason string) func(*etrReader) any {
	return func(r *etrReader) any {
		panic(fmt.Errorf("unsupported type: %s", reason))
	}
}

const (
	etrVersion = 131
)

func Decode(buf []byte) (value any, err error) {
	r := etrReader{bufio.NewReader(bytes.NewBuffer(buf))}
	defer func() {
		perr := recover()
		if perr != nil {
			err = fmt.Errorf("invalid stream: %v", perr)
		}
	}()
	switch b := r.mustByte(); b {
	case etrVersion:
		value = mustDecodeEtr(&r)
	default:
		return nil, fmt.Errorf("unsupported version: %d", b)
	}
	return
}

func (r *etrReader) mustByte() byte {
	b, err := r.ReadByte()
	if err != nil {
		panic(err)
	}
	return b
}

func (r *etrReader) eof() bool {
	return r.Buffered() == 0
}

func mustDecodeEtr(r *etrReader) any {
	if r.eof() {
		return nil
	}
	tag := r.mustByte()
	switch tag {
	}
}
