package serializer_compress_go

func encodeVarint(x uint64) []byte {
	var buf [maxVarintBytes]byte
	var n int
	for n = 0; x > 127; n++ {
		buf[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	buf[n] = uint8(x)
	n++
	return buf[0:n]
}

func decodeVarint(buf []byte) (x uint64, n int) {
	for shift := uint(0); shift < 64; shift += 7 {
		if n >= len(buf) {
			return 0, 0
		}
		b := uint64(buf[n])
		n++
		x |= (b & 0x7F) << shift
		if (b & 0x80) == 0 {
			return x, n
		}
	}

	// The number is too large to represent in a 64-bit value.
	return 0, 0
}

func encodeInt64(value int64) ([]byte, error) {
	if value >= 0 {
		value = value * 2
	} else {
		value = -value
		value = value*2 - 1
	}

	encodedData := encodeVarint(uint64(value))
	return encodedData, nil
}

func decodeInt64(data []byte) (int64, int, error) {
	x, n := decodeVarint(data)
	if x%2 == 0 {
		return int64(x / 2), n, nil
	}
	return -(int64(x) + 1) / 2, n, nil
}
