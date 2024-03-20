package utils

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"uniswap-transaction/pkg/idgenerator/snowflake"
)

func SafeInt(value string, defaultValue int) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}

func SafeInt64(value string, defaultValue int64) int64 {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

func CastEndTime(endTimeValue string) string {
	if endTimeValue != "" {
		endTime, err := time.Parse(DateFormat, endTimeValue)
		if err != nil {
			return ""
		}
		endTime = endTime.AddDate(0, 0, 1)
		return endTime.Format(DateFormat)
	}
	return endTimeValue
}

func ParseUrl(rawUrl string) string {
	if u, err := url.Parse(rawUrl); err == nil {
		return strings.TrimLeft(u.Path, "/")
	}
	return ""
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func FormatTime(val interface{}) string {
	switch t := val.(type) {
	case *time.Time:
		if t == nil {
			return ""
		}
		return t.Format(TimeFormat)
	case time.Time:
		return t.Format(TimeFormat)
	default:
		return ""
	}
}

func GenerateId() int64 {
	node, _ := snowflake.NewNode(0)
	snowId := node.Generate()
	return snowId.Int64()
}

func ConvertInt64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func ConvertFloatToString(value float64) string {
	return strconv.FormatFloat(value, 'E', -1, 64)
}
