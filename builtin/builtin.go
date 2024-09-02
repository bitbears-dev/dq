package builtin

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var supportedKnownTimeFormats = []string{
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	time.Kitchen,
	time.Stamp,
	time.StampMilli,
	time.StampMicro,
	time.StampNano,
}

func Guess(v interface{}, args []interface{}) interface{} {
	if len(args) == 1 {
		v = args[0]
	}

	if isLikelyUnix(v) {
		return FromUnix(v, args)
	}

	if isLikelyUnixMilli(v) {
		return FromUnixMilli(v, args)
	}

	if isLikelyUnixMicro(v) {
		return FromUnixMicro(v, args)
	}

	if isLikelyUnixNano(v) {
		return FromUnixNano(v, args)
	}

	if s, ok := v.(string); ok {
		for _, f := range supportedKnownTimeFormats {
			t := timeFromString(f, s)
			if _, ok := t.(error); !ok {
				return t
			}
		}
	}

	return errors.New("unable to guess")
}

var reAllDigits = regexp.MustCompile("^[[:digit:]]+$")

func isLikelyUnix(v interface{}) bool {
	lenNow := len(fmt.Sprintf("%d", time.Now().Unix()))
	switch x := v.(type) {
	case int:
		if lenNow == len(fmt.Sprintf("%d", x)) {
			return true
		}
	case string:
		if lenNow == len(x) && reAllDigits.MatchString(x) {
			return true
		}
	}

	return false
}

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

func isLikelyUnixMilli(v interface{}) bool {
	lenNow := len(fmt.Sprintf("%d", time.Now().UnixMilli()))
	switch x := v.(type) {
	case int:
		if lenNow == len(fmt.Sprintf("%d", x)) { // diff is smaller than 100 years from now
			return true
		}
	case string:
		if lenNow == len(x) && reAllDigits.MatchString(x) {
			return true
		}
	}

	return false
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

func isLikelyUnixMicro(v interface{}) bool {
	lenNow := len(fmt.Sprintf("%d", time.Now().UnixMicro()))
	switch x := v.(type) {
	case int:
		if lenNow == len(fmt.Sprintf("%d", x)) { // diff is smaller than 100 years from now
			return true
		}
	case string:
		if lenNow == len(x) && reAllDigits.MatchString(x) {
			return true
		}
	}

	return false
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

func isLikelyUnixNano(v interface{}) bool {
	lenNow := len(fmt.Sprintf("%d", time.Now().UnixNano()))
	switch x := v.(type) {
	case int:
		if lenNow == len(fmt.Sprintf("%d", x)) { // diff is smaller than 100 years from now
			return true
		}
	case float64:
		return true
	case string:
		if lenNow == len(x) && reAllDigits.MatchString(x) {
			return true
		}
	}

	return false
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

func AddDate(v interface{}, args []interface{}) interface{} {
	t, ok := DecapTime(v)
	if !ok {
		return errors.Errorf("expected time, but found unexpected type: %T", v)
	}
	if len(args) < 3 {
		return errors.New("insufficient arguments")
	}

	years, err := interpretAsInt(args[0])
	if err != nil {
		return err
	}
	months, err := interpretAsInt(args[1])
	if err != nil {
		return err
	}
	days, err := interpretAsInt(args[2])
	if err != nil {
		return err
	}

	return EncapTime(t.AddDate(years, months, days))
}

func Add(v interface{}, args []interface{}) interface{} {
	t, ok := DecapTime(v)
	if !ok {
		return errors.Errorf("expected time as input, but found unexpected type: %T", v)
	}
	if len(args) < 1 {
		return errors.New("insufficient arguments")
	}

	d, ok := DecapDuration(args[0])
	if !ok {
		return errors.Errorf("expected duration as the first argument, but found unexpected type: %T", args[0])
	}

	return EncapTime(t.Add(*d))
}

func interpretAsInt(arg interface{}) (int, error) {
	switch d := arg.(type) {
	case int:
		return d, nil
	default:
		return 0, errors.Errorf("unexpected argument type: %T", d)
	}
}

func Clock(v interface{}, _ []interface{}) interface{} {
	t, ok := DecapTime(v)
	if !ok {
		return errors.Errorf("expected time, but found unexpected type: %T", v)
	}
	h, m, s := t.Clock()
	return []interface{}{h, m, s}
}

func Date(v interface{}, _ []interface{}) interface{} {
	t, ok := DecapTime(v)
	if !ok {
		return errors.Errorf("expected time, but found unexpected type: %T", v)
	}
	y, m, d := t.Date()
	return []interface{}{y, int(m), d}
}

func UTC(v interface{}, _ []interface{}) interface{} {
	t, ok := DecapTime(v)
	if ok {
		return EncapTime(t.UTC())
	}

	return errors.New("unexpected type")
}

func Local(v interface{}, _ []interface{}) interface{} {
	t, ok := DecapTime(v)
	if ok {
		return EncapTime(t.Local())
	}

	return errors.New("unexpected type")
}

func Hours(v interface{}, args []interface{}) interface{} {
	if len(args) == 1 {
		v = args[0]
	}
	return convertToDuration(v, time.Hour)
}

func Minutes(v interface{}, args []interface{}) interface{} {
	if len(args) == 1 {
		v = args[0]
	}
	return convertToDuration(v, time.Minute)
}

func Seconds(v interface{}, args []interface{}) interface{} {
	if len(args) == 1 {
		v = args[0]
	}
	return convertToDuration(v, time.Second)
}

func Milliseconds(v interface{}, args []interface{}) interface{} {
	if len(args) == 1 {
		v = args[0]
	}
	return convertToDuration(v, time.Millisecond)
}

func Microseconds(v interface{}, args []interface{}) interface{} {
	if len(args) == 1 {
		v = args[0]
	}
	return convertToDuration(v, time.Microsecond)
}

func Nanoseconds(v interface{}, args []interface{}) interface{} {
	if len(args) == 1 {
		v = args[0]
	}
	return convertToDuration(v, time.Nanosecond)
}

func Today(_ interface{}, _ []interface{}) interface{} {
	t := time.Now()
	return EncapTime(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local))
}

func TodayUTC(_ interface{}, _ []interface{}) interface{} {
	t := time.Now().UTC()
	return EncapTime(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC))
}

func Yesterday(_ interface{}, _ []interface{}) interface{} {
	t := time.Now().AddDate(0, 0, -1)
	return EncapTime(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local))
}

