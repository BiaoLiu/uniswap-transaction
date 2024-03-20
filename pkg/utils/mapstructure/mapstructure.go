package mapstructure

import (
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"

	"kratos-sharp/pkg/utils/xtime"
)

func Decode(input interface{}, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:   nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(DecodeTimeHookFunc(), DecodeXTimeHookFunc()),
		Result:     &result,
		TagName:    "json",
	})
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

func DecodeTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		var tTime time.Time
		var err error
		if !reflect.DeepEqual(t, reflect.TypeOf(time.Time{})) {
			return data, nil
		}
		switch f.Kind() {
		case reflect.String:
			tXTime, err := xtime.Parse(xtime.TimeFormat, data.(string))
			if err != nil {
				return nil, err
			}
			tTime = tXTime.Time
		default:
			return data, nil
		}
		if err != nil {
			return nil, err
		}
		return tTime, nil
	}
}

func DecodeXTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		var tTime xtime.Time
		var err error
		if !reflect.DeepEqual(t, reflect.TypeOf(xtime.Time{})) {
			return data, nil
		}
		switch f.Kind() {
		case reflect.String:
			tTime, err = xtime.Parse(xtime.TimeFormat, data.(string))
		default:
			return data, nil
		}
		if err != nil {
			return nil, err
		}
		return tTime, nil
	}
}
