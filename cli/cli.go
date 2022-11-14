package cli

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bitbears-dev/dq/builtin"
	"github.com/itchyny/gojq"
	"github.com/jessevdk/go-flags"
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

	queryString := "."
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
		gojq.WithFunction("fromunix", 0, 1, builtin.FromUnix),
		gojq.WithFunction("from_unix", 0, 1, builtin.FromUnix),
		gojq.WithFunction("fromunixmilli", 0, 1, builtin.FromUnixMilli),
		gojq.WithFunction("from_unixmilli", 0, 1, builtin.FromUnixMilli),
		gojq.WithFunction("fromunixmicro", 0, 1, builtin.FromUnixMicro),
		gojq.WithFunction("from_unixmicro", 0, 1, builtin.FromUnixMicro),
		gojq.WithFunction("fromunixnano", 0, 1, builtin.FromUnixNano),
		gojq.WithFunction("from_unixnano", 0, 1, builtin.FromUnixNano),
		gojq.WithFunction("fromymd", 3, 3, builtin.FromYMD),
		gojq.WithFunction("from_ymd", 3, 3, builtin.FromYMD),
		gojq.WithFunction("fromymdhms", 6, 6, builtin.FromYMDHMS),
		gojq.WithFunction("from_ymdhms", 6, 6, builtin.FromYMDHMS),
		gojq.WithFunction("fromansic", 0, 1, builtin.FromKnownTimeFormat(time.ANSIC)),
		gojq.WithFunction("from_ansic", 0, 1, builtin.FromKnownTimeFormat(time.ANSIC)),
		gojq.WithFunction("toansic", 0, 1, builtin.ToKnownTimeFormat(time.ANSIC)),
		gojq.WithFunction("to_ansic", 0, 1, builtin.ToKnownTimeFormat(time.ANSIC)),
		gojq.WithFunction("fromunixdate", 0, 1, builtin.FromKnownTimeFormat(time.UnixDate)),
		gojq.WithFunction("from_unixdate", 0, 1, builtin.FromKnownTimeFormat(time.UnixDate)),
		gojq.WithFunction("tounixdate", 0, 1, builtin.ToKnownTimeFormat(time.UnixDate)),
		gojq.WithFunction("to_unixdate", 0, 1, builtin.ToKnownTimeFormat(time.UnixDate)),
		gojq.WithFunction("fromrubydate", 0, 1, builtin.FromKnownTimeFormat(time.RubyDate)),
		gojq.WithFunction("from_rubydate", 0, 1, builtin.FromKnownTimeFormat(time.RubyDate)),
		gojq.WithFunction("torubydate", 0, 1, builtin.ToKnownTimeFormat(time.RubyDate)),
		gojq.WithFunction("to_rubydate", 0, 1, builtin.ToKnownTimeFormat(time.RubyDate)),
		gojq.WithFunction("fromrfc822", 0, 1, builtin.FromKnownTimeFormat(time.RFC822)),
		gojq.WithFunction("from_rfc822", 0, 1, builtin.FromKnownTimeFormat(time.RFC822)),
		gojq.WithFunction("torfc822", 0, 1, builtin.ToKnownTimeFormat(time.RFC822)),
		gojq.WithFunction("to_rfc822", 0, 1, builtin.ToKnownTimeFormat(time.RFC822)),
		gojq.WithFunction("fromrfc822z", 0, 1, builtin.FromKnownTimeFormat(time.RFC822Z)),
		gojq.WithFunction("from_rfc822z", 0, 1, builtin.FromKnownTimeFormat(time.RFC822Z)),
		gojq.WithFunction("torfc822z", 0, 1, builtin.ToKnownTimeFormat(time.RFC822Z)),
		gojq.WithFunction("fromrfc850", 0, 1, builtin.FromKnownTimeFormat(time.RFC850)),
		gojq.WithFunction("from_rfc850", 0, 1, builtin.FromKnownTimeFormat(time.RFC850)),
		gojq.WithFunction("torfc850", 0, 1, builtin.ToKnownTimeFormat(time.RFC850)),
		gojq.WithFunction("to_rfc850", 0, 1, builtin.ToKnownTimeFormat(time.RFC850)),
		gojq.WithFunction("fromrfc1123", 0, 1, builtin.FromKnownTimeFormat(time.RFC1123)),
		gojq.WithFunction("from_rfc1123", 0, 1, builtin.FromKnownTimeFormat(time.RFC1123)),
		gojq.WithFunction("torfc1123", 0, 1, builtin.ToKnownTimeFormat(time.RFC1123)),
		gojq.WithFunction("to_rfc1123", 0, 1, builtin.ToKnownTimeFormat(time.RFC1123)),
		gojq.WithFunction("fromrfc1123z", 0, 1, builtin.FromKnownTimeFormat(time.RFC1123Z)),
		gojq.WithFunction("from_rfc1123z", 0, 1, builtin.FromKnownTimeFormat(time.RFC1123Z)),
		gojq.WithFunction("torfc1123z", 0, 1, builtin.ToKnownTimeFormat(time.RFC1123Z)),
		gojq.WithFunction("to_rfc1123z", 0, 1, builtin.ToKnownTimeFormat(time.RFC1123Z)),
		gojq.WithFunction("fromrfc3339", 0, 1, builtin.FromKnownTimeFormat(time.RFC3339)),
		gojq.WithFunction("from_rfc3339", 0, 1, builtin.FromKnownTimeFormat(time.RFC3339)),
		gojq.WithFunction("torfc3339", 0, 1, builtin.ToKnownTimeFormat(time.RFC3339)),
		gojq.WithFunction("to_rfc3339", 0, 1, builtin.ToKnownTimeFormat(time.RFC3339)),
		gojq.WithFunction("fromrfc3339nano", 0, 1, builtin.FromKnownTimeFormat(time.RFC3339Nano)),
		gojq.WithFunction("from_rfc3339nano", 0, 1, builtin.FromKnownTimeFormat(time.RFC3339Nano)),
		gojq.WithFunction("torfc3339nano", 0, 1, builtin.ToKnownTimeFormat(time.RFC3339Nano)),
		gojq.WithFunction("to_rfc3339nano", 0, 1, builtin.ToKnownTimeFormat(time.RFC3339Nano)),
		gojq.WithFunction("fromkitchen", 0, 1, builtin.FromKnownTimeFormat(time.Kitchen)),
		gojq.WithFunction("from_kitchen", 0, 1, builtin.FromKnownTimeFormat(time.Kitchen)),
		gojq.WithFunction("tokitchen", 0, 1, builtin.ToKnownTimeFormat(time.Kitchen)),
		gojq.WithFunction("to_kitchen", 0, 1, builtin.ToKnownTimeFormat(time.Kitchen)),
		gojq.WithFunction("fromstamp", 0, 1, builtin.FromKnownTimeFormat(time.Stamp)),
		gojq.WithFunction("from_stamp", 0, 1, builtin.FromKnownTimeFormat(time.Stamp)),
		gojq.WithFunction("tostamp", 0, 1, builtin.ToKnownTimeFormat(time.Stamp)),
		gojq.WithFunction("to_stamp", 0, 1, builtin.ToKnownTimeFormat(time.Stamp)),
		gojq.WithFunction("fromstampmilli", 0, 1, builtin.FromKnownTimeFormat(time.StampMilli)),
		gojq.WithFunction("from_stampmilli", 0, 1, builtin.FromKnownTimeFormat(time.StampMilli)),
		gojq.WithFunction("tostampmilli", 0, 1, builtin.ToKnownTimeFormat(time.StampMilli)),
		gojq.WithFunction("to_stampmilli", 0, 1, builtin.ToKnownTimeFormat(time.StampMilli)),
		gojq.WithFunction("fromstampmicro", 0, 1, builtin.FromKnownTimeFormat(time.StampMicro)),
		gojq.WithFunction("from_stampmicro", 0, 1, builtin.FromKnownTimeFormat(time.StampMicro)),
		gojq.WithFunction("tostampmicro", 0, 1, builtin.ToKnownTimeFormat(time.StampMicro)),
		gojq.WithFunction("to_stampmicro", 0, 1, builtin.ToKnownTimeFormat(time.StampMicro)),
		gojq.WithFunction("fromstampnano", 0, 1, builtin.FromKnownTimeFormat(time.StampNano)),
		gojq.WithFunction("from_stampnano", 0, 1, builtin.FromKnownTimeFormat(time.StampNano)),
		gojq.WithFunction("tostampnano", 0, 1, builtin.ToKnownTimeFormat(time.StampNano)),
		gojq.WithFunction("to_stampnano", 0, 1, builtin.ToKnownTimeFormat(time.StampNano)),
		gojq.WithFunction("add_date", 3, 3, builtin.AddDate),
		gojq.WithFunction("clock", 0, 0, builtin.Clock),
		gojq.WithFunction("date", 0, 0, builtin.Date),
		gojq.WithFunction("utc", 0, 0, builtin.UTC),
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
