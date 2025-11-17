package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// NormalizePossiblyTextualBytes 尝试多轮将数字/列表/JSON文本还原为二进制或SVG/XML
func NormalizePossiblyTextualBytes(data []byte, _ string) []byte {
	const maxPass = 4
	for pass := 1; pass <= maxPass; pass++ {
		lower := bytes.ToLower(data)
		if bytes.Contains(lower, []byte("<svg")) || bytes.HasPrefix(lower, []byte("<?xml")) {
			return data
		}
		if !(looksLikeByteListHeader(data) || bytes.HasPrefix(lower, []byte("[")) || bytes.HasPrefix(lower, []byte("{"))) {
			return data
		}
		if bin, ok := parseJSONByteArray(data); ok {
			data = bin
		} else if bin, ok := parseAsciiByteList(data); ok {
			data = bin
		} else if sub, ok := extractBracketedByteList(data); ok {
			if bin, ok := parseAsciiByteList(sub); ok {
				data = bin
			} else if bin, ok := parseFlexibleByteArray(sub); ok {
				data = bin
			} else if bin, ok := convertNumericTextToBytes(sub); ok {
				data = bin
			} else {
				return data
			}
		} else if bin, ok := parseFlexibleByteArray(data); ok {
			data = bin
		} else if bin, ok := convertNumericTextToBytes(data); ok {
			data = bin
		} else {
			return data
		}
	}
	return data
}

// convertNumericTextToBytes 无门槛转换：提取数据中所有十进制数字并按字节拼接
func convertNumericTextToBytes(data []byte) ([]byte, bool) {
	out := make([]byte, 0, 1024)
	val := 0
	hasDigit := false
	appended := false
	flush := func() {
		if !hasDigit {
			return
		}
		if val < 0 {
			val = 0
		}
		if val > 255 {
			val = val % 256
		}
		out = append(out, byte(val))
		val = 0
		hasDigit = false
		appended = true
	}
	for _, b := range data {
		if b >= '0' && b <= '9' {
			hasDigit = true
			val = val*10 + int(b-'0')
			continue
		}
		flush()
	}
	flush()
	if !appended || len(out) == 0 {
		return nil, false
	}
	return out, true
}

// parseJSONByteArray 尝试将 JSON 数组（如 [255,216,255,...]）还原为二进制
func parseJSONByteArray(data []byte) ([]byte, bool) {
	var arr []int
	if err := json.Unmarshal(data, &arr); err != nil {
		return nil, false
	}
	out := make([]byte, 0, len(arr))
	for _, v := range arr {
		if v < 0 || v > 255 {
			return nil, false
		}
		out = append(out, byte(v))
	}
	return out, true
}

// parseAsciiByteList 尝试将形如 "[60 63 120 ...]" 的文本还原为二进制
func parseAsciiByteList(data []byte) ([]byte, bool) {
	s := strings.TrimSpace(string(data))
	if len(s) < 2 || s[0] != '[' || s[len(s)-1] != ']' {
		return nil, false
	}
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return nil, false
	}
	out := make([]byte, 0, len(fields))
	for _, f := range fields {
		f = strings.TrimSuffix(f, ",")
		if f == "" {
			continue
		}
		var v int
		n, err := fmt.Sscanf(f, "%d", &v)
		if err != nil || n != 1 || v < 0 || v > 255 {
			return nil, false
		}
		out = append(out, byte(v))
	}
	return out, true
}

// parseFlexibleByteArray 更鲁棒的解析器：允许空白/逗号/换行，忽略非数字符号，提取0-255
func parseFlexibleByteArray(data []byte) ([]byte, bool) {
	if !looksLikeByteListHeader(data) {
		return nil, false
	}
	out := make([]byte, 0, 1024)
	n := 0
	val := 0
	hasDigit := false
	flush := func() bool {
		if !hasDigit {
			return true
		}
		if val < 0 {
			val = 0
		}
		if val > 255 {
			val = val % 256
		}
		out = append(out, byte(val))
		n++
		val = 0
		hasDigit = false
		return true
	}
	for _, b := range data {
		switch {
		case b >= '0' && b <= '9':
			hasDigit = true
			val = val*10 + int(b-'0')
		case b == ' ' || b == ',' || b == '\n' || b == '\r' || b == '\t' || b == '[' || b == ']':
			if !flush() {
				return nil, false
			}
		default:
			if !flush() {
				return nil, false
			}
		}
	}
	_ = flush()
	if len(out) == 0 {
		return nil, false
	}
	return out, true
}

func looksLikeByteListHeader(data []byte) bool {
	max := len(data)
	if max > 128 {
		max = 128
	}
	score := 0
	for i := 0; i < max; i++ {
		b := data[i]
		if b == '[' || b == ']' || b == ' ' || b == ',' || b == '\n' || b == '\r' || b == '\t' || (b >= '0' && b <= '9') {
			score++
		}
	}
	return score*100/max >= 80
}

// 提取文本中首个方括号数组（含 [ 和 ]）的子串，适配诸如 "Uint8Array(...) [60, 63, ...]" 的格式
func extractBracketedByteList(data []byte) ([]byte, bool) {
	start := bytes.IndexByte(data, '[')
	if start == -1 {
		return nil, false
	}
	end := bytes.LastIndexByte(data, ']')
	if end == -1 || end <= start {
		return nil, false
	}
	return data[start : end+1], true
}
