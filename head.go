package serializer_compress_go

const (
	DataHead = 99
)

func encodeHead(idx, t int) byte {
	return byte(idx<<3 | t)
}

func decodeHead(x byte) (idx, t int) {
	return int(x >> 3), int(x & 7)
}
