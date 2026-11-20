package cborx

import (
	"fmt"
	"io"
	"os"
)

// ReadInput reads input bytes from path when it is non-empty, otherwise from r.
func ReadInput(path string, r io.Reader) ([]byte, error) {
	if path != "" {
		data, err := os.ReadFile(path) // #nosec G304 -- path is explicitly provided by the user
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrReadInput, err)
		}
		return data, nil
	}
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadInput, err)
	}
	return data, nil
}

// ReadStringInput returns the string to encode. If s is non-empty it is used;
// otherwise the file at path (or stdin) is read and interpreted as text.
func ReadStringInput(s, path string, r io.Reader) (string, error) {
	if s != "" {
		return s, nil
	}
	data, err := ReadInput(path, r)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", ErrEmptyInput
	}
	return string(data), nil
}

// WriteOutput writes data to path when it is non-empty, otherwise to w.
func WriteOutput(path string, data []byte, w io.Writer) error {
	if path != "" {
		if err := os.WriteFile(path, data, 0o600); err != nil {
			return fmt.Errorf("%w: %v", ErrWriteOutput, err)
		}
		return nil
	}
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("%w: %v", ErrWriteOutput, err)
	}
	return nil
}