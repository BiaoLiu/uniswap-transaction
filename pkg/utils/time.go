package utils

import "time"

const (
	TimeFormat = "2006-01-02 15:04:05"
	DateFormat = "2006-01-02"

	TimeFormatSimplify   = "20060102150405"
	MinuteFormatSimplify = "200601021504"
	DateFormatSimplify   = "20060102"
)

func Format(t time.Time, layout string) string {
	return t.Format(layout)
}
