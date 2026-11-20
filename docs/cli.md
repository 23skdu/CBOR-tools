# CBOR tool CLI reference

`cbor-tool` is a small command line tool for encoding, decoding, and tagging
CBOR (Concise Binary Object Representation, RFC 8949) data.

## Usage

```bash
cbor-tool -a <action> [options]
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

## Examples

Encode a string:

```bash
printf 'hello' | cbor-tool -a Encode
cbor-tool -a Encode -s 'hello' -o hello.cbor
```

Decode CBOR to JSON:

```bash
cbor-tool -a Decode -f hello.cbor
```

Tag a value:

```bash
cbor-tool -a Tag -t 42 -s 'hello'
cbor-tool -a Tag -t 7 -f hello.cbor -o tagged.cbor
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
