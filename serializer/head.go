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
		return encodeHead(idx, String), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return encodeHead(idx, VarInt), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return encodeHead(idx, VarUint), nil
	case reflect.Float32, reflect.Float64:
		return encodeHead(idx, Float), nil
	case reflect.Struct:
		return encodeHead(idx, Struct), nil
	case reflect.Slice:
		return encodeHead(idx, Slice), nil
	default:
		return 0, fmt.Errorf("unsupported kind: %v", k)
	}
}
