package main

import (
	"os"

	"github.com/23skdu/cbor-tools/internal/cli"
)

var osExit = os.Exit

func main() {
	osExit(run())
}

func run() int {
	return cli.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)
}