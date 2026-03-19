package app

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/maxanderson95/md2pdf/internal/browser"
	"github.com/maxanderson95/md2pdf/internal/cli"
	"github.com/maxanderson95/md2pdf/internal/document"
	"github.com/maxanderson95/md2pdf/internal/input"
	"github.com/maxanderson95/md2pdf/internal/markdown"
	"github.com/maxanderson95/md2pdf/internal/pdf"
	"github.com/maxanderson95/md2pdf/internal/version"
)

func Run(ctx context.Context, argv []string, stdout, stderr io.Writer, build version.BuildInfo) error {
	args, err := cli.Parse(argv)
	if err != nil {
		return err
	}

	switch args.Mode {
	case cli.ModeHelp:
		return cli.WriteHelp(stdout)
	case cli.ModeVersion:
		return cli.WriteVersion(stdout, build)
	}

	paths, err := input.Resolve(args.Input)
	if err != nil {
		return cli.ExitError{Code: 3, Err: err}
	}

	content, err := os.ReadFile(paths.InputPath)
	if err != nil {
		return cli.ExitError{Code: 3, Err: fmt.Errorf("read input markdown: %w", err)}
	}

	rendered, err := markdown.New().Render(content)
	if err != nil {
		return cli.ExitError{Code: 4, Err: fmt.Errorf("render markdown: %w", err)}
	}

	builder, err := document.NewBuilder()
	if err != nil {
		return cli.ExitError{Code: 4, Err: fmt.Errorf("load document assets: %w", err)}
	}

	doc, err := builder.Build(rendered, paths.InputPath)
	if err != nil {
		return cli.ExitError{Code: 4, Err: fmt.Errorf("build printable document: %w", err)}
	}

	if err := validateLocalImages(paths.InputDir, doc.ImageRefs); err != nil {
		return cli.ExitError{Code: 4, Err: err}
	}

	browserInfo, err := browser.Find()
	if err != nil {
		return cli.ExitError{Code: 5, Err: err}
	}

	result, err := pdf.Generator{Browser: browserInfo}.Write(ctx, doc, paths.OutputPath)
	if err != nil {
		return cli.ExitError{Code: 6, Err: err}
	}

	for _, warning := range result.Warnings {
		fmt.Fprintf(stderr, "warning: failed to load remote asset: %s\n", warning)
	}

	_, err = fmt.Fprintln(stdout, paths.OutputPath)
	return err
}
