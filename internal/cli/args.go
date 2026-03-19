package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/maxanderson95/md2pdf/internal/version"
)

const usageText = `Usage:
  md2pdf <input.md>

Arguments:
  input.md    Markdown file to convert to PDF

Options:
  -h, --help     Show help
  --version      Show version

Environment:
  MD2PDF_BROWSER   Override the browser executable path
  MD2PDF_INTEGRATION  Enable browser-backed integration tests`

type Mode int

const (
	ModeRun Mode = iota
	ModeHelp
	ModeVersion
)

type Args struct {
	Mode  Mode
	Input string
}

type ExitError struct {
	Code int
	Err  error
}

func (e ExitError) Error() string {
	return e.Err.Error()
}

func (e ExitError) Unwrap() error {
	return e.Err
}

func (e ExitError) ExitCode() int {
	return e.Code
}

func Parse(argv []string) (Args, error) {
	normalized := make([]string, 0, len(argv))
	for _, arg := range argv {
		switch arg {
		case "--help":
			normalized = append(normalized, "-h")
		default:
			normalized = append(normalized, arg)
		}
	}

	fs := flag.NewFlagSet("md2pdf", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	showVersion := fs.Bool("version", false, "show version")

	if err := fs.Parse(normalized); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return Args{Mode: ModeHelp}, nil
		}
		return Args{}, ExitError{Code: 2, Err: fmt.Errorf("invalid arguments: %w", err)}
	}

	if *showVersion {
		if len(fs.Args()) > 0 {
			return Args{}, ExitError{Code: 2, Err: errors.New("--version does not accept positional arguments")}
		}
		return Args{Mode: ModeVersion}, nil
	}

	if len(fs.Args()) == 0 {
		return Args{Mode: ModeHelp}, nil
	}
	if len(fs.Args()) > 1 {
		return Args{}, ExitError{Code: 2, Err: fmt.Errorf("expected exactly one input file, got %d", len(fs.Args()))}
	}

	input := strings.TrimSpace(fs.Arg(0))
	if input == "" {
		return Args{}, ExitError{Code: 2, Err: errors.New("input markdown file cannot be empty")}
	}

	return Args{Mode: ModeRun, Input: input}, nil
}

func WriteHelp(w io.Writer) error {
	_, err := fmt.Fprintln(w, usageText)
	return err
}

func WriteVersion(w io.Writer, info version.BuildInfo) error {
	_, err := fmt.Fprintln(w, info.String())
	return err
}
