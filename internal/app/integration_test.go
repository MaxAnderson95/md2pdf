package app

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maxanderson95/md2pdf/internal/browser"
	"github.com/maxanderson95/md2pdf/internal/version"
)

func TestRunIntegrationCreatesPDF(t *testing.T) {
	if os.Getenv("MD2PDF_INTEGRATION") == "" {
		t.Skip("set MD2PDF_INTEGRATION=1 to run browser-backed integration tests")
	}
	if _, err := browser.Find(); err != nil {
		t.Skipf("supported browser not available: %v", err)
	}

	inputBytes, err := os.ReadFile(filepath.Join("..", "..", "testdata", "markdown", "basic.md"))
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	dir := t.TempDir()
	inputPath := filepath.Join(dir, "basic.md")
	if err := os.WriteFile(inputPath, inputBytes, 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	err = Run(context.Background(), []string{inputPath}, &stdout, &stderr, version.Info())
	if err != nil {
		t.Fatalf("Run returned error: %v; stderr=%s", err, stderr.String())
	}

	outputPath := strings.TrimSpace(stdout.String())
	pdfBytes, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("ReadFile PDF failed: %v", err)
	}
	if !bytes.HasPrefix(pdfBytes, []byte("%PDF-")) {
		t.Fatalf("output is not a PDF: %q", pdfBytes[:min(len(pdfBytes), 16)])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
