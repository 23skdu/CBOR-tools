package main

import (
	"bytes"
	"os"
	"testing"
)

func TestRunHelp(t *testing.T) {
	old := os.Args
	os.Args = []string{"cbor-tool"}
	defer func() { os.Args = old }()
	osExit = func(int) {}
	main()
}

func TestMainCallsExit(t *testing.T) {
	old := os.Args
	os.Args = []string{"cbor-tool", "-a", "Encode", "-s", "x"}
	defer func() { os.Args = old }()
	called := false
	osExit = func(code int) { called = true }
	main()
	if !called {
		t.Fatal("osExit not called")
	}
}

func TestRunVersionFlagNA(t *testing.T) {
	old := os.Args
	os.Args = []string{"cbor-tool", "-a", "Encode", "-s", "x"}
	defer func() { os.Args = old }()
	// redirect stdout
	oldOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	rc := run()
	_ = w.Close()
	os.Stdout = oldOut
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	if rc != 0 {
		t.Fatalf("rc=%d", rc)
	}
}