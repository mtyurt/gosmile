package gosmile

import (
	"bytes"
	"reflect"
)

type EncodeConf struct {
	ContainsRawBinary         bool
	SharedStringValueEnabled  bool
	SharedPropertyNameEnabled bool
	Version                   int
	IncludeHeader             bool
	content                   bytes.Buffer
}

func Marshal(e *EncodeConf, v interface{}) ([]byte, error) {
	e.content.Reset()
	if e.IncludeHeader {
		e.encodeHeader()
	}
	err := marshal(e, v)

	if err != nil {
		return nil, err
	}
	return e.content.Bytes(), nil
}

func NewEncodeConf() *EncodeConf {
	e := &EncodeConf{}
	e.init()
	return e
}

func (e *EncodeConf) init() {
	e.ContainsRawBinary = false
	e.SharedPropertyNameEnabled = true
	e.SharedStringValueEnabled = false
	e.IncludeHeader = true
	e.Version = 0
}

func marshal(e *EncodeConf, v interface{}) error {
	t := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encodeInt(e, rv)
	}
	return nil
}

func (e *EncodeConf) encodeHeader() {
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

func encodeInt(e *EncodeConf, rv reflect.Value) error {
	n := int(rv.Int())
	n = zigzagEncodeInt(n)
	if n <= 0x3F && n >= 0 {
		if n <= 0x1F {
			e.content.WriteByte(byte(token_prefix_small_int + n))
			return nil
		}
	}
	return nil
}

func zigzagEncodeInt(n int) int {
	return (n << 1) ^ (n >> 31)
}

func zigzagEncodeLong(n int64) int64 {
	return (n << 1) ^ (n >> 63)
}
