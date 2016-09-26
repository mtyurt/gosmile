package gosmile

import (
	"math"
	"testing"
)

func TestZigZagDecodeInt(t *testing.T) {
	expect(0, zigzagDecodeInt(0), t, "zigzagdecode0")
	expect(-1, zigzagDecodeInt(1), t, "zigzagdecode-1")
	expect(1, zigzagDecodeInt(2), t, "zigzagdecode1")
	expect(math.MaxInt32, zigzagDecodeInt(0xFFFFFFFE), t, "zigzagdecodeintmax")
	expect(math.MinInt32, zigzagDecodeInt(0xFFFFFFFF), t, "zigzagdecodeintmin")
}

func TestZigZagDecodeLong(t *testing.T) {
	expect(int64(math.MaxInt64), zigzagDecodeLong(int64(-2)), t, "zigzagdecode-2")
	expect(int64(math.MinInt64), zigzagDecodeLong(int64(-1)), t, "zigzagdecode-1")
}
