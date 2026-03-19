package markdown

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderMatchesGoldenHTML(t *testing.T) {
	t.Parallel()

	renderer := New()
	markdownBytes, err := os.ReadFile(filepath.Join("..", "..", "testdata", "markdown", "basic.md"))
	if err != nil {
		t.Fatalf("ReadFile markdown failed: %v", err)
	}
	goldenBytes, err := os.ReadFile(filepath.Join("..", "..", "testdata", "html", "basic.golden.html"))
	if err != nil {
		t.Fatalf("ReadFile golden failed: %v", err)
	}

	got, err := renderer.Render(markdownBytes)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if strings.TrimSpace(got) != strings.TrimSpace(string(goldenBytes)) {
		t.Fatalf("render output mismatch\n--- got ---\n%s\n--- want ---\n%s", got, string(goldenBytes))
	}
}

func TestRenderAddsSyntaxHighlightingClasses(t *testing.T) {
	t.Parallel()

	renderer := New()
	markdownBytes, err := os.ReadFile(filepath.Join("..", "..", "testdata", "markdown", "tables-and-code.md"))
	if err != nil {
		t.Fatalf("ReadFile markdown failed: %v", err)
	}

	html, err := renderer.Render(markdownBytes)
	if err != nil {
		t.Fatalf("Render returned error: %v", err)
	}
	if !strings.Contains(html, "class=\"chroma\"") {
		t.Fatalf("render output missing chroma classes: %s", html)
	}
}
