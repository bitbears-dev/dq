package builtin

import (
	"log"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

func Now(_ interface{}, _ []interface{}) interface{} {
	return EncapTime(time.Now())
}

func FromUnix(v interface{}, args []interface{}) interface{} {
	var u int64
	if len(args) == 1 {
		v = args[0]
	}
	switch x := v.(type) {
	case int:
		u = int64(x)
	case string:
		var err error
		u, err = strconv.ParseInt(x, 10, 64)
		if err != nil {
			return err
		}
	default:
		return errors.Errorf("unexpected type: %T", v)
	}

	return EncapTime(time.Unix(int64(u), 0))
}

func FromUnixMilli(v interface{}, args []interface{}) interface{} {
	var u int64
	if len(args) == 1 {
		v = args[0]
	}

	switch x := v.(type) {
	case int:
		u = int64(x)
	case string:
		var err error
		u, err = strconv.ParseInt(x, 10, 64)
		if err != nil {
			return err
		}
	default:
		return errors.Errorf("unexpected type: %T", v)
	}
	return EncapTime(time.Unix(0, int64(u)*1000000))
}

func FromYMD(_ interface{}, args []interface{}) interface{} {
	if len(args) < 3 {
		log.Printf("args: %v", args)
		return errors.New("insufficient arguments")
	}

	year, ok := args[0].(int)
	if !ok {
		return errors.New("unexpected argument type")
	}

	month, ok := args[1].(int)
	if !ok {
		return errors.New("unexpected argument type")
	}

	day, ok := args[2].(int)
	if !ok {
		return errors.New("unexpected argument type")
	}

	return EncapTime(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local))
}

func FromYMDHMS(_ interface{}, args []interface{}) interface{} {
	if len(args) < 6 {
		return errors.New("insufficient arguments")
	}

	year, ok := args[0].(int)
	if !ok {
		return errors.New("unexpected argument type")
	}

	month, ok := args[1].(int)
	if !ok {
		return errors.New("unexpected argument type")
	}

	day, ok := args[2].(int)
	if !ok {
		return errors.New("unexpected argument type")
	}

	hour, ok := args[3].(int)
	if !ok {
		return errors.New("unexpected argument type")
	}

	minute, ok := args[4].(int)
	if !ok {
		return errors.New("unexpected argument type")
	}

	second, ok := args[5].(int)
	if !ok {
		return errors.New("unexpected argument type")
	}

	return EncapTime(time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local))
}

func FromRFC3339(v interface{}, args []interface{}) interface{} {
	if len(args) == 1 {
		v = args[0]
	}

	if s, ok := v.(string); ok {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return errors.New("unable to parse as rfc3339")
		}
		return EncapTime(t)
	}

	return errors.Errorf("unexpected type: %T", v)
}

func ToRFC3339(v interface{}, args []interface{}) interface{} {
	if len(args) == 1 {
		v = args[0]
	}

	if t, ok := DecapTime(v); ok {
		return t.Format(time.RFC3339)
	}

	return errors.Errorf("unexpected type: %T", v)
}

func AddDay(v interface{}, args []interface{}) interface{} {
	t, ok := DecapTime(v)
	if !ok {
		return errors.Errorf("unexpected type: %T", v)
	}
	if len(args) < 1 {
		return errors.New("insufficient arguments")
	}
	switch d := args[0].(type) {
	case int:
		return EncapTime(t.Add(time.Duration(d) * 24 * time.Hour))
	default:
		return errors.Errorf("unexpected argument type: %T", d)
	}
}

func UTC(v interface{}, _ []interface{}) interface{} {
	t, ok := DecapTime(v)
	if ok {
		return EncapTime(t.UTC())
	}

	return errors.New("unexpected type")
}

func EncapTime(t time.Time) map[string]interface{} {
	zoneName, offset := t.Zone()
	return map[string]interface{}{
		"_source":     t,
		"unixNano":    t.UnixNano(),
		"unixMicro":   t.UnixMicro(),
		"unixMilli":   t.UnixMilli(),
		"unix":        t.Unix(),
		"year":        t.Year(),
		"month":       int(t.Month()),
		"day":         t.Day(),
		"hour":        t.Hour(),
		"hour12":      t.Hour() % 12,
		"minute":      t.Minute(),
		"second":      t.Second(),
		"millisecond": t.Nanosecond() / 1000000,
		"microsecond": t.Nanosecond() / 1000,
		"nanosecond":  t.Nanosecond(),
		"am":          t.Hour() < 12,
		"timezone": map[string]interface{}{
			"short":         zoneName,
			"offsetSeconds": offset,
		},
		"weekday": map[string]interface{}{
			"name": t.Weekday().String(),
		},
		"dayOfYear": t.YearDay(),
	}
}

func DecapTime(v interface{}) (*time.Time, bool) {
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, false
	}

	source, ok := m["_source"]
	if !ok {
		return nil, false
	}

	t, ok := source.(time.Time)
	if !ok {
		return nil, false
	}

	return &t, true
}
