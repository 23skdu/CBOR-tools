package cborx

import (
	"encoding/json"

	cbor "github.com/fxamacker/cbor/v2"
)

// EncodeString encodes s as a CBOR text string.
func EncodeString(s string) ([]byte, error) {
	return cbor.Marshal(s)
}

// EncodeValue encodes v as CBOR.
func EncodeValue(v any) ([]byte, error) {
	out, err := cbor.Marshal(v)
	if err != nil {
		return nil, ErrEncodeFailed
	}
	return out, nil
}

// Decode decodes CBOR data into a Go value.
func Decode(data []byte) (any, error) {
	if len(data) == 0 {
		return nil, ErrEmptyInput
	}
	var v any
	if err := cbor.Unmarshal(data, &v); err != nil {
		return nil, ErrDecodeFailed
	}
	return v, nil
}

// DecodeJSON decodes CBOR data and returns its pretty JSON representation.
func DecodeJSON(data []byte) ([]byte, error) {
	v, err := Decode(data)
	if err != nil {
		return nil, err
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, ErrDecodeFailed
	}
	out = append(out, '\n')
	return out, nil
}

// DecodeCompactJSON decodes CBOR data and returns its compact JSON representation.
func DecodeCompactJSON(data []byte) ([]byte, error) {
	v, err := Decode(data)
	if err != nil {
		return nil, err
	}
	out, err := json.Marshal(v)
	if err != nil {
		return nil, ErrDecodeFailed
	}
	out = append(out, '\n')
	return out, nil
}
