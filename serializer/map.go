package serializer

import (
	"fmt"
	"reflect"
)

func encodeMap(m any) ([]byte, error) {
	t := reflect.TypeOf(m)
	if t.Kind() != reflect.Map {
		return nil, fmt.Errorf("CompressSerializer: unsupported type %d", t.Kind())
	}

	v := reflect.ValueOf(m)
	if v.Len() == 0 {
		return encodeVarint(0), nil
	}

	ret := make([]byte, 0)

	keyType := t.Key()
	if keyType.Kind() == reflect.Ptr {
		keyType = keyType.Elem()
	}
	head, err := encodedKind2Head(DataHead, keyType.Kind())
	if err != nil {
		return nil, fmt.Errorf("CompressSerializer: %s", err)
	}
	ret = append(ret, head)

	valueType := t.Elem()
	if valueType.Kind() == reflect.Ptr {
		valueType = valueType.Elem()
	}
	head, err = encodedKind2Head(DataHead, valueType.Kind())
	if err != nil {
		return nil, fmt.Errorf("CompressSerializer: %s", err)
	}
	ret = append(ret, head)

	for _, key := range v.MapKeys() {
		itemData := make([]byte, 0)

		if key.Kind() == reflect.Ptr {
			key = key.Elem()
		}

		switch key.Kind() {
		case reflect.String:
			data, err := encodeString(key.String())
			if err != nil {
				return nil, fmt.Errorf("CompressSerializer: %s", err)
			}
			itemData = append(itemData, data...)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			data := encodeVarint(uint64(key.Int()))
			itemData = append(itemData, data...)
		case reflect.Struct:
			data, err := encodeStruct(key.Interface())
			if err != nil {
				return nil, fmt.Errorf("CompressSerializer: %s", err)
			}
			itemData = append(itemData, data...)
		case reflect.Slice:
			data, err := encodeSlice(key.Interface())
			if err != nil {
				return nil, fmt.Errorf("CompressSerializer: %s", err)
			}
			itemData = append(itemData, data...)
		default:
			return nil, fmt.Errorf("CompressSerializer: unsupported type %d", key.Kind())
		}

		value := v.MapIndex(key)

		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}

		switch value.Kind() {
		case reflect.String:
			data, err := encodeString(value.String())
			if err != nil {
				return nil, fmt.Errorf("CompressSerializer: %s", err)
			}
			itemData = append(itemData, data...)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			data := encodeVarint(uint64(value.Int()))
			itemData = append(itemData, data...)
		case reflect.Struct:
			data, err := encodeStruct(value.Interface())
			if err != nil {
				return nil, fmt.Errorf("CompressSerializer: %s", err)
			}
			itemData = append(itemData, data...)
		case reflect.Slice:
			data, err := encodeSlice(value.Interface())
			if err != nil {
				return nil, fmt.Errorf("CompressSerializer: %s", err)
			}
			itemData = append(itemData, data...)
		default:
			return nil, fmt.Errorf("CompressSerializer: unsupported type %d", value.Kind())
		}

		dataSize := len(itemData)

		ret = append(ret, encodeVarint(uint64(dataSize))...)
		ret = append(ret, itemData...)

	}

	dataSize := len(ret)
	ret = append(encodeVarint(uint64(dataSize)), ret...)

	return ret, nil
}

func decodeMap(data []byte, target any) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("CompressSerializer: empty data")
	}

	dataSize, n := decodeVarint(data)
	data = data[n:]

	if dataSize == 0 {
		return n, nil
	}

	data = data[:dataSize]

	_, keyTid := decodeHead(data[0])
	data = data[1:]
	_, valueTid := decodeHead(data[0])
	data = data[1:]

	t := reflect.TypeOf(target)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Map {
		return 0, fmt.Errorf("CompressSerializer: unsupported type %d", t.Kind())
	}

	v := reflect.ValueOf(target)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	v.Set(reflect.MakeMap(t))

	keyType := t.Key()
	valueType := t.Elem()

	for len(data) > 0 {
		itemSize, n := decodeVarint(data)
		data = data[n:]

		itemData := data[:itemSize]
		data = data[itemSize:]

		key := reflect.New(keyType).Elem()
		value := reflect.New(valueType).Elem()

		if keyType.Kind() == reflect.Ptr {
			key = reflect.New(keyType.Elem()).Elem()
		}

		if valueType.Kind() == reflect.Ptr {
			value = reflect.New(valueType.Elem()).Elem()
		}

		if key.Kind() == reflect.Ptr {
			key = key.Elem()
		}

		switch keyTid {
		case String:
			str, n, err := decodeString(itemData)
			if err != nil {
				return 0, fmt.Errorf("CompressSerializer: %s", err)
			}
			itemData = itemData[n:]
			key = reflect.ValueOf(str)
		case VarInt:
			num, n := decodeVarint(itemData)
			itemData = itemData[n:]
			key.Set(reflect.ValueOf(num))
		default:
			return 0, fmt.Errorf("CompressSerializer: unsupported type %d", keyTid)
		}

		switch valueTid {
		case String:
			str, n, err := decodeString(itemData)
			if err != nil {
				return 0, fmt.Errorf("CompressSerializer: %s", err)
			}
			itemData = itemData[n:]
			value = reflect.ValueOf(str)
		case VarInt:
			num, n := decodeVarint(itemData)
			itemData = itemData[n:]
			switch valueType.Kind() {
			case reflect.Int:
				value.Set(reflect.ValueOf(int(num)))
			case reflect.Int8:
				value.Set(reflect.ValueOf(int8(num)))
			case reflect.Int16:
				value.Set(reflect.ValueOf(int16(num)))
			case reflect.Int32:
				value.Set(reflect.ValueOf(int32(num)))
			case reflect.Int64:
				value.Set(reflect.ValueOf(int64(num)))
			case reflect.Uint:
				value.Set(reflect.ValueOf(uint(num)))
			case reflect.Uint8:
				value.Set(reflect.ValueOf(uint8(num)))
			case reflect.Uint16:
				value.Set(reflect.ValueOf(uint16(num)))
			case reflect.Uint32:
				value.Set(reflect.ValueOf(uint32(num)))
			case reflect.Uint64:
				value.Set(reflect.ValueOf(uint64(num)))
			default:
				return 0, fmt.Errorf("CompressSerializer: unsupported type %d", valueType.Kind())
			}
		case Struct:
			n, err := decodeStruct(itemData, value.Addr().Interface())
			if err != nil {
				return 0, fmt.Errorf("CompressSerializer: %s", err)
			}
			itemData = itemData[n:]
		case Slice:
			n, err := decodeSlice(itemData, value.Addr().Interface())
			if err != nil {
				return 0, fmt.Errorf("CompressSerializer: %s", err)
			}
			itemData = itemData[n:]
		default:
			return 0, fmt.Errorf("CompressSerializer: unsupported type %d", valueTid)
		}

		if keyType.Kind() == reflect.Ptr {
			key = key.Addr()
		}

		if valueType.Kind() == reflect.Ptr {
			value = value.Addr()
		}

		v.SetMapIndex(key, value)

	}

	return n + int(dataSize), nil
}
