package common

import (
	"database/sql/driver"
	"fmt"
	"pixelpunk/pkg/errors"
	"time"
)

// JSONTime 是一个自定义的时间类型，用于格式化JSON输出中的时间
type JSONTime time.Time

// MarshalJSON 实现json.Marshaler接口，定制JSON格式
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(stamp), nil
}

// Value 实现 driver.Valuer 接口，用于保存到数据库
func (t JSONTime) Value() (driver.Value, error) {
	tTime := time.Time(t)
	return tTime, nil
}

// Scan 实现 sql.Scanner 接口，用于从数据库中读取
func (t *JSONTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = JSONTime(value)
		return nil
	}
	return errors.New(errors.CodeInvalidParameter, fmt.Sprintf("无法将 %v 转换为 JSONTime", v))
}

// JSONTimeNow 返回当前时间的JSONTime版本
func JSONTimeNow() JSONTime {
	return JSONTime(time.Now())
}

func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}
