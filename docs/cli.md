# CBOR tool CLI reference

`cbor-tool` is a small command line tool for encoding, decoding, and tagging
CBOR (Concise Binary Object Representation, RFC 8949) data.

## Usage

```bash
cbor-tool -a <action> [options]
cbor-tool -completion <shell>
cbor-tool -man
```

Actions:

- `Encode` - encode an input string as CBOR.
- `Decode` - decode CBOR input and print it as JSON.
- `Tag` - wrap a CBOR value in a CBOR tag.

## Options

| Flag | Description |
| --- | --- |
| `-a` | Action to perform: `Encode`, `Decode`, or `Tag` (required) |
| `-f` | Input file path (default: stdin) |
| `-s` | Input string (used by `Encode` and `Tag`) |
| `-t` | CBOR tag number (required for `Tag`) |
| `-o` | Output file path (default: stdout) |
| `-format` | Output format for `Decode`: `pretty` (default) or `json` |
| `-root` | Confine file access to a directory (see [Sandbox](#sandbox)) |
| `-completion` | Print a shell completion script for `bash`, `zsh`, or `fish` |
| `-man` | Print the manual page and exit |
| `-q`, `-quiet` | Suppress diagnostics on standard error |
| `-version`, `-v` | Print the version and exit |

Flags accept one or two leading dashes, so `-format` and `--format` are
equivalent. `-completion` and `-man` do not need `-a`.

### Atomic output

Output written to `-o` is atomic: the data is staged in a temporary file in
the destination directory and renamed into place, so an interrupted run never
leaves a partially written file behind. Devices, sockets, and pipes named by
`-o` (for example `/dev/null`) are written directly, because they cannot be
replaced by a rename. A successful file write is reported on standard error
as `wrote N bytes to PATH` unless `-q` is given.

### Sandbox

`-root DIRECTORY` confines every path argument to that directory tree using
[`os.Root`](https://pkg.go.dev/os#Root). Input and output paths are resolved
inside the directory, and the operating system rejects anything that would
escape it, including absolute paths and `..` traversal:

```bash
cbor-tool -root /srv/cbor -a Decode -f incoming.cbor
cbor-tool -root /srv/cbor -a Encode -s 'hello' -o outgoing.cbor
```

Without `-root`, paths are resolved normally, which is what you want when
passing paths from your own shell.

## Examples

Encode a string:

```bash
printf 'hello' | cbor-tool -a Encode
cbor-tool -a Encode -s 'hello' -o hello.cbor
```

Decode CBOR to JSON:

```bash
cbor-tool -a Decode -f hello.cbor
cbor-tool -a Decode -f hello.cbor -format json
```

Tag a value:

```bash
cbor-tool -a Tag -t 42 -s 'hello'
cbor-tool -a Tag -t 7 -f hello.cbor -o tagged.cbor
```

Run a script silently and branch on the exit code:

```bash
if ! cbor-tool -q -a Decode -f broken.cbor > /dev/null; then
  echo "not valid CBOR"
fi
```

More scripts live in [`examples/`](../examples/).

## Shell completion

```bash
cbor-tool -completion bash > completions/cbor-tool.bash
cbor-tool -completion zsh  > completions/_cbor-tool
cbor-tool -completion fish > completions/cbor-tool.fish
```

The same scripts are committed in [`completions/`](../completions/) and are
regenerated with `make completions`.

- bash: `source completions/cbor-tool.bash`, or copy it to
  `/etc/bash_completion.d/`.
- zsh: add the directory containing `_cbor-tool` to `$fpath`.
- fish: copy `cbor-tool.fish` into `~/.config/fish/completions/`.

## Manual page

```bash
cbor-tool -man              # print to stdout
cbor-tool -man -o cbor-tool.1
make man                    # regenerate docs/cbor-tool.1
```

## Environment

No environment variables are required.

## Exit codes

- `0` - success
- `1` - runtime error (invalid input, I/O failure, ...)
- `2` - usage / flag parsing error

## Docker

```bash
docker build -t cbor-tool:latest .
docker run --rm cbor-tool -a Encode -s 'hello'
```

Published images are available from GitHub Container Registry:

```bash
docker pull ghcr.io/23skdu/cbor-tools:latest
docker run --rm ghcr.io/23skdu/cbor-tools:latest -a Decode -f /data/hello.cbor
```

## Input limits

Input read from a file or from standard input is capped at 10 MiB; larger
input fails with exit code `1` instead of being buffered.
