package serializer

import (
	"fmt"
	"reflect"
)

const maxVarintBytes = 10

const (
	VarInt int = iota
	VarUint
	Float
	String
	Struct
	Slice
	Map
)

type Serializer interface {
	Encode(source any) ([]byte, error)
	Decode(data []byte, target any) error
}

type serializer struct {
}

func NewSerializer() Serializer {
	return &serializer{}
}

func (s *serializer) Encode(source any) ([]byte, error) {
	if source == nil {
		return nil, nil
	}

	t := reflect.TypeOf(source)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	ret := make([]byte, 0)
	var err error

	switch t.Kind() {
	case reflect.Struct:
		head := encodeHead(DataHead, Struct)
		ret = append(ret, head)
		data, eErr := encodeStruct(source)
		err = eErr
		ret = append(ret, data...)
	case reflect.Slice:
		head := encodeHead(DataHead, Slice)
		ret = append(ret, head)
		data, eErr := encodeSlice(source)
		err = eErr
		ret = append(ret, data...)
	case reflect.Map:
		head := encodeHead(DataHead, Map)
		ret = append(ret, head)
		data, eErr := encodeMap(source)
		err = eErr
		ret = append(ret, data...)
	default:
		return nil, fmt.Errorf("serializer: unsupported type %s", t.Kind())
	}

	return ret, err
}

func (s *serializer) Decode(data []byte, target any) error {
	if len(data) == 0 {
		return nil
	}

	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("serializer: invalid type %s", t.Kind())
	}

	t = t.Elem()

	_, tId := decodeHead(data[0])
	data = data[1:]

	switch tId {
	case Struct:
		if t.Kind() != reflect.Struct {
			return fmt.Errorf("serializer: invalid type %s", t.Kind())
		}
		_, err := decodeStruct(data, target)
		return err
	case Slice:
		if t.Kind() != reflect.Slice {
			return fmt.Errorf("serializer: invalid type %s", t.Kind())
		}
		_, err := decodeSlice(data, target)
		return err
	case Map:
		if t.Kind() != reflect.Map {
			return fmt.Errorf("serializer: invalid type %s", t.Kind())
		}
		_, err := decodeMap(data, target)
		return err
	default:
		return fmt.Errorf("serializer: unsupported type %d", tId)
	}
}
