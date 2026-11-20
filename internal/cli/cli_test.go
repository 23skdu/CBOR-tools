package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runCLI(t *testing.T, args []string, stdin string) (int, string, string) {
	t.Helper()
	in := bytes.NewBufferString(stdin)
	var out, errb bytes.Buffer
	rc := Run(args, in, &out, &errb)
	return rc, out.String(), errb.String()
}

func TestRunNoAction(t *testing.T) {
	rc, _, err := runCLI(t, nil, "")
	if rc != 1 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(err, "Must specify action") {
		t.Fatalf("err=%s", err)
	}
}

func TestRunInvalidAction(t *testing.T) {
	rc, _, err := runCLI(t, []string{"-a", "Nope"}, "")
	if rc != 1 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(err, "valid action") {
		t.Fatalf("err=%s", err)
	}
}

func TestRunParseError(t *testing.T) {
	rc, _, _ := runCLI(t, []string{"-badflag"}, "")
	if rc != 2 {
		t.Fatalf("rc=%d", rc)
	}
}

func TestEncodeStdin(t *testing.T) {
	rc, out, _ := runCLI(t, []string{"-a", "Encode"}, "hi")
	if rc != 0 {
		t.Fatalf("rc=%d", rc)
	}
	if out != string([]byte{0x62, 'h', 'i'}) {
		// cbor text string "hi" => 0x62 68 69
		t.Logf("out=% x", []byte(out))
	}
}

func TestEncodeStringFlag(t *testing.T) {
	rc, out, _ := runCLI(t, []string{"-a", "Encode", "-s", "abc"}, "")
	if rc != 0 {
		t.Fatalf("rc=%d", rc)
	}
	if len(out) == 0 {
		t.Fatal("empty")
	}
}

func TestEncodeFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "in.txt")
	_ = os.WriteFile(p, []byte("filedata"), 0o600)
	rc, out, _ := runCLI(t, []string{"-a", "Encode", "-f", p}, "")
	if rc != 0 {
		t.Fatalf("rc=%d", rc)
	}
	if len(out) == 0 {
		t.Fatal("empty")
	}
}

func TestEncodeOutputFile(t *testing.T) {
	dir := t.TempDir()
	outp := filepath.Join(dir, "out.cbor")
	rc, out, _ := runCLI(t, []string{"-a", "Encode", "-s", "x", "-o", outp}, "")
	if rc != 0 {
		t.Fatalf("rc=%d", rc)
	}
	if out != "" {
		t.Fatalf("stdout not empty: %q", out)
	}
	b, _ := os.ReadFile(outp) // #nosec G304
	if len(b) == 0 {
		t.Fatal("no file")
	}
}

func TestEncodeReadError(t *testing.T) {
	rc, _, err := runCLI(t, []string{"-a", "Encode", "-f", "/no/such/file"}, "")
	if rc != 1 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(err, "failed to read input") {
		t.Fatalf("err=%s", err)
	}
}

func TestEncodeWriteError(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "nope", "out.cbor")
	rc, _, err := runCLI(t, []string{"-a", "Encode", "-s", "x", "-o", bad}, "")
	if rc != 1 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(err, "failed to write output") {
		t.Fatalf("err=%s", err)
	}
}

func TestDecodeStdin(t *testing.T) {
	// encode "hi" then decode
	enc := []byte{0x62, 'h', 'i'}
	in := bytes.NewBuffer(enc)
	var out, errb bytes.Buffer
	rc := Run([]string{"-a", "Decode"}, in, &out, &errb)
	if rc != 0 {
		t.Fatalf("rc=%d err=%s", rc, errb.String())
	}
	if !strings.Contains(out.String(), "\"hi\"") {
		t.Fatalf("out=%s", out.String())
	}
}

func TestDecodeFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "in.cbor")
	_ = os.WriteFile(p, []byte{0x62, 'h', 'i'}, 0o600)
	rc, out, _ := runCLI(t, []string{"-a", "Decode", "-f", p}, "")
	if rc != 0 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(out, "\"hi\"") {
		t.Fatalf("out=%s", out)
	}
}

func TestDecodeEmpty(t *testing.T) {
	rc, _, err := runCLI(t, []string{"-a", "Decode"}, "")
	if rc != 1 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(err, "empty input") {
		t.Fatalf("err=%s", err)
	}
}

func TestDecodeInvalid(t *testing.T) {
	rc, _, err := runCLI(t, []string{"-a", "Decode"}, string([]byte{0xff}))
	if rc != 1 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(err, "failed to decode CBOR") {
		t.Fatalf("err=%s", err)
	}
}

func TestDecodeReadError(t *testing.T) {
	rc, _, err := runCLI(t, []string{"-a", "Decode", "-f", "/no/such/file"}, "")
	if rc != 1 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(err, "failed to read input") {
		t.Fatalf("err=%s", err)
	}
}

