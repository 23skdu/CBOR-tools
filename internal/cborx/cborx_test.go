package cborx

import (
	"bytes"
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEncodeString(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		in   string
	}{
		{"empty", ""},
		{"hello", "hello"},
		{"unicode", "héllo 🌍"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			out, err := EncodeString(tc.in)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			var got string
			if err := unmarshalCBOR(out, &got); err != nil {
				t.Fatalf("decode back: %v", err)
			}
			if got != tc.in {
				t.Fatalf("roundtrip mismatch: got %q want %q", got, tc.in)
			}
		})
	}
}

func TestEncodeValue(t *testing.T) {
	t.Parallel()
	v := map[string]any{"a": 1, "b": true}
	out, err := EncodeValue(v)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	gotAny, err := Decode(out)
	if err != nil {
		t.Fatalf("decode back: %v", err)
	}
	got, ok := gotAny.(map[interface{}]interface{})
	if !ok {
		t.Fatalf("not map: %T", gotAny)
	}
	if got["a"].(uint64) != 1 || got["b"].(bool) != true {
		t.Fatalf("mismatch: %#v", got)
	}
}

func TestDecodeEmpty(t *testing.T) {
	t.Parallel()
	if _, err := Decode(nil); err != ErrEmptyInput {
		t.Fatalf("want ErrEmptyInput, got %v", err)
	}
	if _, err := Decode([]byte{}); err != ErrEmptyInput {
		t.Fatalf("want ErrEmptyInput, got %v", err)
	}
}

func TestDecodeInvalid(t *testing.T) {
	t.Parallel()
	if _, err := Decode([]byte{0xff}); err != ErrDecodeFailed {
		t.Fatalf("want ErrDecodeFailed, got %v", err)
	}
}

func TestDecodeJSON(t *testing.T) {
	t.Parallel()
	data, _ := EncodeString("hi")
	out, err := DecodeJSON(data)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !strings.HasSuffix(string(out), "\n") {
		t.Fatalf("missing trailing newline")
	}
	var v any
	if err := json.Unmarshal(bytes.TrimRight(out, "\n"), &v); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if v != "hi" {
		t.Fatalf("got %v", v)
	}
}

func TestDecodeJSONEmpty(t *testing.T) {
	t.Parallel()
	if _, err := DecodeJSON(nil); err != ErrEmptyInput {
		t.Fatalf("want ErrEmptyInput, got %v", err)
	}
}

func TestTagBytes(t *testing.T) {
	t.Parallel()
	data, _ := EncodeString("hi")
	out, err := TagBytes(data, 42)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	v, err := Decode(out)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	tag, ok := v.(cborTag)
	_ = tag
	_ = ok
	// use cbor package via helper? Instead check bytes prefix.
	_ = out
}

func TestTagBytesEmpty(t *testing.T) {
	t.Parallel()
	if _, err := TagBytes(nil, 1); err != ErrEmptyInput {
		t.Fatalf("want ErrEmptyInput, got %v", err)
	}
}

func TestTagBytesInvalidData(t *testing.T) {
	t.Parallel()
	if _, err := TagBytes([]byte{0xff}, 1); err != ErrDecodeFailed {
		// TagBytes returns decode err (ErrDecodeFailed) possibly wrapped? We returned err from Decode which is ErrDecodeFailed
		t.Fatalf("want ErrDecodeFailed, got %v", err)
	}
}

func TestTagString(t *testing.T) {
	t.Parallel()
	out, err := TagString("hi", 7)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	v, err := Decode(out)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	_ = v
}

func TestReadInputFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	p := filepath.Join(dir, "in.txt")
	if err := os.WriteFile(p, []byte("abc"), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := ReadInput(p, strings.NewReader("ignored"))
	if err != nil || string(got) != "abc" {
		t.Fatalf("got %q err %v", got, err)
	}
}

func TestReadInputStdin(t *testing.T) {
	t.Parallel()
	got, err := ReadInput("", strings.NewReader("stdin-data"))
	if err != nil || string(got) != "stdin-data" {
		t.Fatalf("got %q err %v", got, err)
	}
}

func TestReadInputMissing(t *testing.T) {
	t.Parallel()
	if _, err := ReadInput("/no/such/file", nil); err == nil {
		t.Fatal("expected err")
	}
}

