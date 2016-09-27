package gosmile

import "bytes"

type Encoder struct {
	ContainsRawBinary         bool
	SharedStringValueEnabled  bool
	SharedPropertyNameEnabled bool
	Version                   int
	content                   bytes.Buffer
}

func (e *Encoder) Marshal(v interface{}) ([]byte, error) {
	e.content.Reset()
	e.encodeHeader()
	return e.content.Bytes(), nil
}

func NewEncoder() *Encoder {
	e := &Encoder{}
	e.init()
	return e
}

func (e *Encoder) init() {
	e.ContainsRawBinary = false
	e.SharedPropertyNameEnabled = true
	e.SharedStringValueEnabled = false
	e.Version = 0
}

func (e *Encoder) encodeHeader() {
	e.content.WriteString(":)\n")
	varByte := 0 //  reserved bit is reserved
	varByte = varByte | (e.Version << 4)
	if e.ContainsRawBinary {
		varByte = varByte | 0x04
	}
	if e.SharedStringValueEnabled {
		varByte = varByte | 0x02
	}
	if e.SharedPropertyNameEnabled {
		varByte = varByte | 0x01
	}
	e.content.WriteByte(byte(varByte))

}

func zigzagEncodeInt(n int) int {
	return (n << 1) ^ (n >> 31)
}

func zigzagEncodeLong(n int64) int64 {
	return (n << 1) ^ (n >> 63)
}
