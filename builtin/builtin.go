package builtin

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func FromUnix(v interface{}, args []interface{}) interface{} {
	if len(args) == 1 {
		v = args[0]
	}

	var u int64
	switch x := v.(type) {
	case int:
		u = int64(x)
	case float64:
		t, err := fromUnixF(x)
		if err != nil {
			return err
		}
		u = (*t).Unix()
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
	if len(args) == 1 {
		v = args[0]
	}

	var u int64
	switch x := v.(type) {
	case int:
		u = int64(x)
	case float64:
		t, err := fromUnixF(x)
		if err != nil {
			return err
		}
		u = (*t).UnixMilli()
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

func FromUnixMicro(v interface{}, args []interface{}) interface{} {
	if len(args) == 1 {
		v = args[0]
	}

	var u int64
	switch x := v.(type) {
	case int:
		u = int64(x)
	case float64:
		t, err := fromUnixF(x)
		if err != nil {
			return err
		}
		u = (*t).UnixMicro()
	case string:
		var err error
		u, err = strconv.ParseInt(x, 10, 64)
		if err != nil {
			return err
		}
	default:
		return errors.Errorf("unexpected type: %T", v)
	}
	return EncapTime(time.Unix(0, int64(u)*1000))
}

func FromUnixNano(v interface{}, args []interface{}) interface{} {
	if len(args) == 1 {
		v = args[0]
	}

	var u int64
	switch x := v.(type) {
	case int:
		u = int64(x)
	case float64:
		t, err := fromUnixF(x)
		if err != nil {
			return err
		}
		u = (*t).UnixNano()
	case string:
		var err error
		u, err = strconv.ParseInt(x, 10, 64)
		if err != nil {
			return err
		}
	default:
		return errors.Errorf("unexpected type: %T", v)
	}
	return EncapTime(time.Unix(0, int64(u)))
}

func fromUnixF(f float64) (*time.Time, error) {
	s := fmt.Sprintf("%.9f", f)
	dotPos := strings.Index(s, ".")
	intPart := int64(math.Floor(f))
	decPart, err := strconv.ParseInt(s[dotPos+1:dotPos+10], 10, 64)
	if err != nil {
		return nil, err
	}
	//log.Printf("s = %s, intPart = %d, decPart = %d", s, intPart, decPart)
	t := time.Unix(intPart, decPart)
	return &t, nil
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

type BuiltinFn func(interface{}, []interface{}) interface{}

func FromKnownTimeFormat(layout string) BuiltinFn {
	return func(v interface{}, args []interface{}) interface{} {
		if s, ok := getStringArg(v, args); ok {
			return timeFromString(layout, s)
		}
		return errors.Errorf("unexpected type: %T", v)
	}
}

func ToKnownTimeFormat(layout string) BuiltinFn {
	return func(v interface{}, args []interface{}) interface{} {
		if t, ok := getTimeArg(v, args); ok {
			return t.Format(layout)
		}
		return errors.Errorf("unexpected type: %T", v)
	}
}

func getStringArg(v interface{}, args []interface{}) (string, bool) {
	if len(args) == 1 {
		v = args[0]
	}
	s, ok := v.(string)
	return s, ok
}

func getTimeArg(v interface{}, args []interface{}) (*time.Time, bool) {
	if len(args) == 1 {
		v = args[0]
	}

	return DecapTime(v)
}

func timeFromString(layout, value string) interface{} {
	t, err := time.Parse(layout, value)
	if err != nil {
		return errors.New("unable to parse using the specified format")
	}
	return EncapTime(t)
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
		"__dq__source":    t,
		"unixNano":        t.UnixNano(),
		"unixNanoString":  fmt.Sprintf("%d", t.UnixNano()),
		"unixMicro":       t.UnixMicro(),
		"unixMicroString": fmt.Sprintf("%d", t.UnixMicro()),
		"unixMilli":       t.UnixMilli(),
		"unixMilliString": fmt.Sprintf("%d", t.UnixMilli()),
		"unix":            t.Unix(),
		"unixString":      fmt.Sprintf("%d", t.Unix()),
		"year":            t.Year(),
		"month":           int(t.Month()),
		"day":             t.Day(),
		"hour":            t.Hour(),
		"hour12":          t.Hour() % 12,
		"minute":          t.Minute(),
		"second":          t.Second(),
		"millisecond":     t.Nanosecond() / 1000000,
		"microsecond":     t.Nanosecond() / 1000,
		"nanosecond":      t.Nanosecond(),
		"am":              t.Hour() < 12,
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

	source, ok := m["__dq__source"]
	if !ok {
		return nil, false
	}

	t, ok := source.(time.Time)
	if !ok {
		return nil, false
	}

	return &t, true
}
