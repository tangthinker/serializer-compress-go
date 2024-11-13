package serializer

import (
	"fmt"
	"reflect"
)

func encodeSlice(source any) ([]byte, error) {
	t := reflect.TypeOf(source)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Slice {
		return nil, fmt.Errorf("serializer: invalid type %s", t.Kind())
	}

	ret := make([]byte, 0)

	v := reflect.ValueOf(source)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Len() == 0 {
		return encodeVarint(0), nil
	}

	itemType := v.Index(0).Type()
	if itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}

	head, err := encodedKind2Head(DataHead, itemType.Kind())
	if err != nil {
		return nil, fmt.Errorf("serializer: encode error %w", err)
	}
	ret = append(ret, head)

	length := v.Len()
	ret = append(ret, encodeVarint(uint64(length))...)

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		switch itemType.Kind() {
		case reflect.Struct:
			data, err := encodeStruct(item.Interface())
			if err != nil {
				return nil, fmt.Errorf("serializer: encode error %w", err)
			}
			ret = append(ret, data...)
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			data, err := encodeInt64(item.Int())
			if err != nil {
				return nil, fmt.Errorf("serializer: encode error %w", err)
			}
			ret = append(ret, data...)
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			data := encodeVarint(item.Uint())
			ret = append(ret, data...)
		case reflect.Float32, reflect.Float64:
			data := encodeFloat64(item.Float())
			ret = append(ret, data...)
		case reflect.String:
			data, err := encodeString(item.String())
			if err != nil {
				return nil, fmt.Errorf("serializer: encode error %w", err)
			}
			ret = append(ret, data...)
		default:
			return nil, fmt.Errorf("serializer: unsupported type %s", itemType.Kind())
		}
	}

	size := len(ret)
	if size > 0 {
		ret = append(encodeVarint(uint64(size)), ret...)
	}

	return ret, nil
}

func decodeSlice(data []byte, target any) (int, error) {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Ptr {
		return 0, fmt.Errorf("CompressSerializer: invalid type %s", t.Kind())
	}
	t = t.Elem()
	if t.Kind() != reflect.Slice {
		return 0, fmt.Errorf("CompressSerializer: invalid type %s", t.Kind())
	}

	dataSize, n := decodeVarint(data)
	data = data[n:]

	if dataSize == 0 {
		return n, nil
	}

	data = data[:dataSize]

	_, eleType := decodeHead(data[0])
	data = data[1:]

	size, n := decodeVarint(data)
	data = data[n:]

	v := reflect.ValueOf(target).Elem()
	if v.Len() != int(size) {
		newSlice := reflect.MakeSlice(t, int(size), int(size))
		v.Set(newSlice)
	}

	for i := 0; i < int(size); i++ {
		item := v.Index(i)
		itemType := item.Type()
		if itemType.Kind() == reflect.Ptr {
			newItem := reflect.New(itemType.Elem()).Elem()
			item.Set(newItem.Addr())
			item = item.Elem()
		}
		switch eleType {
		case Struct:
			n, err := decodeStruct(data, item.Addr().Interface())
			if err != nil {
				return 0, fmt.Errorf("CompressSerializer: decode struct error %w", err)
			}
			data = data[n:]
		case VarInt:
			value, n, err := decodeInt64(data)
			if err != nil {
				return 0, fmt.Errorf("CompressSerializer: decode error %w", err)
			}
			data = data[n:]
			item.SetInt(value)
		case VarUint:
			value, n := decodeVarint(data)
			data = data[n:]
			item.SetUint(value)
		case Float:
			value, n := decodeFloat64(data)
			data = data[n:]
			item.SetFloat(value)
		case String:
			value, n, err := decodeString(data)
			if err != nil {
				return 0, fmt.Errorf("CompressSerializer: decode error %w", err)
			}
			data = data[n:]
			item.SetString(value)
		default:
			return 0, fmt.Errorf("CompressSerializer: unsupported type %d", eleType)
		}
	}

	return n + int(dataSize), nil

}
