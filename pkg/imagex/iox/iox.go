package iox

import (
	"errors"
	"io"
)

// DefaultMaxReadBytes 默认最大读取字节数（64MB）
const DefaultMaxReadBytes int64 = 64 << 20

// ReadAllWithLimit 读取全部数据，若超出上限返回错误
func ReadAllWithLimit(r io.Reader, maxBytes int64) ([]byte, error) {
	if maxBytes <= 0 {
		return io.ReadAll(r)
	}
	lr := &io.LimitedReader{R: r, N: maxBytes + 1}
	b, err := io.ReadAll(lr)
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > maxBytes {
		return nil, errors.New("read limit exceeded")
	}
	return b, nil
}
