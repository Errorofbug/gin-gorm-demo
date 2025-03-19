package util

import (
	"fmt"
	"time"
)

// DefaultTimeFormat 默认的时间格式
const DefaultTimeFormat = "2006-01-02 15:04:05"

// ParseTime 将时间字符串解析为time.Time
func ParseTime(timeStr string, loc *time.Location) (time.Time, error) {
	t, err := time.ParseInLocation(DefaultTimeFormat, timeStr, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("解析时间字符串失败: %v", err)
	}
	return t, nil
}

// FormatTimestamp 将时间戳格式化为字符串
func FormatTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(DefaultTimeFormat)
}

// FormatTime 将time.Time格式化为字符串
func FormatTime(t time.Time) string {
	return t.Format(DefaultTimeFormat)
}

// ParseToTimestamp 将时间字符串解析为时间戳
func ParseToTimestamp(timeStr string, loc *time.Location) (int64, error) {
	t, err := time.ParseInLocation(DefaultTimeFormat, timeStr, loc)
	if err != nil {
		return 0, fmt.Errorf("解析时间字符串失败: %v", err)
	}
	return t.UTC().Unix(), nil
}

// TimestampToTime 将时间戳转换为time.Time
func TimestampToTime(timestamp int64) time.Time {
	return time.Unix(timestamp, 0)
}
