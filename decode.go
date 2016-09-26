package gosmile

func Unmarshal(data []byte, v interface{}) error {
	return nil
}

func zigzagDecodeInt(n int) int {
	return int(uint32(n)>>1) ^ (-(n & 1))
}

func zigzagDecodeLong(n int64) int64 {
	return int64(uint64(n)>>1) ^ (-(n & 1))
}
