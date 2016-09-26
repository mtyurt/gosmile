package gosmile

type Encoder struct {
	ContainsRawBinary         bool
	SharedStringValueEnabled  bool
	SharedPropertyNameEnabled bool
	content                   []byte
}

func (e *Encoder) Marshal(v interface{}) ([]byte, error) {
	return []byte{}, nil
}

func zigzagEncodeInt(n int) int {
	return (n << 1) ^ (n >> 31)
}

func zigzagEncodeLong(n int64) int64 {
	return (n << 1) ^ (n >> 63)
}
