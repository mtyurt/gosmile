package gosmile

import (
	"math"
	"testing"
)

func TestZigZagInt(t *testing.T) {
	expect(0, zigzagEncodeInt(0), t, "zigzag1")
	expect(1, zigzagEncodeInt(-1), t, "zigzag1")
	expect(2, zigzagEncodeInt(1), t, "zigzag2")
	expect(int(0xffffffff), zigzagEncodeInt(math.MinInt32), t, "zigzagintmin")
	expect(int(0xfffffffe), zigzagEncodeInt(math.MaxInt32), t, "zigzagintmax")
}

func TestZigZagLong(t *testing.T) {
	expect(int64(0), zigzagEncodeLong(0), t, "zigzag0")
	expect(int64(-2), zigzagEncodeLong(int64(math.MaxInt64)), t, "zigzag-longmax")
	expect(int64(-1), zigzagEncodeLong(math.MinInt64), t, "zigzag-longmin")
}

func TestEncodeHeader(t *testing.T) {
	e := NewEncoder()
	e.Version = 3
	content, err := e.Marshal(1)
	if err != nil || len(content) < 4 {
		t.Fatal("err here", err, "content size:", len(content))
	}
	expect(byte(':'), content[0], t, "testencodeheader-1")
	expect(byte(')'), content[1], t, "testencodeheader-2")
	expect(byte('\n'), content[2], t, "testencodeheader-3")
	varByte := content[3]
	expect(byte(3)<<4, varByte&0xf0, t, "testencodeheader-version")
	expect(byte(0x00), varByte&0x04, t, "testencodeheader-rawbinary")
	expect(byte(0x00), varByte&0x02, t, "testencodeheader-sharedstringvalue")
	expect(byte(0x01), varByte&0x01, t, "testencodeheader-sharedpropname")

	e = NewEncoder()
	e.SharedStringValueEnabled = true
	e.SharedPropertyNameEnabled = false

	content, err = e.Marshal(1)
	if err != nil || len(content) < 4 {
		t.Fatal("err here", err, "content size:", len(content))
	}
	varByte = content[3]
	expect(byte(1)<<1, varByte&0x02, t, "testencodeheader-sharedstringvalue2")
	expect(byte(0x00), varByte&0x01, t, "testencodeheader-sharedpropname2")

}

func expect(expected interface{}, got interface{}, t *testing.T, test string) {
	if got != expected {
		t.Fatal(test, "expected:", expected, "got:", got)
	}
}
