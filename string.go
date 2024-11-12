package serializer_compress_go

func encodeString(value string) ([]byte, error) {
	size := len(value)
	encodedData := encodeVarint(uint64(size))
	encodedData = append(encodedData, []byte(value)...)
	return encodedData, nil
}

func decodeString(data []byte) (string, int, error) {
	x, n := decodeVarint(data)
	return string(data[n : n+int(x)]), n + int(x), nil
}
