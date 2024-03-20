package xtime

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = t.AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	*t = Time{now}
	return
}

func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}

func (t *Time) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = Time{value}
		return nil
	}
	return fmt.Errorf("can not convert %v to time", v)
}

func (t *Time) ToString() string {
	if t == nil {
		return Time{}.Format(TimeFormat)
	}
	return t.Format(TimeFormat)
}

func (t *Time) ToRFC3399() string {
	if t == nil {
		return Time{}.Format(time.RFC3339)
	}
	return t.Format(time.RFC3339)
}

func Now() Time {
	return Time{time.Now()}
}

func Parse(layout string, value string) (Time, error) {
	t, err := time.ParseInLocation(layout, value, time.Local)
	if err != nil {
		return Time{}, err
	}
	return Time{t}, nil
}
