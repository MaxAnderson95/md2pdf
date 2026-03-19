package document

import (
	"bytes"
	"fmt"
	stdhtml "html"
	"html/template"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/maxanderson95/md2pdf/internal/assets"
	xhtml "golang.org/x/net/html"
)

type Builder struct {
	tmpl         *template.Template
	markdownCSS  string
	printCSS     string
	highlightCSS string
}

type Document struct {
	HTML       string
	ImageRefs  []ImageRef
	Warnings   []string
	BaseHref   string
	InputTitle string
}

type ImageRef struct {
	Source    string
	Resolved  string
	IsRemote  bool
	IsDataURI bool
}

type templateData struct {
	BaseHref     template.URL
	MarkdownCSS  template.CSS
	PrintCSS     template.CSS
	HighlightCSS template.CSS
	Content      template.HTML
	Title        string
}

func NewBuilder() (Builder, error) {
	tmpl, err := template.New("document").Parse(assets.TemplateHTML)
	if err != nil {
		return Builder{}, fmt.Errorf("parse HTML template: %w", err)
	}

	return Builder{
		tmpl:         tmpl,
		markdownCSS:  assets.MarkdownCSS,
		printCSS:     assets.PrintCSS,
		highlightCSS: assets.HighlightCSS,
	}, nil
}

func (b Builder) Build(contentHTML string, inputPath string) (Document, error) {
	baseHref, err := fileBaseHref(filepath.Dir(inputPath))
	if err != nil {
		return Document{}, err
	}

	imageRefs, err := extractImageRefs(contentHTML, baseHref)
	if err != nil {
		return Document{}, fmt.Errorf("extract image references: %w", err)
	}

	title := deriveTitle(inputPath, contentHTML)
	data := templateData{
		BaseHref:     template.URL(baseHref),
		MarkdownCSS:  template.CSS(b.markdownCSS),
		PrintCSS:     template.CSS(b.printCSS),
		HighlightCSS: template.CSS(b.highlightCSS),
		Content:      template.HTML(contentHTML),
		Title:        title,
	}

	var buf bytes.Buffer
	if err := b.tmpl.Execute(&buf, data); err != nil {
		return Document{}, fmt.Errorf("render HTML template: %w", err)
	}

	return Document{
		HTML:       buf.String(),
		ImageRefs:  imageRefs,
		BaseHref:   baseHref,
		InputTitle: title,
	}, nil
}

func fileBaseHref(dir string) (string, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("resolve base directory: %w", err)
	}
	baseURL := &url.URL{Scheme: "file", Path: filepath.ToSlash(absDir) + "/"}
	return baseURL.String(), nil
}

func deriveTitle(inputPath, htmlContent string) string {
	base := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	tokenizer := xhtml.NewTokenizer(strings.NewReader(htmlContent))
	for {
		tt := tokenizer.Next()
		switch tt {
		case xhtml.ErrorToken:
			return base
		case xhtml.StartTagToken:
			tn, _ := tokenizer.TagName()
			if string(tn) != "h1" {
				continue
			}
			if tokenizer.Next() == xhtml.TextToken {
				text := strings.TrimSpace(stdhtml.UnescapeString(string(tokenizer.Text())))
				if text != "" {
					return text
				}
			}
		}
	}
}

func extractImageRefs(fragmentHTML string, baseHref string) ([]ImageRef, error) {
	root, err := xhtml.Parse(strings.NewReader("<!doctype html><html><body>" + fragmentHTML + "</body></html>"))
	if err != nil {
		return nil, err
	}
	baseURL, err := url.Parse(baseHref)
	if err != nil {
		return nil, err
	}

	refs := make([]ImageRef, 0)
	var walk func(*xhtml.Node)
	walk = func(node *xhtml.Node) {
		if node.Type == xhtml.ElementNode && node.Data == "img" {
			for _, attr := range node.Attr {
				if attr.Key != "src" {
					continue
				}
				src := strings.TrimSpace(attr.Val)
				if src == "" {
					continue
				}
				resolved := src
				parsed, parseErr := url.Parse(src)
				isRemote := parseErr == nil && (parsed.Scheme == "http" || parsed.Scheme == "https")
				isData := strings.HasPrefix(src, "data:")
				if parseErr == nil && !isRemote && !isData {
					resolved = baseURL.ResolveReference(parsed).String()
				}
				refs = append(refs, ImageRef{Source: src, Resolved: resolved, IsRemote: isRemote, IsDataURI: isData})
			}
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(root)
	return refs, nil
}
