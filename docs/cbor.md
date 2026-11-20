# CBOR overview

CBOR (Concise Binary Object Representation) is a binary data format defined in
RFC 8949. It is designed to be small, fast to parse, and a superset of JSON's
data model, while also supporting byte strings, tags, and other types.

## Why CBOR

- Compact binary encoding compared to text formats like JSON.
- Self-describing data items.
- Extensible via tags (for example, timestamps, bignums, or custom semantics).
- Suitable for constrained devices and network protocols.

## How this tool is used

`cbor-tool` helps inspect and transform CBOR data from the command line:

- **Encode** text into CBOR for testing or fixtures.
- **Decode** CBOR payloads into readable JSON.
- **Tag** values to experiment with or produce tagged CBOR items.

Typical workflow:

```bash
# Encode a string
cbor-tool -a Encode -s 'hello' -o hello.cbor

# Decode it back
cbor-tool -a Decode -f hello.cbor

# Wrap it in a tag
cbor-tool -a Tag -t 42 -f hello.cbor -o tagged.cbor
```

See [docs/cli.md](cli.md) for the full CLI reference.
