package cborx

import (
	cbor "github.com/fxamacker/cbor/v2"
)

// TagBytes decodes the CBOR data, wraps it in a CBOR tag with the given
// number, and re-encodes the result.
func TagBytes(data []byte, tag uint64) ([]byte, error) {
	if len(data) == 0 {
		return nil, ErrEmptyInput
	}
	v, err := Decode(data)
	if err != nil {
		return nil, err
	}
	t := cbor.Tag{Number: tag, Content: v}
	return cbor.Marshal(t)
}

// TagString encodes s as a CBOR text string, wraps it in a CBOR tag with the
// given number, and returns the encoded bytes.
func TagString(s string, tag uint64) ([]byte, error) {
	t := cbor.Tag{Number: tag, Content: s}
	return cbor.Marshal(t)
}
