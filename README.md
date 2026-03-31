# md2pdf

`md2pdf` is a Go CLI that converts a Markdown file into a PDF with the same basename.

## Installation

Download a pre-compiled binary for your platform from the [latest release](https://github.com/MaxAnderson95/md2pdf/releases/latest), extract it, and place it somewhere on your `PATH`.

## Usage

```bash
md2pdf path/to/file.md
```

- The generated PDF is written beside the source Markdown file.
- If `file.pdf` already exists, the tool writes `file (01).pdf`, `file (02).pdf`, and so on.
- `--help` shows usage information.
- `--version` prints build metadata.

## Browser requirement

Version 1 uses an installed Chromium-based browser to render PDFs.

Supported browsers:
- Google Chrome
- Chromium
- Microsoft Edge

To override automatic detection:

```bash
MD2PDF_BROWSER=/path/to/browser md2pdf docs/guide.md
```

## Supported Markdown features

- CommonMark basics
- GitHub-flavored tables, task lists, and strikethrough
- Syntax highlighting for fenced code blocks
- Local relative images
- Remote images over HTTP and HTTPS

## Current limitations

- YAML front matter is ignored
- Page numbers and custom templates are not supported
- Missing local images fail the command
- Missing remote images warn on stderr and the PDF is still produced

## Development

```bash
go test ./...
go build ./cmd/md2pdf
```

Optional browser-backed integration tests:

```bash
MD2PDF_INTEGRATION=1 go test ./...
```
