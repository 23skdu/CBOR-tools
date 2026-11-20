package main

import (
	"os"

	"github.com/23skdu/cbor-tools/internal/cli"
)

var osExit = os.Exit

// version is set at build time via -ldflags.
var version = "0.1.0"

func main() {
	osExit(run())
}

func run() int {
	cli.Version = version
	return cli.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
}
