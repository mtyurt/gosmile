package gosmile

import (
	"log"
	"math"
	"runtime"
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
	e := NewEncodeConf()
	e.Version = 3
	content, err := Marshal(e, 1)
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

	e = NewEncodeConf()
	e.SharedStringValueEnabled = true
	e.SharedPropertyNameEnabled = false

	content, err = Marshal(e, 1)
	if err != nil || len(content) < 4 {
		t.Fatal("err here", err, "content size:", len(content))
	}
	varByte = content[3]
	expect(byte(1)<<1, varByte&0x02, t, "testencodeheader-sharedstringvalue2")
	expect(byte(0x00), varByte&0x01, t, "testencodeheader-sharedpropname2")

}
func TestEncodeSmallInt(t *testing.T) {
	e := NewEncodeConf()
	e.IncludeHeader = false

	assertSmallInt := func(n int, name string) {
		c, err := Marshal(e, n)
		if err != nil {
			t.Fatal("error", err)
		}
		expect(1, len(c), t, name+"length")
		expect(byte(token_prefix_small_int+zigzagEncodeInt(n)), c[0], t, name)
	}
	assertSmallInt(3, "encode3")
	assertSmallInt(0, "encode0")
	assertSmallInt(-6, "encode-6")
	assertSmallInt(15, "encode15")
	assertSmallInt(-16, "encode-16")
}

func TestOtherInts(t *testing.T) {
	e := NewEncodeConf()
	e.IncludeHeader = false

	assertInt := func(n int, length int, name string) {
		c, err := Marshal(e, n)
		if err != nil {
			t.Fatal(err)
		}
		expect(length, len(c), t, name+"length")
	}
	assertInt(16, 2, "encode16")
	assertInt(-17, 2, "encode-17")
	assertInt(0xfff, 3, "encode0xfff")
	assertInt(-4096, 3, "encode-4096")
	assertInt(0x1000, 4, "encode0x1000")
	assertInt(500000, 4, "encode500000")
	assertInt(math.MaxInt32, 6, "encodemaxint")
	assertInt(math.MinInt32, 6, "encodeminint")

}

func TestFloat32(t *testing.T) {
	e := NewEncodeConf()
	e.IncludeHeader = false

	c, err := Marshal(e, float32(0.125))
	if err != nil || c == nil {
		t.Fatal("encode float failed, content:", c, "error:", err)
	}
	if len(c) != 6 {
		t.Fatal("encode float failed, expected length is 6 but it was:", len(c))
	}
}

func TestFloat64(t *testing.T) {
	e := NewEncodeConf()
	e.IncludeHeader = false

	c, err := Marshal(e, float64(0.125))
	if err != nil || c == nil {
		t.Fatal("encode float64 failed, content:", c, "error:", err)
	}
	if len(c) != 11 {
		t.Fatal("encode float64 failed, expected length is 11 but it was:", len(c))
	}
}

func TestBool(t *testing.T) {
	e := NewEncodeConf()
	e.IncludeHeader = false

	c, err := Marshal(e, false)
	if err != nil || c == nil || len(c) != 1 || c[0] != token_literal_false {
		t.Fatal("encode false failed, content:", c, "error:", err)
	}

	c, err = Marshal(e, true)
	if err != nil || c == nil || len(c) != 1 || c[0] != token_literal_true {
		t.Fatal("encode true failed, content:", c, "error:", err)
	}
}

func TestTinyString(t *testing.T) {
	e := NewEncodeConf()
	e.IncludeHeader = false

	originalVal := "strtest"
	val := originalVal
	c, err := Marshal(e, val)
	if err != nil || c == nil || len(c) != 8 || (c[0] != (token_prefix_tiny_ascii+7)) && string(c[1:]) != val {
		t.Fatal("encode short ascii string failed, err:", err, "content:", c)
	}
	val = originalVal + "œ"
	c, err = Marshal(e, val)
	if err != nil || c == nil || len(c) != 10 || (c[0] != (token_prefix_tiny_unicode+10)) && string(c[1:]) != val {
		t.Fatal("encode short unicode 2-byte string failed, err:", err, "content:", c)
	}
	val = originalVal + "≈"
	c, err = Marshal(e, val)
	if err != nil || c == nil || len(c) != 11 || (c[0] != (token_prefix_tiny_unicode+11)) && string(c[1:]) != val {
		t.Fatal("encode short unicode 3-byte string failed, err:", err, "content:", c)
	}
}

func TestLongString(t *testing.T) {
	e := NewEncodeConf()
	e.IncludeHeader = false

	originalVal := "benim adim insanlarin hizasina yazilmistir\n" +
		"her gun yepyeni ruyalarla odenebilecek bir ceza bu\n"
	val := originalVal
	c, err := Marshal(e, val)
	if err != nil || c == nil || len(c) != 96 || (c[0] != token_byte_long_string_ascii) && string(c[1:]) != val && c[len(c)-1] != byte_marker_end_of_string {
		t.Fatal("encode long ascii string failed, err:", err, "content len:", len(c), "first byte:", c[0], "last byte:", c[len(c)-1])
	}
	val = originalVal + "œ"
	c, err = Marshal(e, val)
	if err != nil || c == nil || len(c) != 98 || (c[0] != token_misc_long_text_unicode) && string(c[1:]) != val && c[len(c)-1] != byte_marker_end_of_string {
		t.Fatal("encode long unicode 2-byte string failed, err:", err, "content len:", len(c), "first byte:", c[0], "last byte:", c[len(c)-1])
	}
	val = originalVal + "≈"
	c, err = Marshal(e, val)
	if err != nil || c == nil || len(c) != 99 || (c[0] != token_misc_long_text_unicode) && string(c[1:]) != val && c[len(c)-1] != byte_marker_end_of_string {
		t.Fatal("encode long unicode 3-byte string failed, err:", err, "content len:", len(c), "first byte:", c[0], "last byte:", c[len(c)-1])
	}
}

func TestSimpleSlice(t *testing.T) {
	e := NewEncodeConf()
	e.IncludeHeader = false

	val := []int{1, 2, 16}

	c, err := Marshal(e, val)
	if err != nil || len(c) != 6 || c[0] != token_literal_start_array || c[5] != token_literal_end_array {
		t.Fatal("encode int slice failed, err:", err, "len(c):", len(c))
	}

	strSlice := []string{"strtest", "strtestœ"}
	c, err = Marshal(e, strSlice)
	if err != nil || len(c) != 20 || c[0] != token_literal_start_array || c[19] != token_literal_end_array {
		t.Fatal("encode string slice failed, err:", err, "len(c):", len(c))
	}
}

func TestSimpleArray(t *testing.T) {
	e := NewEncodeConf()
	e.IncludeHeader = false

	val := [3]int{1, 2, 16}

	c, err := Marshal(e, val)
	if err != nil || len(c) != 6 || c[0] != token_literal_start_array || c[5] != token_literal_end_array {
		t.Fatal("encode int array failed, err:", err, "len(c):", len(c))
	}
}

func expect(expected interface{}, got interface{}, t *testing.T, test string) {
	if got != expected {
		var stack [4096]byte
		runtime.Stack(stack[:], false)
		log.Printf("%s\n", stack[:])
		t.Fatal(test, "expected:", expected, "got:", got)
	}
}