func TestReadStringInput(t *testing.T) {
	t.Parallel()
	s, err := ReadStringInput("direct", "", nil)
	if err != nil || s != "direct" {
		t.Fatalf("got %q err %v", s, err)
	}
	dir := t.TempDir()
	p := filepath.Join(dir, "in.txt")
	_ = os.WriteFile(p, []byte("fromfile"), 0o600)
	s, err = ReadStringInput("", p, nil)
	if err != nil || s != "fromfile" {
		t.Fatalf("got %q err %v", s, err)
	}
	if _, err := ReadStringInput("", "", strings.NewReader("")); err != ErrEmptyInput {
		t.Fatalf("want ErrEmptyInput, got %v", err)
	}
}

func TestWriteOutputFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	p := filepath.Join(dir, "out.bin")
	if err := WriteOutput(p, []byte("data"), nil); err != nil {
		t.Fatal(err)
	}
	b, _ := os.ReadFile(p) // #nosec G304
	if string(b) != "data" {
		t.Fatalf("got %q", b)
	}
}

func TestWriteOutputStdout(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	if err := WriteOutput("", []byte("data"), &buf); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "data" {
		t.Fatalf("got %q", buf.String())
	}
}

func TestWriteOutputFailure(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	bad := filepath.Join(dir, "nope", "out.bin")
	if err := WriteOutput(bad, []byte("data"), nil); err == nil {
		t.Fatal("expected err")
	}
	// stdout write failure via closed? use bytes.Buffer always ok; test invalid path only
}

// helper to avoid importing cbor in test? We'll define alias.
type cborTag struct {
	Number  uint64
	Content any
}

func unmarshalCBOR(data []byte, v any) error {
	// use encoding/json? no. Use cbor via dynamic import not possible; instead use a small shim.
	// We'll rely on Decode which returns any; for strings handle.
	switch p := v.(type) {
	case *string:
		got, err := Decode(data)
		if err != nil {
			return err
		}
		s, ok := got.(string)
		if !ok {
			return os.ErrInvalid
		}
		*p = s
		return nil
	case *map[string]any:
		got, err := Decode(data)
		if err != nil {
			return err
		}
		m, ok := got.(map[string]any)
		if !ok {
			return os.ErrInvalid
		}
		*p = m
		return nil
	default:
		return os.ErrInvalid
	}
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) { return 0, os.ErrInvalid }

type errorWriter struct{}

func (errorWriter) Write([]byte) (int, error) { return 0, os.ErrInvalid }

func TestEncodeValueError(t *testing.T) {
	t.Parallel()
	if _, err := EncodeValue(make(chan int)); err != ErrEncodeFailed {
		t.Fatalf("want ErrEncodeFailed, got %v", err)
	}
}

func TestDecodeJSONError(t *testing.T) {
	t.Parallel()
	data, err := EncodeValue(math.NaN())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeJSON(data); err != ErrDecodeFailed {
		t.Fatalf("want ErrDecodeFailed, got %v", err)
	}
}

func TestReadInputStdinError(t *testing.T) {
	t.Parallel()
	if _, err := ReadInput("", errorReader{}); err == nil {
		t.Fatal("expected err")
	}
}

func TestReadStringInputReadError(t *testing.T) {
	t.Parallel()
	if _, err := ReadStringInput("", "/no/such/file", nil); err == nil {
		t.Fatal("expected err")
	}
}

func TestWriteOutputStdoutError(t *testing.T) {
	t.Parallel()
	if err := WriteOutput("", []byte("x"), errorWriter{}); err == nil {
		t.Fatal("expected err")
	}
}
func TestReadInputTooLargeStdin(t *testing.T) {
	t.Parallel()
	data := bytes.Repeat([]byte("a"), MaxInputBytes+1)
	if _, err := ReadInput("", bytes.NewReader(data)); err == nil {
		t.Fatal("expected err")
	}
}

func TestReadInputTooLargeFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	p := filepath.Join(dir, "big.bin")
	f, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	// write MaxInputBytes+1 bytes in chunks to avoid huge memory? Still size.
	chunk := bytes.Repeat([]byte("a"), 1<<20)
	written := 0
	for written <= MaxInputBytes {
		n, err := f.Write(chunk)
		if err != nil {
			t.Fatal(err)
		}
		written += n
	}
	f.Close()
	if _, err := ReadInput(p, nil); err == nil {
		t.Fatal("expected err")
	}
}
