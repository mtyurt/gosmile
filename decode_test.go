package gosmile

import (
	"math"
	"strings"
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

func TestInvalidHeader(t *testing.T) {
	err := Unmarshal([]byte("asdfhi\n"), 3)
	if err == nil || !strings.Contains(err.Error(), "invalid header: asd") {
		t.Fatal("should return an error, error:", err)
	}
}

func TestInitConfig(t *testing.T) {
	content := []byte{':', ')', '\n', 0xa7}
	d := &decoder{data: content}
	err := d.init()
	if err != nil {
		t.Fatal("shouldn't trigger any error", err)
	}
	if d.version != 0xa || !d.sharedStringValue || !d.sharedPropertyName || !d.rawBinary {
		t.Fatal("wrong config:", d, "version should be 10, all flags should be true")
	}
}
