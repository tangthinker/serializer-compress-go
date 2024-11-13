package serializer

import (
	"fmt"
	"reflect"
	"strconv"
)

func encodeStruct(source any) ([]byte, error) {
	t := reflect.TypeOf(source)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("serializer: invalid type %s", t.Kind())
	}

	m := make(map[int]reflect.StructField)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("compress")
		s, err := strconv.Atoi(tag)
		if err != nil {
			return nil, fmt.Errorf("serializer: invalid tag %s", tag)
		}
		m[s] = field
	}

	v := reflect.ValueOf(source)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	ret := make([]byte, 0)
	for idx, field := range m {
		fieldValue := v.FieldByName(field.Name)
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			header := encodeHead(idx, VarInt)
			ret = append(ret, header)
			data, err := encodeInt64(fieldValue.Int())
			if err != nil {
				return nil, fmt.Errorf("serializer: encode error %w", err)
			}
			ret = append(ret, data...)
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			header := encodeHead(idx, VarUint)
			ret = append(ret, header)
			data := encodeVarint(fieldValue.Uint())
			ret = append(ret, data...)
		case reflect.Float32, reflect.Float64:
			header := encodeHead(idx, Float)
			ret = append(ret, header)
			data := encodeFloat64(fieldValue.Float())
			ret = append(ret, data...)
		case reflect.String:
			header := encodeHead(idx, String)
			ret = append(ret, header)
			data, err := encodeString(fieldValue.String())
			if err != nil {
				return nil, fmt.Errorf("serializer: encode error %w", err)
			}
			ret = append(ret, data...)
		case reflect.Slice:
			header := encodeHead(idx, Slice)
			ret = append(ret, header)
			data, err := encodeSlice(fieldValue.Interface())
			if err != nil {
				return nil, fmt.Errorf("serializer: encode error %w", err)
			}
			ret = append(ret, data...)
		case reflect.Struct:
			header := encodeHead(idx, Struct)
			ret = append(ret, header)
			data, err := encodeStruct(fieldValue.Interface())
			if err != nil {
				return nil, fmt.Errorf("serializer: encode error %w", err)
			}
			ret = append(ret, data...)
		case reflect.Map:
			header := encodeHead(idx, Map)
			ret = append(ret, header)
			data, err := encodeMap(fieldValue.Interface())
			if err != nil {
				return nil, fmt.Errorf("serializer: encode error %w", err)
			}
			ret = append(ret, data...)
		default:
			return nil, fmt.Errorf("serializer: unsupported type %s", field.Type.Kind())

		}

	}

	size := len(ret)
	ret = append(encodeVarint(uint64(size)), ret...)

	return ret, nil

}

func decodeStruct(data []byte, target any) (int, error) {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Ptr {
		return 0, fmt.Errorf("CompressSerializer: invalid type %s", t.Kind())
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		return 0, fmt.Errorf("CompressSerializer: invalid type %s", t.Kind())
	}

	m := make(map[int]reflect.StructField)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("compress")
		s, err := strconv.Atoi(tag)
		if err != nil {
			return 0, fmt.Errorf("CompressSerializer: invalid tag %s", tag)
		}
		m[s] = field
	}

	x, n := decodeVarint(data)
	data = data[n:]
	data = data[:x]

	for len(data) > 0 {
		idx, typeId := decodeHead(data[0])
		data = data[1:]
		field, ok := m[idx]
		if !ok {
			return 0, fmt.Errorf("CompressSerializer: field not found %d", idx)
		}

		switch typeId {
		case VarInt:
			value, n, err := decodeInt64(data)
			if err != nil {
				return 0, fmt.Errorf("CompressSerializer: decode error %w", err)
			}
			data = data[n:]
			reflect.ValueOf(target).Elem().FieldByName(field.Name).SetInt(value)
		case VarUint:
			value, n := decodeVarint(data)
			data = data[n:]
			reflect.ValueOf(target).Elem().FieldByName(field.Name).SetUint(value)
		case Float:
			value, n := decodeFloat64(data)
			data = data[n:]
			reflect.ValueOf(target).Elem().FieldByName(field.Name).SetFloat(value)
		case String:
			value, n, err := decodeString(data)
			if err != nil {
				return 0, fmt.Errorf("CompressSerializer: decode error %w", err)
			}
			data = data[n:]
			reflect.ValueOf(target).Elem().FieldByName(field.Name).SetString(value)
		case Slice:
			n, err := decodeSlice(data, reflect.ValueOf(target).Elem().FieldByName(field.Name).Addr().Interface())
			if err != nil {
				return 0, fmt.Errorf("CompressSerializer: decode error %w", err)
			}
			data = data[n:]
		case Struct:
			n, err := decodeStruct(data, reflect.ValueOf(target).Elem().FieldByName(field.Name).Addr().Interface())
			if err != nil {
				return 0, fmt.Errorf("CompressSerializer: decode error %w", err)
			}
			data = data[n:]
		case Map:
			n, err := decodeMap(data, reflect.ValueOf(target).Elem().FieldByName(field.Name).Addr().Interface())
			if err != nil {
				return 0, fmt.Errorf("CompressSerializer: decode error %w", err)
			}
			data = data[n:]
		default:
			return 0, fmt.Errorf("CompressSerializer: unsupported type %d", typeId)
		}
	}

	return n + int(x), nil
}
