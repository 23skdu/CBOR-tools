# CBOR-tools

[![CI](https://github.com/23skdu/CBOR-tools/actions/workflows/ci.yml/badge.svg)](https://github.com/23skdu/CBOR-tools/actions/workflows/ci.yml)

A Go command line tool for encoding, decoding, and tagging CBOR (Concise
Binary Object Representation, RFC 8949) data.

## Features

- Encode strings as CBOR.
- Decode CBOR into JSON.
- Wrap CBOR values in CBOR tags.
- Streaming-friendly stdin/stdout and file I/O.
- Single static binary, no runtime dependencies.

## Requirements

- Go 1.27.4+

## Build

```bash
make build
```

Or directly:

```bash
go build -o bin/cbor-tool ./cmd/cbor-tool
```

## Usage

```bash
cbor-tool -a Encode -s 'hello'
cbor-tool -a Decode -f data.cbor
cbor-tool -a Tag -t 42 -s 'hello'
```

See [docs/cli.md](docs/cli.md) for the full CLI reference and [docs/cbor.md](docs/cbor.md) for a CBOR overview.

## Docker

```bash
docker build -t cbor-tool:latest .
docker run --rm cbor-tool -a Encode -s 'hello'
```

## Development

```bash
make test       # run unit tests
make fuzz       # run fuzz tests
make cover      # coverage report
make lint       # static analysis
make security   # security scans
make clean      # remove build artifacts
```

## License

MIT
