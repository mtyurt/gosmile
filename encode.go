package gosmile

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

// EncodeConf: Configuration struct for encoding
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
	rv := reflect.ValueOf(v)
	err := marshal(e, rv)

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

func marshal(e *EncodeConf, rv reflect.Value) error {
	fmt.Println("encode value", rv.Kind())

	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encodeInt(e, rv)
	case reflect.Float32:
		return encodeFloat32(e, rv)
	case reflect.Float64:
		return encodeFloat64(e, rv)
	case reflect.Bool:
		return encodeBool(e, rv)
	case reflect.String:
		return encodeString(e, rv)
	case reflect.Slice:
		return encodeSlice(e, rv)
	case reflect.Array:
		return encodeSlice(e, rv)
	case reflect.Map:
		return encodeMap(e, rv)
	case reflect.Struct:
		return encodeStruct(e, rv)
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
func encodeStruct(e *EncodeConf, rv reflect.Value) error {
	e.content.WriteByte(token_literal_start_object)
	for i := 0; i < rv.NumField(); i++ {
		value := rv.Field(i)
		name := rv.Type().Field(i).Name
		if err := writeFieldName(e, name); err != nil {
			return err
		}
		if err := marshal(e, value); err != nil {
			return err
		}

	}
	e.content.WriteByte(token_literal_end_object)
	return nil
}
func encodeMap(e *EncodeConf, rv reflect.Value) error {
	e.content.WriteByte(token_literal_start_object)
	keys := rv.MapKeys()
	for _, key := range keys {
		if key.Kind() != reflect.String {
			return errors.New("Map keys have to be string!")
		}
		if err := writeFieldName(e, key.String()); err != nil {
			return err
		}

		if err := marshal(e, rv.MapIndex(key)); err != nil {
			return err
		}

	}
	e.content.WriteByte(token_literal_end_object)
	return nil
}

func writeFieldName(e *EncodeConf, field string) error {

	flen := len(field)
	if flen == 0 {
		e.content.WriteByte(token_key_empty_string)
		return nil
	}
	//TODO max shared name length
	//TODO shared name

	encoded := []byte(field)
	byteLen := len(encoded)
	if byteLen == flen {
		if byteLen <= max_short_name_ascii_bytes {
			e.content.WriteByte(byte(token_prefix_key_ascii - 1 + byteLen))
			e.content.Write(encoded)
		} else {
			e.content.WriteByte(byte(token_key_long_string))
			e.content.Write(encoded)
			e.content.WriteByte(byte(byte_marker_end_of_string))
		}
	} else {
		if byteLen <= max_short_name_unicode_bytes {
			e.content.WriteByte(byte(token_prefix_key_unicode - 2 + byteLen))
			e.content.Write(encoded)
		} else {
			e.content.WriteByte(byte(token_key_long_string))
			e.content.Write(encoded)
			e.content.WriteByte(byte(byte_marker_end_of_string))
		}
	}
	//TODO add seen name
	return nil
}

func encodeSlice(e *EncodeConf, rv reflect.Value) error {
	e.content.WriteByte(token_literal_start_array)

	for i := 0; i < rv.Len(); i++ {
		iv := rv.Index(i)
		if err := marshal(e, iv); err != nil {
			return err
		}

	}

	e.content.WriteByte(token_literal_end_array)
	return nil
}

func encodeString(e *EncodeConf, rv reflect.Value) error {
	val := rv.String()
	vlen := len(val)
	if vlen == 0 {
		e.content.WriteByte(token_literal_empty_string)
		return nil
	}
	//TODO max shared string length
	//TODO shared string

	encoded := []byte(val)
	byteLen := len(encoded)
	if byteLen <= max_short_value_string_bytes {
		//TODO add seen string if necessary

		if byteLen == vlen {
			e.content.WriteByte(byte(token_prefix_tiny_ascii - 1 + byteLen))
		} else {
			e.content.WriteByte(byte(token_prefix_tiny_unicode - 2 + byteLen))
		}
		e.content.Write(encoded)
	} else {
		token := token_misc_long_text_unicode
		if byteLen == vlen {
			token = token_byte_long_string_ascii
		}
		e.content.WriteByte(byte(token))
		e.content.Write(encoded)
		e.content.WriteByte(byte(byte_marker_end_of_string))
	}
	return nil
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

func encodeFloat64(e *EncodeConf, rv reflect.Value) error {
	n := rv.Float()
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, n)
	if err != nil {
		return err
	}
	i := binary.BigEndian.Uint64(buf.Bytes())
	e.content.WriteByte(token_byte_float_64)
	// first 5 bytes
	hi5 := i >> 35
	byte4 := byte(hi5 & 0x7F)
	hi5 = hi5 >> 7
	byte3 := byte(hi5 & 0x7F)
	hi5 = hi5 >> 7
	byte2 := byte(hi5 & 0x7F)
	hi5 = hi5 >> 7
	byte1 := byte(hi5 & 0x7F)
	hi5 = hi5 >> 7
	byte0 := byte(hi5 & 0x7F)
	e.content.Write([]byte{byte0, byte1, byte2, byte3, byte4})

	//split byte
	e.content.WriteByte(byte((i >> 28) & 0x7F))

	//last 4 bytes

	lo4 := i
	byte3 = byte(lo4 & 0x7F)
	lo4 = lo4 >> 7
	byte2 = byte(lo4 & 0x7F)
	lo4 = lo4 >> 7
	byte1 = byte(lo4 & 0x7F)
	lo4 = lo4 >> 7
	byte0 = byte(lo4 & 0x7F)

	e.content.Write([]byte{byte0, byte1, byte2, byte3})
	return nil
}

func encodeBool(e *EncodeConf, rv reflect.Value) error {
	if rv.Bool() {
		e.content.WriteByte(token_literal_true)
	} else {
		e.content.WriteByte(token_literal_false)
	}
	return nil
}

func zigzagEncodeInt(n int) int {
	return (n << 1) ^ (n >> 31)
}

func zigzagEncodeLong(n int64) int64 {
	return (n << 1) ^ (n >> 63)
}
