package metadata

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Duration struct {
	time.Duration
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)

		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}

		return nil
	default:
		return errors.New("invalid duration")
	}
}

func toTimeDurationHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if t != reflect.TypeOf(Duration{}) && t != reflect.TypeOf(time.Duration(0)) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.TypeOf(time.Duration(0)).Kind():
			return data.(time.Duration), nil
		case reflect.String:
			var val time.Duration
			if data.(string) != "" {
				var err error
				val, err = time.ParseDuration(data.(string))
				if err != nil {
					// If we can't parse the duration, try parsing it as int64 seconds
					seconds, errParse := strconv.ParseInt(data.(string), 10, 0)
					if errParse != nil {
						return nil, errors.Join(err, errParse)
					}
					val = time.Duration(seconds * int64(time.Second))
				}
			}
			if t != reflect.TypeOf(Duration{}) {
				return val, nil
			}
			return Duration{Duration: val}, nil
		case reflect.Float64:
			val := time.Duration(data.(float64))
			if t != reflect.TypeOf(Duration{}) {
				return val, nil
			}
			return Duration{Duration: val}, nil
		case reflect.Int64:
			val := time.Duration(data.(int64))
			if t != reflect.TypeOf(Duration{}) {
				return val, nil
			}
			return Duration{Duration: val}, nil
		default:
			return data, nil
		}
	}
}
