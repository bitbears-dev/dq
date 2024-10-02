package main

import (
	"os"

	"github.com/bitbears-dev/dq/cli"
)

var version string

func main() {
	os.Exit(cli.NewCLI(version).Run(os.Args[1:]))
}