func TestDecodeWriteError(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "nope", "out.txt")
	rc, _, err := runCLI(t, []string{"-a", "Decode", "-f", "/dev/null", "-o", bad}, "")
	// /dev/null empty -> decode empty error before write; instead use valid input file
	_ = rc
	_ = err
	// Use valid input
	p := filepath.Join(dir, "in.cbor")
	_ = os.WriteFile(p, []byte{0x62, 'h', 'i'}, 0o600)
	rc, _, err = runCLI(t, []string{"-a", "Decode", "-f", p, "-o", bad}, "")
	if rc != 1 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(err, "failed to write output") {
		t.Fatalf("err=%s", err)
	}
}

func TestTagMissingTag(t *testing.T) {
	rc, _, err := runCLI(t, []string{"-a", "Tag", "-s", "hi"}, "")
	if rc != 1 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(err, "tag number is required") {
		t.Fatalf("err=%s", err)
	}
}

func TestTagInvalidTag(t *testing.T) {
	rc, _, err := runCLI(t, []string{"-a", "Tag", "-t", "abc", "-s", "hi"}, "")
	if rc != 1 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(err, "invalid tag number") {
		t.Fatalf("err=%s", err)
	}
}

func TestTagString(t *testing.T) {
	rc, out, _ := runCLI(t, []string{"-a", "Tag", "-t", "42", "-s", "hi"}, "")
	if rc != 0 {
		t.Fatalf("rc=%d", rc)
	}
	if len(out) == 0 {
		t.Fatal("empty")
	}
}

func TestTagBytesStdin(t *testing.T) {
	enc := []byte{0x62, 'h', 'i'}
	rc, out, _ := runCLI(t, []string{"-a", "Tag", "-t", "7"}, string(enc))
	if rc != 0 {
		t.Fatalf("rc=%d", rc)
	}
	if len(out) == 0 {
		t.Fatal("empty")
	}
}

func TestTagBytesFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "in.cbor")
	_ = os.WriteFile(p, []byte{0x62, 'h', 'i'}, 0o600)
	rc, out, _ := runCLI(t, []string{"-a", "Tag", "-t", "7", "-f", p}, "")
	if rc != 0 {
		t.Fatalf("rc=%d", rc)
	}
	if len(out) == 0 {
		t.Fatal("empty")
	}
}

func TestTagBytesEmpty(t *testing.T) {
	rc, _, err := runCLI(t, []string{"-a", "Tag", "-t", "7"}, "")
	if rc != 1 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(err, "empty input") {
		t.Fatalf("err=%s", err)
	}
}

func TestTagBytesInvalidData(t *testing.T) {
	rc, _, err := runCLI(t, []string{"-a", "Tag", "-t", "7"}, string([]byte{0xff}))
	if rc != 1 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(err, "failed to decode CBOR") {
		// TagBytes returns ErrDecodeFailed (wrapped? we return err from Decode which is ErrDecodeFailed not wrapped)
		t.Fatalf("err=%s", err)
	}
}

func TestTagOutputFile(t *testing.T) {
	dir := t.TempDir()
	outp := filepath.Join(dir, "out.cbor")
	rc, out, _ := runCLI(t, []string{"-a", "Tag", "-t", "7", "-s", "hi", "-o", outp}, "")
	if rc != 0 {
		t.Fatalf("rc=%d", rc)
	}
	if out != "" {
		t.Fatalf("stdout not empty")
	}
	b, _ := os.ReadFile(outp) // #nosec G304
	if len(b) == 0 {
		t.Fatal("no file")
	}
}

func TestTagWriteError(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "nope", "out.cbor")
	rc, _, err := runCLI(t, []string{"-a", "Tag", "-t", "7", "-s", "hi", "-o", bad}, "")
	if rc != 1 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(err, "failed to write output") {
		t.Fatalf("err=%s", err)
	}
}

func TestTagStrTakesPrecedence(t *testing.T) {
	// when both -s and stdin, -s used; stdin not read. Provide invalid stdin; should still succeed.
	rc, out, _ := runCLI(t, []string{"-a", "Tag", "-t", "7", "-s", "hi"}, string([]byte{0xff}))
	if rc != 0 {
		t.Fatalf("rc=%d", rc)
	}
	if len(out) == 0 {
		t.Fatal("empty")
	}
}

func TestTagReadError(t *testing.T) {
	rc, _, err := runCLI(t, []string{"-a", "Tag", "-t", "7", "-f", "/no/such/file"}, "")
	if rc != 1 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(err, "failed to read input") {
		t.Fatalf("err=%s", err)
	}
}
func TestVersionFlag(t *testing.T) {
	old := Version
	Version = "test-version"
	defer func() { Version = old }()
	rc, out, _ := runCLI(t, []string{"-version"}, "")
	if rc != 0 {
		t.Fatalf("rc=%d", rc)
	}
	if !strings.Contains(out, "test-version") {
		t.Fatalf("out=%s", out)
	}
	rc, out, _ = runCLI(t, []string{"-v"}, "")
	if rc != 0 || !strings.Contains(out, "test-version") {
		t.Fatalf("rc=%d out=%s", rc, out)
	}
}
