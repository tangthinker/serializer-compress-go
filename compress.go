package serializer_compress_go

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

type CompressSerializer struct {
	baseSerializer Serializer
}

func NewCompressSerializer() Serializer {
	return &CompressSerializer{
		baseSerializer: NewSerializer(),
	}
}

func (c CompressSerializer) Encode(source any) ([]byte, error) {
	data, err := c.baseSerializer.Encode(source)
	if err != nil {
		return nil, fmt.Errorf("CompressSerializer: encode error %w", err)
	}

	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err = w.Write(data)
	if err != nil {
		_ = w.Close()
		return nil, fmt.Errorf("CompressSerializer: encode error %w", err)
	}
	_ = w.Close()
	return buf.Bytes(), nil

}

func (c CompressSerializer) Decode(data []byte, target any) error {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("ProtoBuffSerializer: decode error %w", err)
	}

	data, err = io.ReadAll(r)
	if err != nil {
		_ = r.Close()
		return fmt.Errorf("ProtoBuffSerializer: decode error %w", err)
	}
	_ = r.Close()

	return c.baseSerializer.Decode(data, target)
}
