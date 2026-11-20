package cli

import (
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/23skdu/cbor-tools/internal/cborx"
)

// Run parses args and executes the requested action, writing to the provided
// streams. It returns a process exit code.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("cbor-tool", flag.ContinueOnError)
	fs.SetOutput(stderr)

	action := fs.String("a", "", "action to perform: Encode, Decode, or Tag")
	file := fs.String("f", "", "input file path (default: stdin)")
	str := fs.String("s", "", "input string (for Encode)")
	tag := fs.String("t", "", "CBOR tag number (for Tag)")
	out := fs.String("o", "", "output file path (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *action == "" {
		_, _ = fmt.Fprintln(stderr, "ERR: Must specify action(-a)")
		_, _ = fmt.Fprintln(stderr, "Actions: Encode Decode Tag")
		return 1
	}
	a := strings.TrimSpace(*action)
	if a != "Encode" && a != "Decode" && a != "Tag" {
		_, _ = fmt.Fprintln(stderr, "ERR: Must specify a valid action")
		_, _ = fmt.Fprintln(stderr, "Actions: Encode Decode Tag")
		return 1
	}

	rc := 1
	switch a {
	case "Encode":
		input, err := cborx.ReadStringInput(*str, *file, stdin)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "Error: %v\n", err)
			return 1
		}
		data, _ := cborx.EncodeString(input)
		if err := cborx.WriteOutput(*out, data, stdout); err != nil {
			_, _ = fmt.Fprintf(stderr, "Error: %v\n", err)
			return 1
		}
		rc = 0

	case "Decode":
		data, err := cborx.ReadInput(*file, stdin)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "Error: %v\n", err)
			return 1
		}
		jsonOut, err := cborx.DecodeJSON(data)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "Error: %v\n", err)
			return 1
		}
		if err := cborx.WriteOutput(*out, jsonOut, stdout); err != nil {
			_, _ = fmt.Fprintf(stderr, "Error: %v\n", err)
			return 1
		}
		rc = 0

	case "Tag":
		if *tag == "" {
			_, _ = fmt.Fprintln(stderr, "Error: tag number is required (-t)")
			return 1
		}
		n, err := strconv.ParseUint(*tag, 10, 64)
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "Error: invalid tag number: %v\n", err)
			return 1
		}
		var data []byte
		if *str != "" {
			out, _ := cborx.TagString(*str, n)
			data = out
		} else {
			input, err := cborx.ReadInput(*file, stdin)
			if err != nil {
				_, _ = fmt.Fprintf(stderr, "Error: %v\n", err)
				return 1
			}
			out, err := cborx.TagBytes(input, n)
			if err != nil {
				_, _ = fmt.Fprintf(stderr, "Error: %v\n", err)
				return 1
			}
			data = out
		}
		if err := cborx.WriteOutput(*out, data, stdout); err != nil {
			_, _ = fmt.Fprintf(stderr, "Error: %v\n", err)
			return 1
		}
		rc = 0
	}

	return rc
}
