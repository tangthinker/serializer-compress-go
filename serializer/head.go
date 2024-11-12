package serializer

import (
	"fmt"
	"reflect"
)

const (
	DataHead = 99
)

func encodeHead(idx, t int) byte {
	return byte(idx<<3 | t)
}

func decodeHead(x byte) (idx, t int) {
	return int(x >> 3), int(x & 7)
}

func encodedKind2Head(idx int, k reflect.Kind) (byte, error) {
	switch k {
	case reflect.String:
		return encodeHead(DataHead, String), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return encodeHead(DataHead, VarInt), nil
	case reflect.Struct:
		return encodeHead(DataHead, Struct), nil
	case reflect.Slice:
		return encodeHead(DataHead, Slice), nil
	default:
		return 0, fmt.Errorf("unsupported kind: %v", k)
	}
}
