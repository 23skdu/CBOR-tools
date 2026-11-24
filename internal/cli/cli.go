package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/23skdu/cbor-tools/internal/cborx"
)

// Version is the build version. It may be overridden by the main package.
var Version = "dev"

// ErrUnknownShell is returned when an unsupported shell is requested.
var ErrUnknownShell = errors.New("unsupported shell")

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
	format := fs.String("format", "pretty", "output format for Decode: pretty or json")
	root := fs.String("root", "", "confine file access to this directory")
	completion := fs.String("completion", "", "print a shell completion script: bash, zsh, or fish")
	man := fs.Bool("man", false, "print the manual page and exit")
	quiet := fs.Bool("q", false, "suppress diagnostics on stderr")
	quietLong := fs.Bool("quiet", false, "suppress diagnostics on stderr (same as -q)")
	version := fs.Bool("version", false, "print version and exit")
	versionShort := fs.Bool("v", false, "print version and exit (shorthand)")

	if err := fs.Parse(args); err != nil {
		return 2
	}
	silent := *quiet || *quietLong

	if *version || *versionShort {
		_, _ = fmt.Fprintf(stdout, "cbor-tool %s\n", Version)
		return 0
	}

	// Open the sandbox before any file access so -f, -o, -completion and -man
	// are all confined to it.
	var sandbox *os.Root
	if *root != "" {
		r, err := os.OpenRoot(*root)
		if err != nil {
			diag(stderr, silent, "Error: %v\n", err)
			return 1
		}
		sandbox = r
		defer func() { _ = r.Close() }()
	}

	if *completion != "" {
		script, err := CompletionScript(*completion)
		if err != nil {
			diag(stderr, silent, "Error: %v\n", err)
			return 1
		}
		return emit(sandbox, *out, []byte(script), stdout, stderr, silent)
	}
	if *man {
		return emit(sandbox, *out, []byte(ManPage()), stdout, stderr, silent)
	}

	if *action == "" {
		diag(stderr, silent, "ERR: Must specify action(-a)\n")
		diag(stderr, silent, "Actions: Encode Decode Tag\n")
		return 1
	}
	a := strings.TrimSpace(*action)

	switch a {
	case "Encode":
		input, err := cborx.ReadStringInput(sandbox, *str, *file, stdin)
		if err != nil {
			diag(stderr, silent, "Error: %v\n", err)
			return 1
		}
		data, _ := cborx.EncodeString(input)
		return emit(sandbox, *out, data, stdout, stderr, silent)

	case "Decode":
		if *format != "pretty" && *format != "json" {
			diag(stderr, silent, "Error: invalid format: %s\n", *format)
			return 1
		}
		data, err := cborx.ReadInput(sandbox, *file, stdin)
		if err != nil {
			diag(stderr, silent, "Error: %v\n", err)
			return 1
		}
		var jsonOut []byte
		if *format == "json" {
			jsonOut, err = cborx.DecodeCompactJSON(data)
		} else {
			jsonOut, err = cborx.DecodeJSON(data)
		}
		if err != nil {
			diag(stderr, silent, "Error: %v\n", err)
			return 1
		}
		return emit(sandbox, *out, jsonOut, stdout, stderr, silent)

	case "Tag":
		if *tag == "" {
			diag(stderr, silent, "Error: tag number is required (-t)\n")
			return 1
		}
		n, err := strconv.ParseUint(*tag, 10, 64)
		if err != nil {
			diag(stderr, silent, "Error: invalid tag number: %v\n", err)
			return 1
		}
		var data []byte
		if *str != "" {
			data, _ = cborx.TagString(*str, n)
		} else {
			input, err := cborx.ReadInput(sandbox, *file, stdin)
			if err != nil {
				diag(stderr, silent, "Error: %v\n", err)
				return 1
			}
			data, err = cborx.TagBytes(input, n)
			if err != nil {
				diag(stderr, silent, "Error: %v\n", err)
				return 1
			}
		}
		return emit(sandbox, *out, data, stdout, stderr, silent)

	default:
		diag(stderr, silent, "ERR: Must specify a valid action\n")
		diag(stderr, silent, "Actions: Encode Decode Tag\n")
		return 1
	}
}

// emit writes data to path (or to stdout when path is empty) and reports the
// write on stderr unless silent mode is enabled. It returns a process exit
// code.
func emit(root *os.Root, path string, data []byte, stdout, stderr io.Writer, silent bool) int {
	if err := cborx.WriteOutput(root, path, data, stdout); err != nil {
		diag(stderr, silent, "Error: %v\n", err)
		return 1
	}
	if path != "" && !silent {
		_, _ = fmt.Fprintf(stderr, "wrote %d bytes to %s\n", len(data), path)
	}
	return 0
}

// diag writes a diagnostic line to w unless silent mode is enabled.
func diag(w io.Writer, silent bool, format string, a ...any) {
	if silent {
		return
	}
	_, _ = fmt.Fprintf(w, format, a...)
}
