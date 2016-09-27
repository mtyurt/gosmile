package gosmile

import "errors"

type decoder struct {
	version            int
	sharedStringValue  bool
	sharedPropertyName bool
	rawBinary          bool
	data               []byte
}

func Unmarshal(content []byte, v interface{}) error {
	d := &decoder{data: content}
	err := d.init()
	if err != nil {
		return err
	}
	return nil
}

func (d *decoder) init() error {
	if len(d.data) < 4 {
		return errors.New("not enough bytes for header ")
	}

	if d.data[0] != byte(':') || d.data[1] != byte(')') || d.data[2] != byte('\n') {
		return errors.New("invalid header: " + string(d.data[0:3]))
	}
	varByte := d.data[3]
	d.version = int((varByte & 0xf0) >> 4)
	d.rawBinary = (varByte & 0x04) == 0x04
	d.sharedStringValue = (varByte & 0x02) == 0x02
	d.sharedPropertyName = (varByte & 0x01) == 0x01
	return nil
}

func zigzagDecodeInt(n int) int {
	return int(uint32(n)>>1) ^ (-(n & 1))
}

func zigzagDecodeLong(n int64) int64 {
	return int64(uint64(n)>>1) ^ (-(n & 1))
}
