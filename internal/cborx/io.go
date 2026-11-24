package cborx

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// MaxInputBytes is the maximum number of bytes ReadInput will read from a
// file or from stdin.
const MaxInputBytes = 10 << 20

// errTooLarge reports input exceeding MaxInputBytes.
var errTooLarge = errors.New("input too large")

// ReadInput reads input bytes from path when it is non-empty, otherwise from
// r. When root is non-nil, path is resolved inside that directory tree and
// cannot reference locations outside it (see os.OpenRoot).
func ReadInput(root *os.Root, path string, r io.Reader) ([]byte, error) {
	if path != "" {
		data, err := readFile(root, path)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrReadInput, err)
		}
		return data, nil
	}
	data, err := readStream(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadInput, err)
	}
	return data, nil
}

// ReadStringInput returns the string to encode. If s is non-empty it is used;
// otherwise the file at path (or stdin) is read and interpreted as text.
func ReadStringInput(root *os.Root, s, path string, r io.Reader) (string, error) {
	if s != "" {
		return s, nil
	}
	data, err := ReadInput(root, path, r)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", ErrEmptyInput
	}
	return string(data), nil
}

// WriteOutput writes data to path when it is non-empty, otherwise to w.
//
// File writes are atomic: the data is staged in a temporary file in the
// destination directory and renamed into place, so an interrupted run never
// leaves a partially written file behind. Devices, sockets, and pipes cannot
// be replaced by a rename and are written directly.
func WriteOutput(root *os.Root, path string, data []byte, w io.Writer) error {
	if path != "" {
		if err := writeFileAtomic(root, path, data); err != nil {
			return fmt.Errorf("%w: %v", ErrWriteOutput, err)
		}
		return nil
	}
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("%w: %v", ErrWriteOutput, err)
	}
	return nil
}

// readFile reads at most MaxInputBytes+1 bytes from path.
func readFile(root *os.Root, path string) ([]byte, error) {
	var f *os.File
	var err error
	if root != nil {
		f, err = root.Open(path)
	} else {
		f, err = os.Open(path) // #nosec G304 -- path is explicitly provided by the user
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	return readStream(f)
}

// readStream reads at most MaxInputBytes+1 bytes from r.
func readStream(r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, MaxInputBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > MaxInputBytes {
		return nil, errTooLarge
	}
	return data, nil
}

// writeFileAtomic writes data to path relative to root, or relative to the
// directory containing path when root is nil.
func writeFileAtomic(root *os.Root, path string, data []byte) error {
	if root == nil {
		// Anchor a root at the destination directory so the temporary file is
		// created next to the target and the rename stays on one filesystem.
		r, err := os.OpenRoot(filepath.Dir(path))
		if err != nil {
			return err
		}
		defer func() { _ = r.Close() }()
		return writeInto(r, filepath.Base(path), data)
	}
	dir, name := filepath.Split(path)
	if dir == "" {
		dir = "."
	}
	d, err := root.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer func() { _ = d.Close() }()
	return writeInto(d, name, data)
}

// writeInto writes data to name inside the directory rooted at d.
func writeInto(d *os.Root, name string, data []byte) error {
	if fi, err := d.Stat(name); err == nil && fi.Mode()&(os.ModeDevice|os.ModeNamedPipe|os.ModeSocket) != 0 {
		f, err := d.OpenFile(name, os.O_WRONLY, 0)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()
		_, err = f.Write(data)
		return err
	}
	// The temporary name is hidden and unique: a stale one left by a crash is
	// never mistaken for the destination.
	tmp := "." + name + ".tmp-" + rand.Text()
	if err := d.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	if err := d.Rename(tmp, name); err != nil {
		_ = d.Remove(tmp)
		return err
	}
	return nil
}
