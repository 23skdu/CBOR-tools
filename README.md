# CBOR-tools

Go CLI for encoding, decoding, and tagging CBOR data.

## Build

```bash
go build -o bin/cbor-tool ./cmd/cbor-tool
```

## Usage

```bash
cbor-tool -a Encode -s 'hello'
cbor-tool -a Decode -f data.cbor
cbor-tool -a Tag -t 42 -s 'hello'
```

See [docs/cli.md](docs/cli.md).

## Test

```bash
make test
```

## License

This project is released into the public domain under the [Unlicense](https://unlicense.org). See [LICENSE](LICENSE) for details.
