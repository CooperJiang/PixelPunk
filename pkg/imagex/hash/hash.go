package hash

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// FromReader 计算Reader的MD5（读取到EOF，不复位seek）
func FromReader(r io.Reader) (string, error) {
	h := md5.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// FromFile 计算文件路径的MD5
func FromFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return FromReader(f)
}
