package cli

var options struct {
	Version       bool `short:"v" long:"version" description:"print version"`
	InputRaw      bool `short:"R" long:"raw-input" description:"read input as raw strings"`
	InputSlurp    bool `short:"s" long:"slurp" description:"read all inputs into an array"`
	InputStream   bool `long:"stream" description:"parse input in stream fashion"`
	OutputCompact bool `short:"c" long:"compact-output" description:"compact output"`
	OutputRaw     bool `short:"r" long:"raw-output" description:"output raw strings"`
	OutputJoin    bool `short:"j" long:"join-output" description:"stop printing a new line after each output"`
	OutputNul     bool `short:"0" long:"nul-output" description:"print NUL after each output"`
	OutputColor   bool `short:"C" long:"color-output" description:"colorize output even if piped"`
	OutputMono    bool `short:"M" long:"monochrome-output" description:"stop colorizing output"`
	OutputYAML    bool `long:"yaml-output" description:"output by YAML"`
	OutputIndent  *int `long:"indent" description:"number of spaces for indentation"`
	OutputTab     bool `long:"tab" description:"use tabs for indentation"`
}