func YesterdayUTC(_ interface{}, _ []interface{}) interface{} {
	t := time.Now().UTC().AddDate(0, 0, -1)
	return EncapTime(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC))
}

func Tomorrow(_ interface{}, _ []interface{}) interface{} {
	t := time.Now().AddDate(0, 0, 1)
	return EncapTime(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local))
}

func TomorrowUTC(_ interface{}, _ []interface{}) interface{} {
	t := time.Now().UTC().AddDate(0, 0, 1)
	return EncapTime(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC))
}

func convertToDuration(v interface{}, unit time.Duration) interface{} {
	x, ok := v.(int)
	if !ok {
		return errors.Errorf("expected integer, but found unexpected type: %T", v)
	}

	return EncapDuration(time.Duration(x) * unit)
}

func EncapTime(t time.Time) map[string]interface{} {
	zoneName, offset := t.Zone()
	year := t.Year()
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
		"year":            year,
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
			"dst":           t.IsDST(),
		},
		"weekday": map[string]interface{}{
			"name": t.Weekday().String(),
		},
		"dayOfYear":   t.YearDay(),
		"daysInMonth": getDaysInMonth(t),
		"rfc3339":     t.Format(time.RFC3339),
		"leapYear":    year%4 == 0 && (year%100 != 0 || year%400 == 0),
	}
}

func getDaysInMonth(t time.Time) int {
	// https://brandur.org/fragments/go-days-in-month
	// > The reason it works is that we generate a date one month on from the target one (m+1),
	// > but set the day of month to 0. Days are 1-indexed, so this has the effect of rolling back
	// > one day to the last day of the previous month (our target month of m)
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
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

func EncapDuration(d time.Duration) map[string]interface{} {
	return map[string]interface{}{
		"__dq__source": d,
		"hours":        d.Hours(),
		"minutes":      d.Minutes(),
		"seconds":      d.Seconds(),
		"milliseconds": d.Milliseconds(),
		"microseconds": d.Microseconds(),
		"nanoseconds":  d.Nanoseconds(),
	}
}

func DecapDuration(v interface{}) (*time.Duration, bool) {
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil, false
	}

	source, ok := m["__dq__source"]
	if !ok {
		return nil, false
	}

	d, ok := source.(time.Duration)
	if !ok {
		return nil, false
	}

	return &d, true
}
