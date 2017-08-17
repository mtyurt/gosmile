package gosmile

import (
	"errors"
	"reflect"
)

type DecoderConf struct {
	ContainsRawBinary         bool
	SharedStringValueEnabled  bool
	SharedPropertyNameEnabled bool
	Version                   int
	IncludeHeader             bool
	RawBinary                 bool
}
type decoder struct {
	conf  *DecoderConf
	data  []byte
	index int
}

func NewDecoderConf() *DecoderConf {
	c := &DecoderConf{}
	c.ContainsRawBinary = false
	c.SharedPropertyNameEnabled = true
	c.SharedStringValueEnabled = false
	c.IncludeHeader = true
	c.Version = 0
	return c
}

func (d *decoder) init() error {
	if len(d.data) < 4 {
		return errors.New("not enough bytes for header ")
	}

	if d.data[0] != byte(':') || d.data[1] != byte(')') || d.data[2] != byte('\n') {
		return errors.New("invalid header: " + string(d.data[0:3]))
	}
	conf := &DecoderConf{}
	varByte := d.data[3]
	conf.Version = int((varByte & 0xf0) >> 4)
	conf.RawBinary = (varByte & 0x04) == 0x04
	conf.SharedStringValueEnabled = (varByte & 0x02) == 0x02
	conf.SharedPropertyNameEnabled = (varByte & 0x01) == 0x01
	conf.IncludeHeader = true
	d.conf = conf
	d.index = 4
	return nil
}

func Unmarshal(content []byte, v interface{}) error {
	conf := &DecoderConf{}
	d := &decoder{data: content, conf: conf}
	err := d.init()
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("input value should be a pointer")
	}
	return nil
}
func Parse(contant []byte) (*interface{}, error) {
	conf := &DecoderConf{}
	d := &decoder{data: content, conf: conf}
	err := d.init()
	if err != nil {
		return nil, err
	}
	return decode(d)
}
func decode(d *decoder) (*interface{}, error) {
	token := d.data[d.index] & 0xff
	d.index++
	switch token >> 5 {
	case 0: //shared string
		break
	case 1:
		{
			typeBits := ch & 0x1F
			if typeBits < 4 {
				switch typeBits {
				case 0x00:
					return decodeString(d, token)
				case 0x01: //null value
					return nil, nil
				case 0x02:
					return false, nil
				case 0x03:
					return true, nil
				}
			}
			if typeBits == 4 {
				return decodeInt(d)
			}
		}
	}
	return nil
}
func decodeString(d *decoder, token int) (string, error) {
	tokentype := (token >> 5)
	if tokentype == 2 || tokentype == 3 { // tiny & short ASCII
		return decodeShortAsciiString(1 + (token & 0x3F)), nil
	}
	if tokentype == 4 || tokentype == 5 { // tiny & short Unicode
		// short unicode; note, lengths 2 - 65 (off-by-one compared to ASCII)
		return decodeShortUnicodeString(2 + (token & 0x3F)), nil
	}
}
func decodeShortAsciiString(d *decoder, length int) string {
	val := string(decoder.data[d.index : d.index+length])
	d.index += length
	return val
}
func decodeShortUnicodeString(d *decoder, length int) string {
	return ""
}

type token struct {
	tokentype int
	value     interface{}
}

func zigzagDecodeInt(n int) int {
	return int(uint32(n)>>1) ^ (-(n & 1))
}

func zigzagDecodeLong(n int64) int64 {
	return int64(uint64(n)>>1) ^ (-(n & 1))
}
