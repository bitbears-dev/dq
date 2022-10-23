package cli

import (
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/itchyny/gojq"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

type CLI struct {
	version string
}

func NewCLI(version string) *CLI {
	return &CLI{
		version: version,
	}
}

func (c *CLI) Run(args []string) int {
	err := c.run(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		if ex, ok := err.(Exiter); ok {
			return ex.ExitCode()
		}
		return 1
	}
	return 0
}

func (c *CLI) run(args []string) error {
	queryAndInputFiles, err := flags.ParseArgs(&options, args)
	if err != nil {
		return err
	}

	if options.Version {
		fmt.Println(c.version)
		return nil
	}

	queryString := "mynow" // to avoid built-in 'now' function
	if len(queryAndInputFiles) > 0 {
		queryString = queryAndInputFiles[0]
	}
	inputFiles := []string{}
	if len(queryAndInputFiles) > 1 {
		inputFiles = queryAndInputFiles[1:]
	}

	query, err := gojq.Parse(queryString)
	if err != nil {
		return err
	}

	iter := c.createInputIter(queryString, inputFiles)
	defer iter.Close()

	code, err := gojq.Compile(query,
		gojq.WithFunction("mynow", 0, 0, c.now),
		gojq.WithFunction("from_unix", 0, 1, c.fromUnix),
		gojq.WithFunction("from_unixmilli", 0, 1, c.fromUnixMilli),
		gojq.WithFunction("from_ymd", 3, 3, c.fromYMD),
		gojq.WithFunction("from_ymd_hms", 6, 6, c.fromYMDHMS),
		gojq.WithFunction("from_rfc3339", 0, 1, c.fromRFC3339),
		gojq.WithFunction("add_day", 1, 1, c.addDay),
		gojq.WithFunction("utc", 0, 0, c.utc),
	)
	if err != nil {
		return err
	}

	return c.process(iter, code)
}

func (c *CLI) createInputIter(query string, args []string) (iter inputIter) {
	if !isStdinConnectedToPipe() && len(args) == 0 {
		return newGuessedInputIter(time.Now())
	}
	var newIter func(io.Reader, string) inputIter
	switch {
	case options.InputRaw:
		if options.InputSlurp {
			newIter = newReadAllIter
		} else {
			newIter = newRawInputIter
		}
	case options.InputStream:
		newIter = newStreamInputIter
	default:
		newIter = newJSONInputIter
	}
	if options.InputSlurp {
		defer func() {
			if options.InputRaw {
				iter = newSlurpRawInputIter(iter)
			} else {
				iter = newSlurpInputIter(iter)
			}
		}()
	}
	if len(args) == 0 {
		return newIter(os.Stdin, "<stdin>")
	}
	return newFilesInputIter(newIter, args, os.Stdin)
}

func isStdinConnectedToPipe() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		panic("unable to get stat on stdin")
	}
	return (stat.Mode() & os.ModeNamedPipe) != 0
}

func (c *CLI) process(iter inputIter, code *gojq.Code) error {
	var err error
	for {
		v, ok := iter.Next()
		if !ok {
			return err
		}
		if er, ok := v.(error); ok {
			c.printError(er)
			continue
		}
		if er := c.printValues(code.Run(v)); er != nil {
			c.printError(er)
		}
	}
}

func (c *CLI) printValues(iter gojq.Iter) error {
	m := c.createMarshaler()
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return err
		}

		if err := m.marshal(v, os.Stdout); err != nil {
			return err
		}
		if !options.OutputJoin {
			if options.OutputNul {
				os.Stdout.Write([]byte{'\x00'})
			} else {
				os.Stdout.Write([]byte{'\n'})
			}
		}
	}
	return nil
}

func (c *CLI) createMarshaler() marshaler {
	indent := 2
	if options.OutputCompact {
		indent = 0
	} else if options.OutputTab {
		indent = 1
	} else if i := options.OutputIndent; i != nil {
		indent = *i
	}
	f := newEncoder(options.OutputTab, indent)
	if options.OutputRaw || options.OutputJoin || options.OutputNul {
		return &rawMarshaler{f}
	}
	return f
}

func (c *CLI) printError(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
}

func (c *CLI) now(v interface{}, _ []interface{}) interface{} {
	return encap(time.Now())
}

func (c *CLI) fromUnix(v interface{}, args []interface{}) interface{} {
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

	return encap(time.Unix(int64(u), 0))
}

func (c *CLI) fromUnixMilli(v interface{}, args []interface{}) interface{} {
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
	return encap(time.Unix(0, int64(u)*1000000))
}

func (c *CLI) fromYMD(_ interface{}, args []interface{}) interface{} {
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

	return encap(time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local))
}

func (c *CLI) fromYMDHMS(_ interface{}, args []interface{}) interface{} {
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

	return encap(time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local))
}

func (c *CLI) fromRFC3339(v interface{}, args []interface{}) interface{} {
	if len(args) == 1 {
		v = args[0]
	}

	if s, ok := v.(string); ok {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return errors.New("unable to parse as rfc3339")
		}
		return encap(t)
	}

	return errors.Errorf("unexpected type: %T", v)
}

func (c *CLI) addDay(v interface{}, args []interface{}) interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		source, ok := m["_source"]
		if !ok {
			return errors.New("_source not found")
		}
		t, ok := source.(time.Time)
		if !ok {
			return errors.New("unexpected _source type")
		}

		if len(args) < 1 {
			return errors.New("insufficient arguments")
		}
		switch d := args[0].(type) {
		case int:
			return encap(t.Add(time.Duration(d) * 24 * time.Hour))
		default:
			return errors.Errorf("unexpected argument type: %T", d)
		}
	}

	return errors.Errorf("unexpected type: %T", v)
}

func (c *CLI) utc(v interface{}, _ []interface{}) interface{} {
	if m, ok := v.(map[string]interface{}); ok {
		source, ok := m["_source"]
		if !ok {
			return errors.New("_source not found")
		}
		t, ok := source.(time.Time)
		if !ok {
			return errors.New("unexpected _source type")
		}
		return encap(t.UTC())
	}

	return errors.New("unexpected type")
}

func encap(t time.Time) map[string]interface{} {
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
