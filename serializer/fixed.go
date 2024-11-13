package serializer

import "unsafe"

func encodeFixed64(data uint64) []byte {
	var buf [8]byte
	buf[0] = byte(data)
	buf[1] = byte(data >> 8)
	buf[2] = byte(data >> 16)
	buf[3] = byte(data >> 24)
	buf[4] = byte(data >> 32)
	buf[5] = byte(data >> 40)
	buf[6] = byte(data >> 48)
	buf[7] = byte(data >> 56)
	return buf[:]
}

func decodeFixed64(data []byte) (uint64, int) {
	var ret uint64
	ret = uint64(data[0])
	ret |= uint64(data[1]) << 8
	ret |= uint64(data[2]) << 16
	ret |= uint64(data[3]) << 24
	ret |= uint64(data[4]) << 32
	ret |= uint64(data[5]) << 40
	ret |= uint64(data[6]) << 48
	ret |= uint64(data[7]) << 56
	return ret, 8
}

func encodeFixed32(data uint32) []byte {
	var buf [4]byte
	buf[0] = byte(data)
	buf[1] = byte(data >> 8)
	buf[2] = byte(data >> 16)
	buf[3] = byte(data >> 24)
	return buf[:]
}

func decodeFixed32(data []byte) (uint32, int) {
	var ret uint32
	ret = uint32(data[0])
	ret |= uint32(data[1]) << 8
	ret |= uint32(data[2]) << 16
	ret |= uint32(data[3]) << 24
	return ret, 4
}

func encodeFloat32(data float32) []byte {
	var buf [4]byte
	bits := *(*uint32)(unsafe.Pointer(&data))
	buf[0] = byte(bits)
	buf[1] = byte(bits >> 8)
	buf[2] = byte(bits >> 16)
	buf[3] = byte(bits >> 24)
	return buf[:]
}

func decodeFloat32(data []byte) (float32, int) {
	bits := uint32(data[0])
	bits |= uint32(data[1]) << 8
	bits |= uint32(data[2]) << 16
	bits |= uint32(data[3]) << 24
	return *(*float32)(unsafe.Pointer(&bits)), 4
}

func encodeFloat64(data float64) []byte {
	var buf [8]byte
	bits := *(*uint64)(unsafe.Pointer(&data))
	buf[0] = byte(bits)
	buf[1] = byte(bits >> 8)
	buf[2] = byte(bits >> 16)
	buf[3] = byte(bits >> 24)
	buf[4] = byte(bits >> 32)
	buf[5] = byte(bits >> 40)
	buf[6] = byte(bits >> 48)
	buf[7] = byte(bits >> 56)
	return buf[:]
}

func decodeFloat64(data []byte) (float64, int) {
	bits := uint64(data[0])
	bits |= uint64(data[1]) << 8
	bits |= uint64(data[2]) << 16
	bits |= uint64(data[3]) << 24
	bits |= uint64(data[4]) << 32
	bits |= uint64(data[5]) << 40
	bits |= uint64(data[6]) << 48
	bits |= uint64(data[7]) << 56
	return *(*float64)(unsafe.Pointer(&bits)), 8
}
