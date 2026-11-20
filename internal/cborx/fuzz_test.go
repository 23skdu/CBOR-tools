package cborx

import (
	"testing"
	"unicode/utf8"

	cbor "github.com/fxamacker/cbor/v2"
)

func FuzzEncodeDecodeRoundtrip(f *testing.F) {
	f.Add("")
	f.Add("hello")
	f.Add("héllo 🌍")
	f.Fuzz(func(t *testing.T, s string) {
		if !utf8.ValidString(s) {
			t.Skip()
		}
		data, err := EncodeString(s)
		if err != nil {
			t.Fatalf("encode: %v", err)
		}
		got, err := Decode(data)
		if err != nil {
			t.Fatalf("decode: %v", err)
		}
		if got.(string) != s {
			t.Fatalf("mismatch")
		}
	})
}

func FuzzDecodeNoPanic(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0x65, 'h', 'e', 'l', 'l', 'o'})
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = Decode(data)
		_, _ = DecodeJSON(data)
	})
}

func FuzzTagRoundtrip(f *testing.F) {
	f.Add("hello", uint64(42))
	f.Fuzz(func(t *testing.T, s string, tag uint64) {
		data, err := TagString(s, tag)
		if err != nil {
			t.Fatalf("tag: %v", err)
		}
		v, err := Decode(data)
		if err != nil {
			t.Skip()
		}
		gotTag, ok := v.(cbor.Tag)
		if !ok {
			t.Fatalf("not cbor.Tag: %T", v)
		}
		if gotTag.Number != tag {
			t.Fatalf("tag number mismatch")
		}
		if gotTag.Content.(string) != s {
			t.Fatalf("content mismatch")
		}
	})
}
