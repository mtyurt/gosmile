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

func expect(expected interface{}, got interface{}, t *testing.T, test string) {
	if got != expected {
		t.Fatal(test, "expected:", expected, "got:", got)
	}
}
