package cli

type Exiter interface {
	ExitCode() int
}
