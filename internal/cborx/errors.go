package cborx

import "errors"

// Sentinel errors returned by the CBOR helpers.
var (
	ErrInvalidAction  = errors.New("invalid action")
	ErrEmptyInput     = errors.New("empty input")
	ErrInvalidTag     = errors.New("invalid tag number")
	ErrReadInput      = errors.New("failed to read input")
	ErrWriteOutput    = errors.New("failed to write output")
	ErrEncodeFailed   = errors.New("failed to encode CBOR")
	ErrDecodeFailed   = errors.New("failed to decode CBOR")
	ErrTagFailed      = errors.New("failed to tag CBOR")
)
