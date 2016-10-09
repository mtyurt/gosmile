package gosmile

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

// Configuration struct for encoding
// For default values, please use NewEncodeConf() function
type EncodeConf struct {
	ContainsRawBinary         bool
	SharedStringValueEnabled  bool
	SharedPropertyNameEnabled bool
	Version                   int
	IncludeHeader             bool
	content                   bytes.Buffer
}

//Marshals given value with encoding configuration and returns byte slice.
//For encoding conf, refer to EncodeConf struct and also NewEncodeConf() function
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
	case reflect.Float32:
		return encodeFloat32(e, rv)
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
			return e.content.WriteByte(byte(token_prefix_small_int + n))
		}
		e.content.Write([]byte{token_byte_int_32, byte(0x80 + n)})
		return nil
	}

	b0 := byte(0x80 + (n & 0x3F))
	n = int(uint(n) >> 6)
	if n <= 0x7F { //13 bits are enough, 3 total bytes
		e.content.Write([]byte{token_byte_int_32, byte(n), b0})
		return nil
	}
	b1 := byte(n & 0x7F)
	n = n >> 7
	if n <= 0x7F {
		e.content.Write([]byte{token_byte_int_32, byte(n), b1, b0})
		return nil
	}
	b2 := byte(n & 0x7F)
	n = n >> 7
	if n <= 0x7F {
		e.content.Write([]byte{token_byte_int_32, byte(n), b2, b1, b0})
		return nil
	}
	b3 := byte(n & 0x7F)
	n = n >> 7
	e.content.Write([]byte{token_byte_int_32, byte(n), b3, b2, b1, b0})
	return nil
}

func encodeFloat32(e *EncodeConf, rv reflect.Value) error {

	n := float32(rv.Float())
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, n)
	if err != nil {
		return err
	}
	i := binary.BigEndian.Uint32(buf.Bytes())
	e.content.WriteByte(token_byte_float_32)
	byte4 := byte(i & 0x7F)
	i = i >> 7
	byte3 := byte(i & 0x7F)
	i = i >> 7
	byte2 := byte(i & 0x7F)
	i = i >> 7
	byte1 := byte(i & 0x7F)
	i = i >> 7
	byte0 := byte(i & 0x7F)
	e.content.Write([]byte{byte0, byte1, byte2, byte3, byte4})
	return nil
}

func zigzagEncodeInt(n int) int {
	return (n << 1) ^ (n >> 31)
}

func zigzagEncodeLong(n int64) int64 {
	return (n << 1) ^ (n >> 63)
}
