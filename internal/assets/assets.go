package assets

import _ "embed"

var (
	//go:embed template.html
	TemplateHTML string

	//go:embed markdown.css
	MarkdownCSS string

	//go:embed print.css
	PrintCSS string

	//go:embed highlight.css
	HighlightCSS string
)
