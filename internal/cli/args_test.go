package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/maxanderson95/md2pdf/internal/version"
)

func TestParseRunInput(t *testing.T) {
	t.Parallel()

	args, err := Parse([]string{"notes.md"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if args.Mode != ModeRun {
		t.Fatalf("unexpected mode: %v", args.Mode)
	}
	if args.Input != "notes.md" {
		t.Fatalf("unexpected input: %q", args.Input)
	}
}

func TestParseHelp(t *testing.T) {
	t.Parallel()

	args, err := Parse([]string{"--help"})
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if args.Mode != ModeHelp {
		t.Fatalf("unexpected mode: %v", args.Mode)
	}
}

func TestParseVersionRejectsPositionalArguments(t *testing.T) {
	t.Parallel()

	_, err := Parse([]string{"--version", "notes.md"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var exitErr ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("expected ExitError, got %T", err)
	}
	if exitErr.ExitCode() != 2 {
		t.Fatalf("unexpected exit code: %d", exitErr.ExitCode())
	}
}

func TestWriteHelp(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := WriteHelp(&buf); err != nil {
		t.Fatalf("WriteHelp returned error: %v", err)
	}
	if !strings.Contains(buf.String(), "md2pdf <input.md>") {
		t.Fatalf("help output missing usage text: %q", buf.String())
	}
}

func TestWriteVersion(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	info := version.BuildInfo{Version: "1.0.0", Commit: "abc123", Date: "2026-03-17"}
	if err := WriteVersion(&buf, info); err != nil {
		t.Fatalf("WriteVersion returned error: %v", err)
	}
	if !strings.Contains(buf.String(), "1.0.0") {
		t.Fatalf("version output missing version: %q", buf.String())
	}
}
